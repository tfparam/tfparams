# Installation

## go install

```bash
go install github.com/tfkit/tfparams@latest
```

## Homebrew (tap)

```bash
brew install tfkit/tap/tfparams
```

## Pre-built binaries

Download the archive for your OS/arch from the
[GitHub Releases](https://github.com/tfkit/tfparams/releases) page and extract `tfparams`.

## Docker

```bash
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/tfkit/tfparams:latest \
  --plan-json plan.json --docs-json docs.json
```

Images are published to `ghcr.io/tfkit/tfparams` for `linux/amd64` and `linux/arm64`.
