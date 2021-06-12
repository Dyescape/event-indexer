package cmd

import (
	"os"

	"github.com/Dyescape/event-indexer/internal/consumer/kafka"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	consume = &cobra.Command{
		Use:   "consume",
		Short: "Start the consumer",
		Run: func(_ *cobra.Command, _ []string) {

			kaf := &kafka.Kafka{
				Brokers: []string{viper.GetString("consumer.kafka.address")},
				Topic:   viper.GetString("consumer.kafka.topic"),
				Group:   viper.GetString("consumer.kafka.group"),
			}
			kaf.Consume()

			os.Exit(0)
		},
	}
)

func init() {
	rootCmd.AddCommand(consume)
}
