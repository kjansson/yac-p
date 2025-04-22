package types

// Package types defines interfaces and constants used in yac-p

import (
	yace "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

type Logger interface {
	Init(Config) error
	Log(level string, msg string, args ...any)
}

type MetricCollector interface {
	Init(Config) error
	CollectMetrics(Logger, YaceConfig) error
	ExportMetrics(Logger) ([]*io_prometheus_client.MetricFamily, error)
	GetRegistry() *prometheus.Registry
}

type MetricPersister interface {
	Init(Config) error
	PersistMetrics([]prompb.TimeSeries, Logger) error
}

type YaceConfig interface {
	Init(Config) error
	GetYaceOptions(logger Logger) ([]yace.OptionsFunc, error)
}

type Config struct {
	Debug                                             bool
	RemoteWriteURL                                    string
	AuthType                                          string
	AuthToken                                         string
	Username                                          string
	Password                                          string
	Region                                            string
	PrometheusRegion                                  string
	AWSRoleARN                                        string
	YaceCloudwatchConcurrencyPerApiLimitEnabled       string
	YaceCloudwatchConcurrencyListMetricsLimit         string
	YaceCloudwatchConcurrencyGetMetricDataLimit       string
	YaceCloudwatchConcurrencyGetMetricStatisticsLimit string
	YaceMetricsPerQuery                               string
	YaceTaggingAPIConcurrency                         string
	YaceCloudwatchConcurrency                         string
	ConfigFileLoader                                  func() ([]byte, error)
}
