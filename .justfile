#!/usr/bin/env -S just --justfile

set quiet := true
set shell := ['bash', '-eu', '-o', 'pipefail', '-c']

[private]
default:
  just --list

[doc('Build and test an app locally')]
local-build app:
  #!/usr/bin/env bash
  set -euo pipefail
  build_dir="$(mktemp -d)"
  trap 'rm -rf "$build_dir"' EXIT
  rsync -aqIP {{ justfile_dir() }}/include/ {{ justfile_dir() }}/apps/{{ app }}/ "$build_dir"/
  cd "$build_dir"
  docker buildx bake --no-cache --metadata-file docker-bake.json --set=*.output=type=docker --load
  TEST_IMAGE="$(jq -r '."image-local"."image.name" | sub("^docker.io/library/"; "")' docker-bake.json)" \
    go test -v {{ justfile_dir() }}/apps/{{ app }}/...
