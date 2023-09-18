#!/usr/bin/env bash

set -ex

go build -o /dev/null ./cmd/go-build
go build -o ./bin/main_go -gcflags=all="-N -l" cmd/go/main.go
