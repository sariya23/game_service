package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func MustNewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(fmt.Sprintf("cannot connect to kafka; err = %v", err))
	}
	return &KafkaProducer{producer: producer, topic: topic}
}

func (kafka *KafkaProducer) SendMessage(message string) error {
	panic("impl me")
}
