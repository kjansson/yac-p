// Package types defines interfaces and constants used in yac-p
package types

import (
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

// Logger is an interface for logging
type Logger interface {
	Log(level string, msg string, args ...any)
}

// MetricCollector is an interface for collecting metrics and export them in Prometheus format
type MetricCollector interface {
	CollectMetrics(Logger) error
	ExportMetrics(Logger) ([]*io_prometheus_client.MetricFamily, error)
}

// MetricConverter is an interface for converting Prometheus metrics to timeseries format
type MetricConverter interface {
	ConvertMetrics([]*io_prometheus_client.MetricFamily, Logger) ([]prompb.TimeSeries, error)
}

// MetricPersister is an interface for persisting Prometheus time series data
type MetricPersister interface {
	PersistMetrics([]prompb.TimeSeries, Logger) error
}

type Controller struct {
	Logger    Logger          // Logger component
	Collector MetricCollector // Collector component
	Converter MetricConverter // Converter component
	Persister MetricPersister // Persister component
}

// Log extends the logger interface
func (c *Controller) Log(level string, msg string, args ...any) {
	c.Logger.Log(level, msg, args...)
}

// GetRegistry extends the underlying method triggers metrics collection in the Collector component
func (c *Controller) CollectMetrics() error {
	return c.Collector.CollectMetrics(c.Logger)
}

// GetRegistry extendes the underlying method and extracts the metrics from the prometheus registry in the Collector component
func (c *Controller) ExportMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	metrics, err := c.Collector.ExportMetrics(c.Logger)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// ConvertMetrics extends the underlying method and converts metrics to timeseries format using the Converter component
func (c *Controller) ConvertMetrics(metrics []*io_prometheus_client.MetricFamily) ([]prompb.TimeSeries, error) {
	timeSeries, err := c.Converter.ConvertMetrics(metrics, c.Logger)
	if err != nil {
		return nil, err
	}
	return timeSeries, nil
}

// PersistMetrics extends the underlying method and persists timeseries to the remote write endpoint using the Persister component
func (c *Controller) PersistMetrics(timeSeries []prompb.TimeSeries) error {
	return c.Persister.PersistMetrics(timeSeries, c.Logger)
}
