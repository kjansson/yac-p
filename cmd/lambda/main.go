package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kjansson/defcon"
	"github.com/kjansson/yac-p/v2/pkg/controller"
	"github.com/kjansson/yac-p/v2/pkg/prom"
	"github.com/kjansson/yac-p/v2/pkg/types"
)

func main() {
	lambda.Start(HandleRequest) // Start the AWS Lambda function
}

func HandleRequest() {

	config := &types.Config{
		ConfigFileLoader: GetS3Loader(), // Use the S3 loader for Lambda implementation
		RemoteWriteURL:   os.Getenv("PROMETHEUS_REMOTE_WRITE_URL"),
		AuthType:         os.Getenv("AUTH_TYPE"),
		AuthToken:        os.Getenv("TOKEN"),
		Username:         os.Getenv("USERNAME"),
		Password:         os.Getenv("PASSWORD"),
		Region:           os.Getenv("AWS_REGION"),
		PrometheusRegion: os.Getenv("PROMETHEUS_REGION"),
		AWSRoleARN:       os.Getenv("AWS_ROLE_ARN"),
	}

	err := defcon.CheckConfigStruct(config) // Check the config struct for required fields
	if err != nil {
		panic(err)
	}

	c, err := controller.NewController(*config) // Create a new controller instance
	if err != nil {
		panic(err)
	}

	c.Logger.Log("debug", "Starting yac-p lambda function") // Log the start of the function

	c.Logger.Log("debug", "Collecting metrics")
	// Gather cloudwatch metrics
	err = c.CollectMetrics()
	if err != nil {
		panic(err)
	}

	c.Logger.Log("debug", "Extracting metrics")
	// Extract the metrics from the prometheus registry
	metrics, err := c.ExportMetrics()
	if err != nil {
		panic(err)
	}

	c.Logger.Log("debug", "Processing metrics")
	// Process the metrics into timeseries format
	timeSeries, err := prom.ProcessMetrics(metrics, c.Logger)
	if err != nil {
		panic(err)
	}

	c.Logger.Log("debug", "Persisting metrics")
	// Persist the metrics to the remote write endpoint
	err = c.PersistMetrics(timeSeries) // Send the timeseries to the remote write endpoint
	if err != nil {
		panic(err)
	}
}
