package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/snappy"
	yace "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg"
	client "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/clients/v2"
	config "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/config"
	logging "github.com/prometheus-community/yet-another-cloudwatch-exporter/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

const (
	configFilePath = "/tmp/config.yaml"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest() {

	ctx := context.Background()
	logger := logging.NewLogger("", false)

	if os.Getenv("PROMETHEUS_REMOTE_WRITE_URL") == "" {
		panic("PROMETHEUS_REMOTE_WRITE_URL is required")
	}

	configS3Path, configS3Bucket := os.Getenv("CONFIG_S3_PATH"), os.Getenv("CONFIG_S3_BUCKET")

	if configS3Bucket != "" && configS3Path != "" {
		sess, err := createAWSSession()
		if err != nil {
			panic(err)
		}
		s3svc := s3.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
		obj, err := s3svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(configS3Bucket),
			Key:    aws.String(configS3Path),
		})
		if err != nil {
			panic(err)
		}
		content, err := io.ReadAll(obj.Body)
		if err != nil {
			panic(err)
		}
		// Write the configuration to a ephemeral storage, this is needed since the config package expects a file path
		err = os.WriteFile(configFilePath, content, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		panic("CONFIG_S3_BUCKET and CONFIG_S3_PATH is required.")
	}

	config := config.ScrapeConf{} // Create a new scrape config
	conf, err := config.Load(configFilePath, logger)
	if err != nil {
		panic(err)
	}

	registry := prometheus.NewRegistry() // Create a new prometheus registry

	// Create a new yace client factory
	f, err := client.NewFactory(logger, conf, false)
	if err != nil {
		panic(err)
	}

	// Query metrics and resources and update the prometheus registry
	err = yace.UpdateMetrics(ctx, logger, conf, registry, f)
	if err != nil {
		panic(err)
	}

	// Create prometheus timestamp
	newTimestamp := time.Now().UnixNano() / int64(time.Millisecond)

	metrics, err := registry.Gather() // Gather the metrics from the prometheus registry
	if err != nil {
		panic(err)
	}

	timeSeries := []prompb.TimeSeries{} // Create a slice of prometheus time series
	var oldTimestamp int64
	timestamped := false
	// Process metrics into timeseries format that remote write expects
	for _, family := range metrics { // Range through metric types
		metricName, metricType := family.GetName(), family.GetType() // Extraxt the metric type and name to use in prometheus time series

		for _, metric := range family.GetMetric() { // Range through the metrics of the metric type
			ts := prompb.TimeSeries{}

			for _, label := range metric.GetLabel() {
				ts.Labels = append(ts.Labels, prompb.Label{Name: label.GetName(), Value: label.GetValue()}) // Create prometheus time series labels
			}
			// This one is special, we need to add the metric name in the special label that prometheus expects
			ts.Labels = append(ts.Labels, prompb.Label{Name: "__name__", Value: metricName})

			value, err := getValue(metricType, metric) // Extract the value of the metric based on the metric type
			if err != nil {
				panic(err)
			}

			timestamp := metric.GetTimestampMs() // Extract the timestamp of the metric
			// Metrics can have timestamps from Cloudwatch if YACE is configured to use them.
			// If the metric does not have a timestamp, it's either a helper metric created by YACE or YACE is configured to ignore Cloudwatch timestamps.
			// We store the timestamp of the first metric if it's non-zero and use it for the helper metrics,
			// if the first metric is not timestamped we assume that YACE is configured to ignore Cloduwatch timestamps and generate our own.
			if timestamp == 0 {
				if timestamped {
					timestamp = oldTimestamp
				} else {
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
	err = sendRequest(timeSeries) // Send the timeseries to the remote write endpoint
	if err != nil {
		panic(err)
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

// sendRequest sends the timeseries to the remote write endpoint
func sendRequest(ts []prompb.TimeSeries) error {

	authType := os.Getenv("AUTH_TYPE")

	r := &prompb.WriteRequest{
		Timeseries: ts,
	}
	tsProto, err := r.Marshal()
	if err != nil {
		return err
	}

	encoded := snappy.Encode(nil, tsProto)
	body := bytes.NewReader(encoded)

	req, err := http.NewRequest("POST", os.Getenv("PROMETHEUS_REMOTE_WRITE_URL"), body)
	if err != nil {
		return err
	}

	switch authType {
	case "AWS":
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(os.Getenv("AWS_REGION")),
		})
		if err != nil {
			return err
		}

		roleArn := os.Getenv("AWS_ROLE_ARN")
		var awsCredentials *credentials.Credentials
		if roleArn != "" {
			awsCredentials = stscreds.NewCredentials(sess, roleArn, func(p *stscreds.AssumeRoleProvider) {
				host, err := os.Hostname()
				if err != nil {
					host = "unknown"
				}
				p.RoleSessionName = "aws-sigv4-proxy-" + host
			})
		} else {
			awsCredentials = sess.Config.Credentials
		}

		_, err = v4.NewSigner(awsCredentials).Sign(req, body, "aps", os.Getenv("PROMETHEUS_REGION"), time.Now())
		if err != nil {
			return err
		}
	case "BASIC":
		req.SetBasicAuth(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	case "TOKEN":
		req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN"))
	}

	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	response, err := http.DefaultClient.Do(req)
	if err != nil && response.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

// createAWSSession creates a new AWS session
func createAWSSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return sess, err
	}
	return sess, err
}
