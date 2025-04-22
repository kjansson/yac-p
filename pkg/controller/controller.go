package controller

import (
	"github.com/kjansson/yac-p/pkg/types"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

type Controller struct {
	Logger    types.Logger
	Config    types.Config
	Gatherer  types.MetricGatherer
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
	if err := c.Gatherer.Init(configFileLoader); err != nil { // Initialize the metric gatherer with the config file loader
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

// GetRegistry returns the prometheus registry from the gatherer component
func (c *Controller) CollectMetrics() error {
	return c.Gatherer.CollectMetrics(c.Logger, c.Config)
}

// GetRegistry extracts the metrics from the prometheus registry in the gatherer component
func (c *Controller) ExtractMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	metrics, err := c.Gatherer.ExtractMetrics(c.Logger)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// PersistMetrics extends the underlying method and persists timeseries to the remote write endpoint
func (c *Controller) PersistMetrics(timeSeries []prompb.TimeSeries) error {
	return c.Persister.PersistMetrics(timeSeries, c.Logger)
}
