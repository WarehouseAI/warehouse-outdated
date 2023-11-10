package broker

import (
	"os"

	"github.com/IBM/sarama"
)

func InitKafkaProducer() sarama.SyncProducer {
	brokerCfg := sarama.NewConfig()

	brokerCfg.Producer.RequiredAcks = sarama.WaitForAll
	brokerCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_ADDR")}, brokerCfg)
	if err != nil {
		panic(err)
	}

	return producer
}
