// Kafka producer
package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

// Define kafka producer
type Producer struct {
	brokers      []string
	syncProducer sarama.SyncProducer
}

// Configure and create new sync kafka producer
func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()

	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1
	cfg.Producer.CompressionLevel = sarama.CompressionLevelDefault
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Compression = sarama.CompressionGZIP

	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka producer")
	}

	return syncProducer, nil
}

// Create new kafka producer instance
func NewProducer(brokers []string) (*Producer, error) {
	p, err := newSyncProducer(brokers)
	if err != nil {
		return nil, errors.Wrap(err, "error creating producer")
	}

	return &Producer{
		brokers:      brokers,
		syncProducer: p,
	}, nil
}

// Send message to kafka broker
func (k *Producer) SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return k.syncProducer.SendMessage(message)
}

// Send pack of messages
func (k *Producer) SendSyncMessages(messages []*sarama.ProducerMessage) error {
	err := k.syncProducer.SendMessages(messages)
	if err != nil {
		return errors.Wrap(err, "kafka.Connector.Send Messages error")
	}

	return nil
}

// Close kafka producer
func (k *Producer) Close() error {
	err := k.syncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "kafka.Connection.Close")
	}

	return nil
}
