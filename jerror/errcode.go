package jerror

import "fmt"

//错误码定义  程序内部错误相关
const (
	Succeed = iota * (-1) // (0 * -1)
	Failed
	OperationErr
	InvalidArgErr
	ServerInternalErr
	JsonEncodeErr
	JsonDecodeErr
	RandNumberErr
	KeyNotExistErr
	FilePathNotExistErr
)

var errTextMap = map[int]string{
	Succeed:             "succeed",
	Failed:              "failed",
	OperationErr:        "operation err",
	InvalidArgErr:       "invalid arg err",
	ServerInternalErr:   "server internal err",
	JsonEncodeErr:       "json encode err",
	JsonDecodeErr:       "json decode err",
	RandNumberErr:       "generate rand number err",
	KeyNotExistErr:      "key not exist err",
	FilePathNotExistErr: "file path not exist err",
}

// GetErrCodeText 获取错误码文本描述
func GetErrCodeText(code int) string {
	if v, ok := errTextMap[code]; ok {
		return fmt.Sprintf("(%d) - %s", code, v)
	} else {
		v := "unknown error desc"
		return fmt.Sprintf("(%d) - %s", code, v)
	}
}
