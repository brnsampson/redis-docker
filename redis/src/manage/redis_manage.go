package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisInstance struct {
	rc *redis.Client
}

func newRedisInstance(redisAddress string) *redisInstance {
	rc := redis.NewClient(&redis.Options{
		Addr:         redisAddress,
		DialTimeout:  1 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	return &redisInstance{rc}
}

func (ri *redisInstance) isRedisReady() error {
	// Basic connectivity check.
	ping, err := ri.rc.Ping().Result()
	if err != nil {
		return err
	}
	if ping != "PONG" {
		return fmt.Errorf("unexpected value from ping: %s", ping)
	}

	// Check to get instance info
	rawInfo, err := ri.rc.Info().Result()

	if err != nil {
		return err
	}

	info := parseRedisInfo(rawInfo)

	// Check for ongoing loading from existing rdb or aof backup.
	if loading, ok := info["loading"]; ok && loading == "1" {
		return nil
	} else if !ok {
		return fmt.Errorf("loading status of persistant data not found")
	}

	// Check for ongoing SYNC from a master.
	if _, ok := info["master_sync_left_bytes"]; ok {
		return fmt.Errorf("instance still syncing from master")
	}

	return nil
}

func (ri *redisInstance) isRedisMaster() (bool, error) {
	rawInfo, err := ri.rc.Info().Result()

	if err != nil {
		return false, err
	}

	info := parseRedisInfo(rawInfo)

	if role, ok := info["role"]; ok && role == "master" {
		return true, nil
	} else if ok && role != "master" {
		return false, nil
	}

	return false, fmt.Errorf("instance role could not be read")
}

func (ri *redisInstance) makeMaster() error {
	result := ri.rc.SlaveOf("NO", "ONE")
	return result.Err()
}

func (ri *redisInstance) makeSlave(addr string) error {
	addrParts := strings.Split(addr, ":")

	if len(addrParts) != 2 {
		return errors.Errorf("invalid address given, expected format ip:port, got %v", addr)
	}
	return ri.rc.SlaveOf(addrParts[0], addrParts[1]).Err()
}

func parseRedisInfo(in string) map[string]string {
	out := make(map[string]string)
	lines := strings.Split(in, "\r\n")
	for _, line := range lines {
		trimmed := strings.TrimFunc(line, unicode.IsSpace)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.Split(trimmed, ":")

		if len(parts) < 2 {
			continue
		}

		out[parts[0]] = parts[1]
	}
	return out
}
