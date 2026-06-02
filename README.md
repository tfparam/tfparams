<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/public/logo-horizontal-dark.svg">
    <img src="docs/public/logo-horizontal.svg" alt="tfparams" width="320">
  </picture>
</p>

<p align="center">
  Generate Markdown parameter sheets from Terraform <strong>plan</strong> values and variable metadata.
</p>

---

tfparams merges the **applied values from a Terraform plan** (`terraform show -json <planfile>`)
with **variable metadata from terraform-docs** and renders a Markdown parameter sheet —
for a single environment, for a shared module, or compared side by side across environments.

> [!IMPORTANT]
> Input variable values live only in a **plan** file. `terraform.tfstate` (and a bare
> `terraform show -json`) does not carry them. Always feed a plan:
> `terraform plan -out=tfplan && terraform show -json tfplan`.

📖 **Documentation:** <https://tfparam.github.io/tfparams>

## Installation

```bash
# go install
go install github.com/tfparam/tfparams@latest

# Homebrew tap
brew install tfparam/tap/tfparams

# Docker
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/tfparam/tfparams:latest --plan-json plan.json --docs-json docs.json
```

Pre-built binaries are on the [Releases](https://github.com/tfparam/tfparams/releases) page.

## Usage

### Basic

```bash
cd environments/production/
terraform plan -out=tfplan
terraform show -json tfplan | tfparams --docs-json <(terraform-docs json .)
```

```markdown
# Parameter Sheet

**Environment**: production
**Scope**: root
**Source**: terraform show -json tfplan (plan)

## Variables

| Name | Description | Type | Default | Applied Value | Required |
|------|-------------|------|---------|---------------|----------|
| instance_type | EC2 instance type | `string` | `t3.medium` | `t3.xlarge` | - |
| replica_count | RDS replica count | `number` | `1` | `3` | - |
| db_password | Database password | `string` | - | `(sensitive)` | ✓ |
```

### Module-level view

Show the values an environment passes **into a shared module**:

```bash
tfparams --plan-json plan.json --scope module --module app \
  --docs-json <(terraform-docs json ../../modules/app/)
```

### Write to a file

```bash
tfparams --plan-json plan.json --docs-json docs.json --out PARAMETERS.md
```

`--out` overwrites the file; without it, output goes to stdout.

### Compare environments

```bash
tfparams compare \
  --docs-json <(terraform-docs json ./modules/app/) \
  --env dev=s3://my-bucket/env/dev/plan.json \
  --env stg=s3://my-bucket/env/stg/plan.json \
  --env prd=s3://my-bucket/env/prd/plan.json
```

```markdown
# Environment Comparison

| Name | Description | dev | stg | prd | Diff |
|------|-------------|-----|-----|-----|------|
| instance_type | EC2 instance type | `t3.small` | `t3.medium` | `t3.xlarge` | ⚠️ |
| replica_count | RDS replica count | `1` | `2` | `3` | ⚠️ |
| db_password | Database password | `(sensitive)` | `(sensitive)` | `(sensitive)` | - |
```

`--env` values are **plan JSON** files; the scheme (`s3://` / `gs://` / `azblob://` / local)
only selects how the bytes are fetched, using each cloud SDK's default credentials.

### CI/CD

```yaml
# GitHub Actions (composite action)
- uses: tfparam/tfparams@v0.1.0
  with:
    plan-json: plan.json
    docs-json: docs.json
    output-file: PARAMETERS.md
```

```yaml
# pre-commit
repos:
  - repo: https://github.com/tfparam/tfparams
    rev: "v0.1.0"
    hooks:
      - id: tfparams
        args: ["--out", "PARAMETERS.md"]
```

## Output formats

`--format table` (default, Markdown) · `--format csv` · `--format json`.

## Configuration (`.tfparams.yml`)

Searched in order (first match wins) unless `--config` is given:
`./.tfparams.yml` → `./.config/.tfparams.yml` → `$HOME/.tfparams.d/.tfparams.yml`. CLI flags override the file.

```yaml
format: table            # table / csv / json
env: production
scope: root              # root / module
module: ""               # module call name when scope: module (empty = auto)
output:
  file: PARAMETERS.md       # overwritten if it exists
columns:
  show: [name, description, type, default, applied_value, required]
sort:
  by: required              # required (required first, then name) / name
sensitive:
  show: false
recursive:
  enabled: false
  path: .                 # scan root (env dir by default)
  plan_file: tfplan.json  # plan JSON filename per subdirectory
```

See the [configuration reference](https://tfparam.github.io/tfparams/reference/config-file) for every key.

## License

[MIT](LICENSE)
