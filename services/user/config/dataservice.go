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

func NewUserDatabaseCfg() DatabaseCfg {
	return DatabaseCfg{
		Host:     os.Getenv("USERS_DB_HOST"),
		Name:     os.Getenv("USERS_DB_NAME"),
		User:     os.Getenv("USERS_DB_USER"),
		Password: os.Getenv("USERS_DB_PASS"),
		Port:     os.Getenv("USERS_DB_PORT"),
	}
}
