terraform {
  required_version = "~> 1.1.2"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.72.0"
    }
  }

  # backend = "s3"
  # config = {
  #   bucket         = var.tf_backend
  #   key            = "terraform.tfstate"
  #   region         = var.region
  #   encrypt        = true
  #   dynamodb_table = "terraform-state-lock-dynamo"
  # }
}


provider "aws" {
  region = var.region

  default_tags {
    tags = { Environment = var.app_name }
  }
}


locals {
  vpc_cidr           = "10.0.0.0/16"
  availability_zones = ["ap-northeast-1a", "ap-northeast-1c"]
  ingress_cidrs      = ["10.0.0.0/24", "10.0.1.0/24"]
  private_cidrs      = ["10.0.10.0/24", "10.0.11.0/24"]
  egress_cidrs       = ["10.0.20.0/24", "10.0.21.0/24"]
}

# =========================================
# VPC
# =========================================
resource "aws_vpc" "main" {
  cidr_block = local.vpc_cidr

  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "${var.app_name}-vpc"
  }
}

# =========================================
# Internet Gateway
# =========================================
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.app_name}-igw"
  }
}

# =========================================
# Security Group
# =========================================
resource "aws_security_group" "ingress" {
  name        = "${var.app_name}-sg-ingress"
  description = "Allow inbound traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    description      = "HTTPS access"
    from_port        = 443
    to_port          = 443
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "${var.app_name}-sg-ingress"
  }
}

resource "aws_security_group" "private" {
  name        = "${var.app_name}-sg-private"
  description = "access to ecs task"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "access to ecs task"
    from_port       = 0
    to_port         = 65535
    protocol        = "tcp"
    security_groups = [aws_security_group.ingress.id]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "${var.app_name}-sg-private"
  }
}

resource "aws_security_group" "egress" {
  name        = "${var.app_name}-sg-egress"
  description = "for vpc endpoint"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 443
    to_port         = 443
    protocol        = "tcp"
    security_groups = [aws_security_group.private.id]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "${var.app_name}-sg-egress"
  }
}

# =========================================
# Subnet
# =========================================
resource "aws_subnet" "ingress" {
  for_each = { for idx, az in local.availability_zones : idx => az }

  vpc_id                  = aws_vpc.main.id
  cidr_block              = local.ingress_cidrs[each.key]
  availability_zone       = each.value
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.app_name}-sbn-ingress"
  }
}

resource "aws_subnet" "private" {
  for_each = { for idx, az in local.availability_zones : idx => az }

  vpc_id                  = aws_vpc.main.id
  cidr_block              = local.private_cidrs[each.key]
  availability_zone       = each.value
  map_public_ip_on_launch = false

  tags = {
    Name = "${var.app_name}-sbn-private"
  }
}

resource "aws_subnet" "egress" {
  for_each = { for idx, az in local.availability_zones : idx => az }

  vpc_id                  = aws_vpc.main.id
  cidr_block              = local.egress_cidrs[each.key]
  availability_zone       = each.value
  map_public_ip_on_launch = false

  tags = {
    Name = "${var.app_name}-sbn-egress"
  }
}

# =========================================
# Route Table
# =========================================
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "${var.app_name}-rt-public"
  }
}

resource "aws_route_table_association" "ingress" {
  count          = length(local.ingress_cidrs)
  subnet_id      = aws_subnet.ingress[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.app_name}-rt-private"
  }
}

resource "aws_route_table_association" "private" {
  count          = length(local.private_cidrs)
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private.id
}

resource "aws_route_table_association" "egress" {
  count          = length(local.egress_cidrs)
  subnet_id      = aws_subnet.egress[count.index].id
  route_table_id = aws_route_table.private.id
}

# =========================================
# VPC Endpoint
# =========================================
locals {
  interface_endpoint_map = {
    "ecr-api" : "com.amazonaws.ap-northeast-1.ecr.api",
    "ecr-dkr" : "com.amazonaws.ap-northeast-1.ecr.dkr",
    "logs" : "com.amazonaws.ap-northeast-1.logs"
  }
}

resource "aws_vpc_endpoint" "interface" {
  for_each          = local.interface_endpoint_map
  vpc_id            = aws_vpc.main.id
  service_name      = each.value
  vpc_endpoint_type = "Interface"

  security_group_ids = [
    aws_security_group.egress.id,
  ]

  subnet_ids          = [for subnet in aws_subnet.egress : subnet.id]
  private_dns_enabled = true

  tags = {
    Name = "${var.app_name}-vpce-${each.key}"
  }
}

resource "aws_vpc_endpoint" "gateway" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.ap-northeast-1.s3"
  vpc_endpoint_type = "Gateway"

  route_table_ids = [aws_route_table.private.id]

  tags = {
    Name = "${var.app_name}-vpce-s3"
  }
}

# =========================================
# ALB
# =========================================
locals {
  target_group_map = {
    "blue" : 80,
    "green" : 10080
  }
}

resource "aws_lb_target_group" "application" {
  for_each = local.target_group_map
  name     = "${var.app_name}-alb-tg-${each.key}"

  port             = each.value
  protocol         = "HTTP"
  target_type      = "ip"
  vpc_id           = aws_vpc.main.id
  protocol_version = "HTTP1"

  health_check {
    protocol            = "HTTP"
    path                = "/"
    port                = "traffic-port"
    enabled             = true
    healthy_threshold   = 5
    unhealthy_threshold = 2
    timeout             = 5
    interval            = 30
    matcher             = "200"
  }

  tags = {
    Name = "${var.app_name}-alb-tg-${each.key}"
  }
}

resource "aws_lb" "application" {
  name               = "${var.app_name}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.ingress.id]
  subnets            = [for subnet in aws_subnet.ingress : subnet.id]
  ip_address_type    = "ipv4"

  # access_logs {
  #   bucket  = "s3bucket"
  #   prefix  = "prefix"
  #   enabled = true
  # }

  tags = {
    Name = "${var.app_name}-alb"
  }
}

locals {
  public_port_map = {
    "blue" : 443,
    "green" : 10443
  }
}

resource "aws_lb_listener" "application" {
  for_each          = aws_lb_target_group.application
  load_balancer_arn = aws_lb.application.arn
  port              = local.public_port_map[each.key]
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = var.local_certificate_arn

  default_action {

    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      status_code  = "403"
    }
  }
}

resource "aws_lb_listener_rule" "application" {
  for_each     = aws_lb_target_group.application
  listener_arn = aws_lb_listener.application[each.key].arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = each.value.arn
  }

  condition {
    host_header {
      values = [aws_route53_record.alb.name]
    }
  }

  condition {
    http_header {
      http_header_name = var.alb_access_header_name
      values           = [var.alb_access_header_value]
    }
  }
  lifecycle {
    ignore_changes = [action]
  }
}

# =========================================
# Route53
# =========================================
data "aws_route53_zone" "host" {
  name = var.host_zone_name
}

resource "aws_route53_record" "alb" {
  zone_id = data.aws_route53_zone.host.zone_id
  name    = "alb.${var.sub_domain_name}"
  type    = "A"

  alias {
    name                   = aws_lb.application.dns_name
    zone_id                = aws_lb.application.zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "cloud_front" {
  zone_id = data.aws_route53_zone.host.zone_id
  name    = var.sub_domain_name
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.application.domain_name
    zone_id                = aws_cloudfront_distribution.application.hosted_zone_id
    evaluate_target_health = true
  }
}

# =========================================
# CloudFront
# =========================================
data "aws_cloudfront_cache_policy" "managed_caching_disabled" {
  name = "Managed-CachingDisabled"
}

resource "aws_cloudfront_distribution" "application" {
  aliases         = [var.sub_domain_name]
  enabled         = true
  is_ipv6_enabled = true
  web_acl_id      = aws_wafv2_web_acl.application.arn

  origin {
    connection_attempts = 3
    connection_timeout  = 10
    domain_name         = aws_route53_record.alb.name
    origin_id           = aws_lb.application.dns_name

    custom_header {
      name  = var.alb_access_header_name
      value = var.alb_access_header_value
    }

    custom_origin_config {
      http_port                = 80
      https_port               = 443
      origin_keepalive_timeout = 5
      origin_protocol_policy   = "https-only"
      origin_read_timeout      = 30
      origin_ssl_protocols = [
        "TLSv1",
        "TLSv1.1",
        "TLSv1.2",
      ]
    }
  }

  default_cache_behavior {
    target_origin_id       = aws_lb.application.dns_name
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    viewer_protocol_policy = "redirect-to-https"
    cache_policy_id        = data.aws_cloudfront_cache_policy.managed_caching_disabled.id
    compress               = true
    smooth_streaming       = false
    default_ttl            = 0
    min_ttl                = 0
    max_ttl                = 0
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn = var.global_certificate_arn
    ssl_support_method  = "sni-only"
  }

  tags = {
    Name = "${var.app_name}-cloudfront"
  }
}

# =========================================
# ECS
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
  family                   = "${var.app_name}task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 1024
  memory                   = 2048

  # task_role_arn      = ""
  execution_role_arn = aws_iam_role.task_execution.arn

  container_definitions = format("[%s]", templatefile(
    "${path.module}/container_definitions.json",
    {
      container_name = "${local.container_name}-test2"
      region         = var.region
      image_arn      = var.image_arn
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
  name              = "/ecs/${var.app_name}task"
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
# IAM
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

# =========================================
# WAF
# =========================================
provider "aws" {
  region = "us-east-1"
  alias  = "global"

  default_tags {
    tags = { Environment = var.app_name }
  }
}

resource "aws_wafv2_web_acl" "application" {
  provider = aws.global
  name     = "${var.app_name}-waf"
  scope    = "CLOUDFRONT"

  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name                = "${var.app_name}-waf"
    sampled_requests_enabled   = true
  }

  default_action {
    allow {}
  }

  rule {
    name     = "AWS-AWSManagedRulesCommonRuleSet"
    priority = 0

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesCommonRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "AWS-AWSManagedRulesCommonRuleSet"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "AWS-AWSManagedRulesKnownBadInputsRuleSet"
    priority = 1

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesKnownBadInputsRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "AWS-AWSManagedRulesKnownBadInputsRuleSet"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "AWS-AWSManagedRulesAnonymousIpList"
    priority = 2

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesAnonymousIpList"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "AWS-AWSManagedRulesAnonymousIpList"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "AWS-AWSManagedRulesSQLiRuleSet"
    priority = 3

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesSQLiRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "AWS-AWSManagedRulesSQLiRuleSet"
      sampled_requests_enabled   = true
    }
  }

  tags = {
    Name = "${var.app_name}-waf"
  }
}
