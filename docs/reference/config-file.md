# Config File Reference

`.tfparams.yml` configures tfparams. CLI flags always override file values.

## Search order

1. `--config <path>` (explicit)
2. `./.tfparams.yml`
3. `./.config/.tfparams.yml`
4. `$HOME/.tfparams.d/.tfparams.yml`

## Keys

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `format` | string | `table` | Output format: `table` / `csv` / `json` |
| `env` | string | - | Environment name in the header |
| `scope` | string | `root` | `root` / `module` |
| `module` | string | `""` | Module call name when `scope: module` (empty = auto) |
| `output.file` | string | - | Output file path |
| `output.mode` | string | `standalone` | `standalone` / `inject` / `replace` |
| `columns.show` | list | all | Columns to render, in order |
| `sort.enabled` | bool | `false` | Whether to sort rows |
| `sort.by` | string | `name` | `name` / `required` / `type` |
| `sensitive.show` | bool | `false` | Show sensitive values unmasked |
| `sensitive.mask` | string | `(sensitive)` | Mask text |
| `recursive.enabled` | bool | `false` | Recursive mode |
| `recursive.path` | string | `.` | Scan root |
| `recursive.plan_file` | string | `tfplan.json` | Plan JSON filename per subdirectory |

## Example

```yaml
format: table
env: production
scope: module
module: app
output:
  file: PARAMETERS.md
  mode: inject
sensitive:
  show: false
```
