package controller

// Package controller provides the Controller struct and its methods to manage the components of the application

import (
	"github.com/kjansson/yac-p/pkg/types"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

type Controller struct {
	Logger    types.Logger
	Config    types.Config
	Collector types.MetricCollector
	Persister types.MetricPersister
}

// Init initializes all components of the controller
func (c *Controller) Init(configFileLoader func() ([]byte, error)) error {
	if err := c.Logger.Init(); err != nil {
		return err
	}
	if err := c.Config.Init(); err != nil {
		return err
	}
	if err := c.Collector.Init(configFileLoader); err != nil { // Initialize the metric Collector with the config file loader
		return err
	}
	if err := c.Persister.Init(); err != nil {
		return err
	}
	return nil
}

// Log extends the logger interface
func (c *Controller) Log(level string, msg string, args ...any) {
	c.Logger.Log(level, msg, args...)
}

// GetRegistry returns the prometheus registry from the Collector component
func (c *Controller) CollectMetrics() error {
	return c.Collector.CollectMetrics(c.Logger, c.Config)
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
