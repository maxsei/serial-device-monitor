#!/usr/bin/env bash

set -ex

clang -o ./bin/main_c -g -O0 `pkg-config libudev --cflags --libs` cmd/c/main.c
go build -o ./bin/main_go -gcflags=all="-N -l" cmd/go/main.go
