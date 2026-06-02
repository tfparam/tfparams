# CLI Reference

## `tfparams`

Generate a parameter sheet from a plan JSON and terraform-docs metadata.

::: tip Input is a plan JSON, not `*.tfvars`
tfparams reads a **plan JSON** (`terraform show -json <planfile>`), not `*.tfvars`
files. A tfvars file alone is incomplete — it omits `TF_VAR_*` environment
variables, `-var` CLI flags, variable defaults, and computed values. The plan is
the only source carrying the **fully resolved** applied values, so produce one
first:

```bash
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
```
:::

| Flag | Default | Description |
|------|---------|-------------|
| `--plan-json` | stdin | Plan JSON file (`terraform show -json <planfile>`) |
| `--docs-json` | - | terraform-docs JSON, e.g. `<(terraform-docs json .)` (required, repeatable) |
| `--scope` | `root` | `root` (root variables) / `module` (module-passed values) |
| `--module` | (auto) | Module call name when `--scope module` |
| `--out` | stdout | Output file path (overwritten if it exists) |
| `--format` | `markdown` | `markdown` / `csv` / `json` |
| `--env` | - | Environment name shown in the header |
| `--show-sensitive` | false | Show sensitive values unmasked |
| `--no-default-col` | false | Hide the Default column |
| `--sort-by` | `required` | `required` (required first, then name) / `name` |
| `--recursive` | false | Process subdirectories recursively |
| `--recursive-path` | `.` | Root to scan in recursive mode |
| `--config` | (search) | Config file path |

## `tfparams compare`

Compare applied values across environments and highlight differences.

| Flag | Default | Description |
|------|---------|-------------|
| `--env` | - | `name=<uri-or-path>` to a plan JSON (repeat; at least two) |
| `--docs-json` | - | terraform-docs json file (required, repeatable) |
| `--scope` | `root` | `root` / `module` |
| `--module` | (auto) | Module call name when `--scope module` |
| `--out` | stdout | Output file path (overwritten if it exists) |
| `--format` | `markdown` | `markdown` |
| `--highlight-diff` | true | Highlight differing rows with ⚠️ |
| `--show-sensitive` | false | Show sensitive values unmasked |
| `--sort-by` | `required` | `required` (required first, then name) / `name` |

## `--env` URI schemes

The target is each environment's **plan JSON**. The scheme only selects how bytes
are fetched.

| Scheme | Backend | Example |
|--------|---------|---------|
| `s3://` | AWS S3 | `s3://my-bucket/env/prd/plan.json` |
| `gs://` | Google Cloud Storage | `gs://my-bucket/env/prd/plan.json` |
| `azblob://` | Azure Blob Storage | `azblob://account@container/prd/plan.json` |
| (none) | Local file | `./plan.json` |

Credentials come from each cloud SDK's default credential chain. S3 region resolves
from `AWS_REGION`/`~/.aws/config`; GCS project from Application Default Credentials;
the Azure account from the URI or `AZURE_STORAGE_ACCOUNT`.

## Examples

### Generate a sheet (Markdown, the default)

```bash
terraform plan -out=tfplan
tfparams --plan-json <(terraform show -json tfplan) \
  --docs-json <(terraform-docs json .) \
  --env production
```

tfparams writes Markdown to stdout:

```markdown
# Parameter Sheet

**Environment**: production
**Scope**: root
**Source**: terraform show -json tfplan (plan)

## Variables

| Name | Description | Type | Default | Applied Value | Required |
| --- | --- | --- | --- | --- | --- |
| db_password | Database password | `string` | - | (sensitive) | ✓ |
| instance_type | EC2 instance type | `string` | `t3.medium` | `t3.xlarge` | - |
| multi_az | Enable Multi-AZ | `bool` | `false` | `true` | - |
| replica_count | RDS replica count | `number` | `1` | `3` | - |
```

Rendered, that table looks like this:

| Name | Description | Type | Default | Applied Value | Required |
| --- | --- | --- | --- | --- | --- |
| db_password | Database password | `string` | - | (sensitive) | ✓ |
| instance_type | EC2 instance type | `string` | `t3.medium` | `t3.xlarge` | - |
| multi_az | Enable Multi-AZ | `bool` | `false` | `true` | - |
| replica_count | RDS replica count | `number` | `1` | `3` | - |

::: tip
Variables flagged `sensitive` in the terraform-docs metadata render as `(sensitive)`.
The `Required` column shows `✓` when the variable has no default. Pass
`--show-sensitive` to reveal the applied value, or `--no-default-col` to drop the
`Default` column.
:::

### CSV (`--format csv`)

```bash
tfparams --plan-json plan.json --docs-json docs.json --format csv
```

```csv
Name,Description,Type,Default,Applied Value,Required
db_password,Database password,string,,(sensitive),true
instance_type,EC2 instance type,string,t3.medium,t3.xlarge,false
multi_az,Enable Multi-AZ,bool,false,true,false
replica_count,RDS replica count,number,1,3,false
```

### JSON (`--format json`)

```bash
tfparams --plan-json plan.json --docs-json docs.json --format json
```

```json
{
  "scope": "root",
  "source": "terraform show -json tfplan (plan)",
  "variables": [
    {
      "name": "db_password",
      "description": "Database password",
      "type": "string",
      "default": null,
      "applied_value": "(sensitive)",
      "required": true,
      "sensitive": true
    },
    {
      "name": "instance_type",
      "description": "EC2 instance type",
      "type": "string",
      "default": "t3.medium",
      "applied_value": "t3.xlarge",
      "required": false,
      "sensitive": false
    }
  ]
}
```

### Module-passed values (`--scope module`)

For shared-module setups, show the arguments an environment passes **into** the
module instead of its own root variables:

```bash
tfparams --plan-json plan.json --scope module --module app \
  --docs-json <(terraform-docs json ../../modules/app/)
```

### Compare environments (`tfparams compare`)

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

Rows whose applied value differs across environments are flagged `⚠️` in the `Diff`
column (disable with `--highlight-diff=false`).
