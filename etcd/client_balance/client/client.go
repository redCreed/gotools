package main

import (
	"context"
	"fmt"
	"go-do/etcd/client_balance"
	"go-do/etcd/client_balance/proto"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

var (
	// EtcdEndpoints etcd地址
	EtcdEndpoints = []string{"localhost:2379"}
	// SerName 服务名称
	SerName    = "simple_grpc"
	grpcClient proto.SimpleClient
)

func main() {
	r  := client_balance.NewDiscovery(EtcdEndpoints)
	resolver.Register(r)
	// 连接服务器
	conn, err := grpc.Dial(
		r.Scheme()+"://8.8.8.8/simple_grpc",
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure(),
	)
	//grpc.WithBalancerName("round_robin"),

	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立gRPC连接
	grpcClient = proto.NewSimpleClient(conn)
	for i := 0; i < 100; i++ {
		route(i)
		time.Sleep(1 * time.Second)
	}

}

// route 调用服务端Route方法
func route(i int) {
	// 创建发送结构体
	req := proto.SimpleRequest{
		Data: "grpc " + strconv.Itoa(i),
	}
	// 调用我们的服务(Route方法)
	// 同时传入了一个 context.Context ，在有需要时可以让我们改变RPC的行为，比如超时/取消一个正在运行的RPC
	res, err := grpcClient.Route(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	// 打印返回值
	log.Println(res)
}
