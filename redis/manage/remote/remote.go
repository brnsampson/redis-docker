package remote

import (
	"github.com/hashicorp/consul/api"
)

// Implement shim for third party library for us to mock out in testing.
type RemoteShim struct {
	cs *api.Client
}

func (s *RemoteShim) ReadRemoteNodes(name string) (*RemoteNodes, error) {
}

func (s *RemoteShim) ReadKey(key string) (string, error) {
}

func (s *RemoteShim) Lock(path string) (string, error) {
}

func (s *RemoteShim) Unlock(path string) (string, error) {
}

// Implement types for the RemoteServiceLockUnlocker

type RemoteNodes []*struct{
	Address string
	ServiceAddress string
	ServicePort int
}

type RemoteService struct {
	s RemoteShimmer
	Name string
	KeyPrefix string
	LockPath string
}

func (rs *RemoteService) LockRemoteService() error {
}

func (rs *RemoteService) UnlockRemoteService() error {
}

func (rs *RemoteService) ReadRemoteAddrs() ([]string, error) {
}

func (rs *RemoteService) ReadRemoteKey(key string) (string, error) {
}
