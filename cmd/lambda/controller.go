// Package controller provides the yac-p controller struct and its methods to manage the components of the application
package main

import (
	"os"

	"github.com/kjansson/yac-p/v2/pkg/collector/yace"
	"github.com/kjansson/yac-p/v2/pkg/converter"
	"github.com/kjansson/yac-p/v2/pkg/logger"
	"github.com/kjansson/yac-p/v2/pkg/persister/prom"
	"github.com/kjansson/yac-p/v2/pkg/types"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

type Controller struct {
	Logger    types.Logger
	Collector types.MetricCollector
	Converter types.MetricConverter
	Persister types.MetricPersister
	Config    Config
}

func NewController(config Config) (*Controller, error) {

	logger, err := logger.NewLogger(config.LogDestination, config.LogFormat, config.Debug)
	if err != nil {
		return nil, err
	}

	collector, err := yace.NewYaceClient(
		config.ConfigFileLoader,
		yace.YaceOpts{
			YaceCloudwatchConcurrencyPerApiLimitEnabled:       config.YaceCloudwatchConcurrencyPerApiLimitEnabled,
			YaceCloudwatchConcurrencyListMetricsLimit:         config.YaceCloudwatchConcurrencyListMetricsLimit,
			YaceCloudwatchConcurrencyGetMetricDataLimit:       config.YaceCloudwatchConcurrencyGetMetricDataLimit,
			YaceCloudwatchConcurrencyGetMetricStatisticsLimit: config.YaceCloudwatchConcurrencyGetMetricStatisticsLimit,
			YaceMetricsPerQuery:                               config.YaceMetricsPerQuery,
			YaceTaggingAPIConcurrency:                         config.YaceTaggingAPIConcurrency,
			YaceCloudwatchConcurrency:                         config.YaceCloudwatchConcurrency,
		},
	)
	if err != nil {
		return nil, err
	}

	converter := converter.NewConverter(logger)

	persister := &prom.PromClient{
		RemoteWriteURL:   config.RemoteWriteURL,
		AuthType:         config.AuthType,
		AuthToken:        config.AuthToken,
		Username:         config.Username,
		Password:         config.Password,
		Region:           config.Region,
		PrometheusRegion: config.PrometheusRegion,
		AWSRoleARN:       config.AWSRoleARN,
	}

	c := &Controller{
		Logger:    logger,
		Collector: collector,
		Converter: converter,
		Persister: persister,
	}

	return c, nil
}

// Log extends the logger interface
func (c *Controller) Log(level string, msg string, args ...any) {
	c.Logger.Log(level, msg, args...)
}

// GetRegistry returns the prometheus registry from the Collector component
func (c *Controller) CollectMetrics() error {
	return c.Collector.CollectMetrics(c.Logger)
}

// GetRegistry extracts the metrics from the prometheus registry in the Collector component
func (c *Controller) ExportMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	metrics, err := c.Collector.ExportMetrics(c.Logger)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// PersistMetrics extends the underlying method and persists timeseries to the remote write endpoint
func (c *Controller) PersistMetrics(timeSeries []prompb.TimeSeries) error {
	return c.Persister.PersistMetrics(timeSeries, c.Logger)
}

func (c *Controller) ConvertMetrics(metrics []*io_prometheus_client.MetricFamily) ([]prompb.TimeSeries, error) {
	timeSeries, err := c.Converter.ConvertMetrics(metrics, c.Logger)
	if err != nil {
		return nil, err
	}
	return timeSeries, nil
}

type Config struct {
	Debug                                             bool   `env:"DEBUG"`
	RemoteWriteURL                                    string `env:"PROMETHEUS_REMOTE_WRITE_URL"`
	AuthType                                          string `env:"AUTH_TYPE"`
	AuthToken                                         string `env:"AUTH_TYPE"`
	Username                                          string `env:"USERNAME"`
	Password                                          string `env:"PASSWORD"`
	Region                                            string `env:"AWS_REGION"`
	PrometheusRegion                                  string `env:"PROMETHEUS_REGION"`
	AWSRoleARN                                        string `env:"AWS_ROLE_ARN"`
	YaceCloudwatchConcurrencyPerApiLimitEnabled       string `env:"YACE_CLOUDWATCH_CONCURRENCY_PER_API_LIMIT_ENABLED"`
	YaceCloudwatchConcurrencyListMetricsLimit         string `env:"YACE_CLOUDWATCH_CONCURRENCY_LIST_METRICS_LIMIT"`
	YaceCloudwatchConcurrencyGetMetricDataLimit       string `env:"YACE_CLOUDWATCH_CONCURRENCY_GET_METRIC_DATA_LIMIT"`
	YaceCloudwatchConcurrencyGetMetricStatisticsLimit string `env:"YACE_CLOUDWATCH_CONCURRENCY_GET_METRIC_STATISTICS_LIMIT"`
	YaceMetricsPerQuery                               string `env:"YACE_METRICS_PER_QUERY"`
	YaceTaggingAPIConcurrency                         string `env:"YACE_TAGGING_API_CONCURRENCY"`
	YaceCloudwatchConcurrency                         string `env:"YACE_CLOUDWATCH_CONCURRENCY"`
	ConfigFileLoader                                  func() ([]byte, error)
	LogFormat                                         string `env:"LOG_FORMAT"`
	LogLevel                                          string `env:"LOG_LEVEL"`
	LogDestination                                    *os.File
}
