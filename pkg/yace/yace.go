package yace

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/kjansson/yac-p/pkg/types"
	yace "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg"
	client "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/clients/v2"
	"github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/config"
	"github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"gopkg.in/yaml.v2"
)

type YaceClient struct {
	Registry *prometheus.Registry
	Client   *client.CachingFactory
	Config   model.JobsConfig
	Logger   *slog.Logger
}

// CollectMetrics performs the Cloudwatch metrics collection and updates the prometheus registry
func (y *YaceClient) CollectMetrics(logger types.Logger, config types.Config) error {
	ctx := context.Background()

	opts, err := config.GetYaceOptions(logger) // Get the YACE options from the config
	if err != nil {
		return err
	}
	// Query metrics and resources and update the prometheus registry
	err = yace.UpdateMetrics(ctx, y.Logger, y.Config, y.Registry, y.Client, opts...)
	if err != nil {
		return err
	}
	return nil
}

// ExtractMetrics gathers the metrics from the prometheus registry
func (y *YaceClient) ExtractMetrics(logger types.Logger) ([]*io_prometheus_client.MetricFamily, error) {
	metrics, err := y.Registry.Gather() // Gather the metrics from the prometheus registry
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// GetRegistry returns the prometheus registry
func (y *YaceClient) GetRegistry() *prometheus.Registry {
	return y.Registry
}

// Init initializes the YACE client and loads the configuration
func (y *YaceClient) Init(getConfig func() ([]byte, error)) error {
	var err error

	y.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	y.Registry = prometheus.NewRegistry() // Create a new prometheus registry

	contents, err := getConfig()
	if err != nil {
		return err
	}

	conf := config.ScrapeConf{}
	err = yaml.Unmarshal(contents, conf)
	if err != nil {
		return err
	}

	for _, job := range conf.Discovery.Jobs {
		if len(job.Roles) == 0 {
			job.Roles = []config.Role{{}} // use current IAM role
		}
	}

	for _, job := range conf.CustomNamespace {
		if len(job.Roles) == 0 {
			job.Roles = []config.Role{{}} // use current IAM role
		}
	}

	for _, job := range conf.Static {
		if len(job.Roles) == 0 {
			job.Roles = []config.Role{{}} // use current IAM role
		}
	}

	y.Config, err = conf.Validate(y.Logger)
	if err != nil {
		return err
	}

	for _, metric := range yace.Metrics { // Register YACE internal metrics
		err := y.Registry.Register(metric)
		if err != nil {
			return err
		}
	}

	y.Client, err = client.NewFactory(y.Logger, y.Config, false)
	if err != nil {
		return err
	}
	return nil
}

type YaceOptions struct {
	YaceCloudwatchConcurrencyPerApiLimitEnabled       string
	YaceCloudwatchConcurrencyListMetricsLimit         string
	YaceCloudwatchConcurrencyGetMetricDataLimit       string
	YaceCloudwatchConcurrencyGetMetricStatisticsLimit string
	YaceMetricsPerQuery                               string
	YaceTaggingAPIConcurrency                         string
	YaceCloudwatchConcurrency                         string
}

func (c *YaceOptions) Init() error {
	c.YaceCloudwatchConcurrencyPerApiLimitEnabled = os.Getenv("YACE_CLOUDWATCH_CONCURRENCY_PER_API_LIMIT_ENABLED")
	c.YaceCloudwatchConcurrencyListMetricsLimit = os.Getenv("YACE_CLOUDWATCH_CONCURRENCY_LIST_METRICS_LIMIT")
	c.YaceCloudwatchConcurrencyGetMetricDataLimit = os.Getenv("YACE_CLOUDWATCH_CONCURRENCY_GET_METRIC_DATA_LIMIT")
	c.YaceCloudwatchConcurrencyGetMetricStatisticsLimit = os.Getenv("YACE_CLOUDWATCH_CONCURRENCY_GET_METRIC_STATISTICS_LIMIT")
	c.YaceMetricsPerQuery = os.Getenv("YACE_METRICS_PER_QUERY")
	c.YaceTaggingAPIConcurrency = os.Getenv("YACE_TAG_CONCURRENCY")
	c.YaceCloudwatchConcurrency = os.Getenv("YACE_CLOUDWATCH_CONCURRENCY")
	return nil
}

func (c *YaceOptions) GetYaceOptions(logger types.Logger) ([]yace.OptionsFunc, error) {
	optFuncs := []yace.OptionsFunc{}
	var cloudwatchPerApiConcurrencyLimit bool
	var err error
	if c.YaceCloudwatchConcurrencyPerApiLimitEnabled != "" {
		logger.Log("debug", "Using non-default per API concurrency limit", slog.String("per_api_limit", c.YaceCloudwatchConcurrencyPerApiLimitEnabled))
		cloudwatchPerApiConcurrencyLimit, err = strconv.ParseBool(c.YaceCloudwatchConcurrencyPerApiLimitEnabled)
		if err != nil {
			return nil, err
		}
	}
	if c.YaceMetricsPerQuery != "" {
		logger.Log("debug", "Using non-default metrics per query", slog.String("metrics_per_query", c.YaceMetricsPerQuery))
		val, err := strconv.Atoi(c.YaceMetricsPerQuery)
		if err != nil {
			return nil, err
		}
		optFuncs = append(optFuncs, yace.MetricsPerQuery(val))
	}
	if c.YaceTaggingAPIConcurrency != "" {
		logger.Log("debug", "Using non-default tagging API concurrency", slog.String("tagging_api_concurrency", c.YaceTaggingAPIConcurrency))
		val, err := strconv.Atoi(c.YaceTaggingAPIConcurrency)
		if err != nil {
			return nil, err
		}
		optFuncs = append(optFuncs, yace.TaggingAPIConcurrency(val))
	}

	if !cloudwatchPerApiConcurrencyLimit {
		if c.YaceCloudwatchConcurrency != "" {
			logger.Log("debug", "Using non-default cloudwatch concurrency", slog.String("cloudwatch_concurrency", c.YaceCloudwatchConcurrency))
			val, err := strconv.Atoi(c.YaceCloudwatchConcurrency)
			if err != nil {
				return nil, err
			}
			optFuncs = append(optFuncs, yace.CloudWatchAPIConcurrency(val))
		}
	} else {
		limits := yace.DefaultCloudwatchConcurrency
		if c.YaceCloudwatchConcurrencyListMetricsLimit != "" {
			logger.Log("debug", "Using non-default cloudwatch list metrics concurrency", slog.String("cloudwatch_list_metrics_concurrency", c.YaceCloudwatchConcurrencyListMetricsLimit))
			val, err := strconv.Atoi(c.YaceCloudwatchConcurrencyListMetricsLimit)
			if err != nil {
				return nil, err
			}
			limits.ListMetrics = val
		}
		if c.YaceCloudwatchConcurrencyGetMetricDataLimit != "" {
			logger.Log("debug", "Using non-default cloudwatch get metric data concurrency", slog.String("cloudwatch_get_metric_data_concurrency", c.YaceCloudwatchConcurrencyGetMetricDataLimit))
			val, err := strconv.Atoi(c.YaceCloudwatchConcurrencyGetMetricDataLimit)
			if err != nil {
				return nil, err
			}
			limits.GetMetricData = val
		}
		if c.YaceCloudwatchConcurrencyGetMetricStatisticsLimit != "" {
			logger.Log("debug", "Using non-default cloudwatch get metric statistics concurrency", slog.String("cloudwatch_get_metric_statistics_concurrency", c.YaceCloudwatchConcurrencyGetMetricStatisticsLimit))
			val, err := strconv.Atoi(c.YaceCloudwatchConcurrencyGetMetricStatisticsLimit)
			if err != nil {
				return nil, err
			}
			limits.GetMetricStatistics = val
		}
		optFuncs = append(optFuncs, yace.CloudWatchPerAPILimitConcurrency(limits.ListMetrics, limits.GetMetricData, limits.GetMetricStatistics))
	}
	return optFuncs, nil
}
