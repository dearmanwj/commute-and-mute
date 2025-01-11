variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "eu-north-1"
}

variable "strava_secret" {
  type      = string
  sensitive = true
}

variable "webhook_verify_token" {
  type      = string
  sensitive = true
}
