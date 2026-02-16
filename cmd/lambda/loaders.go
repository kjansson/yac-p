package main

// Loaders are used when initializing the YACE client and allows for flexibility in how the configuration is loaded.
// Custom loaders can be used if doing a custom implementation. A loader should return the YACE config file in a byte array along with any errors.
// Example: func MyCustomLoader() func() ([]byte, error)...

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetS3Loader returns a function that loads the config from S3
func GetS3Loader() func() ([]byte, error) {
	return func() (content []byte, err error) {
		configS3Path, configS3Bucket := os.Getenv("CONFIG_S3_PATH"), os.Getenv("CONFIG_S3_BUCKET")
		if configS3Bucket != "" && configS3Path != "" {
			ctx := context.TODO()

			// Load AWS SDK v2 config
			cfg, err := config.LoadDefaultConfig(ctx,
				config.WithRegion(os.Getenv("AWS_REGION")),
			)
			if err != nil {
				return nil, err
			}

			// Create S3 client
			s3svc := s3.NewFromConfig(cfg)

			// Get object from S3
			obj, err := s3svc.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &configS3Bucket,
				Key:    &configS3Path,
			})
			if err != nil {
				return nil, err
			}
			defer func() {
				if closeErr := obj.Body.Close(); closeErr != nil && err == nil {
					err = closeErr
				}
			}()

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
