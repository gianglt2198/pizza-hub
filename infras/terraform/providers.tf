terraform {
  required_version = ">= 1.5"

  backend "s3" {
    bucket                      = "pizza-hub-terraform-state-1773393334"
    workspace_key_prefix        = "pizza-hub"
    key                         = "terraform.tfstate"
    region                      = "us-east-1"
    dynamodb_table              = "pizza-hub-terraform-locks"
    encrypt                     = true
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    force_path_style            = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.28"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 3.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = local.common_tags
  }
}
