#!/usr/bin/env bash
set -e

GO_VERSION=1.25.5

curl -LO https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
tar -C /tmp -xzf go${GO_VERSION}.linux-amd64.tar.gz
export PATH=/tmp/go/bin:$PATH

go version
go run ./cmd/geode/main.go build -dir docs
# make build-site DIR=docs
