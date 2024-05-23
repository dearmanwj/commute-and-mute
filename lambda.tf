data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "archive_file" "lambda" {
  type        = "zip"
  source_file = "cmd/activity/bootstrap"
  output_path = "cmd/activity/bootstrap.zip"
}

resource "aws_lambda_function" "test_lambda" {
  # If the file is not in the current working directory you will need to include a
  # path.module in the filename.
  filename      = "${path.module}/cmd/activity/bootstrap.zip"
  function_name = "commute-and-mute-lambda"
  role          = aws_iam_role.iam_for_lambda.arn

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
