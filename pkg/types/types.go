package types

import (
	yace "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

type Logger interface {
	Init() error
	Log(level string, msg string, args ...any)
}

type MetricGatherer interface {
	Init() error
	CollectMetrics(Logger, Config) error
	ExtractMetrics(Logger) ([]*io_prometheus_client.MetricFamily, error)
	GetRegistry() *prometheus.Registry
}

type MetricPersister interface {
	Init() error
	PersistMetrics([]prompb.TimeSeries, Logger) error
}

type Config interface {
	Init() error
	GetYaceOptions(logger Logger) ([]yace.OptionsFunc, error)
}
