#!/usr/bin/env bash

set -ex

clang -o ./bin/main_c -g -O0 `pkg-config libudev --cflags --libs` ./main.c
go build -o ./bin/main_go main.go
