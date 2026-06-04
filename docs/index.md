---
layout: home

hero:
  name: tfparams
  text: Terraform Parameter Sheet Generator
  tagline: Merge Terraform plan values with variable metadata and output as Markdown.
  image:
    light: /logo-horizontal.svg
    dark: /logo-horizontal-dark.svg
    alt: tfparams
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/tfkit/tfparams

features:
  - title: Plan-aware
    details: Reads actual applied values from a Terraform plan (terraform show -json), including module-level passed values.
  - title: Environment Comparison
    details: Compare dev / staging / production side-by-side with diff highlighting.
  - title: Multi-backend
    details: Fetch each environment's plan JSON from S3, GCS, Azure Blob, or local files.
  - title: CI/CD Ready
    details: GitHub Actions, pre-commit hooks, and Docker image out of the box.
---

## What you get

tfparams writes a single Markdown file (`PARAMETERS.md` by default) with the
**resolved** values from your plan:

```bash
terraform plan -out=tfplan
tfparams --plan-json <(terraform show -json tfplan) \
  --docs-json <(terraform-docs json .) --env production
```

The complete file looks like this:

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

…and renders as:

| Name | Description | Type | Default | Applied Value | Required |
| --- | --- | --- | --- | --- | --- |
| db_password | Database password | `string` | - | (sensitive) | ✓ |
| instance_type | EC2 instance type | `string` | `t3.medium` | `t3.xlarge` | - |
| multi_az | Enable Multi-AZ | `bool` | `false` | `true` | - |
| replica_count | RDS replica count | `number` | `1` | `3` | - |

Use `--format csv` or `--format json` for machine-readable output, and
`tfparams compare` to diff environments side-by-side. See
[Getting Started](/guide/getting-started) and the [CLI reference](/reference/cli).
