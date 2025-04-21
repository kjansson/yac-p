package main

import (
	"context"
	"log/slog"
	"os"

	yace "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg"
	client "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/clients/v2"
	"github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

type YaceMockClient struct {
	Registry *prometheus.Registry
	Client   *client.CachingFactory
	Config   model.JobsConfig
	Logger   *slog.Logger
}

func (y *YaceMockClient) CollectMetrics(logger Logger, config Config) error {
	var err error
	ctx := context.Background()
	// Query metrics and resources and update the prometheus registry
	err = yace.UpdateMetrics(ctx, y.Logger, y.Config, y.Registry, y.Client, config.GetYaceOptions(logger)...)
	if err != nil {
		panic(err)
	}
	return nil
}

func (y *YaceMockClient) ExtractMetrics(logger Logger) ([]*io_prometheus_client.MetricFamily, error) {
	var err error

	metrics, err := y.Registry.Gather() // Gather the metrics from the prometheus registry
	if err != nil {
		panic(err)
	}
	return metrics, nil
}

func (y *YaceMockClient) GetRegistry() *prometheus.Registry {
	return y.Registry
}

func (y *YaceMockClient) Init() error {

	y.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	y.Registry = prometheus.NewRegistry() // Create a new prometheus registry

	return nil
}
