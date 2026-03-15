data "aws_availability_zones" "available" {
  state = "available"
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

# Random password for RDS
resource "random_password" "master_password" {
  length  = 16
  special = true

  # Exclude characters that might cause issues in connection strings  
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

################################################################################
# KMS Key for RDS encryption  
################################################################################
resource "aws_kms_key" "rds" {
  description             = "${var.identifier} RDS Encryption Key"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = merge(var.tags, {
    Name = "${var.identifier}-rds-kms-key"
  })
}

resource "aws_kms_alias" "rds" {
  name          = "alias/${var.identifier}-rds"
  target_key_id = aws_kms_key.rds.key_id
}

################################################################################
# Secrets Manager for RDS credentials
################################################################################
resource "aws_secretsmanager_secret" "rds_credentials" {
  name                    = "${var.identifier}-rds_credentials"
  description             = "RDS credentials for ${var.identifier}"
  kms_key_id              = aws_kms_key.rds.arn
  recovery_window_in_days = 7

  tags = merge(var.tags, {
    Name = "${var.identifier}-rds_credentials"
  })
}

resource "aws_secretsmanager_secret_version" "rds_credentials" {
  secret_id = aws_secretsmanager_secret.rds_credentials.id

  secret_string = jsonencode({
    username = var.username
    password = random_password.master_password.result
    engine   = "postgres"
    host     = aws_db_instance.main.endpoint
    port     = aws_db_instance.main.port
    dbname   = var.database_name
    # Connection string for application  
    connection_string = "postgres://${var.username}:${random_password.master_password.result}@${aws_db_instance.main.endpoint}:${aws_db_instance.main.port}/${var.database_name}?sslmode=require"

  })

  depends_on = [aws_db_instance.main]
}

################################################################################
# Security Group for RDS
################################################################################
resource "aws_db_subnet_group" "main" {
  name       = "${var.identifier}-subnet-group"
  subnet_ids = var.subnet_ids

  tags = merge(var.tags, {
    Name = "${var.identifier}-db-subnet-group"
  })
}

resource "aws_security_group" "rds" {
  name        = "${var.identifier}-rds-sg"
  description = "Security group for RDS instance"
  vpc_id      = var.vpc_id


  # Allow inbound from application security groups  
  dynamic "ingress" {
    for_each = var.allowed_security_group_ids
    content {
      description     = "PostgreSQL from ${ingress.value}"
      from_port       = 5432
      to_port         = 5432
      protocol        = "tcp"
      security_groups = [ingress.value]
    }
  }

  # Allow inbound from application CIDR blocks  
  dynamic "ingress" {
    for_each = var.allowed_cidr_blocks
    content {
      description = "PostgreSQL from ${ingress.value}"
      from_port   = 5432
      to_port     = 5432
      protocol    = "tcp"
      cidr_blocks = [ingress.value]
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${var.identifier}-rds-sg"
  })
}

################################################################################
# RDS Instance
################################################################################
# RDS Parameter Group  
resource "aws_db_parameter_group" "main" {
  family = "postgres${var.engine_version_major}"
  name   = "${var.identifier}-pg"

  # Performance optimizations for workload  
  dynamic "parameter" {
    for_each = var.db_parameters
    content {
      name  = parameter.value.name
      value = parameter.value.value
    }
  }

  tags = merge(var.tags, {
    Name = "${var.identifier}-pg"
  })
}

# RDS Option Group (for PostgreSQL extensions)  
resource "aws_db_option_group" "main" {
  name                     = "${var.identifier}-option-group"
  option_group_description = "Option group for ${var.identifier}"
  engine_name              = "postgres"
  major_engine_version     = var.engine_version_major

  tags = merge(var.tags, {
    Name = "${var.identifier}-option-group"
  })
}

# RDS Instance  
resource "aws_db_instance" "main" {
  identifier = var.identifier

  # Engine configuration  
  engine                = "postgres"
  engine_version        = var.engine_version
  instance_class        = var.instance_class
  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.max_allocated_storage
  storage_type          = var.storage_type
  storage_encrypted     = true
  kms_key_id            = aws_kms_key.rds.arn

  # Database configuration  
  db_name  = var.database_name
  username = var.username
  password = random_password.master_password.result
  port     = 5432

  # Network configuration  
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  publicly_accessible    = false

  # Parameter and option groups  
  parameter_group_name = aws_db_parameter_group.main.name
  option_group_name    = aws_db_option_group.main.name

  # Backup configuration  
  backup_retention_period  = var.backup_retention_period
  backup_window            = var.backup_window
  maintenance_window       = var.maintenance_window
  delete_automated_backups = false

  # Snapshot configuration  
  copy_tags_to_snapshot     = true
  final_snapshot_identifier = "${var.identifier}-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  skip_final_snapshot       = var.skip_final_snapshot

  # Performance configuration  
  performance_insights_enabled = var.performance_insights_enabled
  monitoring_interval          = var.monitoring_interval
  monitoring_role_arn          = var.monitoring_interval > 0 ? aws_iam_role.rds_enhanced_monitoring[0].arn : null

  # Maintenance and updates  
  auto_minor_version_upgrade  = var.auto_minor_version_upgrade
  allow_major_version_upgrade = false

  # Deletion protection  
  deletion_protection = var.deletion_protection

  tags = merge(var.tags, {
    Name = "${var.identifier}-rds-postgresql"
  })

  depends_on = [aws_db_subnet_group.main]
}

# RDS Read Replica (optional, for read scaling)  
resource "aws_db_instance" "read_replica" {
  count = var.create_read_replica ? 1 : 0

  identifier                 = "${var.identifier}-read-replica"
  replicate_source_db        = aws_db_instance.main.identifier
  instance_class             = var.replica_instance_class
  publicly_accessible        = false
  auto_minor_version_upgrade = var.auto_minor_version_upgrade

  tags = merge(var.tags, {
    Name = "${var.identifier}-read-replica"
    Role = "read-replica"
  })
}

################################################################################
# MONITORING (Enhanced Monitoring with CloudWatch Logs)
################################################################################
# Enhanced Monitoring IAM Role (if monitoring enabled)  
resource "aws_iam_role" "rds_enhanced_monitoring" {
  count = var.monitoring_interval > 0 ? 1 : 0

  name = "${var.identifier}-rds-enhanced-monitoring-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "monitoring.rds.amazonaws.com"
        }
      }
    ]
  })

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "rds_enhanced_monitoring" {
  count = var.monitoring_interval > 0 ? 1 : 0

  role       = aws_iam_role.rds_enhanced_monitoring[0].name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
}

# CloudWatch Log Groups for PostgreSQL logs  
resource "aws_cloudwatch_log_group" "postgresql" {
  name              = "/aws/rds/instance/${var.identifier}/postgresql"
  retention_in_days = var.log_retention_days
  kms_key_id        = aws_kms_key.rds.arn

  tags = var.tags
}  