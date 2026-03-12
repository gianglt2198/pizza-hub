terraform {
  required_version = ">= 1.5"

  backend "s3" {
    bucket = "pizza-hub-terraform-state-1773328652"  
    key = "pizza-hub/terraform.tfstate"
    region = "us-east-1"
    dynamodb_table = "pizza-hub-terraform-locks"
    encrypt = true
  }  

  required_providers {
    aws = {
        source = "hashicorp/aws"
        version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project = "pizza-hub"
      Environment = var.environment
      ManagedBy = "gianglt1"
    }
  }
}