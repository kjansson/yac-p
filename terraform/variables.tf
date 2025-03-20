variable "name_prefix" {
  description = "Prefix to apply to all resources"
  type        = string
}

variable "create_amp_workspace" {
  description = "Create an AMP workspace"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = any
  default     = {}
}

variable "create_lambda_log_group" {
  description = "Create a CloudWatch log group for the Lambda function"
  type        = bool
  default     = true
}

variable "lambda_log_group_retention" {
  description = "Retention period for the Lambda log group"
  type        = number
  default     = 7
}

variable "firehose_log_group_name" {
  description = "Name of the Firehose log group"
  type        = string
  default     = ""
}

variable "lambda_log_group_name" {
  description = "Name of the Lambda log group"
  type        = string
  default     = ""
}

variable "prometheus_endpoint" {
  description = "The endpoint of the Prometheus workspace"
  type        = string
  default     = ""
}

variable "prometheus_region" {
  description = "The region of the Prometheus server"
  type        = string
}

variable "lambda_image_uri" {
  description = "The URI of the Lambda function image"
  type        = string
}

variable "lambda_image_repository_arn" {
  description = "The ARN of the Lambda function image repository"
  type        = string
}

variable "lambda_log_level" {
  description = "The log level for the Lambda function"
  type        = string
  default     = "INFO"
}

variable "lambda_schedule_rate" {
  description = "The rate at which to invoke the Lambda function"
  type        = number
  default     = 5
}

variable "yace_config_file_path" {
  description = "Path to the YACE config file"
  type        = string
}