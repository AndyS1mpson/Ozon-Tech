// Consumer group
package kafka

import (
	"context"
	"encoding/json"
	"route256/notifications/internal/model"
	"route256/notifications/internal/pkg/logger"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

// Define Service for send message to notify user
type MessageSenderService interface {
	Save(ctx context.Context, message model.OrderStatusMessage) (model.MessageID, error)
	NotifyUser(ctx context.Context, message model.OrderStatusMessage) error
}

// Define Consumer group for order status messages
type ConsumerGroupHandler struct {
	ready       chan bool
	readyCloser sync.Once
	service     MessageSenderService
}

// Create a new consumer group
func NewConsumerGroupHandler(Service MessageSenderService) ConsumerGroupHandler {
	return ConsumerGroupHandler{
		ready:   make(chan bool),
		service: Service,
	}
}

// Check if the consumer group is ready
func (cg *ConsumerGroupHandler) Ready() <-chan bool {
	return cg.ready
}

// Starting a new session, before ConsumeClaim
func (cg *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	cg.readyCloser.Do(func() {
		close(cg.ready)
	})
	return nil
}

// Cleanup ends the session, after all ConsumeClaims are finished
func (cg *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Read messaged until the session is over and write it to Messages channel
func (cg *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			pm := model.OrderStatusMessage{}
			err := json.Unmarshal(message.Value, &pm)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal")
			}

			// Save message to storage
			_, err = cg.service.Save(context.Background(), pm)
			if err != nil {
				return errors.Wrapf(err, "failed to save message")
			}

			err = cg.service.NotifyUser(context.Background(), pm)
			if err != nil {
				return errors.Wrapf(err, "failed to send message")
			}
			logger.Info(pm)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
