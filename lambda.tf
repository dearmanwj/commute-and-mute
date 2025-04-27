## Activity Lambda
data "archive_file" "activity_lambda" {
  type        = "zip"
  source_file = "artifacts/activity/bootstrap"
  output_path = "activity.zip"
}

resource "aws_lambda_function" "activity_lambda" {
  filename      = "${path.module}/activity.zip"
  function_name = "cam-activity"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "hello.handler"

  source_code_hash = data.archive_file.activity_lambda.output_base64sha256

  runtime = "provided.al2023"

  environment {
    variables = {
      STRAVA_CLIENT_SECRET = var.strava_secret
      USERS_TABLE_NAME     = aws_dynamodb_table.cam-users.name
      WEBHOOK_VERIFY_TOKEN = var.webhook_verify_token
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

