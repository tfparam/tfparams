# Getting Started

tfparams merges the **applied values from a Terraform plan** with **variable metadata
from terraform-docs** and renders a Markdown parameter sheet.

::: warning Data source
Input variable values live only in a **plan** file. `terraform.tfstate` (and a bare
`terraform show -json`) does not carry them, and `*.tfvars` alone is incomplete (it
omits `TF_VAR_*`, `-var`, defaults, and computed values). Always feed a plan:
`terraform plan -out=tfplan && terraform show -json tfplan`.
:::

## Install

```bash
go install github.com/tfparam/tfparams@latest
```

See [Installation](./installation) for Homebrew, binaries, and Docker.

## Basic usage

```bash
cd environments/production/
terraform plan -out=tfplan
terraform show -json tfplan | tfparams --docs-json <(terraform-docs json .)
```

## Output sample

```markdown
# Parameter Sheet

**Environment**: production
**Scope**: root
**Source**: terraform show -json tfplan (plan)

## Variables

| Name | Description | Type | Default | Applied Value | Required |
| --- | --- | --- | --- | --- | --- |
| instance_type | EC2 instance type | `string` | `t3.medium` | `t3.xlarge` | - |
| replica_count | RDS replica count | `number` | `1` | `3` | - |
| db_password | Database password | `string` | - | (sensitive) | ✓ |
```

Rendered, that table looks like this:

| Name | Description | Type | Default | Applied Value | Required |
| --- | --- | --- | --- | --- | --- |
| instance_type | EC2 instance type | `string` | `t3.medium` | `t3.xlarge` | - |
| replica_count | RDS replica count | `number` | `1` | `3` | - |
| db_password | Database password | `string` | - | (sensitive) | ✓ |

## Module-level view

For shared-module setups, show the values an environment passes **into the module**:

```bash
tfparams --plan-json plan.json --scope module --module app \
  --docs-json <(terraform-docs json ../../modules/app/)
```

## Compare environments

```bash
tfparams compare \
  --docs-json <(terraform-docs json ./modules/app/) \
  --env dev=./dev.plan.json \
  --env stg=./stg.plan.json \
  --env prd=./prd.plan.json
```

```markdown
| Name | Description | dev | stg | prd | Diff |
| --- | --- | --- | --- | --- | --- |
| instance_type | EC2 instance type | `t3.small` | `t3.medium` | `t3.xlarge` | ⚠️ |
| replica_count | RDS replica count | `1` | `2` | `3` | ⚠️ |
| db_password | Database password | (sensitive) | (sensitive) | (sensitive) | - |
```

Rendered:

| Name | Description | dev | stg | prd | Diff |
| --- | --- | --- | --- | --- | --- |
| instance_type | EC2 instance type | `t3.small` | `t3.medium` | `t3.xlarge` | ⚠️ |
| replica_count | RDS replica count | `1` | `2` | `3` | ⚠️ |
| db_password | Database password | (sensitive) | (sensitive) | (sensitive) | - |
