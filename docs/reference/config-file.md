# Config File Reference

`.tfparams.yml` configures tfparams. CLI flags always override file values.

## Search order

1. `--config <path>` (explicit)
2. `./.tfparams.yml` (current directory)

Built-in defaults apply when no file is found.

## Keys

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `format` | string | `markdown` | Output format: `markdown` / `csv` / `json` |
| `env` | string | - | Environment name in the header |
| `scope` | string | `root` | `root` / `module` |
| `module` | string | `""` | Module call name when `scope: module` (empty = auto) |
| `output.file` | string | - | Output file path (overwritten if it exists) |
| `columns.show` | list | all | Columns to render, in order |
| `sort.by` | string | `required` | `required` (required first, then name) / `name` |
| `sensitive.show` | bool | `false` | Show sensitive values unmasked |
| `recursive.enabled` | bool | `false` | Recursive mode |
| `recursive.path` | string | `.` | Scan root |
| `recursive.plan_file` | string | `tfplan.json` | Plan JSON filename per subdirectory |

## Example

```yaml
format: markdown
env: production
scope: module
module: app
output:
  file: PARAMETERS.md
sort:
  by: required
sensitive:
  show: false
```
