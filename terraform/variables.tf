variable "name_prefix" {
  description = "Prefix to apply to all resources."
  type        = string
}

variable "create_amp_workspace" {
  description = "Create an Amazon Managed Prometheus workspace."
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = any
  default     = {}
}

variable "create_lambda_log_group" {
  description = "Create a CloudWatch log group for the Lambda function."
  type        = bool
  default     = true
}

variable "lambda_log_group_retention" {
  description = "Retention period for the Lambda log group"
  type        = number
  default     = 7
}

variable "lambda_log_group_name" {
  description = "Name of the Lambda log group. Setting create_lambda_log_group to true will override this value and use the created log group name."
  type        = string
  default     = ""
}

variable "prometheus_endpoint" {
  description = "The endpoint of the Prometheus workspace. Setting create_amp_workspace to true will override this value and use the created workspace endpoint."
  type        = string
  default     = ""
}

variable "prometheus_region" {
  description = "The region of the Prometheus workspace. Only used for Amazon Managed Prometheus. Defaults to current region."
  type        = string
  default = ""
}

variable "lambda_log_level" {
  description = "The log level for the Lambda function."
  type        = string
  default     = "INFO"
}

variable "lambda_schedule_rate" {
  description = "The rate at which to invoke the Lambda function (in minutes)."
  type        = number
  default     = 5
}

variable "lambda_timeout_seconds" {
  description = "The timeout for the Lambda function in seconds."
  type        = number
  default     = 15
}

variable "config_file_local_path" {
  description = "Path to the local YACE config file to upload and use."
  type        = string
}

variable "create_config_file_bucket" {
  description = "Create an S3 bucket to store the YACE config file. Overrides config_bucket."
  type        = bool
  default     = false
}

variable "config_bucket" {
  description = "Name of existing S3 bucket to store the YACE config file. Ignored if create_config_file_bucket is true."
  type        = string
  default     = ""
}

variable "config_path" {
  description = "Custom S3 path for the YACE config file."
  default     = ""
}

variable "yace_options" {
  description = "Additional options to pass to YACE libraries."
  type        = map(string)
  default     = {}
}

variable "eventbridge_schedule_expression" {
  description = "The schedule expression for the EventBridge rule."
  type        = string
  default     = "rate(5 minutes)"
}

variable "assumable_roles" {
  description = "List of IAM role ARNs to add to IAM policy for Lambda to be able to assume. Used for cross-account access."
  type        = list(string)
  default     = []
}

variable "prometheus_remote_write_role_arn" {
  description = "The ARN of the IAM role to use to remote write to Prometheus (AMP only)."
  type        = string
  default     = ""
}
variable "prometheus_auth_type" {
  description = "The authentication type to use for Prometheus remote write. Use when not using Amazon Managed Prometheus or other AWS authentication."
  type        = string
  default     = ""
}
variable "lambda_runtime" {
  description = "The runtime for the Lambda function."
  type        = string
  default     = "provided.al2"
}