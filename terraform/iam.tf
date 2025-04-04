#################
# Eventbridge role
#################

data "aws_iam_policy_document" "eventbridge_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["events.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "eventbridge" {
  assume_role_policy = data.aws_iam_policy_document.eventbridge_assume_role.json
}

data "aws_iam_policy_document" "eventbridge_invoke_lambda" {
  statement {
    effect = "Allow"
    actions = [
      "lambda:InvokeFunction",
      "lambda:GetFunctionConfiguration",
    ]
    resources = [
      aws_lambda_function.this.arn,
    ]
  }
}

resource "aws_iam_role_policy" "eventbridge_invoke_lambda" {
  name   = format("%s-lambda-invoke", var.name_prefix)
  role   = aws_iam_role.eventbridge.id
  policy = data.aws_iam_policy_document.eventbridge_invoke_lambda.json
}

#################
# Lambda role
#################

data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type = "Service"
      identifiers = [
        "lambda.amazonaws.com",
      ]
    }
    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "lambda" {
  name               = format("%s-lambda", var.name_prefix)
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

data "aws_iam_policy_document" "lambda_exec_policy" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      "${aws_cloudwatch_log_group.lambda.arn}:*"
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
    ]
    resources = [
      "${aws_s3_bucket.this[0].arn}/*"
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "tag:GetResources",
      "cloudwatch:GetMetricData",
      "cloudwatch:GetMetricStatistics",
      "cloudwatch:ListMetrics",
      "apigateway:GET",
      "aps:ListWorkspaces",
      "autoscaling:DescribeAutoScalingGroups",
      "dms:DescribeReplicationInstances",
      "dms:DescribeReplicationTasks",
      "ec2:DescribeTransitGatewayAttachments",
      "ec2:DescribeSpotFleetRequests",
      "shield:ListProtections",
      "storagegateway:ListGateways",
      "storagegateway:ListTagsForResource",
      "iam:ListAccountAliases"
    ]
    resources = ["*"] // Needed for discovery
  }

  dynamic "statement" {
    for_each = length(var.assumable_roles) > 0 || var.prometheus_remote_write_role_arn != "" ? [1] : []
    content {
      effect = "Allow"
      actions = [
        "sts:AssumeRole"
      ]
      # If prometheus_remote_write_role_arn is set, add it to the list of assumable roles if it's not already there
      resources = var.prometheus_remote_write_role_arn != "" && !contains(var.assumable_roles, var.prometheus_remote_write_role_arn) ? concat(var.assumable_roles, [var.prometheus_remote_write_role_arn]) : var.assumable_roles
    }
  }
}

resource "aws_iam_policy" "lambda_exec_policy" {
  name   = format("%s-lambda-exec", var.name_prefix)
  policy = data.aws_iam_policy_document.lambda_exec_policy.json
}

resource "aws_iam_role_policy_attachment" "lambda_exec_policy" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda_exec_policy.arn
}

resource "aws_iam_role_policy_attachment" "lambda_cloudwatch_read" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchReadOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "lambda_prom_write" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonPrometheusRemoteWriteAccess"
}