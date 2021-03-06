//
// @Author: openownworld
// @Email:  openownworld@163.com
// @Date:   create on 2020/12/13 14:35
// @File:   log.go
// @Description:

package jlog

import (
	"fmt"
	"github.com/openownworld/go-utils/jerror"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	PanicLevel int = iota // ( = iota )0      (= iota < 1) 1
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

type LogFile struct {
	level    int
	logTime  int64
	fileName string
	fileFd   *os.File
}

var logFile LogFile

// Write 写文件 log.SetOutput(&logFile) 才会调用
func (l *LogFile) Write(buf []byte) (n int, err error) {
	fmt.Printf("%s", buf)
	if l.fileFd != nil {
		return l.fileFd.Write(buf)
	}
	return len(buf), nil
}

// Deprecated: createLogFile 创建文件
// 1 ./log/applog.log
// 2	/log/applog.log
// 3	./applog.log
// 4	applog.log
func (l *LogFile) createLogFile() {
	//创建目录
	if index := strings.LastIndex(l.fileName, "/"); index != -1 {
		logDir := l.fileName[0:index] + "/"
		os.MkdirAll(logDir, os.ModePerm)
	}
	//
	now := time.Now()
	fileName := fmt.Sprintf("%s_%04d%02d%02d_%02d%02d", l.fileName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	if err := os.Rename(l.fileName, fileName); err == nil {
		//压缩tar
		/*go func() {
			tarCmd := exec.Command("tar", "-zcf", fileName+".tar.gz", fileName, "--remove-files")
			tarCmd.Run()
		}()*/
	}
	for index := 0; index < 3; index++ {
		if fd, err := os.OpenFile(l.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive); nil == err {
			l.fileFd.Sync()
			l.fileFd.Close()
			l.fileFd = fd
			break
		} else {
			fmt.Println("open file failed:", l.fileName)
		}
		l.fileFd = nil
	}
}

// InitLog 初始化配置日志参数
func InitLog(fileName string, level int) {
	logFile.fileName = fileName
	logFile.level = level
	log.SetOutput(&logFile)
	//log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	//log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}

// SetLevel 改变level
func SetLevel(level int) {
	logFile.level = level
}

// getUpperFuncName 获取上层调用者的函数名
func getUpperFuncName() string {
	/*funcName, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Println("Func Name=" + runtime.FuncForPC(funcName).Name())
		fmt.Printf("file: %s    line=%d\n", file, line)
	}*/
	//0为当前栈 1上层调用栈
	pc, _, _, ok := runtime.Caller(2)
	if ok {
		return runtime.FuncForPC(pc).Name()
	}
	return "null"
}

// Debugf 打印
func Debugf(format string, a ...interface{}) {
	if logFile.level >= DebugLevel {
		log.SetPrefix("debug ")
		log.Println(getUpperFuncName() + " - " + fmt.Sprintf(format, a...))
	}
}

// Debug 打印
func Debug(v ...interface{}) {
	if logFile.level >= DebugLevel {
		log.SetPrefix("debug ")
		log.Println(getUpperFuncName() + " - " + fmt.Sprintln(v...))
	}
}

// Infof 打印
func Infof(format string, a ...interface{}) {
	if logFile.level >= InfoLevel {
		log.SetPrefix("info  ")
		//log.Println(fmt.Sprintf(format, a))
		log.Println(getUpperFuncName() + " - " + fmt.Sprintf(format, a...))
	}
}

// Info 打印
func Info(v ...interface{}) {
	if logFile.level >= InfoLevel {
		log.SetPrefix("info  ")
		//log.Println(v...)
		log.Println(getUpperFuncName() + " - " + fmt.Sprintln(v...))
	}
}

// Warn 打印
func Warn(v ...interface{}) {
	if logFile.level >= WarnLevel {
		log.SetPrefix("warn  ")
		//log.Println(v...)
		log.Println(getUpperFuncName() + " - " + fmt.Sprintln(v...))
	}
}

// Errorf 打印
func Errorf(format string, a ...interface{}) {
	if logFile.level >= ErrorLevel {
		log.SetPrefix("error ")
		log.Println(getUpperFuncName() + " - " + fmt.Sprintf(format, a...))
	}
}

// Error 打印
func Error(v ...interface{}) {
	if logFile.level >= ErrorLevel {
		log.SetPrefix("error ")
		log.Println(getUpperFuncName() + " - " + fmt.Sprintln(v...))
	}
}

// Fatal 打印,会退出
func Fatal(v ...interface{}) {
	if logFile.level >= FatalLevel {
		log.SetPrefix("fatal ")
		//log.Fatalln(v...)
		log.Println(getUpperFuncName() + " - " + fmt.Sprintln(v...))
		os.Exit(1)
	}
}

// PrintError 日志打印
func PrintError(err error) {
	if logFile.level >= ErrorLevel {
		log.SetPrefix("error ")
		log.Println(getUpperFuncName() + " - " + fmt.Sprintln(err.Error()))
	}
}

// PrintJError 日志打印
func PrintJError(err jerror.Error) {
	if logFile.level >= ErrorLevel {
		log.SetPrefix("error ")
		log.Println(getUpperFuncName()+" - ", err.Code(), fmt.Sprintln(err.Error()))
	}
}
