package utils

import (
	"math/rand"
	"strings"
)

func RightJust(s string, n int, fill string) string {
	return strings.Repeat(fill, n) + s
}

func LeftJust(s string, n int, fill string) string {
	return s + strings.Repeat(fill, n)
}

func Center(s string, n int, fill string) string {
	div := n / 2
	return strings.Repeat(fill, div) + s + strings.Repeat(fill, div)
}

func Random(min, max int) int {
	return rand.Intn(max-min) + min
}

func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	n := len(rs)
	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= n {
		begin = n
	}
	end := begin + length
	if end > n {
		end = n
	}
	// 返回子串
	return string(rs[begin:end])
}
