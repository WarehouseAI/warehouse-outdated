package mail

import (
	"warehouseai/user/adapter/broker/mail"
	"warehouseai/user/config"

	"github.com/IBM/sarama"
)

func NewMailProducer() *mail.MailProducer {
	config := config.NewMailBrokerCfg()

	brokerCfg := sarama.NewConfig()
	brokerCfg.Producer.RequiredAcks = sarama.WaitForAll
	brokerCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{config.Address}, brokerCfg)
	if err != nil {
		panic(err)
	}

	return &mail.MailProducer{
		Producer: producer,
		Topic:    config.Topic,
	}
}
