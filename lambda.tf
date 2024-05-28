locals {
  strava_secret = "80580d66c7d4514ecf1f904f9a698fb314b2463a"

}

## Activity Lambda
data "archive_file" "activity_lambda" {
  type        = "zip"
  source_file = "cmd/activity/bootstrap"
  output_path = "cmd/activity/bootstrap.zip"
}

resource "aws_lambda_function" "activity_lambda" {
  filename      = "${path.module}/cmd/activity/bootstrap.zip"
  function_name = "cam-activity"
  role          = aws_iam_role.iam_for_lambda.arn
  handler = "hello.handler"

  source_code_hash = data.archive_file.activity_lambda.output_base64sha256

  runtime = "provided.al2023"

  environment {
    variables = {
      STRAVA_CLIENT_SECRET = local.strava_secret
      USERS_TABLE_NAME = aws_dynamodb_table.cam-users.name
      WEBHOOK_VERIFY_TOKEN = "KPHo87W@PLVCZs"
    }
  }
}

resource "aws_lambda_function_url" "activity" {
  function_name      = aws_lambda_function.activity_lambda.function_name
  authorization_type = "NONE"
}

resource "aws_cloudwatch_log_group" "activity" {
  name              = "/aws/lambda/${aws_lambda_function.activity_lambda.function_name}"
  retention_in_days = 1
}

## Onboard lambda
data "archive_file" "onboard_lambda" {
  type        = "zip"
  source_file = "cmd/onboard/bootstrap"
  output_path = "cmd/onboard/bootstrap.zip"
}

resource "aws_lambda_function" "onboard_lambda" {
  filename      = "${path.module}/cmd/onboard/bootstrap.zip"
  function_name = "cam-onboard"
  role          = aws_iam_role.iam_for_lambda.arn
  handler = "hello.handler"

  source_code_hash = data.archive_file.onboard_lambda.output_base64sha256

  runtime = "provided.al2023"

  environment {
    variables = {
      STRAVA_CLIENT_SECRET = local.strava_secret
      USERS_TABLE_NAME = aws_dynamodb_table.cam-users.name
    }
  }
}

resource "aws_lambda_function_url" "onboard" {
  function_name      = aws_lambda_function.onboard_lambda.function_name
  authorization_type = "NONE"
}

resource "aws_cloudwatch_log_group" "onboard" {
  name              = "/aws/lambda/${aws_lambda_function.onboard_lambda.function_name}"
  retention_in_days = 1
}
