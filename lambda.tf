data "archive_file" "lambda" {
  type        = "zip"
  source_file = "cmd/activity/bootstrap"
  output_path = "cmd/activity/bootstrap.zip"
}

resource "aws_lambda_function" "activity_lambda" {
  # If the file is not in the current working directory you will need to include a
  # path.module in the filename.
  filename      = "${path.module}/cmd/activity/bootstrap.zip"
  function_name = "cam-activity"
  role          = aws_iam_role.iam_for_lambda.arn
  handler = "hello.handler"

  source_code_hash = data.archive_file.lambda.output_base64sha256

  runtime = "provided.al2023"

  environment {
    variables = {
      STRAVA_CLIENT_SECRET = "80580d66c7d4514ecf1f904f9a698fb314b2463a"
      USERS_TABLE_NAME = "commute-and-mute-users"
      WEBHOOK_VERIFY_TOKEN = "KPHo87W@PLVCZs"
    }
  }
}

resource "aws_lambda_function_url" "activity" {
  function_name      = aws_lambda_function.activity_lambda.function_name
  authorization_type = "NONE"
}

// create log group in cloudwatch to gather logs of our lambda function
resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.activity_lambda.function_name}"
  retention_in_days = 1
}
