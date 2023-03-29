package cmd

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os/exec"
	"runtime"
)

type Charset string

const (
	windows = "windows"
	linux   = "linux"
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

// GenCommand 执行命令，不阻塞/**
func GenCommand(cmd string) *exec.Cmd {
	var c *exec.Cmd
	if runtime.GOOS == linux {
		c = exec.Command("bash", "-c", cmd)
	} else if runtime.GOOS == windows {
		c = exec.Command("cmd", "/C", cmd)
	}
	return c
}

// RunCommand 执行命令，阻塞/**
func RunCommand(cmd string) (*exec.Cmd, error) {
	var c *exec.Cmd
	if runtime.GOOS == linux {
		c = exec.Command("bash", "-c", cmd)
	} else if runtime.GOOS == windows {
		c = exec.Command("cmd", "/C", cmd)
	}
	return c, c.Run()
}
func ReadStdout(stdout io.ReadCloser) string {
	var outputBuf0 bytes.Buffer
	for {
		tempOutPut := make([]byte, 256)
		n, err := stdout.Read(tempOutPut)
		if err != nil {
			if err == io.EOF { // 读取到内容的最后位置
				break
			} else {
				fmt.Println(err)
				return ""
			}
		}
		if n > 0 {
			outputBuf0.Write(tempOutPut[:n])
		}
	}
	return ConvertByte2String(outputBuf0.Bytes(), GB18030)
}

func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}
