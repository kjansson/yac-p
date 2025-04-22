package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/kjansson/yac-p/pkg/types"
	yace "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg"
)

type YaceConfig struct {
	YaceCloudwatchConcurrencyPerApiLimitEnabled       string
	YaceCloudwatchConcurrencyListMetricsLimit         string
	YaceCloudwatchConcurrencyGetMetricDataLimit       string
	YaceCloudwatchConcurrencyGetMetricStatisticsLimit string
	YaceMetricsPerQuery                               string
	YaceTaggingAPIConcurrency                         string
	YaceCloudwatchConcurrency                         string
}

func (c *YaceConfig) Init() error {
	c.YaceCloudwatchConcurrencyPerApiLimitEnabled = os.Getenv("YACE_CLOUDWATCH_CONCURRENCY_PER_API_LIMIT_ENABLED")
	c.YaceCloudwatchConcurrencyListMetricsLimit = os.Getenv("YACE_CLOUDWATCH_CONCURRENCY_LIST_METRICS_LIMIT")
	c.YaceCloudwatchConcurrencyGetMetricDataLimit = os.Getenv("YACE_CLOUDWATCH_CONCURRENCY_GET_METRIC_DATA_LIMIT")
	c.YaceCloudwatchConcurrencyGetMetricStatisticsLimit = os.Getenv("YACE_CLOUDWATCH_CONCURRENCY_GET_METRIC_STATISTICS_LIMIT")
	c.YaceMetricsPerQuery = os.Getenv("YACE_METRICS_PER_QUERY")
	c.YaceTaggingAPIConcurrency = os.Getenv("YACE_TAG_CONCURRENCY")
	c.YaceCloudwatchConcurrency = os.Getenv("YACE_CLOUDWATCH_CONCURRENCY")
	return nil
}

func (c *YaceConfig) GetYaceOptions(logger types.Logger) ([]yace.OptionsFunc, error) {
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
