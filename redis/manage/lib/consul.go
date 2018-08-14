package lib

import (
	"github.com/hashicorp/consul/api"
)

type ConsulGetter interface {
	GetServiceAddrs() ([]string, error)
	GetConsulKey(string) (string, error)
}

type ConsulService struct {
	cc *api.Client
	Name string
	LockPath string
}

func (cs *ConsulService) GetServiceAddrs() ([]string, error) {
	return []string{}, nil
}

func (cs *ConsulService) GetConsulKey(key string) (string, error) {
	return "", nil
}
