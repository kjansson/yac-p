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

	c.Init() // Initialize all components

	// Gather cloudwatch metrics
	c.Gatherer.CollectMetrics(c.Logger, c.Config)

	// Extract the metrics from the prometheus registry
	metrics, err := c.Gatherer.ExtractMetrics(c.Logger)
	if err != nil {
		panic(err)
	}

	// Process the metrics into timeseries format
	timeSeries, err := processMetrics(metrics, c.Logger)
	if err != nil {
		panic(err)
	}

	// Persist the metrics to the remote write endpoint
	err = c.Persister.PersistMetrics(timeSeries, c.Logger) // Send the timeseries to the remote write endpoint
	if err != nil {
		panic(err)
	}
}
