#!/usr/bin/env bash

set -euo pipefail

docker buildx build --platform linux/amd64 --push -t gcr.io/fcproto/server .
gcloud run deploy server --platform managed --region europe-central2 --image gcr.io/fcproto/server
