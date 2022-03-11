package utils

import (
	"fmt"
	"io/ioutil"
	"os"
)

// IsExist returns whether the given file or directory exists or not
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// no such file or dir
		return false
	}
	if info.IsDir() {
		// it's a directory
		return true
	} else {
		// it's a file
		return false
	}
}

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func WriteFile(path string, buf []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Println("error write permission denied")
		}
		if os.IsNotExist(err) {
			fmt.Println("file does not exist")
		}
		return err
	}
	defer f.Close()
	return ioutil.WriteFile(path, buf, 0666)
}
