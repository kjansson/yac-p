package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleRequest) // Start the AWS Lambda function
}

func HandleRequest() {

	c := &Controller{
		Logger:    &SlogLogger{},
		Config:    &YaceConfig{},
		Gatherer:  &YaceClient{},
		Persister: &PromClient{},
	}

	err := c.Init() // Initialize all components
	if err != nil {
		panic(err)
	}

	// Gather cloudwatch metrics
	err = c.CollectMetrics()
	if err != nil {
		panic(err)
	}

	// Extract the metrics from the prometheus registry
	metrics, err := c.ExtractMetrics()
	if err != nil {
		panic(err)
	}

	// Process the metrics into timeseries format
	timeSeries, err := processMetrics(metrics, c.Logger)
	if err != nil {
		panic(err)
	}

	// Persist the metrics to the remote write endpoint
	err = c.PersistMetrics(timeSeries) // Send the timeseries to the remote write endpoint
	if err != nil {
		panic(err)
	}
}
