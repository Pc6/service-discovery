package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type ServiceInfo struct {
	Name    string `json:"serviceName"`
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Version int    `json:"version,omitempty"`
}

type Service struct {
	Info          *ServiceInfo
	StopHeartBeat chan struct{}
	LeaseID       clientv3.LeaseID
}

var serviceMap = make(map[string]*Service)

func Register(serviceName string, info *ServiceInfo) error {
	if len(serviceName) == 0 || info == nil {
		return errors.New("register arg error")
	}

	if _, ok := serviceMap[serviceName]; ok {
		return errors.New("service has been registered")
	}

	key := "service/" + serviceName
	val, err := json.Marshal(info)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := client.Grant(ctx, 5) // default ttl is 5
	cancel()
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = client.Put(ctx, key, string(val), clientv3.WithLease(resp.ID))
	cancel()
	if err != nil {
		return err
	}

	s := &Service{
		Info:          info,
		StopHeartBeat: make(chan struct{}),
		LeaseID:       resp.ID,
	}
	serviceMap[serviceName] = s

	go s.HeartBeat()

	return nil
}

func Deregister(serviceName string) error {
	if len(serviceName) == 0 {
		return errors.New("deregister arg error")
	}

	s, ok := serviceMap[serviceName]
	if !ok {
		return errors.New("service is not registered")
	}

	key := "service/" + serviceName
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := client.Delete(ctx, serviceName)
	cancel()
	if err != nil {
		return err
	}

	// call the go routine to stop heartBeat
	s.StopHeartBeat <- struct{}{}

	delete(serviceMap, serviceName)

	return nil
}

func (s *Service) HeartBeat() {
	for {
		select {
		case <-time.After(5 * time.Second):
			ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
			_, err := client.KeepAliveOnce(ctx, s.LeaseID)
			cancel()
			if err != nil {
				log.Printf("heartBeat: %s\n", err.Error())
				return
			}
			log.Println("heartBeat...")
		case <-s.StopHeartBeat:
			return
		}
	}
}
