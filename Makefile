BINARY  := tfparams
VERSION := $(shell git describe --tags --always 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build test fmt lint vuln install docker-build docker-push

build:
	go build -ldflags="$(LDFLAGS)" -trimpath -o $(BINARY) ./main.go

test:
	go test -race ./...

fmt:
	gofumpt -w .
	goimports -w .

lint:
	golangci-lint run

vuln:
	govulncheck ./...

install:
	go build -ldflags="$(LDFLAGS)" -trimpath -o /usr/local/bin/$(BINARY) ./main.go

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(BINARY):latest .

docker-push:
	docker push $(BINARY):latest
