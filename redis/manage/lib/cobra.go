package lib

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

// RedisPreStart returns a cobra run command for the prerun verb. We check if
// there is already a master registered with consul and if there is we
// reconfigure ourselves as a replica as appropriate.
func RedisPreStart(rm RedisServiceManager) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// Should probably put this in a loop so we can wait on it's success...
		ready, err := rm.IsRedisReady()
		if err != nil {
			fmt.Println("Error getting redis status")
			os.Exit(1)
		} else if ready != true {
			fmt.Println("Redis not ready to accept connections")
			os.Exit(2)
		}

		if err = rm.ConfigureLocalRole(); err != nil {
			fmt.Println("Error while configuring role")
			os.Exit(3)
		}
	}
}

// RedisPreStart returns a cobra run command for the prerun verb.
func RedisPreStop(rm RedisServiceManager) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// Gracefully close any connections as needed.
		err := rm.Close()
		if err != nil {
			fmt.Println("Error closing redis connections!")
			os.Exit(1)
		}
	}
}
