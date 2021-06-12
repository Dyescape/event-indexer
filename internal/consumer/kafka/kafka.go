package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/Dyescape/event-indexer/internal/elastic"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Kafka struct {
	Brokers []string
	Topic   string
	Group   string
	Elastic *elastic.ElasticSearchClient
}

func (k *Kafka) Consume() {
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               k.Brokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         k.Group,
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	messages, err := subscriber.Subscribe(context.Background(), k.Topic)
	if err != nil {
		panic(err)
	}

	k.process(messages)
}

func (k *Kafka) PublishTestMessages() {

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   k.Brokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)

	if err != nil {
		panic(err)
	}

	for {
		msg := message.NewMessage(watermill.NewUUID(), []byte(`{"hello": "world"}`))

		if err := publisher.Publish(k.Topic, msg); err != nil {
			panic(err)
		}

		time.Sleep(time.Second)
	}
}

func (k *Kafka) process(messages <-chan *message.Message) {
	for msg := range messages {
		err := k.Elastic.Index(string(msg.Payload))
		if err != nil {
			fmt.Println("Failed to index event: " + err.Error())
			continue
		}

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
