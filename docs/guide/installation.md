# Installation

## go install

```bash
go install github.com/tfparam/tfparams@latest
```

## Homebrew (tap)

```bash
brew install tfparam/tap/tfparams
```

## Pre-built binaries

Download the archive for your OS/arch from the
[GitHub Releases](https://github.com/tfparam/tfparams/releases) page and extract `tfparams`.

## Docker

```bash
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/tfparam/tfparams:latest \
  --plan-json plan.json --docs-json docs.json
```

Images are published to `ghcr.io/tfparam/tfparams` for `linux/amd64` and `linux/arm64`.
