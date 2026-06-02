# Configuration

tfparams reads `.tfparams.yml` and lets CLI flags override it. The file is searched
in this order (first match wins), unless `--config` is given:

1. `--config <path>`
2. `./.tfparams.yml`
3. `./.config/.tfparams.yml`
4. `$HOME/.tfparams.d/.tfparams.yml`

## Full schema

```yaml
format: table            # table / csv / json
env: production
scope: root              # root / module
module: ""               # module call name when scope: module (empty = auto-select)
output:
  file: PARAMETERS.md
  mode: standalone        # standalone / inject / replace
columns:
  show:                   # column order; --no-default-col drops `default`
    - name
    - description
    - type
    - default
    - applied_value
    - required
sort:
  enabled: false          # default keeps terraform-docs definition order
  by: name                # name / required / type
sensitive:
  show: false
  mask: "(sensitive)"
recursive:
  enabled: false
  path: .                 # scan root (env dir by default)
  plan_file: tfplan.json  # plan JSON filename to look for in each subdirectory
```

See the [Config File reference](../reference/config-file) for details.
