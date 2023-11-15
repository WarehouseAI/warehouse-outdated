package config

import "os"

type BrokerCfg struct {
	Address string
	Topic   string
}

func NewMailBrokerCfg() BrokerCfg {
	return BrokerCfg{
		Address: os.Getenv("KAFKA_MAIL_ADDRESS"),
		Topic:   os.Getenv("KAFKA_MAIL_TOPIC"),
	}
}
