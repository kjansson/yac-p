package main

import (
	"os"

	"github.com/kjansson/yac-p/v3/pkg/collector/yace"
	"github.com/kjansson/yac-p/v3/pkg/converter"
	"github.com/kjansson/yac-p/v3/pkg/logger"
	"github.com/kjansson/yac-p/v3/pkg/persister/prom"
	"github.com/kjansson/yac-p/v3/pkg/types"
)

func NewController(config Config) (*types.Controller, error) {

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

	persister, err := prom.NewPromClient(
		config.RemoteWriteURL,
		config.AuthType,
		config.AuthToken,
		config.Username,
		config.Password,
		config.Region,
		config.PrometheusRegion,
		config.AWSRoleARN,
	)
	if err != nil {
		return nil, err
	}

	c := &types.Controller{
		Logger:    logger,
		Collector: collector,
		Converter: converter,
		Persister: persister,
	}

	return c, nil
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
