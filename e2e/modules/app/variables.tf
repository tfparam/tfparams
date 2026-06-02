variable "instance_type" {
  type        = string
  description = "EC2 instance type"
  default     = "t3.medium"
}

variable "replica_count" {
  type        = number
  description = "RDS replica count"
  default     = 1
}

variable "multi_az" {
  type        = bool
  description = "Enable Multi-AZ deployment"
  default     = false
}

variable "db_password" {
  type        = string
  description = "Database master password"
  sensitive   = true
}

variable "region" {
  type        = string
  description = "Cloud region"
  default     = "us-east-1"
}

variable "name" {
  type        = string
  description = "Application name"
  default     = "app"
}
