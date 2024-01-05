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
	SubName  string
	PubName  string
}

func NewStatDatabaseCfg() DatabaseCfg {
	return DatabaseCfg{
		Host:     os.Getenv("STAT_DB_HOST"),
		Name:     os.Getenv("STAT_DB_NAME"),
		User:     os.Getenv("STAT_DB_USER"),
		Password: os.Getenv("STAT_DB_PASS"),
		Port:     "5432",
	}
}

func NewUserDatabaseCfg() DatabaseCfg {
	return DatabaseCfg{
		Host:     os.Getenv("USERS_DB_HOST"),
		Name:     os.Getenv("USERS_DB_NAME"),
		User:     os.Getenv("USERS_DB_USER"),
		Password: os.Getenv("USERS_DB_PASS"),
		SubName:  os.Getenv("USERS_DB_SUB_NAME"),
		PubName:  os.Getenv("USERS_DB_PUB_NAME"),
		Port:     "5432",
	}
}

func NewAiDatabaseCfg() DatabaseCfg {
	return DatabaseCfg{
		Host:     os.Getenv("AI_DB_HOST"),
		Name:     os.Getenv("AI_DB_NAME"),
		User:     os.Getenv("AI_DB_USER"),
		Password: os.Getenv("AI_DB_PASS"),
		SubName:  os.Getenv("AI_DB_SUB_NAME"),
		PubName:  os.Getenv("AI_DB_PUB_NAME"),
		Port:     "5432",
	}
}
