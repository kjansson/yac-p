package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kjansson/yac-p/pkg/controller"
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

func TestMetricsProcessing(t *testing.T) {
	c := &controller.Controller{
		Logger:   &logger.SlogLogger{},
		Gatherer: &tests.YaceMockClient{},
	}

	err := c.Gatherer.Init(func() ([]byte, error) { return []byte(""), nil })
	if err != nil {
		t.Fatalf("Failed to initialize gatherer: %v", err)
	}
	err = c.Logger.Init()
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	testGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "test_gauge",
			Help: "This is a test gauge",
		},
	)
	testGauge.Set(1.0)
	c.Gatherer.GetRegistry().MustRegister(testGauge)

	metrics, err := c.ExtractMetrics()
	if err != nil {
		t.Fatalf("Failed to extract metrics: %v", err)
	}

	timeseries, err := prom.ProcessMetrics(metrics, c.Logger)
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

	logger := &logger.SlogLogger{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}

	promClient := &prom.PromClient{
		RemoteWriteURL: svr.URL,
	}

	c := &controller.Controller{
		Logger:    logger,
		Gatherer:  &tests.YaceMockClient{},
		Config:    &yace.YaceOptions{},
		Persister: promClient,
	}

	err := c.Gatherer.Init(func() ([]byte, error) { return []byte(""), nil })
	if err != nil {
		t.Fatalf("Failed to initialize gatherer: %v", err)
	}

	testGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "test_gauge",
			Help: "This is a test gauge",
		},
	)
	testGauge.Set(1.0)
	c.Gatherer.GetRegistry().MustRegister(testGauge)

	metrics, err := c.ExtractMetrics()
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

func TestMetricsPersistingBasicAuth(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		err := checkHeaders(r)
		if err != nil {
			t.Fatalf("Header check failed: %v", err)
		}
		username, password, _ := r.BasicAuth()
		if username != "testuser" {
			t.Fatalf("Expected username testuser, got %s", username)
		}
		if password != "testpassword" {
			t.Fatalf("Expected password testpassword, got %s", password)
		}
	}))

	logger := &logger.SlogLogger{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}

	promClient := &prom.PromClient{
		RemoteWriteURL: svr.URL,
		AuthType:       "BASIC",
		Username:       "testuser",
		Password:       "testpassword",
	}

	c := &controller.Controller{
		Logger:    logger,
		Gatherer:  &tests.YaceMockClient{},
		Config:    &yace.YaceOptions{},
		Persister: promClient,
	}

	err := c.Gatherer.Init(func() ([]byte, error) { return []byte(""), nil })
	if err != nil {
		t.Fatalf("Failed to initialize gatherer: %v", err)
	}

	testGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "test_gauge",
			Help: "This is a test gauge",
		},
	)
	testGauge.Set(1.0)
	c.Gatherer.GetRegistry().MustRegister(testGauge)

	metrics, err := c.ExtractMetrics()
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

func TestMetricsPersistingTokenAuth(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		err := checkHeaders(r)
		if err != nil {
			t.Fatalf("Header check failed: %v", err)
		}
		if r.Header.Get("Authorization") != "Bearer testtoken" {
			t.Fatalf("Expected Authorization Bearer testtoken, got %s", r.Header.Get("Authorization"))
		}
	}))

	logger := &logger.SlogLogger{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}

	promClient := &prom.PromClient{
		RemoteWriteURL: svr.URL,
		AuthType:       "TOKEN",
		AuthToken:      "testtoken",
	}

	c := &controller.Controller{
		Logger:    logger,
		Gatherer:  &tests.YaceMockClient{},
		Config:    &yace.YaceOptions{},
		Persister: promClient,
	}

	err := c.Gatherer.Init(func() ([]byte, error) { return []byte(""), nil })
	if err != nil {
		t.Fatalf("Failed to initialize gatherer: %v", err)
	}

	testGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "test_gauge",
			Help: "This is a test gauge",
		},
	)
	testGauge.Set(1.0)
	c.Gatherer.GetRegistry().MustRegister(testGauge)

	metrics, err := c.ExtractMetrics()
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
