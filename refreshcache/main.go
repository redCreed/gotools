package main

import (
	"context"
	"fmt"
	"gotools/refreshcache/sync2"
	"math/rand"
	"time"
)

func main() {
	rc, err := sync2.NewReadCache(3*time.Second, 6*time.Second, Db)
	if err != nil {
		panic(err)
	}
	go rc.Run(context.Background())
	for i := 0; i < 20; i++ {
		stat, err := rc.Get(context.Background(), time.Now())
		if err != nil {
			panic(err)
		}
		s := stat.(*resp)
		fmt.Println("data:", i, s.ss)
		time.Sleep(1 * time.Second)
	}

	time.Sleep(1 * time.Hour)
}

// Db 模拟从db读取数据
func Db(ctx context.Context) (interface{}, error) {
	ss := make([]int64, 0)
	t1 := rand.Int63n(100)
	t2 := rand.Int63n(100)
	ss = append(ss, t1, t2)
	return &resp{ss: ss}, nil
}

type resp struct {
	ss []int64
}
