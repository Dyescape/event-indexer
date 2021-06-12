package cmd

import (
	"os"

	"github.com/Dyescape/event-indexer/internal/elastic"

	"github.com/Dyescape/event-indexer/internal/consumer/kafka"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	publish = &cobra.Command{
		Use:   "publish",
		Short: "Start the test publisher",
		Run: func(_ *cobra.Command, _ []string) {

			kaf := &kafka.Kafka{
				Brokers: viper.GetStringSlice("consumer.kafka.address"),
				Topic:   viper.GetString("consumer.kafka.topic"),
				Group:   viper.GetString("consumer.kafka.group"),
				Elastic: elastic.NewElasticSearchClient(),
			}
			kaf.PublishTestMessages()

			os.Exit(0)
		},
	}
)

func init() {
	rootCmd.AddCommand(publish)
}
