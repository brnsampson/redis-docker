package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/hashicorp/consul/api"
	"strings"
	"unicode"
	"os"
	"time"
)

type redisInstance struct {
	rc *redis.Client
	cc *api.Client
}

func newRedisInstance() *redisInstance {
	rc := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_ADDR"),
		DialTimeout:  1 * time.Second,
		Password:     os.Getenv("REDIS_PASS"),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	return &redisInstance{rc}
}

func (ri *redisInstance) redisPreStart() error {
	// TODO: get lock here before testing for masters.

	// Gather list of existing nodes in the service.
	nodes, meta, err := ri.cc.Service("redis", "dev", *q)
	master := ""
	for node := range nodes {
		// Set the value of master to addr:port of any master instance found.
	}

	// If master instance exists, configure ourself to be a replica of it.

	// Otherwise, start as a master instance.

	// Release the lock
}

func (ri *redisInstance) isRedisReady() (bool, error) {
	// Basic connectivity check.
	ping, err := ri.rc.Ping().Result()
	if err != nil {
		return false, err
	}
	if ping != "PONG" {
		return false, fmt.Errorf("unexpected value from ping: %s", ping)
	}

	// Check to get instance info
	rawInfo, err := ri.rc.Info().Result()

	if err != nil {
		return false, err
	}

	info := parseRedisInfo(rawInfo)

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

func parseRedisInfo(in string) map[string]string {
	out := make(map[string]string)
	lines := strings.Split(in, "\r\n")
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
	return out
}
