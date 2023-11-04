package s3

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
)

type S3Storage struct {
	Bucket  string
	Host    string
	Session *session.Session
}

func NewS3(host string, bucket string) (*S3Storage, error) {
	endpoint := os.Getenv("S3_HOST")
	accessKey := os.Getenv("S3_ACCESSKEY")
	secretKey := os.Getenv("S3_SECRETKEY")
	region := os.Getenv("S3_REGION")

	sess, err := session.NewSession(
		&aws.Config{
			Endpoint:            aws.String(endpoint),
			Region:              aws.String(region),
			STSRegionalEndpoint: endpoints.RegionalSTSEndpoint,
			Credentials: credentials.NewStaticCredentials(
				accessKey,
				secretKey,
				"",
			),
		},
	)

	if err != nil {
		return nil, err
	}

	return &S3Storage{
		Bucket:  bucket,
		Host:    host,
		Session: sess,
	}, nil
}
