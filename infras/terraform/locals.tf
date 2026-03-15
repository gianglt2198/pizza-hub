locals {
  # Environment-specific overrides
  environment_configs = {
    dev = {
      nodes = {
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

    # Optional: spot instances for cost saving  
    spot = {
      desired_capacity = 1
      max_capacity     = 3
      min_capacity     = 0
      instance_types   = ["t3.medium", "t3.large"]
      capacity_type    = "SPOT"
      disk_size        = 20

      labels = {
        role = "spot"
      }

      taints = [{
        key    = "spot"
        value  = "true"
        effect = "NO_SCHEDULE"
      }]
    }
      }
    }

    prod = {
      nodes = {}
      rds_instance_class = "db.r6g.large"
    }
  }

  current_config = local.environment_configs[terraform.workspace]

  # Common tags
  common_tags = {
    Project     = "pizza-hub"
    Environment = terraform.workspace
    ManagedBy   = "gianglt1"
    Workspace   = terraform.workspace
  }
}
