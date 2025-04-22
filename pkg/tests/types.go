// Package tests provides test types and functions for yac-p
package tests

import (
	"context"
	"log/slog"
	"os"

	"github.com/kjansson/yac-p/pkg/types"
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

func (y *YaceMockClient) CollectMetrics(logger types.Logger, config types.YaceConfig) error {
	var err error
	ctx := context.Background()
	// Query metrics and resources and update the prometheus registry
	opts, err := config.GetYaceOptions(logger)
	if err != nil {
		return err
	}
	err = yace.UpdateMetrics(ctx, y.Logger, y.Config, y.Registry, y.Client, opts...)
	if err != nil {
		panic(err)
	}
	return nil
}

func (y *YaceMockClient) ExportMetrics(logger types.Logger) ([]*io_prometheus_client.MetricFamily, error) {
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

func (y *YaceMockClient) Init(types.Config) error {

	y.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	y.Registry = prometheus.NewRegistry() // Create a new prometheus registry

	return nil
}
