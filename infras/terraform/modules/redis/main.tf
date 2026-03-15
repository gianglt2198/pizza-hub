# Data sources  
data "aws_availability_zones" "available" {
  state = "available"
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

# Random auth token for Redis (if enabled)  
resource "random_password" "auth_token" {
  count = var.auth_token_enabled ? 1 : 0

  length  = 32
  special = true

  # ElastiCache auth token constraints  
  override_special = "!&#$^<>-"
}


################################################################################
# Secret Manager secret for Redis auth token (if enabled) 
################################################################################

resource "aws_secretsmanager_secret" "cache_auth_token" {
  count = var.auth_token_enabled ? 1 : 0

  name                    = "${var.project}-cache-auth-token"
  description             = "Secret for Cache auth token"
  kms_key_id              = aws_kms_key.cache.arn
  recovery_window_in_days = 7

  tags = merge(
    var.tags,
    {
      Name = "${var.project}-cache-auth-token"
    }
  )
}

resource "aws_secretsmanager_secret_version" "cache_auth_token" {
  count = var.auth_token_enabled ? 1 : 0

  secret_id = aws_secretsmanager_secret.cache_auth_token[0].id
  secret_string = jsonencode({
    auth_token = random_password.auth_token[0].result
    endpoint   = aws_elasticache_replication_group.cache.primary_endpoint_address
    port       = aws_elasticache_replication_group.cache.port
    # Connection string for application  
    connection_string = var.transit_encryption_enabled ? "rediss://:${random_password.auth_token[0].result}@${aws_elasticache_replication_group.cache.primary_endpoint_address}:${aws_elasticache_replication_group.cache.port}" : "redis://:${random_password.auth_token[0].result}@${aws_elasticache_replication_group.cache.primary_endpoint_address}:${aws_elasticache_replication_group.cache.port}"
  })
}


################################################################################
# KMS Key for ElastiCache encryption  
################################################################################
resource "aws_kms_key" "cache" {
  description             = "KMS key for ElastiCache encryption"
  deletion_window_in_days = 10
  enable_key_rotation     = true

  tags = merge(
    var.tags,
    {
      Name = "${var.project}-cache-kms-key"
    }
  )
}

resource "aws_kms_alias" "cache" {
  name          = "alias/${var.project}-cache-kms-key"
  target_key_id = aws_kms_key.cache.key_id
}

################################################################################
# Networking for ElastiCache
################################################################################
resource "aws_elasticache_subnet_group" "main" {
  name       = "${var.project}-cache-subnet-group"
  subnet_ids = var.subnet_ids

  tags = merge(
    var.tags,
    {
      Name = "${var.project}-cache-subnet-group"
    }
  )
}

resource "aws_security_group" "cache" {
  name        = "${var.project}-cache-sg"
  description = "Security group for ElastiCache cluster"
  vpc_id      = var.vpc_id

  # Allow inbound from allowed security groups
  dynamic "ingress" {
    for_each = var.allowed_security_group_ids
    content {
      description     = "Cache from ${ingress.value}"
      from_port       = 6379
      to_port         = 6379
      protocol        = "tcp"
      security_groups = [ingress.value]
    }
  }


  # Allow inbound from application CIDR blocks  
  dynamic "ingress" {
    for_each = var.allowed_cidr_blocks
    content {
      description = "Redis from ${ingress.value}"
      from_port   = 6379
      to_port     = 6379
      protocol    = "tcp"
      cidr_blocks = [ingress.value]
    }
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.tags,
    {
      Name = "${var.project}-cache-sg"
    }
  )
}

################################################################################
# ElastiCache Cluster
################################################################################
resource "aws_elasticache_parameter_group" "cache" {
  family      = "redis${var.engine_version}"
  name        = "${var.project}-cache-parameter-group"
  description = "Parameter group for ElastiCache cluster"

  dynamic "parameter" {
    for_each = var.cache_params
    content {
      name  = parameter.value.name
      value = parameter.value.value
    }
  }

  tags = merge(
    var.tags,
    {
      Name = "${var.project}-cache-parameter-group"
    }
  )
}


resource "aws_elasticache_replication_group" "cache" {
  replication_group_id = "${var.project}-cache"
  description          = "ElastiCache replication group for ${var.project}"

  # Redis configuration  
  engine               = var.engine
  engine_version       = var.engine_version
  node_type            = var.node_type
  port                 = 6379
  parameter_group_name = aws_elasticache_parameter_group.cache.name

  # Cluster configuration
  num_cache_clusters          = var.cluster_mode_enabled ? null : var.num_cache_clusters
  preferred_cache_cluster_azs = var.preferred_cache_cluster_azs

  # Network configuration
  subnet_group_name  = aws_elasticache_subnet_group.main.name
  security_group_ids = [aws_security_group.cache.id]

  # Security configuration
  at_rest_encryption_enabled = var.at_rest_encryption_enabled
  transit_encryption_enabled = var.transit_encryption_enabled
  kms_key_id                 = aws_kms_key.cache.arn
  auth_token                 = var.auth_token_enabled ? random_password.auth_token[0].result : null
  auth_token_update_strategy = var.auth_token_enabled ? "ROTATE" : null

  # Backup configuration  
  snapshot_retention_limit = var.snapshot_retention_limit
  snapshot_window          = var.snapshot_window

  # Maintenance configuration  
  maintenance_window         = var.maintenance_window
  auto_minor_version_upgrade = var.auto_minor_version_upgrade

  # Multi-AZ configuration  
  automatic_failover_enabled = var.automatic_failover_enabled
  multi_az_enabled           = var.multi_az_enabled

  # Notification configuration  
  notification_topic_arn = var.notification_topic_arn


  # Logging configuration  
  log_delivery_configuration {
    destination      = aws_cloudwatch_log_group.cache_slow.name
    destination_type = "cloudwatch-logs"
    log_format       = "text"
    log_type         = "slow-log"
  }

  tags = merge(
    var.tags,
    {
      Name = "${var.project}-cache"
    }
  )
}

################################################################################
# Monitoring and Logging
################################################################################
resource "aws_cloudwatch_log_group" "cache_slow" {
  name              = "/aws/elasticache/redis/${var.project}/slow-log"
  retention_in_days = var.log_retention_days
  kms_key_id        = aws_kms_key.cache.arn


  tags = merge(
    var.tags,
    {
      Name = "${var.project}-cache-slow-log"
    }
  )
}

# CloudWatch Alarms for Redis monitoring  
resource "aws_cloudwatch_metric_alarm" "cache_cpu" {
  count = var.enable_cloudwatch_alarms ? 1 : 0

  alarm_name          = "${var.project}-redis-high-cpu"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ElastiCache"
  period              = "120"
  statistic           = "Average"
  threshold           = "80"
  alarm_description   = "This metric monitors Redis CPU utilization"
  alarm_actions       = var.alarm_actions

  dimensions = {
    CacheClusterId = "${var.project}-001"
  }

  tags = merge(var.tags, {
    Name = "${var.project}-redis-high-cpu-alarm"
  })
}

resource "aws_cloudwatch_metric_alarm" "cache_memory" {
  count = var.enable_cloudwatch_alarms ? 1 : 0

  alarm_name          = "${var.project}-redis-high-memory"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "DatabaseMemoryUsagePercentage"
  namespace           = "AWS/ElastiCache"
  period              = "120"
  statistic           = "Average"
  threshold           = "80"
  alarm_description   = "This metric monitors Redis memory usage"
  alarm_actions       = var.alarm_actions

  dimensions = {
    CacheClusterId = "${var.project}-001"
  }

  tags = merge(var.tags, {
    Name = "${var.project}-redis-high-memory-alarm"
  })
}

resource "aws_cloudwatch_metric_alarm" "cache_connections" {
  count = var.enable_cloudwatch_alarms ? 1 : 0

  alarm_name          = "${var.project}-redis-high-connections"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "CurrConnections"
  namespace           = "AWS/ElastiCache"
  period              = "120"
  statistic           = "Average"
  threshold           = "100"
  alarm_description   = "This metric monitors Redis connection count"
  alarm_actions       = var.alarm_actions

  dimensions = {
    CacheClusterId = "${var.project}-001"
  }

  tags = merge(var.tags, {
    Name = "${var.project}-redis-high-connections-alarm"
  })
}  