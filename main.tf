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
    key    = "commute-and-mute-tfstate"
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
