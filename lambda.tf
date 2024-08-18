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
  handler = "hello.handler"

  source_code_hash = data.archive_file.activity_lambda.output_base64sha256

  runtime = "provided.al2023"

  environment {
    variables = {
      STRAVA_CLIENT_SECRET = var.strava_secret
      USERS_TABLE_NAME = aws_dynamodb_table.cam-users.name
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

## Onboard lambda
data "archive_file" "onboard_lambda" {
  type        = "zip"
  source_file = "artifacts/onboard/bootstrap"
  output_path = "onboard.zip"
}

resource "aws_lambda_function" "onboard_lambda" {
  filename      = "${path.module}/onboard.zip"
  function_name = "cam-onboard"
  role          = aws_iam_role.iam_for_lambda.arn
  handler = "hello.handler"

  source_code_hash = data.archive_file.onboard_lambda.output_base64sha256

  runtime = "provided.al2023"

  environment {
    variables = {
      STRAVA_CLIENT_SECRET = var.strava_secret
      USERS_TABLE_NAME = aws_dynamodb_table.cam-users.name
      KMS_CAM_KEY_ID = aws_kms_key.cam_idp.key_id
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

## Users lambda
data "archive_file" "users_lambda" {
  type        = "zip"
  source_file = "artifacts/users/bootstrap"
  output_path = "users.zip"
}

resource "aws_lambda_function" "users_lambda" {
  filename      = "${path.module}/users.zip"
  function_name = "cam-users"
  role          = aws_iam_role.iam_for_lambda.arn
  handler = "hello.handler"

  source_code_hash = data.archive_file.users_lambda.output_base64sha256

  runtime = "provided.al2023"

  environment {
    variables = {
      USERS_TABLE_NAME = aws_dynamodb_table.cam-users.name
    }
  }
}

resource "aws_cloudwatch_log_group" "users" {
  name              = "/aws/lambda/${aws_lambda_function.users_lambda.function_name}"
  retention_in_days = 1
}

resource "aws_lambda_permission" "api_invoke_lambda" {
  principal = "apigateway.amazonaws.com"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.users_lambda.arn
  source_arn = "${aws_apigatewayv2_api.cam_users_api.arn}/*/*/users"
}

## Authorizer lambda
data "archive_file" "authorizer_lambda" {
  type        = "zip"
  source_file = "artifacts/authorizer/bootstrap"
  output_path = "authorizer.zip"
}

resource "aws_lambda_function" "authorizer_lambda" {
  filename      = "${path.module}/authorizer.zip"
  function_name = "cam-authorizer"
  role          = aws_iam_role.iam_for_lambda.arn
  handler = "hello.handler"

  source_code_hash = data.archive_file.authorizer_lambda.output_base64sha256

  runtime = "provided.al2023"

  environment {
    variables = {
      KMS_CAM_KEY_ID = aws_kms_key.cam_idp.key_id
    }
  }

}

resource "aws_cloudwatch_log_group" "authorizer" {
  name              = "/aws/lambda/${aws_lambda_function.authorizer_lambda.function_name}"
  retention_in_days = 1
}

resource "aws_lambda_permission" "authorizer_invoke_lambda" {
  principal = "apigateway.amazonaws.com"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.authorizer_lambda.arn
  source_arn = "${aws_apigatewayv2_api.cam_users_api.execution_arn}/authorizers/${aws_apigatewayv2_authorizer.api.id}"
}
