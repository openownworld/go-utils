//
// @Author: openownworld
// @Email:  openownworld@163.com
// @Date:   create on 2020/12/13 15:17
// @File:   ini_test.go
// @Description:

package ini_utils

import (
	"fmt"
	"testing"
)

type MySQL struct {
	Url string `ini:"url"`
}

// 字段必须大写
type IniInfo struct {
	LogFileName   string `ini:"logFileName"`
	ErrorFileName string `ini:"errorFileName"`
	MySQL         MySQL  `ini:"mysql"`
}

func TestINI(t *testing.T) {
	//读配置文件
	v := &IniInfo{} // a pointer
	filePath := "./test.ini"
	iniParser, err := NewIniParser(filePath)
	if err != nil {
		t.Error(err)
	}
	if err := iniParser.GetFileHandle().MapTo(v); err != nil {
		t.Error(err)
	}
	fmt.Printf("%v\n", v)
	fmt.Printf("%+v\n", v)
	fmt.Printf("%+v\n", *v)
	m := iniParser.GetValueMap()
	fmt.Printf("%+v\n", m)
	fmt.Println(m)
	t.Log(v)
}
