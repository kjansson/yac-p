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
	metrics, err := c.Gatherer.GetMetrics(c.Logger, c.Config)
	if err != nil {
		panic(err)
	}

	// Process the metrics into timeseries format
	timeSeries := processMetrics(metrics, c.Logger)

	// Persist the metrics to the remote write endpoint
	err = c.Persister.PersistMetrics(timeSeries, c.Logger) // Send the timeseries to the remote write endpoint
	if err != nil {
		panic(err)
	}
}
