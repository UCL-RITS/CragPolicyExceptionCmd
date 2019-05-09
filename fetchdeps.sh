#!/usr/bin/env bash
# Set up build environment.

if [ -z "$GOPATH" ]; then
  echo "\$GOPATH is not set."
  exit 1
else
  go get github.com/go-sql-driver/mysql
  go get github.com/olekukonko/tablewriter
  go get gopkg.in/alecthomas/kingpin.v2
  go get gopkg.in/AlecAivazis/survey.v1
  go get github.com/jinzhu/gorm
  go get github.com/mattn/go-sqlite3
  go get golang.org/x/crypto/ssh/terminal
fi
