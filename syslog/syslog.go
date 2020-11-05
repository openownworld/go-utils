package syslog

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	_ "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

//go get -u -v github.com/lestrrat-go/file-rotatelogs

var Log *logrus.Logger

func InitLogger() {
	baseLogPath := path.Join("e://", "syslog")
	Log = NewLogger(baseLogPath, 10, "debug")
	Log.Info("InitLogger ok")
}

func NewLogger(path string, maxFileSize uint, logLevel string) *logrus.Logger {

	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M.log",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		//rotatelogs.WithLinkName(path),
		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(24*time.Hour),
		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationCount(maxFileSize),
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
