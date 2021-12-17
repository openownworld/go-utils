package str

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
