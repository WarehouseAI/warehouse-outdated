package broker

import (
	"context"
	"encoding/json"
	"warehouseai/user/dataservice"
	m "warehouseai/user/model"
	"warehouseai/user/service"

	rmq "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Broker struct {
	Connection *rmq.Connection
	Channel    *rmq.Channel
	MailQueue  rmq.Queue
	SagaQueues map[m.QueueName]rmq.Queue
}

func (b Broker) SendEmail(email m.Email) error {
	message := m.EmailRequest{Data: email}

	messageStr, err := json.Marshal(message)

	if err != nil {
		return err
	}

	if err := b.Channel.PublishWithContext(
		context.Background(),
		"",
		b.MailQueue.Name,
		false,
		false,
		rmq.Publishing{
			ContentType: "text/plain",
			Body:        []byte(messageStr),
		},
	); err != nil {
		return err
	}

	return nil
}

func (b Broker) ReceiveTokenReject(userRepository dataservice.UserInterface, logger *logrus.Logger) {
	messages, err := b.Channel.Consume(
		b.SagaQueues[m.Reject].Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	stop := make(chan bool)

	go func() {
		for message := range messages {
			var rejectEvent m.VerificationTokenRequest
			json.Unmarshal(message.Body, &rejectEvent)

			if err := service.Delete(rejectEvent.UserId, userRepository, logger); err != nil {
				continue
			}
		}
	}()

	<-stop
}
