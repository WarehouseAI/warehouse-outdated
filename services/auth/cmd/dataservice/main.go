package dataservice

import (
	"fmt"
	"warehouseai/auth/config"
	"warehouseai/auth/dataservice/picturedata"
	"warehouseai/auth/dataservice/sessiondata"
	"warehouseai/auth/dataservice/tokendata"
	"warehouseai/auth/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPictureStorage() *picturedata.Storage {
	config := config.NewStorageCfg()

	sess, err := session.NewSession(
		&aws.Config{
			Endpoint:            aws.String(config.Endpoint),
			Region:              aws.String(config.Region),
			STSRegionalEndpoint: endpoints.RegionalSTSEndpoint,
			Credentials: credentials.NewStaticCredentials(
				config.AccessKey,
				config.SecretKey,
				"",
			),
		},
	)

	if err != nil {
		fmt.Println("❌Failed to connect to the S3 storage.")
		panic(err)
	}

	return &picturedata.Storage{
		Bucket:  config.Bucket,
		Domain:  config.Domain,
		Session: sess,
	}
}

func NewSessionDatabase() *sessiondata.Database {
	config := config.NewSessionCfg()

	rClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       0,
	})

	return &sessiondata.Database{
		DB: rClient,
	}
}

func NewResetTokenDatabase() *tokendata.Database {
	cfg := config.NewTokenDatabaseCfg()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.ResetToken{})

	return &tokendata.Database{DB: db}
}
