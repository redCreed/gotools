package etcd

import (
	"fmt"
	"testing"
)

func TestRegister_NewRegister(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ser, err := NewRegister(endpoints,  "/web/node1", "localhost:8000",5)
	if err != nil {
		panic(err)
	}

	fmt.Println(ser)
	select {	}
}
