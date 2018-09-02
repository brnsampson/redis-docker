package remote

import "testing"

type TestRemoteShim struct {
	err error
}

func (s *TestRemoteShim) ReadState(name string) (*RemoteState, error) {
	return &RemoteState{}, err
}

func (s *TestRemoteShim) ReadKey(key string) (string, error) {
	return key, err
}

func (s *TestRemoteShim) Lock(path string) error {
	return err
}

func (s *TestRemoteShim) Unlock(path string) error {
	return err
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
