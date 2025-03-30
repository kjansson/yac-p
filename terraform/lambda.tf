resource "random_pet" "always" { # Dummy resource to trigger go build
  length  = 2
  keepers = {
    always_run = timestamp()
  }
}

resource "null_resource" "build" {
  triggers = {
    always_trigger = random_pet.always.id
  }
  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ../bootstrap ../"
  }
}

data "archive_file" "this" {
  depends_on = [null_resource.build]

  type        = "zip"
  source_file = "../bootstrap"
  output_path = "../yac-p.zip"
}

resource "aws_lambda_function" "this" {
  function_name = format("%s-lambda", var.name_prefix)
  role          = aws_iam_role.lambda.arn
  handler       = "bootstrap"
  filename =    "../yac-p.zip"
  runtime       = "provided.al2"
  architectures = ["arm64"]

  environment {
    variables = merge({
      PROMETHEUS_REMOTE_WRITE_URL = var.create_amp_workspace ? "${aws_prometheus_workspace.this.0.prometheus_endpoint}api/v1/remote_write" : var.prometheus_endpoint
      PROMETHEUS_REGION           = var.prometheus_region == "" ? data.aws_region.current.name : var.prometheus_region
      CONFIG_S3_PATH              = var.config_path == "" ? format("%s-yace-config/config.yaml", var.name_prefix) : var.config_path
      CONFIG_S3_BUCKET            = var.create_config_file_bucket ? aws_s3_bucket.this[0].bucket : var.config_bucket
      AUTH_TYPE                   = var.create_amp_workspace ? "aws4" : var.prometheus_auth_type
      AWS_ROLE_ARN                = var.prometheus_remote_write_role_arn
    }, var.yace_options)
  }
  logging_config {
    log_group        = aws_cloudwatch_log_group.lambda.id
    log_format       = "JSON"
    system_log_level = var.lambda_log_level
  }
  timeout = var.lambda_timeout_seconds
  tags    = var.tags
  depends_on = [ data.archive_file.this ]
}

resource "aws_cloudwatch_log_group" "lambda" {
  name              = format("%s-lambda", var.name_prefix)
  retention_in_days = var.lambda_log_group_retention
}

resource "aws_lambda_permission" "lambda" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "events.amazonaws.com"
  source_arn    = module.scheduler.eventbridge_rule_arns[format("%s-cron", var.name_prefix)]
}

