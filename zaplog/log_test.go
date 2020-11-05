package zaplog

import (
	"go.lbj.pkg/go-utils/tools"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestDemo(t *testing.T) {
	defer PrintPanicLog()
	//panic("ppp")
	t.Log(GetCallerInfo(1))
	Println("ttttttttttt")
	Printfln("%d", 111111)
	confFile := "logConfig.ini"
	InitLoggerFile(tools.GetAbsolutePath(confFile))
	for i := 0; i < 1; i++ {
		//time.Sleep(1 * time.Second)
		time.Sleep(10 * time.Millisecond)
		Debug("debug log", zap.String("s", "ssss"), zap.Int("counter", i))
		Debugf("%d", 1000000)
		Info("info log", zap.String("s", "ssss"), zap.Int("counter", i))
		Error("error log", zap.String("s", "ssss"), zap.Int("counter", i))
		Warn("warn log", zap.String("s", "ssss"), zap.Int("counter", i))
		//Panic("panic log", zap.String("s", "ssss"), zap.Int("counter", i))
	}
}

func TestDefaultCfg(t *testing.T) {
	InitLoggerDefaultCfg()
	Debug("run go, panic---------------------")
}

func TestLoopLog(t *testing.T) {
	//runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//confFile := runDir + "\\logConfig.ini"
	//confFile := "logConfig.ini"
	//InitLoggerFile(confFile)
	for n := 0; n < 0; n++ {
		for i := 0; i < 1; i++ {
			//time.Sleep(10 * time.Second)
			Debug("debug log", zap.String("test", "ssss"), zap.Int("counter", i))
			Info("info log", zap.String("test", "ssss"), zap.Int("counter", i))
			//logger.Error("error log", zap.String("test", "ssss"), zap.Int("counter", i))
			Warn("warn log", zap.String("test", "ssss"), zap.Int("counter", i))
			Error("err log", zap.String("test", "ssss"), zap.Int("counter", i))
			//SetLevel("error")
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
