package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/ssm"
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
	logger := logging.NewLogger("debug", false)

	if os.Getenv("PROMETHEUS_REMOTE_WRITE_URL") == "" {
		panic("PROMETHEUS_REMOTE_WRITE_URL is required")
	}

	configSSMParameter := os.Getenv("CONFIG_SSM_PARAMETER")

	if configSSMParameter != "" {
		sess, err := session.NewSessionWithOptions(session.Options{
			Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			panic(err)
		}
		ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
		param, err := ssmsvc.GetParameters(&ssm.GetParametersInput{
			Names: []*string{&configSSMParameter},
		})
		if err != nil {
			panic(err)
		}
		content := param.Parameters[0].Value
		err = os.WriteFile(configFilePath, []byte(*content), 0644)
		if err != nil {
			panic(err)
		}
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

	gMetrics, err := registry.Gather() // Gather the metrics from the prometheus registry
	if err != nil {
		panic(err)
	}

	timeSeries := []prompb.TimeSeries{} // Create a slice of prometheus time series
	var oldTs int64
	// Process metrics into timeseries format that remote write expects
	for _, fam := range gMetrics { // Range through metric types
		metricName, metricType := fam.GetName(), fam.GetType() // Extraxt the metric type and name to use in prometheus time series

		for _, metric := range fam.GetMetric() { // Range through the metrics of the metric type
			ts := prompb.TimeSeries{}

			labels := metric.GetLabel() // Extract the labels of the metric
			for _, label := range labels {
				lv := label.GetValue()
				ts.Labels = append(ts.Labels, prompb.Label{Name: label.GetName(), Value: lv}) // Create prometheus time series labels
			}
			// This one is special, we need to add the metric name in the special label that prometheus expects
			ts.Labels = append(ts.Labels, prompb.Label{Name: "__name__", Value: metricName})

			value, err := getValue(metricType.String(), metric) // Extract the value of the metric based on the metric type
			if err != nil {
				panic(err)
			}
			timeStamp := metric.GetTimestampMs() // Extract the timestamp of the metric
			if timeStamp == 0 {                  // The helper metrics does not have a timestamp, so we need to create one by storing the metric timestamps
				timeStamp = oldTs
			} else {
				oldTs = timeStamp
			}

			ts.Samples = append(ts.Samples, prompb.Sample{Value: value, Timestamp: timeStamp}) // Create prometheus time series samples
			timeSeries = append(timeSeries, ts)

		}
	}

	err = sendRequest(timeSeries) // Send the timeseries to the remote write endpoint
	if err != nil {
		panic(err)
	}
}

func getValue(valueType string, metric *io_prometheus_client.Metric) (float64, error) {
	switch valueType {
	case "GAUGE":
		return *metric.Gauge.Value, nil
	case "COUNTER":
		return *metric.Counter.Value, nil
	default:
		return 0, fmt.Errorf("unknown metric type: %s", valueType)
	}
}

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
