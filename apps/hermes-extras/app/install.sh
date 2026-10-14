#!/bin/sh
# hermes-extras initializer — runs as a Kubernetes initContainer.
# Mirrors bundled binaries, caveman skills and hermes plugin python deps
# into the shared hermes volume.
set -eu

BIN_DIR="${BIN_DIR:-/opt/data/.local/bin}"
SKILLS_DIR="${SKILLS_DIR:-/opt/data/skills/productivity}"
PYLIB_DIR="${PYLIB_DIR:-/opt/data/.local/lib/hermes-plugins}"

mkdir -p "$BIN_DIR" "$SKILLS_DIR" "$PYLIB_DIR"

# Binaries: overwrite what we ship; never delete foreign files in BIN_DIR.
cp -f /data/* "$BIN_DIR/"

# Skills: exact mirror of the dirs we own
# (same set as caveman's upstream hermes installer, bin/install.js).
for skill in caveman caveman-commit caveman-review caveman-help caveman-stats caveman-compress cavecrew; do
  rm -rf "$SKILLS_DIR/$skill"
  cp -r "/skills/$skill" "$SKILLS_DIR/$skill"
done

# Plugin python deps: exact mirror of the packages we own (loaded by the
# hermes venv via PYTHONPATH; stale copies of a previous pin must not linger)
rm -rf "$PYLIB_DIR"/lark_oapi "$PYLIB_DIR"/lark_oapi-*.dist-info \
       "$PYLIB_DIR"/requests_toolbelt "$PYLIB_DIR"/requests_toolbelt-*.dist-info
cp -r /pylib/. "$PYLIB_DIR/"

echo "hermes-extras initialized: binaries -> $BIN_DIR, caveman skills -> $SKILLS_DIR, plugin deps -> $PYLIB_DIR"
