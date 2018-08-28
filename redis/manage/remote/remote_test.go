package remote

import "testing"

type TestRemoteShim struct {}

func (s *TestRemoteShim) ReadState(name string) (*RemoteState, error) {
}

func (s *TestRemoteShim) ReadKey(key string) (string, error) {
}

func (s *TestRemoteShim) Lock(path string) (string, error) {
}

func (s *TestRemoteShim) Unlock(path string) (string, error) {
}

// Perform tests
func TestLockRemoteService(t *testing.T) {
}

func TestUnlockRemoteService(t *testing.T) {
}

func TestReadRemoteAddrs(t *testing.T) {
}

func TestReadRemoteKey(t *testing.T) {
}
