package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	yace "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg"
	client "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/clients/v2"
	"github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/config"
	"github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

type YaceClient struct {
	registry *prometheus.Registry
	client   *client.CachingFactory
	config   model.JobsConfig
	logger   *slog.Logger
}

func (y *YaceClient) GetMetrics(logger logger, config Config) ([]*io_prometheus_client.MetricFamily, error) {
	var err error
	ctx := context.Background()
	// Query metrics and resources and update the prometheus registry
	err = yace.UpdateMetrics(ctx, y.logger, y.config, y.registry, y.client, config.GetYaceOptions(logger)...)
	if err != nil {
		panic(err)
	}

	metrics, err := y.registry.Gather() // Gather the metrics from the prometheus registry
	if err != nil {
		panic(err)
	}
	return metrics, nil
}

func (y *YaceClient) Init() error {
	var err error

	y.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	y.registry = prometheus.NewRegistry() // Create a new prometheus registry

	configS3Path, configS3Bucket := os.Getenv("CONFIG_S3_PATH"), os.Getenv("CONFIG_S3_BUCKET")
	if configS3Bucket != "" && configS3Path != "" {
		sess, err := createAWSSession()
		if err != nil {
			return err
		}
		y.logger.Debug("Using S3 config", slog.String("bucket", configS3Bucket), slog.String("path", configS3Path))
		s3svc := s3.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
		obj, err := s3svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(configS3Bucket),
			Key:    aws.String(configS3Path),
		})
		if err != nil {
			return err
		}
		content, err := io.ReadAll(obj.Body)
		if err != nil {
			return err
		}
		// Write the configuration to a ephemeral storage, this is needed since the config package expects a file path
		err = os.WriteFile(configFilePath, content, 0644)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("CONFIG_S3_BUCKET and CONFIG_S3_PATH is required")
	}
	config := config.ScrapeConf{}
	y.config, err = config.Load(configFilePath, y.logger)
	if err != nil {
		panic(err)
	}

	for _, metric := range yace.Metrics { // Register YACE internal metrics
		err := y.registry.Register(metric)
		if err != nil {
			return err
		}
	}

	// Create a new yace client factory
	y.client, err = client.NewFactory(y.logger, y.config, false)
	if err != nil {
		return err
	}
	return nil
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
