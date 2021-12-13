package utils

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
)

func GetCurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("Can not get current file info")
	}
	return file
}

func GetCurrentFilePath() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("Can not get current file info")
	}
	return path.Join(path.Dir(file), "")
}

// GetCurFuncName 获取当前调用者的函数名
func GetCurFuncName() string {
	/*funcName, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Println("Func Name=" + runtime.FuncForPC(funcName).Name())
		fmt.Printf("file: %s    line=%d\n", file, line)
	}*/
	//0为当前栈 1上层调用栈
	pc, _, _, ok := runtime.Caller(0)
	if ok {
		return runtime.FuncForPC(pc).Name()
	}
	return "null"
}

// GetUpperFuncName 获取上层调用者的函数名
func GetUpperFuncName(skip int) string {
	/*funcName, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Println("Func Name=" + runtime.FuncForPC(funcName).Name())
		fmt.Printf("file: %s    line=%d\n", file, line)
	}*/
	//0为当前栈 1上层调用栈 2 上上层
	pc, _, _, ok := runtime.Caller(skip)
	if ok {
		return runtime.FuncForPC(pc).Name()
	}
	return "null"
}

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}

func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}
