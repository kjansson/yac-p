package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kjansson/yac-p/pkg/controller"
	"github.com/kjansson/yac-p/pkg/loaders"
	"github.com/kjansson/yac-p/pkg/logger"
	"github.com/kjansson/yac-p/pkg/prom"
	"github.com/kjansson/yac-p/pkg/yace"
)

func main() {
	lambda.Start(HandleRequest) // Start the AWS Lambda function
}

func HandleRequest() {

	c := &controller.Controller{
		Logger:    &logger.SlogLogger{},
		Config:    &yace.YaceOptions{},
		Collector: &yace.YaceClient{},
		Persister: &prom.PromClient{},
	}

	err := c.Init(loaders.GetS3Loader()) // Initialize all components, use the S3 loader for Lambda implementation
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
