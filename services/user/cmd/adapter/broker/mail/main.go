package mail

import (
	"fmt"
	"warehouseai/user/adapter/broker/mail"
	"warehouseai/user/config"

	rmq "github.com/rabbitmq/amqp091-go"
)

func NewMailProducer() *mail.MailProducer {
	config := config.NewMailBrokerCfg()

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

	return &mail.MailProducer{
		Channel:    ch,
		Queue:      queue,
		Connection: conn,
	}
}
