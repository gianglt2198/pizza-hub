output "replication_group_id" {
  description = "ElastiCache replication group ID"
  value       = aws_elasticache_replication_group.cache.id
}

output "replication_group_arn" {
  description = "ElastiCache replication group ARN"
  value       = aws_elasticache_replication_group.cache.arn
}

output "primary_endpoint_address" {
  description = "Primary endpoint address"
  value       = aws_elasticache_replication_group.cache.primary_endpoint_address
}

output "reader_endpoint_address" {
  description = "Reader endpoint address"
  value       = aws_elasticache_replication_group.cache.reader_endpoint_address
}

output "port" {
  description = "Redis port"
  value       = aws_elasticache_replication_group.cache.port
}

output "auth_token" {
  description = "Redis auth token"
  value       = var.auth_token_enabled ? random_password.auth_token[0].result : null
  sensitive   = true
}

output "secret_manager_secret_arn" {
  description = "ARN of the Secrets Manager secret containing Redis credentials"
  value       = var.auth_token_enabled ? aws_secretsmanager_secret.cache_auth_token[0].arn : null
}

output "secret_manager_secret_name" {
  description = "Name of the Secrets Manager secret containing Redis credentials"
  value       = var.auth_token_enabled ? aws_secretsmanager_secret.cache_auth_token[0].name : null
}

output "connection_string" {
  description = "Redis connection string"
  value = var.auth_token_enabled ? (
    var.transit_encryption_enabled ?
    "rediss://:${random_password.auth_token[0].result}@${aws_elasticache_replication_group.cache.primary_endpoint_address}:${aws_elasticache_replication_group.cache.port}" :
    "redis://:${random_password.auth_token[0].result}@${aws_elasticache_replication_group.cache.primary_endpoint_address}:${aws_elasticache_replication_group.cache.port}"
    ) : (
    var.transit_encryption_enabled ?
    "rediss://${aws_elasticache_replication_group.cache.primary_endpoint_address}:${aws_elasticache_replication_group.cache.port}" :
    "redis://${aws_elasticache_replication_group.cache.primary_endpoint_address}:${aws_elasticache_replication_group.cache.port}"
  )
  sensitive = true
}

output "subnet_group_name" {
  description = "ElastiCache subnet group name"
  value       = aws_elasticache_subnet_group.main.name
}

output "security_group_id" {
  description = "Security group ID for Redis"
  value       = aws_security_group.cache.id
}

output "kms_key_id" {
  description = "KMS key ID used for encryption"
  value       = aws_kms_key.cache.key_id
}

output "kms_key_arn" {
  description = "KMS key ARN used for encryption"
  value       = aws_kms_key.cache.arn
}

output "cloudwatch_log_group_name" {
  description = "CloudWatch log group name for Redis slow log"
  value       = aws_cloudwatch_log_group.cache_slow.name
}  