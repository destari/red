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
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/go-redis/redis"

	"github.com/spf13/cobra"
)

func Generate(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	c := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println("Generate thread starting up: ", id)

	for {
		for n := int64(0); n < config.Qty; n++ {
			//c.Send("XADD", config.Stream, "*", config.Key, config.Payload)
			if _, err := c.XAdd(&redis.XAddArgs{
				Stream:       config.Stream,
				MaxLen:       0,
				MaxLenApprox: 0,
				ID:           "*",
				Values:       map[string]interface{}{config.Key: config.Payload},
			}).Result(); err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println("Completed: ", "Thread: ", id, " Qty: ", config.Qty)
		time.Sleep(time.Millisecond * time.Duration(config.Delay))
	}
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generator description here..",
	Long:  `Longer description goes here.. but I have nothing to say.`,
	Args:  cobra.MinimumNArgs(0),

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Generate: ") // + strings.Join(args, ", ")

		// open payload file and load into payload config var
		config.Payload = "{}"
		if config.PayloadFile != "" {
			data, err := ioutil.ReadFile(config.PayloadFile)
			check(err)
			fmt.Print(string(data))
			config.Payload = string(data)
		}

		fmt.Println("Starting threads..")
		var wg sync.WaitGroup
		for t := 0; t < config.Threads; t++ {
			wg.Add(1)
			go Generate(&wg, t)
		}
		wg.Wait()
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	generateCmd.PersistentFlags().StringVarP(&config.PayloadFile, "payload", "p", "", "Filename for sample payload data")
	generateCmd.PersistentFlags().Int64VarP(&config.Delay, "delay", "d", 1000, "Number of milliseconds to delay between runs")
	generateCmd.PersistentFlags().Int64VarP(&config.Qty, "count", "n", 1000, "Number of items per run")
	generateCmd.PersistentFlags().StringVarP(&config.Key, "key", "k", "project-id", "Name of key in stream to use")
	rootCmd.AddCommand(generateCmd)

}
