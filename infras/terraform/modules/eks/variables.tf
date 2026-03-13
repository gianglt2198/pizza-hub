variable "node_ami_type" {
  description = "node ami type"
  type        = string
}

variable "project" {
  description = "Project name used for resource naming"
  type        = string
  default     = "pizza-hub"
}

variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}

variable "vpc_id" {
  description = "VPC ID where the EKS cluster will be deployed"
  type        = string
}

variable "private_subnets" {
  description = "List of private subnet IDs where the EKS cluster nodes will be deployed"
  type        = list(string)
}
