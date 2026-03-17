module "vpc" {
  source = "./modules/vpc"

  project              = var.project
  vpc_cidr             = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs

  tags = local.common_tags
}

module "eks" {
  source = "./modules/eks"

  cluster_name    = "${var.project}-${var.environment}-eks"
  cluster_version = var.eks_cluster_version
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnet_ids

  node_groups         = var.node_groups
  public_access_cidrs = var.eks_public_access_cidrs

  tags = local.common_tags
}

#################### # Optional Modules - Uncomment and configure as needed ####################
# module "rds" {
#   source = "./modules/rds"

#   identifier            = "${var.project}-${var.environment}-rds"
#   engine_version        = var.rds_engine_version
#   engine_version_major  = var.rds_engine_version_major
#   instance_class        = var.rds_instance_class
#   allocated_storage     = var.rds_allocated_storage
#   max_allocated_storage = var.rds_max_allocated_storage

#   database_name = var.database_name
#   username      = var.database_username

#   vpc_id     = module.vpc.vpc_id
#   subnet_ids = module.vpc.private_subnet_ids

#   allowed_security_group_ids = [module.eks.node_security_group_id]

#   backup_retention_period      = var.rds_backup_retention_period
#   skip_final_snapshot          = var.rds_skip_final_snapshot
#   deletion_protection          = var.rds_deletion_protection
#   performance_insights_enabled = var.rds_performance_insights_enabled
#   monitoring_interval          = var.rds_monitoring_interval

#   db_parameters = var.rds_parameters

#   tags = local.common_tags
# }

# module "redis" {
#   source = "./modules/redis"

#   project            = "${var.project}-${var.environment}-redis"
#   engine_version     = var.redis_engine_version
#   node_type          = var.redis_node_type
#   num_cache_clusters = var.redis_num_cache_clusters

#   vpc_id     = module.vpc.vpc_id
#   subnet_ids = module.vpc.private_subnet_ids

#   # Allow access from EKS nodes  
#   allowed_security_group_ids = [module.eks.node_security_group_id]

#   # Development settings  
#   transit_encryption_enabled = var.redis_transit_encryption_enabled
#   auth_token_enabled         = var.redis_auth_token_enabled
#   automatic_failover_enabled = var.redis_automatic_failover_enabled
#   multi_az_enabled           = var.redis_multi_az_enabled
#   snapshot_retention_limit   = var.redis_snapshot_retention_limit
#   enable_cloudwatch_alarms   = var.redis_enable_cloudwatch_alarms

#   # Cache parameters optimized for project
#   cache_params = var.redis_parameters

#   tags = local.common_tags
# }
