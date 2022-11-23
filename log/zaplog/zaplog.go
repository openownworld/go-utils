//
// @Author: openownworld
// @Email:  openownworld@163.com
// @Date:   create on 2020/12/13 14:35
// @File:   zaplog.go
// @Description:

package zaplog

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	ini_utils "github.com/openownworld/go-utils/ini-utils"
	"github.com/openownworld/go-utils/log/errlog"
	"github.com/openownworld/go-utils/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig 封装高性能日志库zap
// 支持日志分级，按最大天数保存，最大文件限制，最大文件数限制
// 支持日志分级warn error能单独复制提取保存文件，支持动态设置日志分级，支持日志压缩
type LogConfig struct {
	//配置文件要通过tag来指定配置文件中的名称
	ServiceName        string `ini:"serviceName"`
	CustomTimeEnable   bool   `ini:"customTimeEnable"`
	LogFileName        string `ini:"logFileName"`
	ErrorFileName      string `ini:"errorFileName"`
	MaxSize            int    `ini:"maxSize"`
	MaxBackups         int    `ini:"maxBackups"`
	MaxDays            int    `ini:"maxDays"`
	Compress           bool   `ini:"compress"`
	Level              string `ini:"level"`
	StacktraceLevel    string `ini:"stacktraceLevel"`
	ErrorFileLevel     string `ini:"errorFileLevel"`
	LevelHttpEnable    bool   `ini:"levelHttpEnable"`
	LevelHttpApi       string `ini:"levelHttpApi"`
	LevelHttpPort      string `ini:"levelHttpPort"`
	SocketLoggerEnable bool   `ini:"socketLoggerEnable"`
	SocketLoggerJSON   bool   `ini:"socketLoggerJSON"`
	SocketType         string `ini:"socketType"`
	SocketIP           string `ini:"socketIP"`
	SocketPort         string `ini:"socketPort"`
	FileLogger         bool   `ini:"fileLogger"`
	ConsoleLogger      bool   `ini:"consoleLogger"`
	FileLoggerJSON     bool   `ini:"fileLoggerJSON"`
	ConsoleLoggerJSON  bool   `ini:"consoleLoggerJSON"`
}

var logConfig LogConfig
var connSocketLogger net.Conn

func GetConfig() *LogConfig {
	return &logConfig
}

var FirstInit bool = false

func init() {
	//确保日志最先初始化
	if FirstInit {
		logConfigFile := "log.ini"
		runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		InitLoggerFile(path.Join(runDir, logConfigFile))
	}
}

func InitConfigByIO(reader io.Reader, v interface{}) error {
	//读配置文件
	iniParser, err := ini_utils.NewIniParserByIO(reader)
	if err != nil {
		return err
	}
	if err := iniParser.GetFileHandle().MapTo(v); err != nil {
		errlog.WriteLogErrorf("MapTo config stream, error %s \n", err.Error())
		return err
	}
	return nil
}

func InitConfig(filePath string, v interface{}) error {
	//读配置文件
	iniParser, err := ini_utils.NewIniParser(filePath)
	if err != nil {
		return err
	}
	if err := iniParser.GetFileHandle().MapTo(v); err != nil {
		errlog.WriteLogErrorf("MapTo config file %s, error %s \n", filePath, err.Error())
		return err
	}
	return nil
}

// 支持文件大小限制，level，文件行号函数名。
// go get -u -v github.com/lestrrat-go/file-rotatelogs
var logger *zap.Logger

// var errorLogger *zap.SugaredLogger
var atomLevel zap.AtomicLevel

// InitLoggerFile 3个接口，1初始化日志库，通过配置文件
func InitLoggerFile(filePath string) error {
	if !utils.IsExist(filePath) {
		Error(filePath + " file path not exist")
		return errors.New(filePath + " file path not exist")
	}
	err := InitConfig(filePath, &logConfig)
	errlog.CheckError(err)
	if err == nil && logger == nil {
		logger = initLogger(logConfig)
	}
	return nil
}

func InitLoggerByIO(reader io.Reader) error {
	err := InitConfigByIO(reader, &logConfig)
	errlog.CheckError(err)
	if err == nil && logger == nil {
		logger = initLogger(logConfig)
	}
	return nil
}

// InitLoggerCfg 3个接口，1初始化日志库，通过结构体
func InitLoggerCfg(logConfig LogConfig) {
	if logger == nil {
		logger = initLogger(logConfig)
	}
}

// InitLoggerDefaultCfg 3个接口，1初始化日志库，通过默认结构体
func InitLoggerDefaultCfg() {
	logConfig = LogConfig{
		LogFileName:        "./log/log.log",
		ErrorFileName:      "./log/error.log",
		MaxSize:            50,
		MaxBackups:         30,
		MaxDays:            30,
		Compress:           true,
		Level:              "debug",
		StacktraceLevel:    "panic",
		ErrorFileLevel:     "warn",
		FileLogger:         true,
		ConsoleLogger:      true,
		FileLoggerJSON:     false,
		ConsoleLoggerJSON:  false,
		SocketLoggerEnable: false,
		SocketLoggerJSON:   false,
		SocketType:         "udp",
		SocketIP:           "127.0.0.1",
		SocketPort:         "9990",
		LevelHttpEnable:    false,
		LevelHttpApi:       "/api/log/level",
		LevelHttpPort:      "9090",
	}
	if logger == nil {
		logger = initLogger(logConfig)
	}
}

const (
	maxStack  = 20
	separator = "capture panic---------------------------------------"
)

// PrintPanicLog 打印堆栈信息
func PrintPanicLog() {
	if err := recover(); err != nil {
		t := time.Now().Format("2006-01-02 15:04:05.000") + " "
		str := fmt.Sprintf("\n%s%s start\nruntime error: %v\ntraceback:\n", t, separator, err)
		i := 2
		for {
			pc, file, line, ok := runtime.Caller(i)
			if !ok || i > maxStack {
				break
			}
			str += fmt.Sprintf("\tstack: %d %v [file: %s:%d] func: %s\n", i-1, ok, file, line, runtime.FuncForPC(pc).Name())
			i++
		}
		str += t + separator + " end\n" + string(debug.Stack())
		Error(str)
		//debug.PrintStack()
	}
}

// CloseConnSocketLogger 3个接口，2关闭Socket推送日志
func CloseConnSocketLogger() {
	if connSocketLogger != nil {
		connSocketLogger.Close()
	}
}

// SetLevel 3个接口，3设置日志级别
func SetLevel(level string) {
	if logger != nil {
		atomLevel.SetLevel(getLevel(level))
	}
}

func getLevel(loglevel string) zapcore.Level {
	var level zapcore.Level
	switch loglevel {
	case "debug":
		// DebugLevel logs are typically voluminous, and are usually disabled in production.
		level = zap.DebugLevel
	case "info":
		// InfoLevel is the default logging priority.
		level = zap.InfoLevel
	case "warn":
		// WarnLevel logs are more important than Info, but don't need individual human review.
		level = zap.WarnLevel
	case "error":
		// ErrorLevel logs are high-priority. If an application is running smoothly,
		// it shouldn't generate any error-level logs.
		level = zap.ErrorLevel
	case "dpanic":
		// DPanicLevel logs are particularly important errors. In development the
		// logger panics after writing the message.
		level = zap.DPanicLevel
	case "panic":
		// PanicLevel logs a message, then panics.
		level = zap.PanicLevel
	case "fatal":
		// FatalLevel logs a message, then calls os.Exit(1).
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	return level
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000000") + "]")
	//enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// logPath 日志文件路径
// loglevel 日志级别
func initLogger(logConfig LogConfig) *zap.Logger {
	//[1]文件log hook MaxBackups和MaxAge 任意达到限制，对应的文件就会被清理
	hookAll := lumberjack.Logger{
		Filename:   logConfig.LogFileName, // 日志文件路径 ./log/log.log
		MaxSize:    logConfig.MaxSize,     // 最大文件大小 M字节
		MaxBackups: logConfig.MaxBackups,  // 最多保留3个备份
		MaxAge:     logConfig.MaxDays,     // 文件最多保存多少天
		Compress:   logConfig.Compress,    // 是否压缩 disabled by default
	}
	hookError := lumberjack.Logger{
		Filename:   logConfig.ErrorFileName, // 日志文件路径 ./log/error.log
		MaxSize:    logConfig.MaxSize,       // 最大文件大小 M字节
		MaxBackups: logConfig.MaxBackups,    // 最多保留3个备份
		MaxAge:     logConfig.MaxDays,       // 文件最多保存多少天
		Compress:   logConfig.Compress,      // 是否压缩 disabled by default
	}
	//[2]设置level 动态level
	atomLevel = zap.NewAtomicLevel()
	atomLevel.SetLevel(getLevel(logConfig.Level))
	// 实现两个判断日志等级的interface,回调函数每次打印日志都会调用
	//allLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl < zapcore.WarnLevel
	//})
	//errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	//return lvl >= zapcore.WarnLevel
	//	return lvl >= zapcore.ErrorLevel
	//})
	var errorLevel zapcore.Level
	if logConfig.ErrorFileLevel == "error" {
		errorLevel = zapcore.ErrorLevel
	} else {
		errorLevel = zapcore.WarnLevel
	}
	//[3]日志输出配置
	//NewDevelopment和NewProduction区别 zap.NewDevelopment() 包含代码中文件信息
	//encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	// Keys can be anything except the empty string.
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "name"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stack"
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	//encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder // 小写编码器
	//encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if logConfig.CustomTimeEnable {
		encoderConfig.EncodeTime = customTimeEncoder //zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
	} else {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 UTC 时间格式
	}
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder //zapcore.FullCallerEncoder,      // 全路径编码器
	//[4]配置多个输出方式
	//TCP 没解决 当服务器断开，重连的问题
	//UDP 不存在这样的问题
	var socketCore zapcore.Core
	if logConfig.SocketLoggerEnable {
		conn, err := net.DialTimeout("udp", logConfig.SocketIP+":"+logConfig.SocketPort, 2*time.Second)
		//conn, err := net.DialTimeout(logConfig.SocketType, logConfig.SocketIP+":"+logConfig.SocketPort, 2*time.Second)
		if err != nil {
			errlog.WriteLogErrorf(logConfig.SocketType, logConfig.SocketIP+":"+logConfig.SocketPort, err.Error())
		} else {
			connSocketLogger = conn
			// read or write on conn
			//defer conn.Close()
			wSocket := zapcore.AddSync(conn)
			if logConfig.SocketLoggerJSON {
				socketEncoder := zapcore.NewJSONEncoder(encoderConfig)
				socketCore = zapcore.NewCore(socketEncoder, wSocket, atomLevel)
			} else {
				socketEncoder := zapcore.NewConsoleEncoder(encoderConfig)
				socketCore = zapcore.NewCore(socketEncoder, wSocket, atomLevel)
			}
		}
	}
	//
	// Assume that we have clients for two Kafka topics. The clients implement
	// zapcore.WriteSyncer and are safe for concurrent use. (If they only
	// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
	// method. If they're not safe for concurrent use, we can add a protecting
	// mutex with zapcore.Lock.)
	//topicDebugging := zapcore.AddSync(ioutil.Discard)
	//topicErrors := zapcore.AddSync(ioutil.Discard)
	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	var consoleWriter zapcore.WriteSyncer
	var allWriter zapcore.WriteSyncer
	var errorWriter zapcore.WriteSyncer
	//consoleWriter := zapcore.Lock(os.Stdout)
	//consoleErrors := zapcore.Lock(os.Stderr)
	if logConfig.ConsoleLogger {
		consoleWriter = zapcore.Lock(os.Stdout)
		//consoleWriter = zapcore.AddSync(os.Stdout)
	} else {
		consoleWriter = zapcore.AddSync(ioutil.Discard)
	}
	if logConfig.FileLogger {
		allWriter = zapcore.AddSync(&hookAll)
		errorWriter = zapcore.AddSync(&hookError)
	} else {
		allWriter = zapcore.AddSync(ioutil.Discard)
		errorWriter = zapcore.AddSync(ioutil.Discard)
	}
	// Optimize the Kafka output for machine consumption and the console output for human operators.
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	if logConfig.FileLoggerJSON {
		fileEncoder = jsonEncoder
	}
	//终端支持彩色打印
	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	if logConfig.ConsoleLoggerJSON {
		consoleEncoder = jsonEncoder
	}
	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	var core zapcore.Core
	if socketCore != nil {
		core = zapcore.NewTee(
			socketCore,
			zapcore.NewCore(fileEncoder, allWriter, atomLevel),
			zapcore.NewCore(fileEncoder, errorWriter, errorLevel),
			zapcore.NewCore(consoleEncoder, consoleWriter, atomLevel),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, allWriter, atomLevel),
			zapcore.NewCore(fileEncoder, errorWriter, errorLevel),
			zapcore.NewCore(consoleEncoder, consoleWriter, atomLevel),
		)
	}
	// From a zapcore.Core, it's easy to construct a Logger.
	//[5]创建日志logger
	// 开启开发模式，堆栈跟踪
	//caller := zap.AddCaller()
	// 开启文件及行号
	//development := zap.Development()
	// 设置初始化字段 service key
	filed := zap.Fields()
	if len(logConfig.ServiceName) > 0 {
		filed = zap.Fields(zap.String("service", logConfig.ServiceName))
	}
	logger := zap.New(core, zap.Development(), zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(getLevel(logConfig.StacktraceLevel)), filed)
	//logger := zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(getLevel(logConfig.StacktraceLevel)))
	defer logger.Sync()
	//
	//logger.Info("")
	//logger.Error("")
	logger.Info("DefaultLogger init success [" + logConfig.Level + "]")
	logger.Error("DefaultLogger init success [" + logConfig.Level + "]")
	//logger.Debug("Debug")
	//logger.Warn("Warn")
	//logger.Error("Error")
	//logger.Panic("Panic")
	//logger.Fatal("Fatal")
	//[6]支持 restful api修改level  http://localhost:9090/api/log/level
	//curl -H "Content-Type:application/json" -X PUT --data "{\"level\":\"error\"}" http://localhost:9090/api/log/level
	if logConfig.LevelHttpEnable {
		levelHttp := logConfig.LevelHttpApi
		levelHttpPort := ":" + logConfig.LevelHttpPort
		http.HandleFunc(levelHttp, atomLevel.ServeHTTP)
		go func() {
			if err := http.ListenAndServe(levelHttpPort, nil); err != nil {
				//panic(err)
				logger.Error("DefaultLogger init success", zap.String("err", err.Error()))
			}
		}()
	}
	return logger
}

//Sprint采用默认格式将其参数格式化，串联所有输出生成并返回一个字符串。如果两个相邻的参数都不是字符串，会在它们的输出之间添加空格
//Sprintln采用默认格式将其参数格式化，串联所有输出生成并返回一个字符串。总是会在相邻参数的输出之间添加空格并在输出结束后添加换行符。

func GetNowTimeMs() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

// GetCallerInfo 0为当前栈 1上层调用栈
func GetCallerInfo(stackNum int) string {
	//0为当前栈 1上层调用栈
	//pc, file, line, ok := runtime.Caller(stackNum)
	_, file, line, ok := runtime.Caller(stackNum)
	if !ok {
		file = "runtime call err"
	}
	//pcName := runtime.FuncForPC(pc).Name() //获取函数名，听说很耗时
	//where:= fmt.Sprintf("%v %s %d %t %s",pc,file,line,ok,pcName)
	//where := fmt.Sprintf("%s %d %s", file, line, pcName)
	where := fmt.Sprintf("%s:%d", file, line) //冒号拼接，goland可以直接点开到文件行数
	return where
}

// Println 打印日志到终端
func Println(args ...interface{}) {
	fmt.Printf(fmt.Sprintf("%s %s %s %s", GetNowTimeMs(), "console", GetCallerInfo(2), fmt.Sprintln(args...)))
}

// Printfln 打印日志到终端
func Printfln(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf("%s %s %s %s", GetNowTimeMs(), "console", GetCallerInfo(2), fmt.Sprintf(format, args...)))
}

// Printf 打印日志到终端
func Printf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf("%s %s %s %s", GetNowTimeMs(), "console", GetCallerInfo(2), fmt.Sprintf(format, args...)))
}

func printlnLog(level string, args ...interface{}) {
	//fmt.Println(fmt.Sprintf("%s %s %s %s", GetNowTimeMs(), level, GetCallerInfo(3), fmt.Sprint(args...)))
	fmt.Print(fmt.Sprintf("%s %s %s %s", GetNowTimeMs(), level, GetCallerInfo(3), fmt.Sprintln(args...)))
}

// Debug logs a message at level Debug on the compatibleLogger.
func Debug(args ...interface{}) {
	if logger != nil {
		logger.Debug(strings.TrimRight(fmt.Sprintln(args...), "\n")) //Sprintln带了换行符，TrimRight从右边去掉“\n”换行符
		//logger.Debug(fmt.Sprint(args...)) //zap 默认带换行符，Sprint不带空格拼接，Sprintln带空格拼接
	} else {
		printlnLog("c-debug", args...)
	}
}

// Debugf logs a message at level Debug on the compatibleLogger.
func Debugf(format string, args ...interface{}) {
	if logger != nil {
		logger.Debug(fmt.Sprintf(format, args...))
	} else {
		printlnLog("c-debug", fmt.Sprintf(format, args...))
	}
}

// Info logs a message at level Info on the compatibleLogger.
func Info(args ...interface{}) {
	if logger != nil {
		logger.Info(strings.TrimRight(fmt.Sprintln(args...), "\n"))
	} else {
		printlnLog("c-info", args...)
	}
}

// Infof logs a message at level Info on the compatibleLogger.
func Infof(format string, args ...interface{}) {
	if logger != nil {
		logger.Info(fmt.Sprintf(format, args...))
	} else {
		printlnLog("c-info", fmt.Sprintf(format, args...))
	}
}

// Warn logs a message at level Warn on the compatibleLogger.
func Warn(args ...interface{}) {
	if logger != nil {
		logger.Warn(strings.TrimRight(fmt.Sprintln(args...), "\n"))
	} else {
		printlnLog("c-warn", args...)
	}
}

// Warnf logs a message at level Warn on the compatibleLogger.
func Warnf(format string, args ...interface{}) {
	if logger != nil {
		logger.Warn(fmt.Sprintf(format, args...))
	} else {
		printlnLog("c-warn", fmt.Sprintf(format, args...))
	}
}

// Error logs a message at level Error on the compatibleLogger.
func Error(args ...interface{}) {
	if logger != nil {
		logger.Error(strings.TrimRight(fmt.Sprintln(args...), "\n"))
	} else {
		printlnLog("c-error", args...)
	}
}

// Errorf logs a message at level Error on the compatibleLogger.
func Errorf(format string, args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprintf(format, args...))
	} else {
		printlnLog("c-error", fmt.Sprintf(format, args...))
	}
}

// Fatal logs a message at level Fatal on the compatibleLogger.
func Fatal(args ...interface{}) {
	if logger != nil {
		logger.Fatal(strings.TrimRight(fmt.Sprintln(args...), "\n"))
	} else {
		printlnLog("c-fatal", args...)
		os.Exit(1)
	}
}

// Fatalf logs a message at level Fatal on the compatibleLogger. followed by a call to os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	if logger != nil {
		logger.Fatal(fmt.Sprintf(format, args...))
	} else {
		printlnLog("c-fatal", fmt.Sprintf(format, args...))
		os.Exit(1)
	}
}

// Panic logs a message at level Painc on the compatibleLogger.  followed by a call to panic().
func Panic(args ...interface{}) {
	if logger != nil {
		logger.Panic(strings.TrimRight(fmt.Sprintln(args...), "\n"))
	} else {
		printlnLog("c-panic", args...)
		panic(fmt.Sprint(args...))
	}
}

// Panicf logs a message at level Painc on the compatibleLogger.
func Panicf(format string, args ...interface{}) {
	if logger != nil {
		logger.Panic(fmt.Sprintf(format, args...))
	} else {
		printlnLog("c-panic", fmt.Sprintf(format, args...))
		panic(fmt.Sprint(args...))
	}
}

// With return a logger with an extra field.
func With(key string, value interface{}) {
	logger.With(zap.Any(key, value))
}

// WithField return a logger with an extra field.
func WithField(key string, value interface{}) {
	logger.With(zap.Any(key, value))
}

// WithFields return a logger with extra fields.
func WithFields(fields map[string]interface{}) {
	first := make([]zap.Field, len(fields))
	i := 0
	for k, v := range fields {
		first[i] = zap.Any(k, v)
		i++
	}
	logger.With(first...)
}

// DebugWithField logs a message at level Info on the compatibleLogger.
func DebugWithField(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Debug(msg, fields...)
	} else {
		printlnLog("c-debug", msg)
	}
}

// InfoWithField logs a message at level Info on the compatibleLogger.
func InfoWithField(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Info(msg, fields...)
	} else {
		printlnLog("c-info", msg)
	}
}

// WarnWithField logs a message at level Info on the compatibleLogger.
func WarnWithField(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Warn(msg, fields...)
	} else {
		printlnLog("c-warn", msg)
	}
}

// ErrorWithField logs a message at level Info on the compatibleLogger.
func ErrorWithField(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Error(msg, fields...)
	} else {
		printlnLog("c-error", msg)
	}
}
