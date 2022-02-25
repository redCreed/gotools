package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var a uint32 = 10
	atomic.AddUint32(&a, 1)
	fmt.Println(a)
	// uint32需要一个非负整数，uint32(int32(-2)), 会被编译器报错，需要一个中间变量b来绕过 结果减2
	b := int32(-2)
	atomic.AddUint32(&a, uint32(b))
	fmt.Println(a)
	// ^uint32(n-1), n为要减去的数
	// 整数在计算机以补码形式存在，这里的异或求出来的补码与b的补码相同
	atomic.AddUint32(&a, ^uint32(3-1))
	fmt.Println(a)

	//atomic.AddUint32(&a, ^uint32(6))
	//fmt.Println(a)
}
