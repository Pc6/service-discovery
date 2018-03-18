package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/coreos/etcd/clientv3"
)

var cacheInfo = make(map[string]*ServiceInfo)

func GetServiceInfo(serviceName string) *ServiceInfo {
	if len(serviceName) == 0 {
		log.Fatalln("getInfo arg error")
	}

	key := "services/" + serviceName

	info, ok := cacheInfo[key]
	if ok {
		return info
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := client.Get(ctx, key)
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

	return s
}

func WatchServices(serviceNames ...string) error {
	for _, serviceName := range serviceNames {
		if len(serviceName) == 0 {
			return errors.New("watch arg error")
		}

		key := "services" + serviceName

		ctx, cancel := context.WithTimeout(context.Background, requestTimeout)
		rch := client.Watch(ctx, key)
		cancel()

		go watchEvents(rch)
	}
	return nil
}

func watchEvents(rch clientv3.WatchChan) {
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case "PUT":
				s := new(ServiceInfo)
				err := json.Unmarshal(ev.Value, s)
				if err != nil {
					log.Println("json unmarshal error")
					continue
				}
				cacheInfo[serviceName] = s
			case "DELETE":
				delete(cacheInfo, serviceName)
			}
		}
	}
}
