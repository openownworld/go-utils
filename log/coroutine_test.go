package log

import (
	"fmt"
	"github.com/openownworld/go-utils/utils"
	"testing"
)

func TestHandlePanic(t *testing.T) {
	func() {
		defer printPanic()
		panic("throw err for TestHandlePanic")
	}()
}

func TestGo(t *testing.T) {
	Go(func() {
		fmt.Println("run go, panic---------------------")
		panic("throw err for TestSafe")
	})
	t.Logf("GetProcName(): %s", utils.GetProcName())
}
