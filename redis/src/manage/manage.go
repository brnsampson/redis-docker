package main

import (
	"flag"
	"fmt"
)

var cmd = flag.String("cmd", "isRedisReady", "command to run")

func main() {
	flag.Parse()
	ri := newRedisInstance()
	switch *cmd {
	case "isRedisReady":
		fmt.Println(ri.isRedisReady())
	case "isRedisMaster":
		fmt.Println(ri.isRedisMaster())
	case "isConsulReady":
		fmt.Println(isConsulReady())
	}
}
