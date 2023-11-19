package config

import "os"

type BrokerCfg struct {
	User     string
	Password string
	Host     string
	Port     string
}

func NewMailBrokerCfg() BrokerCfg {
	return BrokerCfg{
		User:     os.Getenv("RMQ_USER"),
		Password: os.Getenv("RMQ_PASS"),
		Host:     os.Getenv("RMQ_HOST"),
		Port:     os.Getenv("RMQ_PORT"),
	}
}
