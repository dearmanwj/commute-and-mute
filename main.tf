terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.38.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.4.2"
    }
  }

  backend "s3" {
    bucket = "commute-and-mute-tf-backend"
    key = "commute-and-mute-tfstate"
    region = "eu-north-1"
  }

  required_version = "~> 1.2"
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      project = "commute-and-mute"
    }
  }

}

resource "aws_s3_bucket" "authorizer_lambda_bucket" {
  bucket = "commute-and-mute-authorizer-bucket"
}

resource "aws_s3_bucket_ownership_controls" "authorizer_lambda_bucket" {
  bucket = aws_s3_bucket.authorizer_lambda_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "authorizer_lambda_bucket" {
  depends_on = [aws_s3_bucket_ownership_controls.authorizer_lambda_bucket]

  bucket = aws_s3_bucket.authorizer_lambda_bucket.id
  acl    = "private"
}
