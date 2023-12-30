package dataservice

import (
	"fmt"
	"warehouseai/ai/config"
	"warehouseai/ai/dataservice/aidata"
	"warehouseai/ai/dataservice/commanddata"
	"warehouseai/ai/dataservice/picturedata"
	"warehouseai/ai/dataservice/ratingdata"
	"warehouseai/ai/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewAiDatabase() *aidata.Database {
	cfg := config.NewAiDatabaseCfg()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.AI{})

	return &aidata.Database{DB: db}
}

func NewCommandDatabase() *commanddata.Database {
	cfg := config.NewAiDatabaseCfg()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.Command{})

	return &commanddata.Database{DB: db}
}

func NewRatingDatabase() *ratingdata.Database {
	cfg := config.NewAiDatabaseCfg()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.RatingPerUser{})

	return &ratingdata.Database{DB: db}
}

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
