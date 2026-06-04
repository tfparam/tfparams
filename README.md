<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/public/logo-horizontal-dark.svg">
    <img src="docs/public/logo-horizontal.svg" alt="tfparams" width="320">
  </picture>
</p>

<p align="center">
  Generate Markdown / CSV / JSON <strong>parameter sheets</strong> for Terraform by merging
  the <strong>applied input-variable values from a plan</strong> with <strong>variable metadata from terraform-docs</strong>.
</p>

📖 **Documentation:** <https://tfkit.github.io/tfparams>

---

## How it works — and a common misconception

> [!IMPORTANT]
> **`terraform.tfstate` does NOT contain your input variable values.**
> This trips a lot of people up, so it's worth being precise about where each
> value actually lives.

### What each data source really holds

| Source | What it contains | tfparams uses it for |
|--------|------------------|----------------------|
| `terraform.tfstate` | The **attributes of the resources that were created** (the deployed reality) and outputs. **It does not store input variables.** | — (not read directly) |
| `variables.tf` → `terraform-docs json` | Variable **metadata**: type, description, default, required | Name / Type / Description / Default / Required columns |
| **plan JSON** → `terraform show -json <planfile>` | The **resolved input-variable values** (after merging `*.tfvars`, `TF_VAR_*`, `-var`, and defaults), plus the prior state and any drift | **Applied Value** column |

### Why aren't variable values in the state?

Terraform re-reads variables from `*.tfvars` / `TF_VAR_*` / `-var` on **every
run**, so it never needs to persist them. The state file records the **results**
(resource attributes) — not the **inputs** (variables). Drift detection compares
those recorded resource attributes against the real infrastructure; variables
are an input to that computation, not part of it.

So the only place the question *"what value did this variable actually get?"* is
answered is a **plan**. tfparams therefore reads `terraform show -json <planfile>`:

```bash
terraform plan -out=tfplan
terraform show -json tfplan | tfparams --docs-json <(terraform-docs json .)
```

### Yes — you can do this after `apply`

`terraform plan` runs at any time, including on already-applied infrastructure.
A post-apply plan is usually a no-op (`No changes`), but its JSON **still
carries** everything tfparams needs:

| Key in plan JSON | Meaning |
|------------------|---------|
| `variables` | the current **resolved variable values** (what tfparams renders) |
| `prior_state` | the refreshed, **actually-deployed** resource values |
| `resource_changes` | **drift**, if reality has diverged from the state |

“Generate a plan after apply, then run tfparams” is a perfectly normal flow —
it's the same thing drift-detection pipelines do. If you want to skip the live
refresh, use `terraform plan -refresh=false`; for drift only, `-refresh-only`.

---

## Installation

```bash
# go install
go install github.com/tfkit/tfparams@latest

# Homebrew tap
brew install tfkit/tap/tfparams

# Docker
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/tfkit/tfparams:latest --plan-json plan.json --docs-json docs.json
```

Pre-built binaries are on the [Releases](https://github.com/tfkit/tfparams/releases) page.

## Usage

### Basic (root variables)

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
| --- | --- | --- | --- | --- | --- |
| db_password | Database password | `string` | - | (sensitive) | ✓ |
| instance_type | EC2 instance type | `string` | `t3.medium` | `t3.xlarge` | - |
| replica_count | RDS replica count | `number` | `1` | `3` | - |
```

### Module-level view

Show the values an environment passes **into a shared module**:

```bash
tfparams --plan-json plan.json --scope module --module app \
  --docs-json <(terraform-docs json ../../modules/app/)
```

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

### Output formats

`--format markdown` (default) · `--format csv` · `--format json`.

### Write to a file

```bash
tfparams --plan-json plan.json --docs-json docs.json --out PARAMETERS.md
```

`--out` overwrites the file; without it, output goes to stdout.

### CI/CD

```yaml
# GitHub Actions — see github.com/tfkit/tfparams-action
- uses: tfkit/tfparams-action@v1
  with:
    plan-json: plan.json
    docs-json: docs.json
    out: PARAMETERS.md
```

```yaml
# pre-commit
repos:
  - repo: https://github.com/tfkit/tfparams
    rev: "v0.1.0"
    hooks:
      - id: tfparams
        args: ["--out", "PARAMETERS.md"]
```

## Configuration (`.tfparams.yml`)

Read from `./.tfparams.yml` if present (built-in defaults otherwise). Pass
`--config <path>` to load any other file. CLI flags override file values.

```yaml
format: markdown         # markdown / csv / json
env: production
scope: root              # root / module
module: ""               # module call name when scope: module (empty = auto)
output:
  file: PARAMETERS.md     # overwritten if it exists
columns:
  show: [name, description, type, default, applied_value, required]
sort:
  by: required            # required (required first, then name) / name
sensitive:
  show: false
recursive:
  enabled: false
  path: .                 # scan root (env dir by default)
  plan_file: tfplan.json  # plan JSON filename per subdirectory
```

See the [configuration reference](https://tfkit.github.io/tfparams/reference/config-file) for every key.

## Use as a library

The core packages are public and can be imported on their own (the CLI in
`cmd/` is just a thin wrapper over them):

```go
import (
    "github.com/tfkit/tfparams/pkg/parser"    // parse plan JSON / terraform-docs JSON
    "github.com/tfkit/tfparams/pkg/merger"    // merge into a parameter sheet model
    "github.com/tfkit/tfparams/pkg/formatter" // render Markdown / CSV / JSON
    "github.com/tfkit/tfparams/pkg/backend"   // fetch plan JSON from s3/gcs/azblob/local
    "github.com/tfkit/tfparams/pkg/config"    // load .tfparams.yml
)
```

API docs: [pkg.go.dev/github.com/tfkit/tfparams](https://pkg.go.dev/github.com/tfkit/tfparams).

## License

[MIT](LICENSE)
