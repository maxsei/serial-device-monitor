#!/usr/bin/env bash

set -ex

go build -o /dev/null ./cmd/build-test
go build -o ./bin/example -gcflags=all="-N -l" cmd/example/main.go
