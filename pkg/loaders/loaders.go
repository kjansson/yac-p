package loaders

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetS3Loader() func() ([]byte, error) {
	return func() ([]byte, error) {
		var content []byte
		configS3Path, configS3Bucket := os.Getenv("CONFIG_S3_PATH"), os.Getenv("CONFIG_S3_BUCKET")
		if configS3Bucket != "" && configS3Path != "" {
			sess, err := createAWSSession()
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

func createAWSSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return sess, err
	}
	return sess, err
}
