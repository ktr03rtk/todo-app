variable "alb_access_header_name" {}
variable "alb_access_header_value" {}
variable "app_name" {}
variable "db_username" { sensitive = true }
variable "db_password" { sensitive = true }
variable "db_name" {}
variable "global_certificate_arn" {}
variable "host_zone_name" {}
variable "image_arn" {}
variable "management_image_arn" {}
variable "local_certificate_arn" {}
variable "region" {}
variable "sub_domain_name" {}
variable "time_zone" {}
