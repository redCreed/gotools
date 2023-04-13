package sync2

import (
	"fmt"
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {
	//当onDo执行出现超时，会调用onTimeout
	WithTimeout(3*time.Second, onDo, onTimeout)
}

func onTimeout() {
	fmt.Println("func onDo timeout")
}

func onDo() {
	time.Sleep(5 * time.Second)
}