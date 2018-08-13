package lib

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"unicode"
	"time"
)

// RedisPreStart returns a cobra run command for the prerun verb. We check if
// there is already a master registered with consul and if there is we
// reconfigure ourselves as a replica as appropriate.
func RedisPreStart(addr, pass, name, lockPath string) func(cmd *cobra.Command, args []string) {
	rc := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  1 * time.Second,
		Password:     pass,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})
	cc, err := api.NewClient(api.DefaultConfig())

	ri := redisInstance{rc, cc, name, lockPath}

	return func(cmd *cobra.Command, args []string) {
		// Should probably put this in a loop so we can wait on it's success...
		ready, err := ri.isRedisReady()
		if err != nil {
			fmt.Println("Error getting redis status")
			os.exit(1)
		} else if ready != true {
			fmt.Println("Redis not ready to accept connections")
			os.exit(2)
		}

		ri.configureRole()

		// Gather list of existing nodes in the service.
		nodes, meta, err := ri.cc.Catalog().Service("redis", "dev", *q)
		master := ""
		for node := range nodes {
			// Set the value of master to addr:port of any master instance found.
		}

		// If master instance exists, configure ourself to be a replica of it.

		// Otherwise, start as a master instance.

		// Release the lock

	}
}

type redisInstance struct {
	rc *redis.Client
	cc *api.Client
	name string
	lockPath string
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

func (ri *redisInstance) configureRole() error {
	// Acquire lock
	// Get servers in service from consul
	// Check if any of them are masters
	//     If so, become replica of that
	//     If not, become master
	//     If no nodes in service, become master
	return nil
}
