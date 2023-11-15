package cmd

import (
	"fmt"
	"os"
	"warehouseai/mail/cmd/broker"
	"warehouseai/mail/cmd/mail"
	"warehouseai/mail/cmd/server"
	"warehouseai/mail/config"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	file, err := os.OpenFile("./mail.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	config := config.NewMailCfg()
	mailDialer := mail.NewMailDialer(config)

	handlers := map[string]sarama.ConsumerGroupHandler{
		os.Getenv("KAFKA_MAIL_TOPIC"): server.NewMailHandler(mailDialer, log, config.Sender),
	}

	broker.RunConsumers(log, handlers)
}
