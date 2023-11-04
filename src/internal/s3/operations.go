package s3

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (cfg *S3Storage) UploadFile(file multipart.File, fileName string) (string, *S3Error) {
	uploader := s3manager.NewUploader(cfg.Session)

	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cfg.Bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(fileName),
		Body:   file,
	})

	fmt.Println(output)

	if err != nil {
		fmt.Println(err)
		return "", NewS3Error(System, "Can't upload image.", err.Error())
	}

	fileLink := cfg.Host + fileName

	return fileLink, nil
}

func (cfg *S3Storage) DeleteImage(fileName string) *S3Error {
	svc := s3.New(cfg.Session)
	objectInfo := s3.DeleteObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(fileName),
	}

	_, err := svc.DeleteObject(&objectInfo)

	if err != nil {
		return NewS3Error(System, "Can't delete image.", err.Error())
	}

	return nil
}
