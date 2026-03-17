# Basic Configuration
aws_region          = "us-east-1"
project             = "pizza-hub"
environment         = "prod"
eks_cluster_version = "1.35"

# VPC Configuration - Production (different CIDR)
vpc_cidr             = "10.1.0.0/16"
public_subnet_cidrs  = ["10.1.1.0/24", "10.1.2.0/24"]
private_subnet_cidrs = ["10.1.11.0/24", "10.1.12.0/24"]

# EKS Configuration - Production
eks_public_access_cidrs = ["203.0.113.0/24"] # Restricted to office/VPN

# EKS Node Groups - Production scale
node_groups = {
  general = {
    desired_capacity = 3
    max_capacity     = 10
    min_capacity     = 2
    instance_types   = ["t3.large", "t3.xlarge"]
    capacity_type    = "ON_DEMAND"
    disk_size        = 50
    labels = {
      role        = "general"
      environment = "prod"
    }
  }
}

# RDS Configuration - Production grade
rds_instance_class               = "db.r6g.large"
rds_allocated_storage            = 100
rds_max_allocated_storage        = 1000
rds_backup_retention_period      = 7
rds_skip_final_snapshot          = false
rds_deletion_protection          = true
rds_performance_insights_enabled = true
rds_monitoring_interval          = 60

# Production RDS parameters
rds_parameters = [
  {
    name  = "shared_preload_libraries"
    value = "pg_stat_statements"
  },
  {
    name  = "log_statement"
    value = "all" # Full logging for prod
  },
  {
    name  = "max_connections"
    value = "200"
  },
  {
    name  = "work_mem"
    value = "16MB"
  }
]

# Redis Configuration - Production grade
redis_node_type                  = "cache.r6g.large"
redis_num_cache_clusters         = 2
redis_transit_encryption_enabled = true
redis_auth_token_enabled         = true
redis_automatic_failover_enabled = true
redis_multi_az_enabled           = true
redis_snapshot_retention_limit   = 7
redis_enable_cloudwatch_alarms   = true

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
    value = "1000"
  }
]

# Database Configuration
database_name            = "pizzahub"
database_username        = "pizza"
rds_engine_version       = "15.4"
rds_engine_version_major = "15"
redis_engine_version     = "7.0"
