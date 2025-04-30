package prom

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kjansson/yac-p/v2/pkg/logger"
	"github.com/kjansson/yac-p/v2/pkg/types"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
	"google.golang.org/protobuf/proto"
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

func createTestMetricsFamily() []*io_prometheus_client.MetricFamily {
	metric := []*io_prometheus_client.Metric{{
		Label: []*io_prometheus_client.LabelPair{
			//	{Name: proto.String("__name__"), Value: proto.String("test_gauge")},
			{Name: proto.String("label1"), Value: proto.String("value1")},
		},
		Gauge: &io_prometheus_client.Gauge{
			Value: proto.Float64(1.0),
		},
	}}

	return []*io_prometheus_client.MetricFamily{{
		Name:   proto.String("test_gauge"),
		Help:   proto.String("This is a test gauge"),
		Type:   io_prometheus_client.MetricType_GAUGE.Enum(),
		Metric: metric,
	}}
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

func TestMetricsProcessing(t *testing.T) {
	logger := &logger.SlogLogger{}
	err := logger.Init(types.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	timeseries, err := ProcessMetrics(createTestMetricsFamily(), logger)
	if err != nil {
		t.Fatalf("Failed to process metrics: %v", err)
	}

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

	logger := &logger.SlogLogger{}
	err := logger.Init(types.Config{})
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

	logger := &logger.SlogLogger{}
	err := logger.Init(types.Config{})
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

	logger := &logger.SlogLogger{}
	err := logger.Init(types.Config{})
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
