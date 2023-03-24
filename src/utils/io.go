package utils

import "os"

func CheckDirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}
