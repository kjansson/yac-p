package loaders

// Package loaders provides functions to load configuration files from different sources

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetS3Loader returns a function that loads the config from S3
func GetS3Loader() func() ([]byte, error) {
	return func() ([]byte, error) {
		var content []byte
		configS3Path, configS3Bucket := os.Getenv("CONFIG_S3_PATH"), os.Getenv("CONFIG_S3_BUCKET")
		if configS3Bucket != "" && configS3Path != "" {
			sess, err := session.NewSessionWithOptions(session.Options{
				Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
				SharedConfigState: session.SharedConfigEnable,
			})
			if err != nil {
				return nil, err
			}
			s3svc := s3.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
			obj, err := s3svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(configS3Bucket),
				Key:    aws.String(configS3Path),
			})
			if err != nil {
				return nil, err
			}
			content, err = io.ReadAll(obj.Body)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("CONFIG_S3_BUCKET and CONFIG_S3_PATH is required")
		}
		return content, nil
	}
}

// GetLocalFileLoader returns a function that loads the config from a local file
func GetLocalFileLoader() func() (content []byte, err error) {
	return func() (content []byte, err error) {
		configFilePath := os.Getenv("CONFIG_FILE_PATH")
		if configFilePath != "" {
			file, err := os.Open(configFilePath)
			if err != nil {
				return nil, err
			}

			defer func() {
				err = file.Close()
			}()

			content, err = io.ReadAll(file)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("CONFIG_FILE_PATH is required")
		}
		return content, nil
	}
}
