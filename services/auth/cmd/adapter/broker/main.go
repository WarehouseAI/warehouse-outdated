package broker

import (
	"fmt"
	"warehouseai/auth/adapter/broker"
	"warehouseai/auth/config"
	m "warehouseai/auth/model"

	rmq "github.com/rabbitmq/amqp091-go"
)

func NewBroker() *broker.Broker {
	cfg := config.NewMailBrokerCfg()

	conn, err := rmq.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port))

	if err != nil {
		panic(fmt.Sprintf("Unable to open connect to RabbitMQ; %s", err))
	}

	ch, err := conn.Channel()

	if err != nil {
		panic(fmt.Sprintf("Unable to open channel; %s", err))
	}

	var sagaQueues = map[m.QueueName]rmq.Queue{}
	for key, element := range cfg.Queues {
		queue, err := ch.QueueDeclare(
			element,
			false,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			panic(fmt.Sprintf("Unable to create queue in the channel; %s", err))
		}

		sagaQueues[key] = queue
	}

	mailQueue, err := ch.QueueDeclare(
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

	return &broker.Broker{
		Channel:    ch,
		Connection: conn,
		MailQueue:  mailQueue,
		SagaQueues: sagaQueues,
	}
}
