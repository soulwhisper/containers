#!/usr/bin/env bash
# devbox entrypoint — claude-code + git, mise-managed tools.
#
# Backend routing is 100% env-driven (set via k8s `env:` / docker -e); the image
# bakes no provider, so it never forces models.
#   ANTHROPIC_BASE_URL    Anthropic-compatible route
#   ANTHROPIC_AUTH_TOKEN  bearer token (inject from a Secret)
#   ANTHROPIC_MODEL, ANTHROPIC_DEFAULT_{OPUS,SONNET,HAIKU}_MODEL,
#   CLAUDE_CODE_SUBAGENT_MODEL  model ids
#
# No argv  -> runs the CMD (sleep infinity); `kubectl exec -it -- bash` to use.
# With argv -> exec's it as-is (`claude`, `git ...`, CI commands).
set -euo pipefail

exec "$@"
