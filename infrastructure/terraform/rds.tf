# =========================================
# =========================================
# RDS
# =========================================
# =========================================

resource "aws_db_instance" "db" {
  identifier                      = "todo-db"
  allocated_storage               = 20
  auto_minor_version_upgrade      = true
  backup_retention_period         = 1
  copy_tags_to_snapshot           = true
  db_subnet_group_name            = "todo-db"
  delete_automated_backups        = true
  deletion_protection             = true
  enabled_cloudwatch_logs_exports = ["error", "slowquery"]
  engine                          = "mysql"
  engine_version                  = "5.7.36"
  instance_class                  = "db.t3.micro"
  kms_key_id                      = aws_kms_key.db.arn
  multi_az                        = true
  max_allocated_storage           = 100
  monitoring_interval             = 60
  monitoring_role_arn             = aws_iam_role.db.arn
  name                            = var.db_name
  option_group_name               = aws_db_option_group.db.name
  parameter_group_name            = aws_db_parameter_group.db.name
  password                        = var.db_password
  port                            = 3306
  skip_final_snapshot             = true
  storage_encrypted               = true
  storage_type                    = "gp2"
  username                        = var.db_username
  vpc_security_group_ids          = [aws_security_group.db.id]

  lifecycle {
    ignore_changes = [password]
  }
}

resource "aws_db_option_group" "db" {
  name                     = "mysql5-7"
  option_group_description = "mysql5-7 option group"
  engine_name              = "mysql"
  major_engine_version     = "5.7"
}

resource "aws_db_parameter_group" "db" {
  name        = "mysql5-7"
  family      = "mysql5.7"
  description = "mysql5.7 parameter group"

  parameter {
    apply_method = "immediate"
    name         = "character_set_client"
    value        = "utf8mb4"
  }

  parameter {
    apply_method = "immediate"
    name         = "character_set_connection"
    value        = "utf8mb4"
  }

  parameter {
    apply_method = "immediate"
    name         = "character_set_database"
    value        = "utf8mb4"
  }

  parameter {
    apply_method = "immediate"
    name         = "character_set_filesystem"
    value        = "utf8mb4"
  }

  parameter {
    apply_method = "immediate"
    name         = "character_set_results"
    value        = "utf8mb4"
  }

  parameter {
    apply_method = "immediate"
    name         = "character_set_server"
    value        = "utf8mb4"
  }

  parameter {
    apply_method = "immediate"
    name         = "collation_connection"
    value        = "utf8mb4_unicode_ci"
  }

  parameter {
    apply_method = "immediate"
    name         = "collation_server"
    value        = "utf8mb4_unicode_ci"
  }

  parameter {
    apply_method = "immediate"
    name         = "time_zone"
    value        = var.time_zone
  }
}

resource "aws_iam_role" "db" {
  name = "todo-rds-monitoring-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "monitoring.rds.amazonaws.com"
        }
      },
    ]
  })

  managed_policy_arns = [data.aws_iam_policy.db.arn]
}

data "aws_iam_policy" "db" {
  name = "AmazonRDSEnhancedMonitoringRole"
}

resource "aws_kms_key" "db" {
  description             = "KMS for db"
  deletion_window_in_days = 7
}
