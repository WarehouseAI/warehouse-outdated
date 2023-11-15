package broker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

func RunConsumers(logger *logrus.Logger, handlers map[string]sarama.ConsumerGroupHandler) {
	kafkaConsumerGroups := initAllConsumerGroups(logger)

	for topic, group := range kafkaConsumerGroups {
		go func(topic string, group *sarama.ConsumerGroup) {
			defer func() {
				if r := recover(); r != nil {
					logger.WithFields(logrus.Fields{"time": time.Now(), "error": r}).Info("Setup email")
					fmt.Println(r)
				}
			}()

			for {
				err := (*group).Consume(context.Background(), []string{topic}, handlers[topic])
				if err != nil {
					logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Setup email")
					fmt.Println("Consumer group error")
				}
			}
		}(topic, group)
	}
}

func initAllConsumerGroups(logger *logrus.Logger) map[string]*sarama.ConsumerGroup {
	return map[string]*sarama.ConsumerGroup{
		os.Getenv("KAFKA_MAIL_TOPIC"): initGroup(os.Getenv("KAFKA_MAIL_TOPIC"), logger),
	}
}

func initGroup(topic string, logger *logrus.Logger) *sarama.ConsumerGroup {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_3_0_0
	cfg.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup([]string{os.Getenv("KAFKA_ADDR")}, topic, cfg)
	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Setup email")
		fmt.Println("Message hasn't been marshaled.")
		return nil
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.WithFields(logrus.Fields{"time": time.Now(), "error": r}).Info("Setup email")
				fmt.Println(fmt.Sprintf("%s", r))
			}
		}()

		for err := range group.Errors() {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Setup email")
			fmt.Println("Consumer group error")
		}
	}()

	return &group
}
