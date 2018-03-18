package main

import (
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
	endPoints      = []string{"localhost:2379"}
)

func init() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

}
