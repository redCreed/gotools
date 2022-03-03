package etcd

import (
	"log"
	"testing"
	"time"
)

func TestNewDiscovery(t *testing.T) {
	endpoints := []string{"localhost:2379"}
	ser, err :=  NewDiscovery(endpoints)
	if err != nil {
		panic(err)
	}
	defer ser.Close()
	ser.WatchSrv("/web/")
	for {
		select {
		case <-time.Tick(3 * time.Second):
			log.Println(ser.GetSrv())
		}
	}
}