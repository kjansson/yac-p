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
	Registry *prometheus.Registry
	Client   *client.CachingFactory
	Config   model.JobsConfig
	Logger   *slog.Logger
}

func (y *YaceClient) CollectMetrics(logger logger, config Config) error {
	var err error
	ctx := context.Background()
	// Query metrics and resources and update the prometheus registry
	err = yace.UpdateMetrics(ctx, y.Logger, y.Config, y.Registry, y.Client, config.GetYaceOptions(logger)...)
	if err != nil {
		panic(err)
	}
	return nil
}

func (y *YaceClient) ExtractMetrics(logger logger) ([]*io_prometheus_client.MetricFamily, error) {
	var err error

	metrics, err := y.Registry.Gather() // Gather the metrics from the prometheus registry
	if err != nil {
		panic(err)
	}
	return metrics, nil
}

func (y *YaceClient) Init() error {
	var err error

	y.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	y.Registry = prometheus.NewRegistry() // Create a new prometheus registry

	configS3Path, configS3Bucket := os.Getenv("CONFIG_S3_PATH"), os.Getenv("CONFIG_S3_BUCKET")
	if configS3Bucket != "" && configS3Path != "" {
		sess, err := createAWSSession()
		if err != nil {
			return err
		}
		y.Logger.Debug("Using S3 config", slog.String("bucket", configS3Bucket), slog.String("path", configS3Path))
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
	y.Config, err = config.Load(configFilePath, y.Logger)
	if err != nil {
		panic(err)
	}

	for _, metric := range yace.Metrics { // Register YACE internal metrics
		err := y.Registry.Register(metric)
		if err != nil {
			return err
		}
	}

	// Create a new yace client factory
	y.Client, err = client.NewFactory(y.Logger, y.Config, false)
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
