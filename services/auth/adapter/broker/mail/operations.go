package mail

import (
	"context"
	"encoding/json"
	"warehouseai/auth/model"

	rmq "github.com/rabbitmq/amqp091-go"
)

type MailProducer struct {
	Connection *rmq.Connection
	Channel    *rmq.Channel
	Queue      rmq.Queue
}

func (k MailProducer) SendEmail(email model.Email) error {
	message := model.EmailRequest{Data: email}

	messageStr, err := json.Marshal(message)

	if err != nil {
		return err
	}

	if err := k.Channel.PublishWithContext(
		context.Background(),
		"",
		k.Queue.Name,
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
