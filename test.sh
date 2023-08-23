#!/usr/bin/env bash
set -eu -o pipefail

set -x
golangci-lint run
go test