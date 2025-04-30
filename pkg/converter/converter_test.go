package converter

import (
	"os"
	"testing"

	"github.com/kjansson/yac-p/v3/pkg/logger"
	"google.golang.org/protobuf/proto"

	io_prometheus_client "github.com/prometheus/client_model/go"
)

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

func TestMetricsProcessing(t *testing.T) {
	logger, err := logger.NewLogger(
		os.Stdout,
		"text",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	c := NewConverter(logger)

	timeseries, err := c.ConvertMetrics(createTestMetricsFamily(), logger)
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
