# Parameter Sheet

**Environment**: production
**Scope**: root
**Generated at**: 2025-01-15 10:30:00 JST
**Source**: terraform show -json tfplan (plan)

## Variables

| Name | Description | Type | Default | Applied Value | Required |
| --- | --- | --- | --- | --- | --- |
| instance_type | EC2インスタンスタイプ | `string` | `t3.medium` | `t3.xlarge` | - |
| replica_count | RDSレプリカ数 | `number` | `1` | `3` | - |
| db_password | DBパスワード | `string` | - | (sensitive) | ✓ |
| multi_az | Multi-AZ有効化 | `bool` | `false` | `true` | - |
| extra | (no description) | - | - | (computed) | - |
