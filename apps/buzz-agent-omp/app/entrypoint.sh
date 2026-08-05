#!/usr/bin/env bash
# buzz-agent-omp entrypoint — buzz-acp harness with omp (oh-my-pi) as the ACP
# agent subprocess.
#
# Identity, relay and model routing are 100% env-driven (k8s `env:` / secretKeyRef);
# the image bakes no keys, no relay, no provider:
#   BUZZ_PRIVATE_KEY   agent nsec identity (required — launch fails closed)
#   BUZZ_RELAY_URL     relay WebSocket URL (default ws://localhost:3000)
#   BUZZ_AUTH_TAG      NIP-OA owner-auth tag
#   BUZZ_API_TOKEN     relay API token (when the relay enforces token auth)
#   BUZZ_ACP_AGENT_COMMAND / BUZZ_ACP_AGENT_ARGS  ACP agent (default: `omp acp`)
#   omp provider credentials per omp docs (ANTHROPIC_API_KEY, OPENAI_API_KEY, …)
set -euo pipefail

# Match sprig/desktop's URL-scoped git credential configuration without
# installing a helper globally (which would answer for unrelated remotes).
if [[ -n "${BUZZ_RELAY_URL:-}" ]]; then
    relay_http_url="${BUZZ_RELAY_URL/#ws:/http:}"
    relay_http_url="${relay_http_url/#wss:/https:}"
    relay_http_url="${relay_http_url%/}"
    git config --global "credential.${relay_http_url}/git.helper" \
        /usr/local/bin/git-credential-nostr
    git config --global "credential.${relay_http_url}/git.useHttpPath" true
fi

# Fail closed on identity (remote-agents spec, I1): never launch an
# identityless agent body.
if [[ -z "${BUZZ_PRIVATE_KEY:-}" && -z "${BUZZ_ACP_PRIVATE_KEY:-}" ]]; then
    echo "buzz-agent-omp: BUZZ_PRIVATE_KEY (agent nsec) is required; refusing to launch identityless" >&2
    exit 1
fi

# The harness must receive the substrate's termination signal directly.
exec buzz-acp "$@"
