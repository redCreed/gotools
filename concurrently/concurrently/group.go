package concurrently

import (
	"fmt"
	"sync"
)

//errgroup只能捕获一组协程中的第一个错误

//可以记录多个错误，但是暂未实现超时取消
type Group struct {
	wg   sync.WaitGroup
	mu   sync.Mutex
	errs []error
}

//执行一个带所有返回错误的函数
func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				fmt.Println("oops, panic")
			}
		}()
		if err := f(); err != nil {
			g.mu.Lock()
			g.errs = append(g.errs, err)
			g.mu.Unlock()
		}
		g.wg.Done()
	}()
}

func (g *Group) Wait() []error {
	g.wg.Wait()
	return g.errs
}
