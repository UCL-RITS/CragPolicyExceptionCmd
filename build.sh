#!/usr/bin/env bash

# Not strictly necessary, but a little more contained than using 'go install'

mkdir -p bin

for dir in ./cmd/*; do
  for arch_label in "" linux-amd64 linux-arm linux-arm64 darwin-amd64; do
    export GOOS="${arch_label%%-*}"
    export GOARCH="${arch_label##*-}"
    go build \
        -ldflags "-X main.buildDate=$(date +%Y-%m-%d) -X main.commitLabel=$(git rev-parse --short HEAD)" \
        -o "bin/${dir##*/}${arch_label:+.}${arch_label}" \
        "$dir" 
  done
done

