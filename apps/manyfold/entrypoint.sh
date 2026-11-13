#!/bin/ash
set -eu

export TMPDIR="${TMPDIR:-/tmp/manyfold}"
mkdir -p "$TMPDIR"
chmod 0700 "$TMPDIR"

exec "$@"
