terraform {
  required_version = "~> 1.1.2"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.70.0"
    }
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = { Environment = var.app_name }
  }
}

resource "aws_instance" "test" {
  instance_type = "t2.micro"
  ami           = "ami-0218d08a1f9dac831"
}
