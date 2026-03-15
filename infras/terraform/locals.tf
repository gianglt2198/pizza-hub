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
        instance_class               = "db.t3.micro"
        allocated_storage            = 20
        max_allocated_storage        = 100
        backup_retention_period      = 3
        skip_final_snapshot          = true
        deletion_protection          = false
        performance_insights_enabled = false
        monitoring_interval          = 0
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

      cache = {
        number_cache_clusters      = 1     # Single node for dev
        transit_encryption_enabled = false # Disable TLS for dev simplicity  
        auth_token_enabled         = false # No auth for dev  
        automatic_failover_enabled = false # Single node, no failover  
        multi_az_enabled           = false # No multi-AZ for dev  
        snapshot_retention_limit   = 1     # Minimal backup for dev  
        enable_cloudwatch_alarms   = false # No alarms for dev  
        parameters = [
          {
            name  = "maxmemory-policy"
            value = "allkeys-lru" # Evict least recently used keys  
          },
          {
            name  = "timeout"
            value = "300" # 5 minutes timeout  
          },
          {
            name  = "tcp-keepalive"
            value = "300"
          },
          {
            name  = "maxclients"
            value = "100" # Lower for dev  
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
        instance_class               = "db.r6g.large"
        allocated_storage            = 100
        max_allocated_storage        = 1000
        backup_retention_period      = 7
        skip_final_snapshot          = false
        deletion_protection          = true
        performance_insights_enabled = true
        monitoring_interval          = 60
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

      cache = {
        number_cache_clusters      = 2     # Multi-node for high availability
        transit_encryption_enabled = false # Disable TLS for dev simplicity  
        auth_token_enabled         = false # No auth for dev  
        automatic_failover_enabled = false # Single node, no failover  
        multi_az_enabled           = false # No multi-AZ for dev  
        snapshot_retention_limit   = 1     # Minimal backup for dev  
        enable_cloudwatch_alarms   = false # No alarms for dev  
        parameters = [
          {
            name  = "maxmemory-policy"
            value = "allkeys-lru" # Evict least recently used keys  
          },
          {
            name  = "timeout"
            value = "300" # 5 minutes timeout  
          },
          {
            name  = "tcp-keepalive"
            value = "300"
          },
          {
            name  = "maxclients"
            value = "100" # Lower for dev  
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
