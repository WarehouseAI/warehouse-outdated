package config

import "os"

type StorageCfg struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
	Domain    string
	Bucket    string
}

type SessionCfg struct {
	Host     string
	Port     string
	Password string
}

type DatabaseCfg struct {
	Host     string
	Name     string
	User     string
	Password string
	Port     string
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

func NewSessionCfg() SessionCfg {
	return SessionCfg{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
}

func NewTokenDatabaseCfg() DatabaseCfg {
	return DatabaseCfg{
		Host:     os.Getenv("AUTH_DB_HOST"),
		Name:     os.Getenv("AUTH_DB_NAME"),
		User:     os.Getenv("AUTH_DB_USER"),
		Password: os.Getenv("AUTH_DB_PASS"),
		Port:     os.Getenv("AUTH_DB_PORT"),
	}
}
