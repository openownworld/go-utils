package utils

import (
	"fmt"
	"github.com/openownworld/go-utils/zaplog"
	"runtime"
)

const (
	maxStack  = 20
	separator = "---------------------------------------\n"
)

var panicHandler func(string)

// SetPanicHandler 设置异常时处理接口
func SetPanicHandler(h func(string)) {
	panicHandler = func(str string) {
		defer func() {
			printPanic()
		}()
		h(str)
	}
}
func printPanic() {
	if err := recover(); err != nil {
		str := fmt.Sprintf("\n%sruntime error: %v\ntraceback:\n", separator, err)
		i := 2
		for {
			pc, file, line, ok := runtime.Caller(i)
			if !ok || i > maxStack {
				break
			}
			str += fmt.Sprintf("\tstack: %d %v [file: %s:%d] func: %s\n", i-1, ok, file, line, runtime.FuncForPC(pc).Name())
			i++
		}
		str += separator
		zaplog.Error(str)
	}
}

func handlePanic() {
	if err := recover(); err != nil {
		errstr := fmt.Sprintf("\n%sruntime error: %v\ntraceback:\n", separator, err)
		i := 2
		for {
			pc, file, line, ok := runtime.Caller(i)
			if !ok || i > maxStack {
				break
			}
			errstr += fmt.Sprintf("\tstack: %d %v [file: %s:%d] func: %s\n", i-1, ok, file, line, runtime.FuncForPC(pc).Name())
			i++
		}
		errstr += separator
		if panicHandler != nil {
			panicHandler(errstr)
		} else {
			fmt.Println(errstr)
		}
	}
}

// Go 封装的协程支持异常恢复，打印堆栈
func Go(cb func()) {
	go func() {
		defer zaplog.PrintPanicLog()
		cb()
	}()
}
