package logdecode

import (
	"bytes"
	"clam-server/jwt"
	"clam-server/utils/cmd"
	"clam-server/utils/fileio"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

type cmdData struct {
	cmd      *exec.Cmd
	out      *bytes.Buffer
	err      *bytes.Buffer
	filePath string
}

const scriptsPath = "../scripts/decode.sh"
const temporaryFilePath = "../resources/"

var token2CmdDataMap sync.Map

func Router(r *gin.Engine) {
	r.POST("/util/log-decode", logDecode)
	r.GET("/util/log-decode/get-res", getFileRes)
}

func logDecode(c *gin.Context) {
	uid, _ := jwts.ParseToken(c)
	_, ok := token2CmdDataMap.Load(uid)
	// 正在执行
	if ok {
		c.JSON(http.StatusProxyAuthRequired, gin.H{
			"message": "脚本正在执行或文件未取走",
		})
		return
	}
	xLogHandler, err := c.FormFile("xlog")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "未找到请求文件",
		})
		return
	}
	if xLogHandler.Filename[len(xLogHandler.Filename)-5:] != ".xlog" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "上传的文件必须以‘.xlog’结尾",
		})
		return
	}
	xLog, err := xLogHandler.Open()
	if err != nil {
		c.JSON(http.StatusProxyAuthRequired, gin.H{
			"message": err,
		})
		return
	}
	filePath, err := fileio.WriteFile(genTemporaryFilePath(uid), xLogHandler.Filename, xLog)
	if err != nil {
		c.JSON(http.StatusProxyAuthRequired, gin.H{
			"message": err,
		})
		return
	}
	scripts, err := fileio.ReadFile(scriptsPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	d := cmd.GenCommand(string(scripts) + filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	var stdout, stderr bytes.Buffer
	d.Stdout = &stdout
	d.Stderr = &stderr
	err = d.Start()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	token2CmdDataMap.Store(uid, cmdData{d, &stdout, &stderr, filePath})
	c.JSON(http.StatusOK, gin.H{
		"message": "成功，请尝试取走文件",
	})
}

func genTemporaryFilePath(token string) string {
	return temporaryFilePath + token
}

func getFileRes(c *gin.Context) {
	uid, _ := jwts.ParseToken(c)
	d, exit := token2CmdDataMap.Load(uid)
	if !exit {
		c.JSON(http.StatusProxyAuthRequired, gin.H{
			"message": "没有文件正在解码",
		})
		return
	}
	data := d.(cmdData)
	// if (data.out != nil && data.out.String() == "") || (data.err == nil && data.err.String() == "") {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": cmd.ConvertByte2String(data.out.Bytes(), cmd.GB18030),
	// 		"err":     cmd.ConvertByte2String(data.err.Bytes(), cmd.GB18030),
	// 	})
	// 	return
	// }
	_, err := os.Stat(data.filePath + ".log")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.File(data.filePath + ".log")
	token2CmdDataMap.Delete(uid)
}
