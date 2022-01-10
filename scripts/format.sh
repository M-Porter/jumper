#!/usr/bin/env bash

set -eu

THIS_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
PROJ_ROOT=$(cd -- "$(dirname -- "${THIS_DIR}/../.")" &> /dev/null && pwd)

echo "Running gofmt on $PROJ_ROOT"
gofmt -w "$PROJ_ROOT"
