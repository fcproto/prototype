#!/usr/bin/env bash

set -euo pipefail

watch -n 0.5 -t curl -s https://server-ix6omulhiq-lm.a.run.app/status
