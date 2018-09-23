package main

import (
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: /bin/manage isRedisHealthy|isConsulHealthy|redisPreStart|onChange [flags]")
		os.Exit(1)
	}

	// TODO: flag for this
	log.SetLevel(log.DebugLevel)

	publicRedisAddress := getSessionName()
	log.Debugf("Retrieved redis address : %v\n", publicRedisAddress)

	ri := newRedisInstance(publicRedisAddress)
	ci, err := newConsulInstance("/tmp/consulCacheFile")
	if err != nil {
		log.Fatal(err)
	}

	manager, err := newSingleMasterManager(ri, ci, publicRedisAddress)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[1]
	switch cmd {
	case "isRedisHealthy":
		if err := manager.handleHealthCheck(); err != nil {
			log.Errorf("Redis health check failed: %+v\n", err)
			os.Exit(1)
		}
	case "isConsulHealthy":
		if err := ci.isConsulReady(); err != nil {
			log.Errorf("Consul health check failed: %+v\n", err)
			os.Exit(1)
		}
	case "redisPreStart":
		if err := manager.handlePreStart(); err != nil {
			log.Errorf("Redis preStart failed: %+v\n", err)
			os.Exit(1)
		}
	case "onChange":
		if err := manager.handleChange(); err != nil {
			log.Errorf("Redis onChange failed: %+v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command %v\n", cmd)
		os.Exit(1)
	}
}

func getSessionName() string {
	return fmt.Sprintf("%v:%v", getLocalIp(), 6379) // TODO
}

func getLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	panic("wtf")
}
