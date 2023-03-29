package fileio

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func CheckDirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// ReadFile 读取文件
func ReadFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte(""), err
	}
	return data, nil
}

// WriteFile 保存文件（没有就创建，删除并覆盖）
func WriteFile(path string, fileName string, file multipart.File) (filepath string, err error) {
	filepath = path + "/" + fileName
	err = os.MkdirAll(path, 0777)
	if err != nil {
		return "", err
	}
	out, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return "", err
	}
	defer func(out *os.File) {
		err1 := out.Close()
		if err1 != nil {
			err = err1
			return
		}
	}(out)
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}
	return filepath, nil
}

func ConvertToUtf8(in string) string {
	s := []byte(in)
	reg := regexp.MustCompile(`\\[0-7]{3}`)
	out := reg.ReplaceAllFunc(s,
		func(b []byte) []byte {
			i, _ := strconv.ParseInt(string(b[1:]), 8, 0)
			return []byte{byte(i)}
		})
	return string(out)
}

func BadStrToUtf8(input string) string {
	reg, _ := regexp.Compile("\\\\u\\w{4}")
	return reg.ReplaceAllStringFunc(input, func(input string) string {
		replaceU := strings.Replace(input, "\\u", "", -1)
		tmp, _ := strconv.ParseInt(replaceU, 16, 32)
		return fmt.Sprintf("%s", string(rune(tmp)))
	})
}
