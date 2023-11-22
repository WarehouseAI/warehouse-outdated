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

func NewAiDatabaseCfg() DatabaseCfg {
	return DatabaseCfg{
		Host:     os.Getenv("AI_DB_HOST"),
		Name:     os.Getenv("AI_DB_NAME"),
		User:     os.Getenv("AI_DB_USER"),
		Password: os.Getenv("AI_DB_PASS"),
		Port:     "5432",
	}
}
