#!/usr/bin/env bash

# Not strictly necessary, but a little more contained than using 'go install'

mkdir -p bin

cmd_name="exceptions"

go build \
    -ldflags "-X main.buildDate=$(date +%Y-%m-%d) -X main.commitLabel=$(git rev-parse --short HEAD)" \
    -o "bin/${cmd_name##*/}" \
    .

