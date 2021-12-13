//
// @Author: openownworld
// @Email:  openownworld@163.com
// @Date:   create on 2020/12/13 14:35
// @File:   log.go
// @Description:

package logrus

import (
	"bufio"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	_ "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// Logrus不提供日志轮换。日志轮换应由logrotate(8)可以压缩和删除旧日志条目的外部程序（如）完成。它不应该是应用程序级记录器的功能。
// go get -u -v github.com/lestrrat-go/file-rotatelogs

var Log *logrus.Logger

// InitLogger 路径，日志最大文件大小M
func InitLogger(path string, maxFileSize int64, logLevel string) {
	// baseLogPath := path.Join("/Users/start/code-my/go-utils/log/", "logrus")
	// 默认10个文件
	maxFileSize = maxFileSize * 1024 * 1024
	Log = newLogger(path, maxFileSize, 10, logLevel)
	Log.Info("InitLogger ok")
}

func newLogger(path string, maxFileSize int64, maxFileCount uint, logLevel string) *logrus.Logger {
	writer, err := rotatelogs.New(
		path+".%Y%m%d-%H%M.log",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		//rotatelogs.WithLinkName(path),
		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(24*time.Hour),
		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationCount(maxFileCount),
		rotatelogs.WithRotationSize(maxFileSize),
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}
	/*
			//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
		    logrus.SetFormatter(&logrus.JSONFormatter{})
		    //设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
		    logrus.SetOutput(os.Stdout)
		    //设置最低loglevel
			logrus.SetLevel(logrus.InfoLevel)
	*/
	logrus.SetFormatter(&logrus.TextFormatter{})
	/*
	   如果日志级别不是debug就不要打印日志到控制台了
	*/
	switch level := logLevel; level {

	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stderr)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true})
	logrus.AddHook(lfsHook)
	return logrus.StandardLogger()
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	logrus.SetOutput(writer)
}
