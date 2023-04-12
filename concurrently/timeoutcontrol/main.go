package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func hardWork(job interface{}) error {
	time.Sleep(time.Second * 10)
	dd := make([]int, 0)
	fmt.Println("dd:", dd[1])
	fmt.Println("exit:", time.Now())
	return nil
}
func requestWork(ctx context.Context, job interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	done := make(chan error, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		fmt.Println("before:", time.Now())
		//todo hardWork耗时太久 下面select会走ctx.Done
		done <- hardWork(job)
	}()
	select {
	case err := <-done:
		return err
	case p := <-panicChan:
		panic(p)
	case <-ctx.Done():
		return ctx.Err()
	}
}

//最好使用http请求的使用
func main() {
	const total = 10
	var wg sync.WaitGroup
	wg.Add(total)
	now := time.Now()
	for i := 0; i < total; i++ {
		go func() {
			defer func() {
				if p := recover(); p != nil {
					//todo 这里最好要日志记录
					fmt.Println("oops, panic", p)
				}
			}()
			defer wg.Done()
			requestWork(context.Background(), "any")
		}()
	}
	wg.Wait()
	fmt.Println("elapsed:", time.Since(now))
	time.Sleep(time.Second * 20)
	fmt.Println("number of goroutines:", runtime.NumGoroutine())
}
