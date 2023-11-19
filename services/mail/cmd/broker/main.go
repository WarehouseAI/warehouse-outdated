package broker

import (
	"fmt"
	"warehouseai/mail/config"

	rmq "github.com/rabbitmq/amqp091-go"
)

type MailConsumer struct {
	Channel    *rmq.Channel
	Connection *rmq.Connection
	Queue      rmq.Queue
}

func NewMailConsumer() *MailConsumer {
	config := config.NewMailBrokerCfg()

	fmt.Println(config)
	conn, err := rmq.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.User, config.Password, config.Host, config.Port))

	if err != nil {
		panic(fmt.Sprintf("Unable to open connect to RabbitMQ; %s", err))
	}

	ch, err := conn.Channel()

	if err != nil {
		panic(fmt.Sprintf("Unable to open channel; %s", err))
	}

	queue, err := ch.QueueDeclare(
		"mail",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(fmt.Sprintf("Unable to create queue in the channel; %s", err))
	}

	return &MailConsumer{
		Channel:    ch,
		Connection: conn,
		Queue:      queue,
	}
}
