terraform {
  required_version = ">= 1.5"
}

variable "instance_type" {
  type        = string
  description = "EC2 instance type for this environment"
}

variable "replica_count" {
  type        = number
  description = "RDS replica count for this environment"
}

variable "multi_az" {
  type        = bool
  description = "Enable Multi-AZ for this environment"
}

variable "db_password" {
  type        = string
  description = "Database master password"
  sensitive   = true
}

module "app" {
  source = "../../modules/app"

  instance_type = var.instance_type
  replica_count = var.replica_count
  multi_az      = var.multi_az
  db_password   = var.db_password
  region        = "ap-northeast-1"          # constant expression
  name          = "${var.instance_type}-app" # computed expression
}

output "app_summary" {
  value = module.app.summary
}
