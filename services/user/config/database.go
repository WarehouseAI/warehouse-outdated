package config

import (
	"os"
)

type DatabaseCfg struct {
	Host     string
	Name     string
	User     string
	Password string
	Port     string
}

func NewDatabaseCfg() DatabaseCfg {
	return DatabaseCfg{
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Port:     os.Getenv("DB_PORT"),
	}
}
