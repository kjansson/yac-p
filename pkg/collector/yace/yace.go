// Package yace provides a client for the Yet Another Cloudwatch Exporter (YACE) packages for collecting metrics from AWS Cloudwatch as well as managing
// It implements the types.MetricCollector interface.
package yace

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/kjansson/yac-p/v2/pkg/types"
	yace "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg"
	client "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/clients/v2"
	yace_config "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/config"
	"github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"gopkg.in/yaml.v2"
)

type YaceOpts struct {
	YaceCloudwatchConcurrencyPerApiLimitEnabled       string
	YaceCloudwatchConcurrencyListMetricsLimit         string
	YaceCloudwatchConcurrencyGetMetricDataLimit       string
	YaceCloudwatchConcurrencyGetMetricStatisticsLimit string
	YaceMetricsPerQuery                               string
	YaceTaggingAPIConcurrency                         string
	YaceCloudwatchConcurrency                         string
}

type YaceClient struct {
	Registry         *prometheus.Registry
	Client           *client.CachingFactory
	JobConfig        model.JobsConfig
	Logger           *slog.Logger
	YaceOpts         YaceOpts
	ConfigFileLoader func() ([]byte, error)
}

func NewYaceClient(
	configFileLoader func() ([]byte, error),
	yaceOpts YaceOpts,
) (*YaceClient, error) {
	var err error

	y := &YaceClient{
		ConfigFileLoader: configFileLoader,
		YaceOpts:         yaceOpts,
	}

	y.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	y.Registry = prometheus.NewRegistry() // Create a new prometheus registry

	contents, err := y.ConfigFileLoader()
	if err != nil {
		return nil, err
	}

	conf := yace_config.ScrapeConf{}
	err = yaml.Unmarshal(contents, &conf)
	if err != nil {
		return nil, err
	}

	for _, job := range conf.Discovery.Jobs {
		if len(job.Roles) == 0 {
			job.Roles = []yace_config.Role{{}} // use current IAM role
		}
	}

	for _, job := range conf.CustomNamespace {
		if len(job.Roles) == 0 {
			job.Roles = []yace_config.Role{{}} // use current IAM role
		}
	}

	for _, job := range conf.Static {
		if len(job.Roles) == 0 {
			job.Roles = []yace_config.Role{{}} // use current IAM role
		}
	}

	y.JobConfig, err = conf.Validate(y.Logger)
	if err != nil {
		return nil, err
	}

	for _, metric := range yace.Metrics { // Register YACE internal metrics
		err := y.Registry.Register(metric)
		if err != nil {
			return nil, err
		}
	}

	y.Client, err = client.NewFactory(y.Logger, y.JobConfig, false)
	if err != nil {
		return nil, err
	}
	return y, nil
}

// CollectMetrics performs the Cloudwatch metrics collection and updates the prometheus registry
func (y *YaceClient) CollectMetrics(logger types.Logger) error {
	ctx := context.Background()

	opts, err := getYaceOptions(y.YaceOpts, logger) // Get the YACE options from the config
	if err != nil {
		return err
	}
	// Query metrics and resources and update the prometheus registry
	err = yace.UpdateMetrics(ctx, y.Logger, y.JobConfig, y.Registry, y.Client, opts...)
	if err != nil {
		return err
	}
	return nil
}

// ExportMetrics exports metrics from the prometheus registry
func (y *YaceClient) ExportMetrics(logger types.Logger) ([]*io_prometheus_client.MetricFamily, error) {
	metrics, err := y.Registry.Gather() // Gather the metrics from the prometheus registry
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// GetRegistry returns the prometheus registry used by YACE, this is mostly used for testing
func (y *YaceClient) GetRegistry() *prometheus.Registry {
	return y.Registry
}

// See YACE documentation for more details on the options
type YaceOptions struct {
	YaceCloudwatchConcurrencyPerApiLimitEnabled       string
	YaceCloudwatchConcurrencyListMetricsLimit         string
	YaceCloudwatchConcurrencyGetMetricDataLimit       string
	YaceCloudwatchConcurrencyGetMetricStatisticsLimit string
	YaceMetricsPerQuery                               string
	YaceTaggingAPIConcurrency                         string
	YaceCloudwatchConcurrency                         string
}

// GetYaceOptions returns the YACE options function based on the config
func getYaceOptions(opts YaceOpts, logger types.Logger) ([]yace.OptionsFunc, error) {

	optFuncs := []yace.OptionsFunc{}
	var cloudwatchPerApiConcurrencyLimit bool
	var err error
	if opts.YaceCloudwatchConcurrencyPerApiLimitEnabled != "" {
		logger.Log("debug", "Using non-default per API concurrency limit", slog.String("per_api_limit", opts.YaceCloudwatchConcurrencyPerApiLimitEnabled))
		cloudwatchPerApiConcurrencyLimit, err = strconv.ParseBool(opts.YaceCloudwatchConcurrencyPerApiLimitEnabled)
		if err != nil {
			return nil, err
		}
	}
	if opts.YaceMetricsPerQuery != "" {
		logger.Log("debug", "Using non-default metrics per query", slog.String("metrics_per_query", opts.YaceMetricsPerQuery))
		val, err := strconv.Atoi(opts.YaceMetricsPerQuery)
		if err != nil {
			return nil, err
		}
		optFuncs = append(optFuncs, yace.MetricsPerQuery(val))
	}
	if opts.YaceTaggingAPIConcurrency != "" {
		logger.Log("debug", "Using non-default tagging API concurrency", slog.String("tagging_api_concurrency", opts.YaceTaggingAPIConcurrency))
		val, err := strconv.Atoi(opts.YaceTaggingAPIConcurrency)
		if err != nil {
			return nil, err
		}
		optFuncs = append(optFuncs, yace.TaggingAPIConcurrency(val))
	}

	if !cloudwatchPerApiConcurrencyLimit {
		if opts.YaceCloudwatchConcurrency != "" {
			logger.Log("debug", "Using non-default cloudwatch concurrency", slog.String("cloudwatch_concurrency", opts.YaceCloudwatchConcurrency))
			val, err := strconv.Atoi(opts.YaceCloudwatchConcurrency)
			if err != nil {
				return nil, err
			}
			optFuncs = append(optFuncs, yace.CloudWatchAPIConcurrency(val))
		}
	} else {
		limits := yace.DefaultCloudwatchConcurrency
		if opts.YaceCloudwatchConcurrencyListMetricsLimit != "" {
			logger.Log("debug", "Using non-default cloudwatch list metrics concurrency", slog.String("cloudwatch_list_metrics_concurrency", opts.YaceCloudwatchConcurrencyListMetricsLimit))
			val, err := strconv.Atoi(opts.YaceCloudwatchConcurrencyListMetricsLimit)
			if err != nil {
				return nil, err
			}
			limits.ListMetrics = val
		}
		if opts.YaceCloudwatchConcurrencyGetMetricDataLimit != "" {
			logger.Log("debug", "Using non-default cloudwatch get metric data concurrency", slog.String("cloudwatch_get_metric_data_concurrency", opts.YaceCloudwatchConcurrencyGetMetricDataLimit))
			val, err := strconv.Atoi(opts.YaceCloudwatchConcurrencyGetMetricDataLimit)
			if err != nil {
				return nil, err
			}
			limits.GetMetricData = val
		}
		if opts.YaceCloudwatchConcurrencyGetMetricStatisticsLimit != "" {
			logger.Log("debug", "Using non-default cloudwatch get metric statistics concurrency", slog.String("cloudwatch_get_metric_statistics_concurrency", opts.YaceCloudwatchConcurrencyGetMetricStatisticsLimit))
			val, err := strconv.Atoi(opts.YaceCloudwatchConcurrencyGetMetricStatisticsLimit)
			if err != nil {
				return nil, err
			}
			limits.GetMetricStatistics = val
		}
		optFuncs = append(optFuncs, yace.CloudWatchPerAPILimitConcurrency(limits.ListMetrics, limits.GetMetricData, limits.GetMetricStatistics))
	}
	return optFuncs, nil
}
