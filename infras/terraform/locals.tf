locals {
  # Environment-specific overrides
  environment_configs = {
    dev = {
      eks = {
        public_access_cidrs = ["0.0.0.0/0"] # Restrict this in production  
      }
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
      rds = {
        instance_class        = "db.t3.micro"
        allocated_storage     = 20
        max_allocated_storage = 100
        backup_retention_period = 3
        skip_final_snapshot = true
        deletion_protection = false
        performance_insights_enabled = false
        monitoring_interval = 0
        parameters = [
          {
            name  = "shared_preload_libraries"
            value = "pg_stat_statements"
          },
          {
            name  = "log_statement"
            value = "ddl" # Log only DDL in dev  
          },
          {
            name  = "max_connections"
            value = "50" # Lower for dev environment  
          },
          {
            name  = "work_mem"
            value = "4MB"
          }
        ]
      }
    }

    prod = {
      eks = {
        public_access_cidrs = ["0.0.0.0/0"] # Restrict this in production  
      }
      nodes = {}
      rds = {
        instance_class        = "db.r6g.large"
        allocated_storage     = 100
        max_allocated_storage = 1000
        backup_retention_period = 7
        skip_final_snapshot = false
        deletion_protection = true
        performance_insights_enabled = true
        monitoring_interval = 60
        parameters = [
          {
            name  = "shared_preload_libraries"
            value = "pg_stat_statements"
          },
          {
            name  = "max_connections"
            value = "50" # Lower for dev environment  
          },
          {
            name  = "work_mem"
            value = "4MB"
          }
        ]
      }
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
