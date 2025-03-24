name_prefix          = "demo"
create_amp_workspace = true // If you want to use Amazon Managed Prometheus

tags = {
  "Environment" = "demo"
}

create_lambda_log_group = true // If you want to create a CloudWatch log group for the Lambda function

lambda_image_repository_arn = "arn:aws:ecr:<region>:<account>:repository/<rep>"
lambda_image_uri            = "<image_uri>"

prometheus_region         = "<region>" # If using Amazon Managed Prometheus
config_file_local_path    = "/path/config.yaml"
create_config_file_bucket = true // If you want to create an S3 bucket for the config file

eventbridge_schedule_expression = "rate(5 minute)"