package prom

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kjansson/yac-p/v2/pkg/logger"
	"github.com/prometheus/prometheus/prompb"
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

func createTestTimeSeries() []prompb.TimeSeries {
	return []prompb.TimeSeries{
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "test_gauge"},
				{Name: "label1", Value: "value1"},
			},
			Samples: []prompb.Sample{
				{Value: 1.0, Timestamp: 1234567890},
			},
		},
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

	logger, err := logger.NewLogger(
		os.Stdout,
		"debug",
		"text",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	p := &PromClient{
		RemoteWriteURL: svr.URL,
	}

	err = p.PersistMetrics(createTestTimeSeries(), logger)
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

	logger, err := logger.NewLogger(
		os.Stdout,
		"debug",
		"text",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	p := &PromClient{
		RemoteWriteURL: svr.URL,
		AuthType:       "BASIC",
		Username:       "testuser",
		Password:       "testpassword",
	}

	err = p.PersistMetrics(createTestTimeSeries(), logger)
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

	logger, err := logger.NewLogger(
		os.Stdout,
		"debug",
		"text",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	p := &PromClient{
		RemoteWriteURL: svr.URL,
		AuthType:       "TOKEN",
		AuthToken:      "testtoken",
	}

	err = p.PersistMetrics(createTestTimeSeries(), logger)
	if err != nil {
		t.Fatalf("Failed to persist metrics: %v", err)
	}

	defer svr.Close()

}
