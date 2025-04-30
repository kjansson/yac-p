// Package converter provides methods to convert metrics to timeseries format for Prometheus remote write. It implements the types.MetricConverter interface.
package converter

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/kjansson/yac-p/v3/pkg/types"
	"github.com/prometheus/prometheus/prompb"

	io_prometheus_client "github.com/prometheus/client_model/go"
)

type Converter struct {
	Logger types.Logger // Logger instance
}

func NewConverter(logger types.Logger) *Converter {
	return &Converter{
		Logger: logger,
	}
}

// getValue extracts the value of the metric based on the metric type
func getValue(valueType io_prometheus_client.MetricType, metric *io_prometheus_client.Metric) (float64, error) {
	switch valueType {
	case io_prometheus_client.MetricType_GAUGE:
		return *metric.Gauge.Value, nil
	case io_prometheus_client.MetricType_COUNTER:
		return *metric.Counter.Value, nil
	default:
		return 0, fmt.Errorf("unknown metric type: %s", valueType)
	}
}

// ConvertMetrics accepts Prometheus metrics gathered from a Prometheus registry, converts and returns them in timeseries format suitable for the Prometheus remote write API
func (c *Converter) ConvertMetrics(metrics []*io_prometheus_client.MetricFamily, logger types.Logger) ([]prompb.TimeSeries, error) {

	newTimestamp := time.Now().UnixNano() / int64(time.Millisecond)
	timeSeries := []prompb.TimeSeries{} // Create a slice of prometheus time series
	var oldTimestamp int64
	timestamped := false
	// Process metrics into timeseries format that remote write expects
	for _, family := range metrics { // Range through metric types
		metricName, metricType := family.GetName(), family.GetType() // Extraxt the metric type and name to use in prometheus time series
		logger.Log("debug", "Processing metric", slog.String("metric_name", metricName), slog.String("metric_type", metricType.String()))
		for _, metric := range family.GetMetric() { // Range through the metrics of the metric type
			ts := prompb.TimeSeries{}

			// This one is special, we need to add the metric name in the special label that prometheus expects
			ts.Labels = append(ts.Labels, prompb.Label{Name: "__name__", Value: metricName})
			for _, label := range metric.GetLabel() {
				ts.Labels = append(ts.Labels, prompb.Label{Name: label.GetName(), Value: label.GetValue()}) // Create prometheus time series labels
			}

			value, err := getValue(metricType, metric) // Extract the value of the metric based on the metric type
			if err != nil {
				return nil, err
			}

			timestamp := metric.GetTimestampMs() // Extract the timestamp of the metric
			// Metrics can have timestamps from Cloudwatch if YACE is configured to use them.
			// If the metric does not have a timestamp, it's either a helper metric created by YACE or YACE is configured to ignore Cloudwatch timestamps.
			// We store the timestamp of the first metric if it's non-zero and use it for the helper metrics,
			// if the first metric is not timestamped we assume that YACE is configured to ignore Cloduwatch timestamps and generate our own.
			if timestamp == 0 {
				if timestamped {
					logger.Log("debug", "Using stored timestamp from previous metric", slog.Int64("timestamp", oldTimestamp))
					timestamp = oldTimestamp
				} else {
					logger.Log("debug", "Using generated timestamp", slog.Int64("timestamp", newTimestamp))
					timestamp = newTimestamp
				}
			} else {
				oldTimestamp = timestamp
				timestamped = true
			}

			ts.Samples = append(ts.Samples, prompb.Sample{Value: value, Timestamp: timestamp}) // Create prometheus time series samples
			timeSeries = append(timeSeries, ts)
		}
	}
	return timeSeries, nil
}
