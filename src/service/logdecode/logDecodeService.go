package logdecode

import (
	"bytes"
	"clam-server/component/jwt"
	"clam-server/config"
	"clam-server/serverlogger"
	"clam-server/utils/cmd"
	"clam-server/utils/fileio"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

type cmdData struct {
	cmd        *exec.Cmd
	out        *bytes.Buffer
	err        *bytes.Buffer
	filePath   string
	createTime time.Time
}

var (
	token2CmdDataMap sync.Map
)

func Router(r *gin.Engine) {
	r.LoadHTMLFiles("../../web/index.html")
	r.POST("/util/log-decode", logDecode)
	r.GET("/util/log-decode/get-res", getFileRes)
	r.GET("/main", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Center",
		})
	})
	r.Static("./web", "../../web")
	go deleteTemporaryFile()
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
	scripts, err := fileio.ReadFile(config.GetConfig().Decoder.ScriptsPath)
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
	token2CmdDataMap.Store(uid, cmdData{d, &stdout, &stderr, filePath, time.Now()})
	token, _ := c.Get(config.GetConfig().Jwt.DefaultContextKey)
	jwtStr, _ := jwts.TokenToString(token.(*jwt.Token))
	c.JSON(http.StatusOK, gin.H{
		"message": "成功，请尝试取走文件",
		"jwt":     jwtStr,
	})
}

func genTemporaryFilePath(token string) string {
	return config.GetConfig().Decoder.TemporaryFilePath + token
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

func deleteTemporaryFile() {
	for {
		serverlogger.Warn("删除文件任务开始执行")
		now := time.Now()
		dirs, err := os.ReadDir(config.GetConfig().Decoder.TemporaryFilePath)
		if err != nil {
			serverlogger.Warn("删除文件失败", zap.Error(err))
		}
		for _, entry := range dirs {
			data, ok := token2CmdDataMap.Load(entry.Name())
			if !ok {
				err := os.RemoveAll(genTemporaryFilePath(entry.Name()))
				if err != nil {
					serverlogger.Warn("删除文件失败", zap.Error(err))
				}
			} else {
				if data.(cmdData).createTime.Add(time.Duration(config.GetConfig().Decoder.FileTimeOut) * time.Minute).Before(now) {
					err := os.RemoveAll(genTemporaryFilePath(entry.Name()))
					if err != nil {
						serverlogger.Warn("删除文件失败", zap.Error(err))
					}
					token2CmdDataMap.Delete(entry.Name())
				}
			}
		}
		serverlogger.Warn("删除文件任务执行结束")
		t := time.NewTimer(time.Duration(config.GetConfig().Decoder.DeleteFilePeriod) * time.Hour)
		<-t.C
	}
}
