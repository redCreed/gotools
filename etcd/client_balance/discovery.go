package client_balance

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"log"
	"sync"
	"time"
)
const schema = "grpclb"

//相关文档: https://www.cnblogs.com/FireworksEasyCool/p/12912839.html
//客户端负载均衡: 需要实现Builder和Resolver接口
//服务发现
type Discovery struct {
	Client *clientv3.Client
	cc resolver.ClientConn
	srvList sync.Map //服务列表
}

func NewDiscovery(endpoints []string) resolver.Builder{
	var (
		client *clientv3.Client
		err error
	)
	client, err  = clientv3.New(clientv3.Config{
		Endpoints:            endpoints,
		DialTimeout:          5*time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Discovery{Client: client}
}


func (d *Discovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
 	d.cc = cc
 	preFix := "/"+target.Scheme+"/"+target.Endpoint+"/"

 	//根据前缀获取key
 	getResp, err := d.Client.Get(context.TODO(), preFix, clientv3.WithPrefix())
 	if err != nil {
 		return nil, err
	}
	adds := make([]resolver.Address, 0)
	for _, v := range getResp.Kvs {
		adds = append(adds, resolver.Address{
			Addr:      string(v.Value) ,
		})
	}

	if err := d.cc.UpdateState(resolver.State{Addresses: adds}); err != nil {
		return d, err
	}

	//监视前缀，修改变更的server
	go d.watcher(preFix)
	return d, nil
}

func (d *Discovery) Scheme() string {
	return schema
}

func (d *Discovery) ResolveNow(options resolver.ResolveNowOptions) {
	panic("implement me")
}

func (d *Discovery) Close() {
	d.Client.Close()
}

func  (d *Discovery)setSrvList(key, value string) {
	d.srvList.Store(key, resolver.Address{Addr: value})
	//获取所有地址
	if err := d.cc.UpdateState(resolver.State{Addresses: d.getSrvList()}); err != nil {
		log.Fatal(err)
	}
}

func (d *Discovery)getSrvList() []resolver.Address  {
	addrs := make([]resolver.Address, 0)
	d.srvList.Range(func(key, value interface{}) bool {
		addrs = append(addrs, value.(resolver.Address))
		return true
	})

	return addrs
}

//DelServiceList 删除服务地址
func (s *Discovery) DelSrvList(key string) {
	s.srvList.Delete(key)
	s.cc.UpdateState(resolver.State{Addresses: s.getSrvList()})
	log.Println("del key:", key)
}

//watcher 监听前缀
func (s *Discovery) watcher(prefix string) {
	rch := s.Client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: //新增或修改
				s.setSrvList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				s.DelSrvList(string(ev.Kv.Key))
			}
		}
	}
}