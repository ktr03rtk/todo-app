# =========================================
# =========================================
# todo-app infrastructure
# =========================================
# =========================================

terraform {
  required_version = "~> 1.1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.74.1"
    }
  }
}


provider "aws" {
  region = var.region

  default_tags {
    tags = { Environment = var.app_name }
  }
}
