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

type Controller struct {
	Logger    Logger
	Collector MetricCollector
	Converter MetricConverter
	Persister MetricPersister
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
