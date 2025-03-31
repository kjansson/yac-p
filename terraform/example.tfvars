name_prefix             = "demo"
create_amp_workspace    = true // If you want to use Amazon Managed Prometheus
create_lambda_log_group = true // If you want to create a CloudWatch log group for the Lambda function

config_file_local_path    = "/path/config.yaml"
create_config_file_bucket = true // If you want to create an S3 bucket for the config file
tags = {
  "Environment" = "demo"
}
