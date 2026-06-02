# Configuration

tfparams reads `./.tfparams.yml` (if present) and lets CLI flags override it.
Pass `--config <path>` to load a file from any other location.

## Full schema

```yaml
format: markdown         # markdown / csv / json
env: production
scope: root              # root / module
module: ""               # module call name when scope: module (empty = auto-select)
output:
  file: PARAMETERS.md       # overwritten if it exists
columns:
  show:                   # column order; --no-default-col drops `default`
    - name
    - description
    - type
    - default
    - applied_value
    - required
sort:
  by: required            # required (required first, then name) / name
sensitive:
  show: false
recursive:
  enabled: false
  path: .                 # scan root (env dir by default)
  plan_file: tfplan.json  # plan JSON filename to look for in each subdirectory
```

See the [Config File reference](../reference/config-file) for details.
