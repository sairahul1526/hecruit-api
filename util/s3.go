package util

import (
	"fmt"
	"strconv"
	"time"

	CONFIG "hecruit-backend/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GeneratePUTSignedURL(path, fileType string) (string, string, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(CONFIG.S3ID, CONFIG.S3Secret, ""),
		Endpoint:         aws.String(CONFIG.S3Endpoint),
		Region:           aws.String(CONFIG.S3Region),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	fileName := path + strconv.FormatInt(time.Now().Unix(), 10) + "." + fileType

	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(CONFIG.S3Bucket),
		Key:    aws.String(fileName),
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		fmt.Println("GeneratePUTSignedURL", err)
		return "", "", err
	}

	return fileName, url, nil
}
