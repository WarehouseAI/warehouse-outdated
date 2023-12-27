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

func NewStatDatabaseCfg() DatabaseCfg {
	return DatabaseCfg{
		Host:     os.Getenv("STAT_DB_HOST"),
		Name:     os.Getenv("STAT_DB_NAME"),
		User:     os.Getenv("STAT_DB_USER"),
		Password: os.Getenv("STAT_DB_PASS"),
		Port:     "5432",
	}
}

// TODO: Этот моментик мне не нравится, но пока я не могу придумать как делать по другому
func UserDatabaseCfg() DatabaseCfg {
	return DatabaseCfg{
		Host:     os.Getenv("USERS_DB_HOST"),
		Name:     os.Getenv("USERS_DB_NAME"),
		User:     os.Getenv("USERS_DB_USER"),
		Password: os.Getenv("USERS_DB_PASS"),
		Port:     "5432",
	}
}

func AiDatabaseCfg() DatabaseCfg {
	return DatabaseCfg{
		Host:     os.Getenv("AI_DB_HOST"),
		Name:     os.Getenv("AI_DB_NAME"),
		User:     os.Getenv("AI_DB_USER"),
		Password: os.Getenv("AI_DB_PASS"),
		Port:     "5432",
	}
}
