terraform {
  required_version = "~> 1.1.2"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.71.0"
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



# =========================================
# vpc
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
# igw
# =========================================
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.app_name}-igw"
  }
}

# =========================================
# sg
# =========================================
resource "aws_security_group" "ingress" {
  name        = "${var.app_name}-sg-ingress"
  description = "Allow inbound traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    description      = "HTTP access"
    from_port        = 80
    to_port          = 80
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
# subnet
# =========================================
locals {
  vpc_cidr           = "10.0.0.0/16"
  availability_zones = ["ap-northeast-1a", "ap-northeast-1c"]
  ingress_cidrs      = ["10.0.0.0/24", "10.0.1.0/24"]
  private_cidrs      = ["10.0.10.0/24", "10.0.11.0/24"]
  egress_cidrs       = ["10.0.20.0/24", "10.0.21.0/24"]
}

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
# route table
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

  # route {
  #   cidr_block = local.vpc_cidr
  #   gateway_id = aws_internet_gateway.igw.id
  #   vpc_endpoint_id
  # }

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
# vpc endpoint
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

resource "aws_lb_listener" "application" {
  for_each          = aws_lb_target_group.application
  load_balancer_arn = aws_lb.application.arn
  port              = each.value.port
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = each.value.arn
  }
}
