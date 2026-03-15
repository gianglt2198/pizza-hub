variable "aws_region" {
  description = "The AWS region to deploy resources in"
  type        = string
  default     = "us-east-1"
}

variable "eks_cluster_name" {
  description = "The name of EKS in AWS"
  type        = string
  default     = "pizza-hub-eks-cluster"
}

variable "eks_cluster_version" {
  description = "the Kubernetes version of EKS in AWS"
  type        = string
  default     = "1.35"
}

variable "project" {
  description = "Project name"
  type        = string
}


variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string

  validation {
    condition     = can(cidrhost(var.vpc_cidr, 0))
    error_message = "VPC CIDR must be a valid IPv4 CIDR block."
  }
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)

  validation {
    condition     = length(var.public_subnet_cidrs) >= 2
    error_message = "At least 2 public subnets are required for high availability."
  }
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)

  validation {
    condition     = length(var.private_subnet_cidrs) >= 2
    error_message = "At least 2 private subnets are required for high availability."
  }
}

variable "database_name" {
  description = "Initial database name"
  type        = string
  default     = "pizzahub"
}

variable "database_username" {
  description = "Username for Database"
  type        = string
  default     = "pizza"
}

variable "rds_engine_version" {
  description = "RDS engine version"
  type        = string
  default     = "15.4"
}

variable "rds_engine_version_major" {
  description = "RDS engine version major"
  type        = string
  default     = "15"
}

