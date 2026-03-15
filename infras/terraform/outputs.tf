
output "cluster_arn" {
  description = "The Amazon Resource Name (ARN) of the cluster"
  value       = module.eks.cluster_arn
}

output "cluster_id" {
  description = "The ID of the EKS cluster. Note: currently a value is returned only for local EKS clusters created on Outposts"
  value       = module.eks.cluster_id
}

output "cluster_name" {
  description = "The name of the EKS cluster"
  value       = module.eks.cluster_name
}


output "rds_endpoint" {
  value = module.rds.db_endpoint
}

output "rds_port" {
  value = module.rds.db_port
}

output "rds_database_name" {
  value = module.rds.db_name
}

output "rds_secret_arn" {
  description = "ARN of the secret containing RDS credentials"
  value       = module.rds.secret_manager_secret_arn
}

output "rds_password" {
  description = "Master password for the RDS instance"
  value       = module.rds.db_password
  sensitive   = true
}

output "rds_connection_string" {
  description = "Connection string for the RDS instance"
  value       = module.rds.db_connection_string
  sensitive   = true
}


output "redis_primary_endpoint" {
  value = module.redis.primary_endpoint_address
}

output "redis_port" {
  value = module.redis.port
}

output "redis_connection_string" {
  description = "Redis connection string for application"
  value       = module.redis.connection_string
  sensitive   = true
}  