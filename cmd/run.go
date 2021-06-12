package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Start the consumer",
		Run: func(_ *cobra.Command, _ []string) {

			address := viper.GetString("consumer.kafka.address")
			topic := viper.GetString("consumer.kafka.topic")
			partition := viper.GetInt("consumer.kafka.partition")

			conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
			if err != nil {
				log.Fatal("failed to dial leader: ", err)
			}

			err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				log.Fatal("failed to set deadline: ", err)
			}

			batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

			b := make([]byte, 10e3) // 10KB max per message
			for {
				_, err := batch.Read(b)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				fmt.Println(string(b))
			}

			if err := batch.Close(); err != nil {
				log.Fatal("failed to close batch: ", err)
			}

			if err := conn.Close(); err != nil {
				log.Fatal("failed to close connection: ", err)
			}

			os.Exit(0)
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}
