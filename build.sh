#!/usr/bin/env bash

# Not strictly necessary, but a little more contained than using 'go install'

mkdir -p bin

for dir in ./cmd/*; do
  # For whatever reason, using go-sqlite3 for this project now requires a cross-compiler
  #  where it didn't before. See: https://github.com/mattn/go-sqlite3/issues/742
  for arch_label in ""; do
    export GOOS="${arch_label%%-*}"
    export GOARCH="${arch_label##*-}"
    go build \
        -ldflags "-X main.buildDate=$(date +%Y-%m-%d) -X main.commitLabel=$(git rev-parse --short HEAD)" \
        -o "bin/${dir##*/}${arch_label:+.}${arch_label}" \
        "$dir" 
  done
done

