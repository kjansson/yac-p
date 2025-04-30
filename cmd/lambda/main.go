// Lambda implementation of yac-p
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	defcon "github.com/kjansson/defcon"
)

func main() {
	lambda.Start(HandleRequest) // Start the AWS Lambda function
}

func HandleRequest() {

	config := Config{}
	err := defcon.CheckConfigStruct(&config) // Validate the config struct
	if err != nil {
		panic(err)
	}
	config.ConfigFileLoader = GetS3Loader() // Set the config file loader to S3 for Lambda

	c, err := NewController(config) // Create a new controller instance
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
	timeSeries, err := c.ConvertMetrics(metrics)
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
