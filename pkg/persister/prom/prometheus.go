// Package prom provides a client for processing and persisting metrics to a Prometheus remote write endpoint.
// It implements the types.MetricPersister interface.
package prom

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/golang/snappy"
	"github.com/kjansson/yac-p/v2/pkg/types"
	"github.com/prometheus/prometheus/prompb"
)

type PromClient struct {
	RemoteWriteURL   string // The URL of the Prometheus remote write endpoint
	AuthType         string // The type of authentication to use (AWS, BASIC, TOKEN)
	AuthToken        string // The token to use for authentication (if using TOKEN auth)
	Username         string // The username to use for authentication (if using BASIC auth)
	Password         string // The password to use for authentication (if using BASIC auth)
	Region           string // The AWS region to use for authentication (if using AWS auth)
	PrometheusRegion string // The AWS region of the Prometheus remote write endpoint (if using Amazon Managed Prometheus)
	AWSRoleARN       string // The ARN of the AWS role to assume for remote write (if using Amazon Managed Prometheus cross-account)
}

func NewPersister(
	remoteWriteURL string,
	authType string,
	authToken string,
	username string,
	password string,
	region string,
	prometheusRegion string,
	awsRoleARN string,
) *PromClient {
	return &PromClient{
		RemoteWriteURL:   remoteWriteURL,
		AuthType:         authType,
		AuthToken:        authToken,
		Username:         username,
		Password:         password,
		Region:           region,
		PrometheusRegion: prometheusRegion,
		AWSRoleARN:       awsRoleARN,
	}
}

// PeristMetrics creates a Prometheus remote write request and sends it to the remote write URL
func (p *PromClient) PersistMetrics(timeSeries []prompb.TimeSeries, logger types.Logger) error {

	logger.Log("debug", "Sending timeseries", slog.Int("timeseries_count", len(timeSeries)))
	logger.Log("debug", "Auth type", slog.String("auth_type", p.AuthType))

	r := &prompb.WriteRequest{
		Timeseries: timeSeries,
	}
	tsProto, err := r.Marshal()
	if err != nil {
		return err
	}

	encoded := snappy.Encode(nil, tsProto)
	body := bytes.NewReader(encoded)

	req, err := http.NewRequest("POST", p.RemoteWriteURL, body)
	if err != nil {
		return err
	}

	switch p.AuthType {
	case "AWS":
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(p.Region),
		})
		if err != nil {
			return err
		}

		var awsCredentials *credentials.Credentials
		if p.AWSRoleARN != "" {
			logger.Log("debug", "Using AWS role", slog.String("role_arn", p.AWSRoleARN))
			awsCredentials = stscreds.NewCredentials(sess, p.AWSRoleARN, func(p *stscreds.AssumeRoleProvider) {
				host, err := os.Hostname()
				if err != nil {
					host = "unknown"
				}
				p.RoleSessionName = "aws-sigv4-proxy-" + host
				logger.Log("debug", "Using AWS role session name", slog.String("role_session_name", p.RoleSessionName))
			})
		} else {
			awsCredentials = sess.Config.Credentials
		}

		_, err = v4.NewSigner(awsCredentials).Sign(req, body, "aps", p.PrometheusRegion, time.Now())
		if err != nil {
			return err
		}
	case "BASIC":
		logger.Log("debug", "Using basic auth")
		req.SetBasicAuth(p.Username, p.Password)
	case "TOKEN":
		logger.Log("debug", "Using token auth")
		req.Header.Set("Authorization", "Bearer "+p.AuthToken)
	}

	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	logger.Log("debug", "Sending request", slog.String("url", p.RemoteWriteURL), slog.Int("body_size", len(encoded)))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send metrics: %s", response.Status)
	}
	logger.Log("debug", "Response", slog.String("status", response.Status), slog.Int("status_code", response.StatusCode))

	return nil
}

// // getValue extracts the value of the metric based on the metric type
// func getValue(valueType io_prometheus_client.MetricType, metric *io_prometheus_client.Metric) (float64, error) {
// 	switch valueType {
// 	case io_prometheus_client.MetricType_GAUGE:
// 		return *metric.Gauge.Value, nil
// 	case io_prometheus_client.MetricType_COUNTER:
// 		return *metric.Counter.Value, nil
// 	default:
// 		return 0, fmt.Errorf("unknown metric type: %s", valueType)
// 	}
// }

// // ProcessMetrics accepts Prometheus metrics gathered from a Prometheus registry, converts and returns them in timeseries format suitable for the Prometheus remote write API
// func ProcessMetrics(metrics []*io_prometheus_client.MetricFamily, logger types.Logger) ([]prompb.TimeSeries, error) {

// 	newTimestamp := time.Now().UnixNano() / int64(time.Millisecond)
// 	timeSeries := []prompb.TimeSeries{} // Create a slice of prometheus time series
// 	var oldTimestamp int64
// 	timestamped := false
// 	// Process metrics into timeseries format that remote write expects
// 	for _, family := range metrics { // Range through metric types
// 		metricName, metricType := family.GetName(), family.GetType() // Extraxt the metric type and name to use in prometheus time series
// 		logger.Log("debug", "Processing metric", slog.String("metric_name", metricName), slog.String("metric_type", metricType.String()))
// 		for _, metric := range family.GetMetric() { // Range through the metrics of the metric type
// 			ts := prompb.TimeSeries{}

// 			// This one is special, we need to add the metric name in the special label that prometheus expects
// 			ts.Labels = append(ts.Labels, prompb.Label{Name: "__name__", Value: metricName})
// 			for _, label := range metric.GetLabel() {
// 				ts.Labels = append(ts.Labels, prompb.Label{Name: label.GetName(), Value: label.GetValue()}) // Create prometheus time series labels
// 			}

// 			value, err := getValue(metricType, metric) // Extract the value of the metric based on the metric type
// 			if err != nil {
// 				return nil, err
// 			}

// 			timestamp := metric.GetTimestampMs() // Extract the timestamp of the metric
// 			// Metrics can have timestamps from Cloudwatch if YACE is configured to use them.
// 			// If the metric does not have a timestamp, it's either a helper metric created by YACE or YACE is configured to ignore Cloudwatch timestamps.
// 			// We store the timestamp of the first metric if it's non-zero and use it for the helper metrics,
// 			// if the first metric is not timestamped we assume that YACE is configured to ignore Cloduwatch timestamps and generate our own.
// 			if timestamp == 0 {
// 				if timestamped {
// 					logger.Log("debug", "Using stored timestamp from previous metric", slog.Int64("timestamp", oldTimestamp))
// 					timestamp = oldTimestamp
// 				} else {
// 					logger.Log("debug", "Using generated timestamp", slog.Int64("timestamp", newTimestamp))
// 					timestamp = newTimestamp
// 				}
// 			} else {
// 				oldTimestamp = timestamp
// 				timestamped = true
// 			}

// 			ts.Samples = append(ts.Samples, prompb.Sample{Value: value, Timestamp: timestamp}) // Create prometheus time series samples
// 			timeSeries = append(timeSeries, ts)
// 		}
// 	}
// 	return timeSeries, nil
// }
