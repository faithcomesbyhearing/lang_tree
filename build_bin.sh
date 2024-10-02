#!/bin/bash
GOOS=linux GOARCH=amd64 go install .
GOOS=darwin GOARCH=arm64 go install .
