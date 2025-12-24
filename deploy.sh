#!/bin/sh

curl -LO https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
tar -C /tmp -xzf go1.22.3.linux-amd64.tar.gz
export PATH=/tmp/go/bin:$PATH

go version
make build-site DIR=docs
