package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kjansson/yac-p/pkg/logger"
	"github.com/kjansson/yac-p/pkg/prom"
	"github.com/kjansson/yac-p/pkg/tests"
	"github.com/kjansson/yac-p/pkg/yace"
	"github.com/prometheus/client_golang/prometheus"
)

func checkHeaders(r *http.Request) error {
	if r.Header.Get("Content-Encoding") != "snappy" {
		return fmt.Errorf("Expected snappy encoding, got %s", r.Header.Get("Content-Encoding"))
	}
	if r.Header.Get("X-Prometheus-Remote-Write-Version") != "0.1.0" {
		return fmt.Errorf("Expected X-Prometheus-Remote-Write-Version 0.1.0, got %s", r.Header.Get("X-Prometheus-Remote-Write-Version"))
	}
	if r.Header.Get("Content-Type") != "application/x-protobuf" {
		return fmt.Errorf("Expected application/x-protobuf, got %s", r.Header.Get("Content-Type"))
	}
	return nil
}

func TestConfigLoad(t *testing.T) {
	c := &Controller{
		Logger:    &logger.SlogLogger{},
		Collector: &tests.YaceMockClient{},
		Config:    &yace.YaceOptions{},
		Persister: &prom.PromClient{},
	}

	os.Setenv("PROMETHEUS_REMOTE_WRITE_URL", "http://localhost:9090/api/v1/write")
	defer os.Unsetenv("PROMETHEUS_REMOTE_WRITE_URL")

	err := c.Init(tests.GetTestConfigLoader()) // Initialize all components
	if err != nil {
		t.Fatalf("Failed to initialize with vaild config: %v", err)
	}
}

func TestMetricsPersistingNoAuth(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		err := checkHeaders(r)
		if err != nil {
			t.Fatalf("Header check failed: %v", err)
		}
		if r.Header.Get("Authorization") != "" {
			t.Fatalf("Expected no Authorization header, got %s", r.Header.Get("Authorization"))
		}
	}))

	os.Setenv("PROMETHEUS_REMOTE_WRITE_URL", svr.URL)
	defer os.Unsetenv("PROMETHEUS_REMOTE_WRITE_URL")
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	c := &Controller{
		Logger:    &logger.SlogLogger{},
		Collector: &tests.YaceMockClient{},
		Config:    &yace.YaceOptions{},
		Persister: &prom.PromClient{},
	}

	err := c.Init(func() ([]byte, error) { return []byte(""), nil })
	if err != nil {
		t.Fatalf("Failed to initialize Collector: %v", err)
	}

	testGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "test_gauge",
			Help: "This is a test gauge",
		},
	)
	testGauge.Set(1.0)
	c.Collector.GetRegistry().MustRegister(testGauge)

	metrics, err := c.ExportMetrics()
	if err != nil {
		t.Fatalf("Failed to extract metrics: %v", err)
	}

	timeseries, err := prom.ProcessMetrics(metrics, c.Logger)
	if err != nil {
		t.Fatalf("Failed to process metrics: %v", err)
	}

	err = c.PersistMetrics(timeseries)
	if err != nil {
		t.Fatalf("Failed to persist metrics: %v", err)
	}

	defer svr.Close()
}

// No further persist tests needed here, this is only to test the extended method
