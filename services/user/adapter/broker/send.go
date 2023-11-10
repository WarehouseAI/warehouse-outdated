package broker

import (
	"encoding/json"
	"warehouseai/user/model"

	"github.com/IBM/sarama"
)

func SendEmail(message model.Email, topic string, kafka sarama.SyncProducer) error {
	msgStr, err := json.Marshal(message)

	if err != nil {
		return err
	}

	producerMsg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(msgStr)}
	_, _, err = kafka.SendMessage(producerMsg)

	if err != nil {
		return err
	}

	return nil
}
