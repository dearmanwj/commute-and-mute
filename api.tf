resource "aws_apigatewayv2_api" "cam_users_api" {
  name          = "cam-users-api"
  protocol_type = "HTTP"
  
  cors_configuration {
    allow_origins = ["*"]
  }
}

resource "aws_apigatewayv2_integration" "cam_users_integration" {
  api_id           = aws_apigatewayv2_api.cam_users_api.id
  integration_type = "AWS_PROXY"

  description            = "CAM User lambda integration"
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.users_lambda.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "route" {
  api_id = aws_apigatewayv2_api.cam_users_api.id
  route_key = "GET /users"
  target = "integrations/${aws_apigatewayv2_integration.cam_users_integration.id}"
}

resource "aws_apigatewayv2_stage" "stage" {
  api_id = aws_apigatewayv2_api.cam_users_api.id
  name   = "$default"
  auto_deploy = "true"
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api.arn
    format = "$context.identity.sourceIp $context.identity.caller $context.identity.user [$context.requestTime] $context.httpMethod $context.resourcePath $context.protocol $context.status $context.responseLength $context.requestId $context.extendedRequestId"
  }
}

resource "aws_cloudwatch_log_group" "api" {
  name              = "API-Gateway-Execution-Logs_${aws_apigatewayv2_api.cam_users_api.id}"
  retention_in_days = 1
}