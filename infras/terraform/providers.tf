terraform {
  required_version = ">= 1.5"

  backend "s3" {
    bucket               = "pizza-hub-terraform-state-1773393334"
    workspace_key_prefix = "pizza-hub"
    key                  = "terraform.tfstate"
    region               = "us-east-1"
    dynamodb_table       = "pizza-hub-terraform-locks"
    encrypt              = true
    profile              = "localstack" # for local development with LocalStack, remove or change for production
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 3.0"
    }
  }
}

provider "aws" {
  region  = var.aws_region
  profile = "localstack" # for local development with LocalStack, remove or change for production

  default_tags {
    tags = {
      Project     = "pizza-hub"
      Environment = var.environment
      ManagedBy   = "gianglt1"
    }
  }
}
