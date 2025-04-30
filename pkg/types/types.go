// Package types defines interfaces and constants used in yac-p
package types

import (
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

type Logger interface {
	Log(level string, msg string, args ...any)
}

// MetricCollector is an interface for collecting metrics and export them in Prometheus format
type MetricCollector interface {
	CollectMetrics(Logger) error
	ExportMetrics(Logger) ([]*io_prometheus_client.MetricFamily, error)
}

type MetricConverter interface {
	ConvertMetrics([]*io_prometheus_client.MetricFamily, Logger) ([]prompb.TimeSeries, error)
}

// MetricPersister is an interface for persisting Prometheus time series data
type MetricPersister interface {
	PersistMetrics([]prompb.TimeSeries, Logger) error
}
