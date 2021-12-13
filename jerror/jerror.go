package jerror

import (
	//"encoding/json"
	_ "errors"
	"fmt"
	"runtime"
)

type jerror struct {
	code  int
	msg   string
	where string
}

type Error interface {
	Code() int
	Msg() string
	Error() string
}

func (err *jerror) Error() string {
	if len(err.where) > 0 {
		return fmt.Sprintf("(%v) - (%s) - %s", err.code, err.where, err.msg)
	} else {
		return fmt.Sprintf("(%v) - %s", err.code, err.msg)
	}
	//bstr, _ := json.Marshal(err)
	//return string(bstr)
}

func (err *jerror) Code() int {
	return err.code
}

func (err *jerror) Msg() string {
	return err.msg
}

// NewSucceed succeed failed
func NewSucceed(code int) Error {
	return &jerror{code, "succeed", ""}
}

func New(code int, msg string) Error {
	if code != 0 {
		//0为当前栈 1上层调用栈
		pc, file, line, ok := runtime.Caller(1)
		if !ok {
			return &jerror{code, msg, "runtime call err"}
		}
		pcName := runtime.FuncForPC(pc).Name() //获取函数名
		//where:= fmt.Sprintf("%v %s %d %t %s",pc,file,line,ok,pcName)
		where := fmt.Sprintf("%s %d %s", file, line, pcName)
		return &jerror{code, msg, where}
	} else {
		return &jerror{code, msg, ""}
	}
}

func NewError(code int, err error) Error {
	if code != 0 {
		//0为当前栈 1上层调用栈
		pc, file, line, ok := runtime.Caller(1)
		if !ok {
			return &jerror{code, err.Error(), "runtime call err"}
		}
		pcName := runtime.FuncForPC(pc).Name() //获取函数名
		//where:= fmt.Sprintf("%v %s %d %t %s",pc,file,line,ok,pcName)
		where := fmt.Sprintf("%s %d %s", file, line, pcName)
		return &jerror{code, err.Error(), where}
	} else {
		return &jerror{code, err.Error(), ""}
	}
}

func NewErrorText(code int, err error, msg string) Error {
	if code != 0 {
		//0为当前栈 1上层调用栈
		pc, file, line, ok := runtime.Caller(1)
		if !ok {
			return &jerror{code, msg + " - " + err.Error(), "runtime call err"}
		}
		pcName := runtime.FuncForPC(pc).Name() //获取函数名
		//where:= fmt.Sprintf("%v %s %d %t %s",pc,file,line,ok,pcName)
		where := fmt.Sprintf("%s %d %s", file, line, pcName)
		return &jerror{code, msg + " - " + err.Error(), where}
	} else {
		return &jerror{code, msg + " - " + err.Error(), ""}
	}
}

func ErrorCode(err Error) int {
	if err != nil {
		return err.Code()
	}
	return -10000
}
