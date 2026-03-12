variable "aws_region" {
  description = "The AWS region to deploy resources in"
  type        = string
}

variable "environment" {  
  description = "Environment name (dev, staging, prod)"  
  type        = string  
  default     = "dev"  
  
  validation {  
    condition     = contains(["dev", "staging", "prod"], var.environment)  
    error_message = "Environment must be one of: dev, staging, prod."  
  }  
}  
