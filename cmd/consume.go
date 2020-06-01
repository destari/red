/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

var ctx = context.Background()

func Consume(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	c := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println("Consume thread starting up: ", id, config.Stream, config.Name)

	c.XGroupCreateMkStream(ctx, config.Stream, "mygroup", "0")

	for {
		// XREADGROUP GROUP $GroupName $ConsumerName BLOCK 2000 COUNT 10 STREAMS mystream >
		results, err := c.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "mygroup",
			Consumer: config.Name,
			Streams:  []string{config.Stream, ">"},
			Count:    config.Qty,
			Block:    time.Duration(config.Delay),
			NoAck:    true,
		}).Result()

		if err != nil {
			fmt.Println(err)
		}
		/*
		[
			{
				teststream: [
					{
						1590965499617-7: map[keyname: PAYLOAD ]
					}
				]
			}
		]

		 */

		for _, res := range results {
			fmt.Println(res.Messages[0].ID)
			c.XDel(ctx, config.Stream, res.Messages[0].ID)

		}

		fmt.Println("Completed: ", "Thread: ", id, " Qty: ", config.Qty)
	}
}

// consumeCmd represents the consume command
var cmdConsume = &cobra.Command{
	Use:   "consume",
	Short: "Consumer description here..",
	Long: `Longer description goes here.. but I have nothing to say.`,
	Args: cobra.MinimumNArgs(0),

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Consume: ")

		fmt.Println("Starting threads..")
		var wg sync.WaitGroup
		for t := 0; t < config.Threads; t++ {
			wg.Add(1)
			go Consume(&wg, t)
		}
		wg.Wait()
	},
}



func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consumeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmdConsume.PersistentFlags().StringVarP(&config.Name, "name", "n", "client", "Name of client consuming stream data")
	cmdConsume.PersistentFlags().Int64VarP(&config.Delay, "delay", "d", 1000, "Number of milliseconds to delay between runs")
	cmdConsume.PersistentFlags().Int64VarP(&config.Qty, "count", "c", 1000, "Number of items to batch receive")
	cmdConsume.PersistentFlags().StringVarP(&config.Key, "key", "k", "keyname", "Name of key in stream to use")
	rootCmd.AddCommand(cmdConsume)
}
