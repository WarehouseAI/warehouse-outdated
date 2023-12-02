package config

import (
	"os"
	m "warehouseai/user/model"
)

type BrokerCfg struct {
	User     string
	Password string
	Host     string
	Port     string
	Queues   map[m.QueueName]string
}

func NewMailBrokerCfg() BrokerCfg {
	return BrokerCfg{
		User:     os.Getenv("RMQ_USER"),
		Password: os.Getenv("RMQ_PASS"),
		Host:     os.Getenv("RMQ_HOST"),
		Port:     os.Getenv("RMQ_PORT"),
		Queues: map[m.QueueName]string{
			m.Reject: os.Getenv("TOKEN_REJECTED_QUEUE"),
		},
	}
}
