# =========================================
# =========================================
# ECS
# =========================================
# =========================================

# =========================================
# ECS for application
# =========================================
locals {
  cluster_name = "${var.app_name}-ecs-cluster"
}

resource "aws_ecs_cluster" "application" {
  name = local.cluster_name

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = "${var.app_name}-ecs-cluster"
  }
}

locals {
  container_name = var.app_name
}
resource "aws_ecs_task_definition" "application" {
  family                   = "${var.app_name}-task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 1024
  memory                   = 2048
  skip_destroy             = true

  task_role_arn      = aws_iam_role.task.arn
  execution_role_arn = aws_iam_role.task_execution.arn

  container_definitions = format("[%s]", templatefile(
    "${path.module}/container_definitions.json",
    {
      container_name = local.container_name
      region         = var.region
      image_arn      = var.image_arn
      logs_group     = aws_cloudwatch_log_group.ecs_task.name
      cpu            = 128
      memory         = 256
      entry_point    = "server"
      db_host        = aws_db_instance.db.address
      db_username    = var.db_username
      db_password    = var.db_password
      db_name        = var.db_name
    }
  ))

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = null
  }

  tags = {
    Name = "${var.app_name}-task-def-ecs"
  }
}

resource "aws_cloudwatch_log_group" "ecs_task" {
  name              = "/ecs/${var.app_name}-task"
  retention_in_days = 30

  tags = {
    Name = "${var.app_name}-lg-ecs"
  }
}

resource "aws_ecs_service" "application" {
  name                              = "${var.app_name}-ecs-service"
  cluster                           = aws_ecs_cluster.application.id
  task_definition                   = aws_ecs_task_definition.application.arn
  launch_type                       = "FARGATE"
  scheduling_strategy               = "REPLICA"
  desired_count                     = 2
  health_check_grace_period_seconds = 60
  enable_execute_command            = true

  deployment_controller {
    type = "CODE_DEPLOY"
  }

  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100

  enable_ecs_managed_tags = true

  load_balancer {
    target_group_arn = aws_lb_target_group.application["blue"].arn
    container_name   = local.container_name
    container_port   = 8080
  }

  network_configuration {
    subnets          = [for subnet in aws_subnet.private : subnet.id]
    security_groups  = [aws_security_group.private.id]
    assign_public_ip = false
  }

  lifecycle {
    ignore_changes = [desired_count, task_definition, load_balancer]
  }

  tags = {
    Name = "${var.app_name}-ecs-service"
  }
}

# =========================================
# ECS for management server
# =========================================
locals {
  management_cluster_name = "${var.app_name}-ecs-management-cluster"
}

resource "aws_ecs_cluster" "management" {
  name = local.management_cluster_name

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = "${var.app_name}-ecs-management-cluster"
  }
}

locals {
  management_container_name = "management"
}
resource "aws_ecs_task_definition" "management" {
  family                   = "${var.app_name}-management-task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 1024
  memory                   = 2048
  skip_destroy             = true

  task_role_arn      = aws_iam_role.task.arn
  execution_role_arn = aws_iam_role.task_execution.arn

  container_definitions = format("[%s]", templatefile(
    "${path.module}/container_definitions.json",
    {
      container_name = local.management_container_name
      region         = var.region
      image_arn      = var.management_image_arn
      logs_group     = aws_cloudwatch_log_group.ecs_management_task.name
      cpu            = 128
      memory         = 1024
      entry_point    = "top"
      db_host        = aws_db_instance.db.address
      db_username    = var.db_username
      db_password    = var.db_password
      db_name        = var.db_name
    }
  ))

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = null
  }

  tags = {
    Name = "${var.app_name}-management-task-def-ecs"
  }
}

resource "aws_cloudwatch_log_group" "ecs_management_task" {
  name              = "/ecs/${var.app_name}-management-task"
  retention_in_days = 30

  tags = {
    Name = "${var.app_name}-lg-ecs-management"
  }
}

resource "aws_ecs_service" "management" {
  name                   = "${var.app_name}-ecs-management-service"
  cluster                = aws_ecs_cluster.management.id
  task_definition        = aws_ecs_task_definition.management.arn
  launch_type            = "FARGATE"
  scheduling_strategy    = "REPLICA"
  desired_count          = 1
  enable_execute_command = true

  deployment_controller {
    type = "ECS"
  }

  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100

  enable_ecs_managed_tags = true

  network_configuration {
    subnets          = [for subnet in aws_subnet.management : subnet.id]
    security_groups  = [aws_security_group.management.id]
    assign_public_ip = false
  }

  lifecycle {
    ignore_changes = [desired_count]
  }

  tags = {
    Name = "${var.app_name}-ecs-management-service"
  }
}

# =========================================
# Code Deploy
# =========================================
resource "aws_codedeploy_app" "application" {
  compute_platform = "ECS"
  name             = "${var.app_name}-deploy-app"
}

resource "aws_codedeploy_deployment_group" "application" {
  app_name               = aws_codedeploy_app.application.name
  deployment_config_name = "CodeDeployDefault.ECSAllAtOnce"
  deployment_group_name  = "${var.app_name}-deployment"
  service_role_arn       = aws_iam_role.ecs_code_deploy_role.arn

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE"]
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout = "CONTINUE_DEPLOYMENT"
    }

    terminate_blue_instances_on_deployment_success {
      action                           = "TERMINATE"
      termination_wait_time_in_minutes = 5
    }
  }

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  ecs_service {
    cluster_name = aws_ecs_cluster.application.name
    service_name = aws_ecs_service.application.name
  }

  load_balancer_info {
    target_group_pair_info {
      prod_traffic_route {
        listener_arns = [aws_lb_listener.application["blue"].arn]
      }

      target_group {
        name = aws_lb_target_group.application["blue"].name
      }

      target_group {
        name = aws_lb_target_group.application["green"].name
      }
    }
  }

  tags = {
    Name = "${var.app_name}-deployment"
  }
}

# =========================================
# Autoscaling
# =========================================
resource "aws_appautoscaling_target" "ecs_target" {
  max_capacity       = 4
  min_capacity       = 1
  resource_id        = "service/${aws_ecs_cluster.application.name}/${aws_ecs_service.application.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "ecs_scale_up" {
  name               = "${var.app_name}-ecs-scale-up"
  policy_type        = "StepScaling"
  resource_id        = aws_appautoscaling_target.ecs_target.resource_id
  scalable_dimension = aws_appautoscaling_target.ecs_target.scalable_dimension
  service_namespace  = aws_appautoscaling_target.ecs_target.service_namespace

  step_scaling_policy_configuration {
    adjustment_type         = "ChangeInCapacity"
    cooldown                = 60
    metric_aggregation_type = "Average"

    step_adjustment {
      metric_interval_lower_bound = 0
      scaling_adjustment          = 1
    }
  }
}

resource "aws_appautoscaling_policy" "ecs_scale_down" {
  name               = "${var.app_name}-ecs-scale-down"
  policy_type        = "StepScaling"
  resource_id        = aws_appautoscaling_target.ecs_target.resource_id
  scalable_dimension = aws_appautoscaling_target.ecs_target.scalable_dimension
  service_namespace  = aws_appautoscaling_target.ecs_target.service_namespace

  step_scaling_policy_configuration {
    adjustment_type         = "ChangeInCapacity"
    cooldown                = 60
    metric_aggregation_type = "Average"

    step_adjustment {
      metric_interval_upper_bound = 0
      scaling_adjustment          = -1
    }
  }
}


resource "aws_cloudwatch_metric_alarm" "ecs_cpu_high" {
  alarm_name          = "${var.app_name}-ecs-cpu-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "80"

  dimensions = {
    ClusterName = aws_ecs_cluster.application.name
    ServiceName = aws_ecs_service.application.name
  }

  alarm_actions = [aws_appautoscaling_policy.ecs_scale_up.arn]

  tags = {
    Name = "${var.app_name}-ecs-cpu-high"
  }
}

resource "aws_cloudwatch_metric_alarm" "ecs_cpu_low" {
  alarm_name          = "${var.app_name}-ecs-cpu-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Average"
  threshold           = "30"

  dimensions = {
    ClusterName = aws_ecs_cluster.application.name
    ServiceName = aws_ecs_service.application.name
  }

  alarm_actions = [aws_appautoscaling_policy.ecs_scale_down.arn]

  tags = {
    Name = "${var.app_name}-ecs-cpu-low"
  }
}

# =========================================
# IAM for ECS
# =========================================
data "aws_iam_policy" "task_execution" {
  name = "AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "task_execution" {
  name = "${var.app_name}-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      },
    ]
  })

  managed_policy_arns = [data.aws_iam_policy.task_execution.arn]
}

resource "aws_iam_role" "task" {
  name = "${var.app_name}-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      },
    ]
  })

  inline_policy {
    name = "ecs_exec_policy"

    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action = [
            "ssmmessages:CreateDataChannel",
            "ssmmessages:OpenDataChannel",
            "ssmmessages:OpenControlChannel",
            "ssmmessages:CreateControlChannel"
          ]
          Effect   = "Allow"
          Resource = "*"
        },
      ]
    })
  }
}

data "aws_iam_policy" "ecs_code_deploy_policy" {
  name = "AWSCodeDeployRoleForECS"
}

resource "aws_iam_role" "ecs_code_deploy_role" {
  name = "${var.app_name}-ecs-code-deploy-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "codedeploy.amazonaws.com"
        }
      },
    ]
  })

  managed_policy_arns = [data.aws_iam_policy.ecs_code_deploy_policy.arn]
}
