package kafka

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"

	"log"
)

type Kafka struct {
	Brokers []string
	Topic   string
	Group   string
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

	process(messages)
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

func process(messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
