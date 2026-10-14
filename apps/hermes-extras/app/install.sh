#!/bin/sh
# hermes-extras initializer — runs as a Kubernetes initContainer.
# Mirrors hermes plugin python deps into the shared hermes volume
# (loaded by the agent venv via PYTHONPATH).
set -eu

PYLIB_DIR="${PYLIB_DIR:-/opt/data/.local/lib/hermes-plugins}"

mkdir -p "$PYLIB_DIR"

# Exact mirror of the packages we own — stale copies of a previous pin must
# not linger next to a newer pin on the same PYTHONPATH entry.
rm -rf "$PYLIB_DIR"/lark_oapi "$PYLIB_DIR"/lark_oapi-*.dist-info \
       "$PYLIB_DIR"/requests_toolbelt "$PYLIB_DIR"/requests_toolbelt-*.dist-info
cp -r /pylib/. "$PYLIB_DIR/"

echo "hermes-extras initialized: plugin deps -> $PYLIB_DIR"
