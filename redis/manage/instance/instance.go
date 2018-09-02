package instance

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"unicode"
)

// Implement shim for third party library for us to mock out while testing.
type InstanceShim struct {
	client *redis.Client
}

func (s *InstanceShim) Ping() error {
	ping, err := s.client.Ping().Result()
	if err != nil {
		return err
	}
	if ping != "PONG" {
		return fmt.Errorf("unexpected value from ping: %s", ping)
	}
	return nil
}

func (s *InstanceShim) GetInfo() (string, error) {
	rawInfo, err := s.client.Info().Result()
	if err != nil {
		return rawInfo, err
	}
	return rawInfo, nil
}

func (s *InstanceShim) SetReplication(host, port string) error {
	_, err := s.client.SlaveOf(host, port).Result()
	if err != nil {
		return err
	}
	return nil
}

//func (s *InstanceShim) ReadConfig(k string) (string, error) {
//}
//
//func (s *InstanceShim) UpdateConfig(k, v string) error {
//}

func (s *InstanceShim) Quit() error {
	result, err := s.client.Quit().Result()
	if err != nil {
		return err
	}
	if result != "OK" {
		return fmt.Errorf("recieved unexpected value from quit command: %q", result)
	}

	return nil
}

// Implement an info parser that we can mock out while testing
type InstanceInfo map[string]string

type InfoParser struct{}

func (ip *InfoParser) ParseInstanceInfo(rawInfo string, info InstanceInfo) error {
	lines := strings.Split(rawInfo, "\r\n")
	for _, line := range lines {
		trimmed := strings.TrimFunc(line, unicode.IsSpace)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.Split(trimmed, ":")
		if len(parts) < 2 {
			continue
		}
		info[parts[0]] = parts[1]
	}
	return nil
}

// Provide types that satisfy the InstanceConfigKeyReaderUpdater interface
type Instance struct {
	InstanceShimmer
	InstanceInfoParser
}

func (i *Instance) GetParsedInfo() (InstanceInfo, error) {
	rawInfo, err := i.GetInfo()
	if err != nil {
		return nil, err
	}

	info := make(InstanceInfo)
	err = i.ParseInstanceInfo(rawInfo, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (i *Instance) IsInstanceReady() (bool, error) {
	// Basic connectivity check.
	err := i.Ping()
	if err != nil {
		return false, err
	}
	// Check to get instance info
	info, err := i.GetParsedInfo()
	if err != nil {
		return false, err
	}

	// Check for ongoing loading from existing rdb or aof backup.
	if loading, ok := info["loading"]; ok && loading == "1" {
		return false, nil
	} else if !ok {
		return false, fmt.Errorf("loading status of persistant data not found")
	}

	// Check for ongoing SYNC from a master.
	if role, ok := info["role"]; ok && role == "master" {
		return true, nil
	} else if !ok {
		return false, fmt.Errorf("instance info did not contain role data")
	}

	syncBytes, ok := info["master_sync_left_bytes"]
	if ok != true {
		return false, fmt.Errorf("replica instance info did not contain sync data")
	}

	b, err := strconv.Atoi(syncBytes)
	if err != nil {
		return false, err
	}

	if b != 0 {
		return false, fmt.Errorf("replica instance still syncing data from master")
	}

	return true, nil
}

func (i *Instance) ReadInstanceMaster() (string, error) {
	info, err := i.GetParsedInfo()
	if err != nil {
		return "", err
	}

	if role, ok := info["role"]; ok != true {
		return "", fmt.Errorf("could not extract role from instance info")
	} else if role == "master" {
		// This instance has no master
		return "", nil
	}

	host, ok := info["master_host"]
	if ok != true {
		return "", fmt.Errorf("could not extract master host from instance info")
	}

	port, ok := info["master_port"]
	if ok != true {
		return "", fmt.Errorf("could not extract master port from instance info")
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	return addr, nil
}

func (i *Instance) UpdateInstanceMaster(addr string) error {
	parts := strings.Split(addr, ":")
	host := parts[0]
	port := parts[1]
	err := i.SetReplication(host, port)
	if err != nil {
		return err
	}

	m, err := i.ReadInstanceMaster()
	if err != nil {
		return err
	}
	if m != addr {
		return fmt.Errorf("instance master is unchanged after update")
	}

	return nil
}

func (i *Instance) ClaimInstanceMaster() error {
	// The special redis command 'SLAVE OF NO ONE' sets an instance as a master
	err := i.SetReplication("NO", "ONE")
	if err != nil {
		return err
	}

	m, err := i.ReadInstanceMaster()
	if err != nil {
		return err
	}
	if m != "" {
		return fmt.Errorf("instance has not been changed to master after update")
	}
	return nil
}
