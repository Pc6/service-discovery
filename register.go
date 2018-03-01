package main

type ServiceInfo struct {
	IP      string
	Port    int
	Version int
}

func Register(serviceName string, info *ServiceInfo) error {
	return nil
}

func Deregister(serviceName string) error {
	return nil
}
