//
// @Author: openownworld
// @Email:  openownworld@163.com
// @Date:   create on 2021/12/13 17:12
// @File:   errlog.go
// @Description:

package errlog

import (
	"fmt"
	"github.com/openownworld/go-utils/utils"
	"os"
	"time"
)

var appErrorLogFile = "./app_error.log"

const AppErrorLogFileSize = 20 * 1024 * 1024

//WriteLogError 只有日志模块的初始话，才使用次函数，其他打印日志，走日志模块
func WriteLogError(v ...interface{}) error {
	fileSize := utils.GetFileSizeBySeek(appErrorLogFile)
	if fileSize > AppErrorLogFileSize {
		err := os.Remove(appErrorLogFile)
		if err != nil {
			fmt.Println("WriteLogError", err.Error())
		}
	}
	currentTime := time.Now().Local()
	timeString := currentTime.Format("2006-01-02 15:04:05.000")
	//以读写方式打开文件，如果不存在，则创建
	file, err := os.OpenFile(appErrorLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Println(err)
		return err
	}
	log := timeString + " - " + utils.GetUpperFuncName(2) + " - " + fmt.Sprintln(v...)
	fmt.Println(log)
	file.WriteString(log)
	file.Close()
	return nil
}

//WriteLogErrorf 只有日志模块的初始话，才使用次函数，其他打印日志，走日志模块
func WriteLogErrorf(format string, a ...interface{}) error {
	fileSize := utils.GetFileSizeBySeek(appErrorLogFile)
	if fileSize > AppErrorLogFileSize {
		err := os.Remove(appErrorLogFile)
		if err != nil {
			fmt.Println("WriteLogErrorf", err.Error())
		}
	}
	currentTime := time.Now().Local()
	timeString := currentTime.Format("2006-01-02 15:04:05.000")
	//以读写方式打开文件，如果不存在，则创建
	file, err := os.OpenFile(appErrorLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Println(err)
		return err
	}
	log := timeString + " - " + utils.GetUpperFuncName(2) + " - " + fmt.Sprintf(format, a...)
	fmt.Println(log)
	file.WriteString(log)
	file.Close()
	return nil
}

//CheckErrorToExit 检查错误，不正常则退出
func CheckErrorToExit(err error) {
	if err != nil {
		WriteLogError("system exit,", err.Error())
		os.Exit(1)
	}
}

//CheckError 检查错误，不正常则退出
func CheckError(err error) {
	if err != nil {
		WriteLogError("error------------,a serious error has occurred\n\n", err.Error())
	}
}
