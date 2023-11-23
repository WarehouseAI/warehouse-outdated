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

type StorageCfg struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
	Domain    string
	Bucket    string
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

func NewStorageCfg() StorageCfg {
	return StorageCfg{
		Endpoint:  os.Getenv("S3_HOST"),
		AccessKey: os.Getenv("S3_ACCESSKEY"),
		SecretKey: os.Getenv("S3_SECRETKEY"),
		Domain:    os.Getenv("S3_LINK"),
		Bucket:    os.Getenv("S3_BUCKET"),
		Region:    os.Getenv("S3_REGION"),
	}
}
