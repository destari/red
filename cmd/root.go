package cmd

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/cobra"
	"os"
)

type Config struct {
	Threads  	int
	Host 		string
	Port 		string
	Redis 		redis.Conn
	Stream 		string
	Key 		string
	PayloadFile string
	Payload 	string
	Name 		string
	Qty 		int64
	Delay		int64
}

var config Config

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var rootCmd = &cobra.Command{
	Use:   "red",
	Short: "Red is a redis stream tester",
	Long: `This is a long description of the magic.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running the main program.")
	},
}

func initConfig() {
	config.Key = "keyname"
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().IntVarP(&config.Threads, "threads", "t", 2, "Number of threads to run in parallel")
	rootCmd.PersistentFlags().StringVarP(&config.Host, "host", "H", "127.0.0.1", "Local host/IP redis is on")
	rootCmd.PersistentFlags().StringVarP(&config.Port, "port", "P", "6379", "Port redis is listening on")
	rootCmd.PersistentFlags().StringVarP(&config.Stream, "stream", "s", "teststream", "Name of stream to use")
	rootCmd.PersistentFlags().StringVarP(&config.Key, "key", "k", "keyname", "Name of key in stream to use")

}


