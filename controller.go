package main

import (
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

type Controller struct {
	Logger    Logger
	Config    Config
	Gatherer  MetricGatherer
	Persister MetricPersister
}

func (c *Controller) Init() error {
	if err := c.Logger.Init(); err != nil {
		return err
	}
	if err := c.Config.Init(); err != nil {
		return err
	}
	if err := c.Gatherer.Init(); err != nil {
		return err
	}
	if err := c.Persister.Init(); err != nil {
		return err
	}
	return nil
}

func (c *Controller) Log(level string, msg string, args ...any) {
	c.Logger.Log(level, msg, args...)
}

func (c *Controller) CollectMetrics() error {
	return c.Gatherer.CollectMetrics(c.Logger, c.Config)
}

func (c *Controller) ExtractMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	metrics, err := c.Gatherer.ExtractMetrics(c.Logger)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (c *Controller) PersistMetrics(timeSeries []prompb.TimeSeries) error {
	return c.Persister.PersistMetrics(timeSeries, c.Logger)
}
