# tfparams

Generate Markdown / CSV / JSON **parameter sheets** for Terraform by merging the
**applied input-variable values from a plan** with **variable metadata from
terraform-docs**.

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
go install github.com/tfparam/tfparams@latest

# Homebrew tap
brew install tfparam/tap/tfparams

# Docker
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/tfparam/tfparams:latest --plan-json plan.json --docs-json docs.json
```

Pre-built binaries are on the [Releases](https://github.com/tfparam/tfparams/releases) page.

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
| instance_type | EC2 instance type | `string` | `t3.medium` | `t3.xlarge` | - |
| replica_count | RDS replica count | `number` | `1` | `3` | - |
| db_password | Database password | `string` | - | (sensitive) | ✓ |
```

### Module-level view

Show the values an environment passes **into a shared module**:

```bash
tfparams --plan-json plan.json --scope module --module app \
  --docs-json <(terraform-docs json ../../modules/app/)
```

### Output formats

`--format table` (default, Markdown) · `--format csv` · `--format json`.

### Write to a file

```bash
tfparams --plan-json plan.json --docs-json docs.json --out PARAMETERS.md
```

`--out` overwrites the file; without it, output goes to stdout.

## Configuration (`.tfparams.yml`)

Searched in order (first match wins) unless `--config` is given:
`./.tfparams.yml` → `./.config/.tfparams.yml` → `$HOME/.tfparams.d/.tfparams.yml`.
CLI flags override the file.

```yaml
format: table            # table / csv / json
env: production
scope: root              # root / module
module: ""               # module call name when scope: module (empty = auto)
output:
  file: PARAMETERS.md     # overwritten if it exists
columns:
  show: [name, description, type, default, applied_value, required]
sort:
  by: name                # name / required / type
sensitive:
  show: false
```

## License

[MIT](LICENSE)
