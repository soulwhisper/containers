#!/bin/sh
# hermes-extras initializer — runs as a Kubernetes initContainer.
# Mirrors bundled binaries and caveman skills into the shared hermes volume.
set -eu

BIN_DIR="${BIN_DIR:-/opt/data/.local/bin}"
SKILLS_DIR="${SKILLS_DIR:-/opt/data/skills/productivity}"

mkdir -p "$BIN_DIR" "$SKILLS_DIR"

# Binaries: overwrite what we ship; never delete foreign files in BIN_DIR.
cp -f /data/* "$BIN_DIR/"

# Skills: exact mirror of the dirs we own
# (same set as caveman's upstream hermes installer, bin/install.js).
for skill in caveman caveman-commit caveman-review caveman-help caveman-stats caveman-compress cavecrew; do
  rm -rf "$SKILLS_DIR/$skill"
  cp -r "/skills/$skill" "$SKILLS_DIR/$skill"
done

echo "hermes-extras initialized: binaries -> $BIN_DIR, caveman skills -> $SKILLS_DIR"
