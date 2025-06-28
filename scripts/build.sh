#!/bin/sh

BASE_OUTPUT=bin

# Build for mac
env GOOS=darwin GOARCH=arm64 go build -v -o $BASE_OUTPUT/dbcli

# Build for windows
env GOOS=windows GOARCH=amd64 go build -v -o $BASE_OUTPUT/dbcli.exe
