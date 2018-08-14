package lib

import "testing"

// Testing functions that return cobra stuff with our redis code stubbed out.

// FakeRedisService is a stub struct that implements the RedisServiceManager
// interface.
type FakeRedisService struct {}

func (rs *FakeRedisService) IsRedisReady() (bool, error) {
	return true, nil
}

func (rs *FakeRedisService) GetRole() (string, error) {
	return "master", nil
}

func (rs *FakeRedisService) ConfigureRole() error {
	return nil
}

func (rs *FakeRedisService) ParsedRedisInfo() (map[string]string, error) {
	m := make(map[string]string)
	return m, nil
}

func (rs *FakeRedisService) Close() error {
	return nil
}

func TestRedisPreRun(t *testing.T) {
}

func TestRedisPostRun(t *testing.T) {
}
