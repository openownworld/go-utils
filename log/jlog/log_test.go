//
// @Author: openownworld
// @Email:  openownworld@163.com
// @Date:   create on 2020/12/13 15:17
// @File:   log_test.go
// @Description:

package jlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	InitLog("./log.log", DebugLevel)
	Debug("abc")
	Error("abc")
}
