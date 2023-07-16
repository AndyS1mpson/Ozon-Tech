// Kafka sender
package sender

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"log"
	"route256/loms/internal/domain"
	"route256/loms/internal/kafka"
)

// Define kafka sender
type KafkaSender struct {
	producer *kafka.Producer
	topic    string
}

// Create new kafka sender
func NewKafkaSender(producer *kafka.Producer, topic string) *KafkaSender {
	return &KafkaSender{
		producer: producer,
		topic:    topic,
	}
}

// Send messsage to Kafka
func (s *KafkaSender) SendMessage(message domain.OrderStatusNotification) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		return errors.Wrap(err, "fail build message")
	}

	partition, offset, err := s.producer.SendSyncMessage(kafkaMsg)

	if err != nil {
		return errors.Wrap(err, "fail send message")
	}

	log.Printf("Partition: %v, Offset: %v, OrderID: %v\n", partition, offset, message.OrderID)
	return nil
}

// Send pack of messages
func (s *KafkaSender) SendMessages(messages []domain.OrderStatusNotification) error {
	var kafkaMsg []*sarama.ProducerMessage
	var message *sarama.ProducerMessage
	var err error

	for _, m := range messages {
		message, err = s.buildMessage(m)
		kafkaMsg = append(kafkaMsg, message)

		if err != nil {
			return errors.Wrap(err, "fail build message")
		}
	}

	err = s.producer.SendSyncMessages(kafkaMsg)

	if err != nil {
		return errors.Wrap(err, "fail send message")
	}

	return nil
}

// Create kafka message from input data
func (s *KafkaSender) buildMessage(message domain.OrderStatusNotification) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)

	if err != nil {
		return nil, errors.Wrap(err, "send message marshal error")
	}

	return &sarama.ProducerMessage{
		Topic: s.topic,
		Value: sarama.ByteEncoder(msg),
		Key:   sarama.StringEncoder(fmt.Sprint(message.OrderID)),
	}, nil
}
