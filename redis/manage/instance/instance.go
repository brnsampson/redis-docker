package instance

import (
	"fmt"
	"github.com/go-redis/redis"
)

// Implement shim for third party library for us to mock out while testing.
type InstanceShim struct{
	Client *redis.Client
}

func(s *InstanceShim) Ping() error {
}

func(s *InstanceShim) GetInfo() (string, error) {
}

func(s *InstanceShim) SetReplication(master string) error {
}

func (c *InstanceShim) ReadConfig(k string) (string, error) {
}

func (c *InstanceShim) UpdateConfig(k, v string) error {
}

func(s *InstanceShim) Quit() error {
}

// Implement an info parser that we can mock out while testing
type InfoParser struct{}

func(ip *InfoParser) ParseInstanceInfo(rawInfo string, c *InstanceInfo) error {
}

// Provide types that satisfy the InstanceConfigKeyReaderUpdater interface
type InstanceInfo map[string]string

type Instance struct{
	InstanceShimmer
	InstanceInfoParser
}

func (i *Instance) IsInstanceReady() (bool, error) {
}

func (i *Instance) ReadInstanceMaster() error {
}

func (i *Instance) UpdateInstanceMaster(addr string) error {
}

func (i *Instance) ClaimInstanceMaster() error {
}
