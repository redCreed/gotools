package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

func doSomething(u string) { // 模拟抓取任务的执行
	time.Sleep(2 * time.Second)
	fmt.Println("sleep", u)
}

const (
	Limit  = 3 // 同時并行运行的goroutine上限
	Weight = 1 // 每个goroutine获取信号量资源的权重
)

func main() {
	urls := []string{
		"http://www.example.com",
		"http://www.example.net",
		"http://www.example.net/foo",
		"http://www.example.net/bar",
		"http://www.example.net/baz",
	}
	s := semaphore.NewWeighted(Limit)
	var w sync.WaitGroup
	t := time.Now()
	for _, u := range urls {
		w.Add(1)
		go func(u string) {
			s.Acquire(context.Background(), Weight)
			doSomething(u)
			s.Release(Weight)
			w.Done()
		}(u)
	}
	w.Wait()
	fmt.Println(time.Since(t))
	fmt.Println("All Done")
}
