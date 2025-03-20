# YAC-p

YAC-p is based on YACE (Yet Another Cloudwatch Exporter)  
https://github.com/prometheus-community/yet-another-cloudwatch-exporter

YAC-P utilizes Go libraries from YACE to collect and convert Cloudwatch metrics into Prometheus format. It runs as a Lambda function and requires no infrastructure to host.   
Where YACE serves these metrics for scraping like an exporter, YAC-p converts them to remote write format and pushes them to a Prometheus remote write endpoint. 

## Deployment

- Build the image and push to ECR
- Write a YACE job config file
- Deploy with Terraform

## Purpose
TODO

