package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	var value int32 = 0
	go func() {
		t  := time.NewTicker(1*time.Second)
		for {
			select {
			case <- t.C:
				fmt.Println("value:", value)
				atomic.AddInt32(&value, 1)
			}
		}
	}()
	ch := make(chan int32)

	go func() {
		for {
			if atomic.CompareAndSwapInt32(&value, 10, 11111) {
				fmt.Println("The second number has gone to zero.")
				ch <- value
				break
			}
			fmt.Println("num:", value)
			time.Sleep(time.Millisecond * 500)
		}
	}()

	for {
		data := <- ch
		fmt.Println("data:", data)
	}
}
