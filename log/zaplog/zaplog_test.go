package zaplog

import (
	"encoding/json"
	"fmt"
	"github.com/openownworld/go-utils/utils"
	"go.uber.org/zap/zapcore"
	"log"
	"path"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	logConfigFile := "logConfig.ini"
	InitLoggerFile(path.Join(utils.GetCurrentFilePath(), logConfigFile))
	defer PrintPanicLog()
	//panic("ppp")
	t.Log(GetCallerInfo(1))
	for i := 0; i < 1; i++ {
		//time.Sleep(1 * time.Second)
		time.Sleep(10 * time.Millisecond)
		Debug("debug log", zap.String("s", "ssss"), zap.Int("counter", i))
		Debugf("%d", 1000000)
		Info("info log", zap.String("s", "ssss"), zap.Int("counter", i))
		Error("error log", zap.String("s", "ssss"), zap.Int("counter", i))
		Warn("warn log", zap.String("s", "ssss"), zap.Int("counter", i))
	}
	Panic("panic log", zap.String("s", "ssss"))
}

func TestDefaultCfg(t *testing.T) {
	InitLoggerDefaultCfg()
	Debug("run go, panic---------------------")
}

func TestZaplog(t *testing.T) {
	FirstInit = false
	fmt.Println("GetCurrentFilePath", utils.GetCurrentFilePath())
	lp := "log"
	lv := "debug"
	isDebug := true
	var js string
	if isDebug {
		js = fmt.Sprintf(`{
	      "level": "%s",
	      "encoding": "json",
	      "outputPaths": ["stdout"],
	      "errorOutputPaths": ["stdout"]
	      }`, lv)
	} else {
		js = fmt.Sprintf(`{
	      "level": "%s",
	      "encoding": "json",
	      "outputPaths": ["%s"],
	      "errorOutputPaths": ["%s"]
	      }`, lv, lp, lp)
	}
	var cfg zap.Config
	if err := json.Unmarshal([]byte(js), &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	logger, err = cfg.Build()
	if err != nil {
		log.Fatal("init logger error: ", err)
	}
	logger.Info("abc")
}
