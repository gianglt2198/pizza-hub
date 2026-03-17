variable "aws_region" {
  description = "The AWS region to deploy resources in"
  type        = string
  default     = "us-east-1"
}

variable "project" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

###################################################################
# VPC
####################################################################
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

####################################################################
# EKS
####################################################################
variable "eks_cluster_name" {
  description = "The name of EKS in AWS"
  type        = string
  default     = "${var.project}-eks-cluster"
}

variable "eks_cluster_version" {
  description = "the Kubernetes version of EKS in AWS"
  type        = string
  default     = "1.35"
}

# EKS Node Groups Configuration  
variable "node_groups" {
  description = "EKS node group configurations"
  type = map(object({
    desired_capacity = number
    max_capacity     = number
    min_capacity     = number
    instance_types   = list(string)
    capacity_type    = optional(string, "ON_DEMAND")
    disk_size        = optional(number, 20)
    labels           = optional(map(string), {})
    taints = optional(list(object({
      key    = string
      value  = string
      effect = string
    })), [])
  }))
  default = {
    general = {
      desired_capacity = 2
      max_capacity     = 4
      min_capacity     = 1
      instance_types   = ["t3.medium"]
      capacity_type    = "ON_DEMAND"
      disk_size        = 20
      labels = {
        role = "general"
      }
    }
  }
}

# EKS Access Configuration  
variable "eks_public_access_cidrs" {
  description = "CIDR blocks that can access the EKS cluster endpoint"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

############################################################################################
# RDS
############################################################################################
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

# RDS Configuration Variables  
variable "rds_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "rds_allocated_storage" {
  description = "RDS allocated storage in GB"
  type        = number
  default     = 20
}

variable "rds_max_allocated_storage" {
  description = "RDS maximum allocated storage for autoscaling"
  type        = number
  default     = 100
}

variable "rds_backup_retention_period" {
  description = "RDS backup retention period in days"
  type        = number
  default     = 7
}

variable "rds_skip_final_snapshot" {
  description = "Whether to skip final snapshot when deleting RDS"
  type        = bool
  default     = false
}

variable "rds_deletion_protection" {
  description = "Enable deletion protection for RDS"
  type        = bool
  default     = true
}

variable "rds_performance_insights_enabled" {
  description = "Enable RDS Performance Insights"
  type        = bool
  default     = false
}

variable "rds_monitoring_interval" {
  description = "RDS enhanced monitoring interval"
  type        = number
  default     = 0
  validation {
    condition     = contains([0, 1, 5, 10, 15, 30, 60], var.rds_monitoring_interval)
    error_message = "Monitoring interval must be 0, 1, 5, 10, 15, 30, or 60 seconds."
  }
}

# RDS Parameters  
variable "rds_parameters" {
  description = "RDS database parameters"
  type = list(object({
    name  = string
    value = string
  }))
  default = [
    {
      name  = "shared_preload_libraries"
      value = "pg_stat_statements"
    },
    {
      name  = "log_statement"
      value = "all"
    },
    {
      name  = "max_connections"
      value = "100"
    }
  ]
}


######################################################################
# ElastiCache
######################################################################
variable "redis_engine_version" {
  description = "Redis engine version for ElastiCache"
  type        = string
  default     = "7.0"
}

variable "redis_node_type" {
  description = "Redis node type for ElastiCache"
  type        = string
  default     = "cache.t3.micro"
}
# Redis Configuration Variables  
variable "redis_node_type" {
  description = "Redis node type"
  type        = string
  default     = "cache.t3.micro"
}

variable "redis_num_cache_clusters" {
  description = "Number of Redis cache clusters"
  type        = number
  default     = 1
}

variable "redis_transit_encryption_enabled" {
  description = "Enable Redis transit encryption"
  type        = bool
  default     = false
}

variable "redis_auth_token_enabled" {
  description = "Enable Redis authentication"
  type        = bool
  default     = false
}

variable "redis_automatic_failover_enabled" {
  description = "Enable Redis automatic failover"
  type        = bool
  default     = false
}

variable "redis_multi_az_enabled" {
  description = "Enable Redis Multi-AZ"
  type        = bool
  default     = false
}

variable "redis_snapshot_retention_limit" {
  description = "Redis snapshot retention limit"
  type        = number
  default     = 3
}

variable "redis_enable_cloudwatch_alarms" {
  description = "Enable Redis CloudWatch alarms"
  type        = bool
  default     = false
}

variable "redis_parameters" {
  description = "Redis cache parameters"
  type = list(object({
    name  = string
    value = string
  }))
  default = [
    {
      name  = "maxmemory-policy"
      value = "allkeys-lru"
    },
    {
      name  = "timeout"
      value = "300"
    },
    {
      name  = "maxclients"
      value = "1000"
    }
  ]
}
