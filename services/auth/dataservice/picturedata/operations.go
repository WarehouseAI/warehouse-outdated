package picturedata

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage struct {
	Bucket  string
	Host    string
	Session *session.Session
}

func (s *Storage) UploadFile(file multipart.File, fileName string) (string, error) {
	uploader := s3manager.NewUploader(s.Session)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	fileLink := s.Host + fileName

	return fileLink, nil
}

func (s *Storage) DeleteImage(fileName string) error {
	svc := s3.New(s.Session)
	objectInfo := s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fileName),
	}

	_, err := svc.DeleteObject(&objectInfo)

	if err != nil {
		return err
	}

	return nil
}
