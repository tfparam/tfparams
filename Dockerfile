# syntax=docker/dockerfile:1

# ---- Build Stage ----
# --platform=$BUILDPLATFORM builds on the host arch and cross-compiles with Go.
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

# alpine lacks git; needed for go mod download / VCS-based deps.
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# buildx sets these per target platform.
ARG TARGETOS
ARG TARGETARCH
# version is computed on the host and passed in (.dockerignore excludes .git).
ARG VERSION=dev

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -trimpath \
    -o /tfparams \
    ./main.go

# ---- Final Stage ----
FROM gcr.io/distroless/static-debian12:nonroot

USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /tfparams /tfparams

ENTRYPOINT ["/tfparams"]
