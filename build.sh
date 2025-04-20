#!/bin/bash
mkdir -p functions
GOOS=linux GOARCH=amd64 go build -o functions/hello main.go
