package main

import (
	yace "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

type logger interface {
	Init() error
	Log(level string, msg string, args ...any)
}

type metricGatherer interface {
	Init() error
	GetMetrics(logger, Config) ([]*io_prometheus_client.MetricFamily, error)
}

type metricPersister interface {
	Init() error
	PersistMetrics([]prompb.TimeSeries, logger) error
}

type Config interface {
	Init() error
	GetYaceOptions(logger logger) []yace.OptionsFunc
}

type Controller struct {
	Logger    logger
	Config    Config
	Gatherer  metricGatherer
	Persister metricPersister
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
