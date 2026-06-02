# CLI Reference

## `tfparams`

Generate a parameter sheet from a plan JSON and terraform-docs metadata.

| Flag | Default | Description |
|------|---------|-------------|
| `--plan-json` | stdin | Plan JSON file (`terraform show -json <planfile>`) |
| `--docs-json` | - | terraform-docs JSON, e.g. `<(terraform-docs json .)` (required, repeatable) |
| `--scope` | `root` | `root` (root variables) / `module` (module-passed values) |
| `--module` | (auto) | Module call name when `--scope module` |
| `--out` | stdout | Output file path (overwritten if it exists) |
| `--format` | `table` | `table` / `csv` / `json` |
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
| `--format` | `table` | `table` |
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
