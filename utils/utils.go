package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

/*
0在C语言以及Go语言不同的表现在大部分情况下不会造成问题，但是当使用io.Reader(b []byte)时，如果传入的字节数组b本身长度大于reader可读到的长度，则会导致末尾被0补齐。当直接使用string(b)强制类型转换时会导致显示上看似无问题，但是实际上字符串并不相同。
要解决这个问题需要对[]byte和string的转换过程进行一个封装。
这里实现了针对两种情况的解决方案，前者是遇到0就结束转换，后者则是忽略所有的0并将剩余部分拼接
*/

// String 将 `[]byte` 转换为 `string`
func String(b []byte) string {
	for idx, c := range b {
		if c == 0 {
			return string(b[:idx])
		}
	}
	return string(b)
}

// ByteNoneZero String 将 `[]byte` 转换为 `string`
func ByteNoneZero(b []byte) []byte {
	for idx, c := range b {
		if c == 0 {
			return b[:idx]
		}
	}
	return b
}

// StringWithoutZero 将 `[]byte` 转换为 `string`
func StringWithoutZero(b []byte) string {
	s := make([]rune, len(b))
	offset := 0
	for i, c := range b {
		if c == 0 {
			offset++
		} else {
			s[i-offset] = rune(c)
		}
	}
	return string(s[:len(b)-offset-1])
}

// ByteToHexString 打印16hex
func ByteToHexString(buf []byte, bufLen int) string {
	// 遍历, 转为16进制
	var hexStr string
	for _, b := range buf[:bufLen] {
		hex := fmt.Sprintf("%x", b)
		if len(hex) == 1 {
			hex = "0" + hex
		}
		hexStr += hex + " "
	}
	return hexStr
}

// BytesCombine 多个[]byte数组合并成一个[]byte
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// GetFileSize "c:/test.bmp"
func GetFileSize(filenamePath string) int64 {
	var result int64
	filepath.Walk(filenamePath, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

func GetFileSizeBySeek(filenamePath string) int64 {
	fin, err := os.Open(filenamePath)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	fin.Seek(0, os.SEEK_SET)
	result, _ := fin.Seek(0, os.SEEK_END)
	fin.Close()
	return result
}

func GetRunPath() string {
	runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return runDir
}

func GetAbsolutePath(addPath string) string {
	runDir := GetRunPath()
	fileDir := ""
	switch runtime.GOOS {
	case "darwin":
		fileDir = runDir + "/" + addPath
	case "windows":
		fileDir = runDir + "\\" + addPath
	case "linux":
		fileDir = runDir + "/" + addPath
	}
	return fileDir
}

func GetWorkPath() string {
	workDir, _ := os.Getwd()
	return workDir
}

func BoolString(b bool) string {
	if b == true {
		return "true"
	} else {
		return "false"
	}
}

func StrToInt(strNumber string, value interface{}) (err error) {
	var number interface{}
	number, err = strconv.ParseInt(strNumber, 10, 64)
	switch v := number.(type) {
	case int64:
		switch d := value.(type) {
		case *int64:
			*d = v
		case *int:
			*d = int(v)
		case *int16:
			*d = int16(v)
		case *int32:
			*d = int32(v)
		case *int8:
			*d = int8(v)
		}
	}
	return
}

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func WriteFile(path string, buf []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Println("error write permission denied")
		}
		if os.IsNotExist(err) {
			fmt.Println("file does not exist")
		}
		return err
	}
	defer f.Close()
	return ioutil.WriteFile(path, buf, 0666)
}

func GetUintByInterface(obj interface{}) uint {
	if it, ok := obj.(float64); ok {
		return uint(it)
	} else {
		return 0
	}
}

func GetIntByInterface(obj interface{}) int {
	if it, ok := obj.(float64); ok {
		return int(it)
	} else {
		return 0
	}
}

func FloatDecimalPoint(v float64, pointNum int) float64 {
	format := fmt.Sprintf("%d", pointNum)
	format = "%." + format + "f" //
	//s := fmt.Sprintf("%.2f", v)
	s := fmt.Sprintf(format, v)
	f, _ := strconv.ParseFloat(s, pointNum)
	return f
}

func GetTimeFormat() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

func ToJson(obj interface{}) string {
	bJson, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(bJson)
}

// DeepCopy 最简单的方法是用json unmarshal，变成字符串，然后再用 json marshal生成新的map。这种方法对结构体也适用。
// 如果是map[string]interface{}和[]interface{}的组合，用代码递归也很简单：
func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}
		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}
		return newSlice
	}
	return value
}

func Typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func CheckMapItemIsNull(m map[string]interface{}) (string, bool) {
	for k, v := range m {
		fmt.Println(k, v)
		switch value := v.(type) {
		case string:
			fmt.Printf("type is string, %s[%v]\n", k, v)
		case bool:
			fmt.Printf("type is bool, %s[%v]\n", k, v)
		case int:
			fmt.Printf("type is int, %s[%v]\n", k, v)
		case float32, float64:
			fmt.Printf("type is float, %s[%v]\n", k, v)
		case []interface{}:
			fmt.Printf("type is []interface, %s[%v]\n", k, v)
		default:
			fmt.Printf("type:unkown type is %T, %s[%v]\n", value, k, v)
			return k, true
		}
	}
	return "", false
}

func CheckStructItemIsNull(obj interface{}) (string, bool) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	fmt.Println("interface Type:", t.Name())
	count := t.NumField()
	// for i := 0; i < count; i++ {
	// 	f := t.Field(i)
	// 	fmt.Printf("index:%d\nName:%s\nPkgPath:%s\nType:%v\nTag:%s\nOffset:%v\nIndex:%v\nAnonymous:%v\n",
	// 		i, f.Name, f.PkgPath, f.Type, f.Tag, f.Offset, f.Index, f.Anonymous)
	// 	fmt.Println("------------------------------")
	// }
	for i := 0; i < count; i++ {
		//注：当结构体中含有非导出字段时，v.Field(i).Interface()会panic,
		if v.Field(i).CanInterface() { //判断是否为可导出字段
			f := t.Field(i)
			val := v.Field(i).Interface()
			fmt.Printf("Fileds: %6s : %v %v\n", f.Name, f.Type, val)
			switch v.Field(i).Type().Kind() {
			case reflect.String:
				if val == "" {
					return f.Name, true
				}
			case reflect.Int:
				if val == 0 {
					return f.Name, true
				}
			case reflect.Int64:
				if val == 0 {
					return f.Name, true
				}
			case reflect.Interface:
				if val == nil {
					return f.Name, true
				}
			case reflect.Struct:
				if val == nil {
					return f.Name, true
				}
			case reflect.Slice:
				if val == nil {
					return f.Name, true
				}
				s := reflect.ValueOf(val)
				if s.Len() <= 0 {
					return f.Name, true
				}
				// if len(val.([]interface{})) <= 0 {
				// 	return f.Name, true
				// }
			case reflect.Map:
				if val == nil {
					return f.Name, true
				}
				if len(val.(map[string]interface{})) <= 0 {
					return f.Name, true
				}
			case reflect.Bool:

			default:
				return f.Name, true
			}
		}
	}
	return "", false
}

// IsExist returns whether the given file or directory exists or not
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// no such file or dir
		return false
	}
	if info.IsDir() {
		// it's a directory
		return true
	} else {
		// it's a file
		return false
	}
}
