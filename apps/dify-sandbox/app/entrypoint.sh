#!/usr/bin/env bash
# dify-sandbox entrypoint — mise + node + officecli.
#
# No argv  -> runs the CMD (sleep infinity); `kubectl exec -it -- bash` to use.
# With argv -> exec's it as-is.
set -euo pipefail

exec "$@"
