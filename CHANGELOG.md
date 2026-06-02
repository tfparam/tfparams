# Changelog

## [0.1.0] - 2026-06-02


### 🐛 Bug Fixes

- Mask sensitive from plan config; parse terraform-docs default as raw value ([74c0e3f](https://github.com/tfparam/tfparams/commit/74c0e3f15e85ccb5cfaad820ba44d8e6a052dc76))
- Use patched Go 1.25.10 toolchain and bump x/net to v0.55.0 ([3e5cc75](https://github.com/tfparam/tfparams/commit/3e5cc75b09af889a2d237d707e3ad825fd2bc784))

### 📚 Documentation

- VitePress documentation site (Track E) ([1616d1a](https://github.com/tfparam/tfparams/commit/1616d1ac2f7d30d74c9cf1105a7b0d8e4bf729fc))
- Comprehensive README (Track F) ([0c16aca](https://github.com/tfparam/tfparams/commit/0c16aca211e0e68dd346462e80becd05cc815a4a))

### 🔧 Maintenance

- Remove accidentally committed terraform artifacts; ignore tf working files ([336bb29](https://github.com/tfparam/tfparams/commit/336bb2933624b4b0b2ef2a2d28b5029b3f98a3e2))
- Tooling, CI, Docker, goreleaser (Track D) ([43e5fe8](https://github.com/tfparam/tfparams/commit/43e5fe892edf0613d53838d1993def50eabfb802))
- Remove accidentally committed docs build artifacts; ignore them ([afa2d98](https://github.com/tfparam/tfparams/commit/afa2d98d85c04c75d19c70a921b21e836a057fa3))
- Bump GitHub Actions to Node24-based majors ([b52ebff](https://github.com/tfparam/tfparams/commit/b52ebff8da9016f3e8b54ced3ab52474f69094fa))

### 🔨 Refactoring

- Address PR #1 review feedback (#9) ([fb74777](https://github.com/tfparam/tfparams/commit/fb74777a2e106c518884f2e88ba3cf2d5eeee331))
- Address PR #2 review (keep plan JSON; cleanups + markdown lib) ([add327c](https://github.com/tfparam/tfparams/commit/add327c66c99c7a08f62d4345ff684dc1855adc8))

### 🚀 Features

- Implement tfparams core (plan/docs parser, merger, formatter, CLI) ([df09810](https://github.com/tfparam/tfparams/commit/df09810584a57c7bece5e633ff14e698818fb56d))
- Add CSV and JSON formatters ([a37e94a](https://github.com/tfparam/tfparams/commit/a37e94a911428265d55e0248815059dee3def45f))
- Recursive mode ([4acf2ec](https://github.com/tfparam/tfparams/commit/4acf2ecbfdaf1522db1587be83d5a18edb11764d))
- Compare subcommand + backend URI dispatch (local) ([a9a5488](https://github.com/tfparam/tfparams/commit/a9a54882356b9c80ce3e8f13c850da54a8e086b5))
- Implement cloud backends (S3/GCS/Azure) for compare (Track C-3) ([f2b79aa](https://github.com/tfparam/tfparams/commit/f2b79aa22d3558f0f01c52b51b3dcd28c02c84d7))
- Tfparams core (plan/docs parser, merger, formatter, CLI) (#1) ([aa03613](https://github.com/tfparam/tfparams/commit/aa036136947a7d0762851c8d374f64ea3ba3e882))

