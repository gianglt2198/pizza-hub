# Basic Configuration
aws_region  = "us-east-1"
project     = "pizza-hub"
environment = "dev"

eks_cluster_version = "1.35"


# VPC Configuration - Development
vpc_cidr             = "10.0.0.0/16"
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.11.0/24", "10.0.12.0/24"]

# EKS Configuration - Development
eks_public_access_cidrs = ["0.0.0.0/0"] # More permissive for dev

# EKS Node Groups - Smaller for dev
node_groups = {
  general = {
    desired_capacity = 1
    max_capacity     = 2
    min_capacity     = 1
    instance_types   = ["t3.small"]
    capacity_type    = "ON_DEMAND"
    disk_size        = 20
    labels = {
      role        = "general"
      environment = "dev"
    }
  }
}

# RDS Configuration - Cost-optimized for dev
rds_instance_class               = "db.t3.micro"
rds_allocated_storage            = 20
rds_max_allocated_storage        = 100
rds_backup_retention_period      = 1
rds_skip_final_snapshot          = true
rds_deletion_protection          = false
rds_performance_insights_enabled = false
rds_monitoring_interval          = 0

# Dev-specific RDS parameters
rds_parameters = [
  {
    name  = "shared_preload_libraries"
    value = "pg_stat_statements"
  },
  {
    name  = "log_statement"
    value = "ddl" # Only DDL for dev
  },
  {
    name  = "max_connections"
    value = "50" # Lower for dev
  }
]

# Redis Configuration - Minimal for dev
redis_node_type                  = "cache.t3.micro"
redis_num_cache_clusters         = 1
redis_transit_encryption_enabled = false
redis_auth_token_enabled         = false
redis_automatic_failover_enabled = false
redis_multi_az_enabled           = false
redis_snapshot_retention_limit   = 1
redis_enable_cloudwatch_alarms   = false

redis_parameters = [
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
    value = "100" # Lower for dev
  }
]

# Database Configuration
database_name            = "pizzahub_dev"
database_username        = "pizza_dev"
rds_engine_version       = "15.4"
rds_engine_version_major = "15"
redis_engine_version     = "7.0"
