package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.

type YaceMockClient struct {
	registry *prometheus.Registry
	//logger   *slog.Logger
}

func (y *YaceMockClient) Init() error {

	//y.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	y.registry = prometheus.NewRegistry() // Create a new prometheus registry

	testGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "test_gauge",
			Help: "This is a test gauge",
		},
	)
	testGauge.Set(1.0)
	y.registry.Register(testGauge)

	return nil
}

func (y *YaceMockClient) CollectMetrics(logger logger, config Config) error {
	return nil
}

func (y *YaceMockClient) ExtractMetrics(logger logger) ([]*io_prometheus_client.MetricFamily, error) {
	metrics, err := y.registry.Gather() // Gather the metrics from the prometheus registry
	if err != nil {
		panic(err)
	}
	return metrics, nil
}

func TestMetricsProcessing(t *testing.T) {
	c := &Controller{
		Logger:   &SlogLogger{},
		Gatherer: &YaceMockClient{},
	}

	c.Gatherer.Init()
	c.Logger.Init()

	metrics, err := c.Gatherer.ExtractMetrics(c.Logger)
	if err != nil {
		t.Fatalf("Failed to extract metrics: %v", err)
	}

	timeseries, err := processMetrics(metrics, c.Logger)
	if err != nil {
		t.Fatalf("Failed to process metrics: %v", err)
	}

	// Check if the timeseries are in the correct format

	for _, ts := range timeseries {
		if len(ts.Labels) == 0 {
			t.Fatalf("Timeseries has no labels")
		}
		if ts.Labels[0].Name != "__name__" {
			t.Fatalf("Timeseries does not have __name__ label")
		}
		if ts.Labels[0].Value != "test_gauge" {

			t.Fatalf("Timeseries does not have correct __name__ label value")
		}
		if len(ts.Samples) == 0 {
			t.Fatalf("Timeseries has no samples")
		}
		if ts.Samples[0].Value != 1.0 {
			t.Fatalf("Timeseries does not have correct sample value")
		}
		if ts.Samples[0].Timestamp == 0 {
			t.Fatalf("Timeseries does not have correct timestamp")
		}
	}
}
