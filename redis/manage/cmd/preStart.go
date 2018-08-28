// Copyright © 2018 NAME HERE brnsampson@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/brnsampson/redis-docker/redis/manage/service"
)

// preStartCmd represents the preStart command
var preStartCmd = &cobra.Command{
	Use:   "preStart",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: RedisPreStart,
}

func init() {
	redisCmd.AddCommand(preStartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// preStartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// preStartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// RedisPreStart returns a cobra run command for the prerun verb. We check if
// there is already a master registered with consul and if there is we
// reconfigure ourselves as a replica as appropriate.
func RedisPreStart(rm RedisServiceManager) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// Should probably put this in a loop so we can wait on it's success...
		s := service.Service{}
		ready, err := s.IsRedisReady()
		if err != nil {
			fmt.Println("Error getting redis status")
			os.Exit(1)
		} else if ready != true {
			fmt.Println("Redis not ready to accept connections")
			os.Exit(2)
		}

		if err = s.ConfigureServiceRoles(); err != nil {
			fmt.Println("Error while configuring role")
			os.Exit(3)
		}
	}
}
