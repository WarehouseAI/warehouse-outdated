package broker

import (
	"context"
	"encoding/json"
	m "warehouseai/auth/model"

	rmq "github.com/rabbitmq/amqp091-go"
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

func (b Broker) SendUserReject(userId string) error {
	message := m.VerificationTokenRequest{UserId: userId}

	messageStr, err := json.Marshal(message)

	if err != nil {
		return err
	}

	if err := b.Channel.PublishWithContext(
		context.Background(),
		"",
		b.SagaQueues[m.Reject].Name,
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
