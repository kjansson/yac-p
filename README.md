# YAC-p

YAC-p (Yet Another Cloudwatch Pusher) is heavily <b>(heavily)</b> utilizing YACE Go packages (Yet Another Cloudwatch Exporter)  
https://github.com/prometheus-community/yet-another-cloudwatch-exporter

<u>Note: This is a work in progress.</u>

YAC-p utilizes Go libraries from YACE to collect and convert Cloudwatch metrics into Prometheus format remote write format and writes to your endpoint of choice.  
It runs as a Lambda function and only requires AWS managed services to run (You don't have to host anything!).   

## Purpose

YAC-p fits in wherever you don't want to do metrics scraping to get access to your Cloudwatch metrics in Prometheus.   
There are multiple scenarios where push-based metrics collection might be more suitable than pull-based;

- <b>Decentralizing</b> - When in a multi-account cloud environment with centralized metric collection, keeping the configuration and responsibility of metrics collection in the scope of the client account simplifies scaling.

- <b>Network access</b> - Scraping through Firewalls or other network access control mechanism can be a hassle. Providing a single endpoint for push-based metrics delivery simplifies things.

- <b>Closer to real-time</b> - Instead of relying on the timing of an exporter and a scraping job, YAC-p delivers as fresh metrics as you want and Cloudwatch can manage to collect.  

## Features

- <b>Nothing to host</b> - Can run on fully managed AWS infrastructure, Eventbridge + Lambda + Amazon managed Prometheus (Works with your self-hosted Prometheus as well)
- <b>Yace compatible</b> - Uses YACE native job configurations and its amazing discovery and metric enrichment features
- <b>Manage Prometheus the way you want (or not)</b> - Authentication options for Amazon Managed Prometheus, self-hosted Prometheus, etc

## Deployment

YAC-p can be deployed using fully managed services. Using Eventbridge to schedule the YAC-p Lambda function it will deliver AWS Cloudwatch metrics to any Prometheus server, but when using Amazon Managed Prometheus it becomes a fully managed collection process.  
The included Terraform example code deploys everything you need including a Amazon Managed Prometheus workspace.

![Deployment](img/deployment.png)

## Try it out
- Build the image and push to ECR
- Write a YACE job config file (https://github.com/prometheus-community/yet-another-cloudwatch-exporter/blob/master/docs/configuration.md) or use the included example
- Deploy with included Terraform code

## Lambda configuration
The included Terraform code will configure the Lambda for you, but if you want to deploy it yourself there are a few environment variables to set;

- PROMETHEUS_REMOTE_WRITE_URL - The URL of the Prometheus remote write endpoint
- PROMETHEUS_REGION - If using AMP, the region needs to be configured
- CONFIG_S3_BUCKET - The S3 bucket where the config file is stored
- CONFIG_S3_PATH - The path of the config file
- AUTH_TYPE - Authentication type to use for the remote write endpoint. Valid options are "AWS", "BASIC", "TOKEN. Leave empty if no authentication is required.

## Deployment comparison to YACE

YACE is awesome, but it usually requires something to host it. That might be something that you are trying to avoid. It can also introduce latency into the metric collection process, since it requires two unsynchronized collections of the metrics before it ends up in Prometheus.  

![YACE](img/YACE.png)

YAC-p fits better where you want to manage as little infrastructure as possible. It can also reduce latency in the metrics collection process, since it queries Cloudwatch at you desired rate and delivers metrics instantly to Prometheus.

![YAC-p](img/YAC-p.png)
