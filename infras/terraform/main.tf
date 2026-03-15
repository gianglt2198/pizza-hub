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

  cluster_name    = var.eks_cluster_name
  cluster_version = var.eks_cluster_version
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnet_ids

  tags = local.common_tags

  node_groups = local.current_config.nodes

  public_access_cidrs = local.current_config.eks.public_access_cidrs # Restrict this in production  
}

module "rds" {
  source = "./modules/rds"

  identifier            = "${var.project}-${terraform.workspace}"
  engine_version        = var.rds_engine_version
  engine_version_major  = var.rds_engine_version_major
  instance_class        = local.current_config.rds.instance_class
  allocated_storage     = local.current_config.rds.allocated_storage
  max_allocated_storage = local.current_config.rds.max_allocated_storage

  database_name = var.database_name
  username      = var.database_username

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnet_ids

  # Allow access from EKS nodes  
  allowed_security_group_ids = [module.eks.node_security_group_id]

  # Development settings  
  backup_retention_period      = local.current_config.rds.backup_retention_period
  skip_final_snapshot          = local.current_config.rds.skip_final_snapshot          # Skip snapshot for dev  
  deletion_protection          = local.current_config.rds.deletion_protection          # Allow deletion in dev  
  performance_insights_enabled = local.current_config.rds.performance_insights_enabled # Disable PI for cost saving  
  monitoring_interval          = local.current_config.rds.monitoring_interval          # Disable enhanced monitoring  

  # Database parameters optimized for pizza-hub[^7]  
  db_parameters = local.current_config.rds.parameters

  tags = local.common_tags
}

module "redis" {
  source = "./modules/redis"

  project            = "${var.project}-${terraform.workspace}"
  engine_version     = var.redis_engine_version
  node_type          = var.redis_node_type
  num_cache_clusters = local.current_config.cache.number_cache_clusters # Single node for dev  

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnet_ids

  # Allow access from EKS nodes  
  allowed_security_group_ids = [module.eks.node_security_group_id]

  # Development settings  
  transit_encryption_enabled = local.current_config.cache.transit_encryption_enabled
  auth_token_enabled         = local.current_config.cache.auth_token_enabled
  automatic_failover_enabled = local.current_config.cache.automatic_failover_enabled
  multi_az_enabled           = local.current_config.cache.multi_az_enabled
  snapshot_retention_limit   = local.current_config.cache.snapshot_retention_limit
  enable_cloudwatch_alarms   = local.current_config.cache.enable_cloudwatch_alarms

  # Cache parameters optimized for pizza-hub[^8]  
  cache_params = local.current_config.cache.parameters

  tags = local.common_tags
}
