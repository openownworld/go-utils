//
// @Author: openownworld
// @Email:  openownworld@163.com
// @Date:   create on 2020/12/13 15:17
// @File:   logrus_test.go
// @Description:

package logrus

import "testing"

func TestLog(t *testing.T) {
	// /Users/start/code-my/go-utils/log/logrus/logrus.20211213-0000.log
	InitLogger("logrus", 10, "debug")
	for i := 0; i < 10; i++ {
		Log.Info("InitLogger ok")
	}
}
