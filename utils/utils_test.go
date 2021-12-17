package utils

import (
	"github.com/openownworld/go-utils/utils/file"
	"testing"
)

func TestDemo(t *testing.T) {
	t.Log("测试是否存在----------", file.IsExist("E:/code3"))
	t.Log("测试是否文件夹----------", file.IsDir("E:/code3"))
	//t.Error("-----") //打印错误，测试就会FAIL
}
