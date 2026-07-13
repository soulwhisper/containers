#!/usr/bin/env bash
# dify-sandbox entrypoint — mise + upstream sandbox server.
#
# Activates mise shims (node, officecli) then runs CMD (/main by default).
# Upstream node extraction is skipped; mise provides node via shims.
set -euo pipefail

eval "$(mise activate bash)"
exec "$@"
