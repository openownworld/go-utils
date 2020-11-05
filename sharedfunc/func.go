package sharedfunc

import "runtime"

//GetCurFuncName 获取当前调用者的函数名
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

//GetUpperFuncName 获取上层调用者的函数名
func GetUpperFuncName() string {
	/*funcName, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Println("Func Name=" + runtime.FuncForPC(funcName).Name())
		fmt.Printf("file: %s    line=%d\n", file, line)
	}*/
	//0为当前栈 1上层调用栈
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		return runtime.FuncForPC(pc).Name()
	}
	return "null"
}
