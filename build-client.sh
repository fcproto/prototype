#!/usr/bin/env bash

set -euo pipefail

GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -ldflags="-extldflags '-static' -s -w" ./cmd/edge-service/
