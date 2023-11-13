package mail

import (
	"warehouseai/auth/adapter/broker/mail"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

func InitKafkaProducer(addr string, topic string, logger *logrus.Logger) *mail.MailProducer {
	brokerCfg := sarama.NewConfig()
	brokerCfg.Producer.RequiredAcks = sarama.WaitForAll
	brokerCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{addr}, brokerCfg)
	if err != nil {
		panic(err)
	}

	return &mail.MailProducer{
		Producer: producer,
		Topic:    topic,
	}
}
