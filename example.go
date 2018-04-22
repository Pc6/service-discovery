package main

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var (
	client *clientv3.Client
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 2 * time.Second
	delayTime      = 500 * time.Millisecond
	endPoints      = []string{"localhost:2379"}
)

const (
	prefix = "service/"
	ttl    = int64(5)
)

func init() {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("before watch")
	WatchServices("test")

	info := &ServiceInfo{
		Name:    "test",
		IP:      "192.168.0.1",
		Port:    9090,
		Version: 1,
	}
	Register("test", info)
	time.Sleep(5 * time.Second)

	info = GetServiceInfo("test")
	fmt.Printf("info: %v\n", *info)

	err := Deregister("test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("deregister service")

	for {
	}
}
