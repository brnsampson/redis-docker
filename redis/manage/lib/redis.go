package lib

import (
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"unicode"
	"time"
)

type RedisServiceManager interface {
	ServiceLogger
	RedisCommunicator
	ConsulGetter
	ConfigureLocalRole() error
	Close() error
}

type RedisCommunicator interface {
	Ping() error
	Info() (string, error)
	ParsedInfo() (map[string]string, error)
	GetRole() (string, error)
	IsRedisReady() (bool, error)
	SetReplication() error
}

type RedisInstance struct {
	rc *redis.Client
	Address string
	Password string
	DialTimeout time.Duration
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

type RedisService struct {
	RedisInstance
	ConsulService
	ServiceLogger
}

// func (ri *RedisInstance) getRedisClient() *redis.Client {
// 	if ri.rc == nil {
// 		if ri.DialTimeout == 0 { ri.DialTimeout = 1 * time.Second }
// 		if ri.ReadTimeout == 0 { ri.ReadTimeout = 1 * time.Second }
// 		if ri.WriteTimeout == 0 { ri.WriteTimeout = 1 * time.Second }
// 
// 		ri.rc = &redis.NewClient(&redis.Options{
// 			Addr:         ri.Address,
// 			DialTimeout:  ri.DialTimeout,
// 			Password:     ri.Password,
// 			ReadTimeout:  ri.ReadTimeout,
// 			WriteTimeout: ri.WriteTimeout,
// 		})
// 	}
// 	return ri.rc
// }

func (ri *RedisInstance) Ping() error {
	ping, err := ri.rc.Ping().Result()

	if err != nil {
		return err
	}

	if ping != "PONG" {
		return fmt.Errorf("unexpected value from ping: %s", ping)
	}
	return nil
}

func (ri *RedisInstance) Info() (string, error) {
	rawInfo, err := ri.rc.Info().Result()

	if err != nil {
		return rawInfo, err
	}
	return rawInfo, nil
}

func (ri *RedisInstance) ParsedInfo() (map[string]string, error) {
	rawInfo, err := ri.Info()
	if err != nil {
		return nil, err
	}

	out := make(map[string]string)
	lines := strings.Split(rawInfo, "\r\n")
	for _, line := range lines {
		trimmed := strings.TrimFunc(line, unicode.IsSpace)
		//if strings.HasPrefix(trimmed, "#") {
		//	continue
		//}

		parts := strings.Split(trimmed, ":")

		if len(parts) < 2 {
			continue
		}

		out[parts[0]] = parts[1]
	}
	return out, nil
}

func (ri *RedisInstance) IsRedisReady() (bool, error) {
	// Basic connectivity check.
	err := ri.Ping()
	if err != nil {
		return false, err
	}

	// Check to get instance info
	info, err := ri.ParsedInfo()
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
	if _, ok := info["master_sync_left_bytes"]; ok != true {
		return false, fmt.Errorf("instance still syncing from master")
	}

	return true, nil
}

func (ri *RedisInstance) SetReplication(addr string) error {
	parts := strings.Split(addr, ":")
	host := parts[0]
	port := parts[1]

	_, err := ri.rc.SlaveOf(host, port).Result()
	if err != nil {
		return err
	}
	return nil
}

func (ri *RedisInstance) GetRole() (string, error) {
	info, err := ri.ParsedInfo()
	if err != nil {
		return "", err
	}

	role, ok := info["role"];

	if ok != true {
		return role, fmt.Errorf("could not read redis role from instance info")
	}

	return role, nil
}

func (rs *RedisService) ConfigureLocalRole() error {
	// TODO: write this
	// Acquire lock
	// Get servers in service from consul
	// Check if any of them are masters
	//     If so, become replica of that
	//     If not, become master
	//     If no nodes in service, become master
	return nil
}

func (rs *RedisService) Close() error {
	//TODO
	return nil
}
