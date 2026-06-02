terraform {
  required_version = ">= 1.5"
}

# No providers/resources so `terraform init` works offline; the plan still
# carries input variables and module-call expressions, which is what tfparams reads.
locals {
  summary = "${var.name} (${var.instance_type} x${var.replica_count}) in ${var.region}, multi_az=${var.multi_az}"
}
