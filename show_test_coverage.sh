#!/usr/bin/env bash
set -eu -o pipefail
rm -f coverage.out
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
rm -f coverage.out
