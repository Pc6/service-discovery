package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

var cacheInfo = make(map[string]*ServiceInfo)

func GetServiceInfo(serviceName string) *ServiceInfo {
	if len(serviceName) == 0 {
		log.Fatalln("getInfo arg error")
	}

	info, ok := cacheInfo[serviceName]
	if ok {
		return info
	}

	key := prefix + serviceName

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := client.Get(ctx, key)
	cancel()
	if err != nil {
		log.Fatalln(err)
	}

	s := new(ServiceInfo)
	for _, ev := range resp.Kvs {
		err = json.Unmarshal(ev.Value, s)
		if err != nil {
			log.Fatalln(err)
		}
	}

	cacheInfo[serviceName] = s

	return s
}

func WatchServices(serviceNames ...string) error {
	for _, serviceName := range serviceNames {
		if len(serviceName) == 0 {
			return errors.New("watch arg error")
		}

		key := prefix + serviceName

		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		rch := client.Watch(ctx, key)
		cancel()

		go watchEvents(serviceName, rch)
	}
	return nil
}

func watchEvents(serviceName string, rch clientv3.WatchChan) {
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				s := new(ServiceInfo)
				err := json.Unmarshal(ev.Kv.Value, s)
				if err != nil {
					log.Println("json unmarshal error")
					continue
				}
				cacheInfo[serviceName] = s
				log.Println("watch update event")
			case mvccpb.DELETE:
				delete(cacheInfo, serviceName)
				log.Println("watch delete event")
			}
		}
	}
}
