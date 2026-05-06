#!/bin/sh
# Stork container entrypoint.
#
# Two execution modes are supported:
#
#   1. Direct (default) — exec stork-server or stork-agent as PID 1.
#      Configuration is supplied either via container env vars
#      (STORK_DATABASE_*, STORK_REST_*, STORK_AGENT_*, ...) or via an env
#      file mounted at:
#         /etc/stork/server.env   (when STORK_MODE=server)
#         /etc/stork/agent.env    (when STORK_MODE=agent)
#      The path can be overridden with STORK_ENV_FILE. Vars from the env
#      file are exported to the binary's environment.
#      This is the preferred mode for Kubernetes and Docker Compose.
#
#   2. Supervisor (opt-in) — set STORK_SUPERVISOR_CONF to a readable
#      supervisord config path and that config will be used instead.
#      Useful when co-locating multiple processes (server + agent,
#      exporters, hooks) in one container, or for parity with the
#      previous image behavior. Sample configs may be shipped under
#      /app/defaults/.
#
# Mode is selected by STORK_MODE: "server" (default) or "agent".
# Any extra CLI arguments are forwarded to the stork binary.

set -eu

STORK_MODE="${STORK_MODE:-server}"

# ---- Supervisor mode -------------------------------------------------------
if [ -n "${STORK_SUPERVISOR_CONF:-}" ]; then
    if [ ! -f "${STORK_SUPERVISOR_CONF}" ]; then
        echo "entrypoint: STORK_SUPERVISOR_CONF=${STORK_SUPERVISOR_CONF} not found" >&2
        exit 1
    fi
    echo "entrypoint: starting supervisord with ${STORK_SUPERVISOR_CONF}"
    exec /usr/bin/supervisord -c "${STORK_SUPERVISOR_CONF}"
fi

# ---- Direct mode -----------------------------------------------------------
case "${STORK_MODE}" in
    server)
        binary="stork-server"
        env_file_default="/etc/stork/server.env"
        ;;
    agent)
        binary="stork-agent"
        env_file_default="/etc/stork/agent.env"
        ;;
    *)
        echo "entrypoint: unknown STORK_MODE='${STORK_MODE}' (expected 'server' or 'agent')" >&2
        exit 1
        ;;
esac

env_file="${STORK_ENV_FILE:-${env_file_default}}"
if [ -f "${env_file}" ]; then
    echo "entrypoint: loading env from ${env_file}"
    set -a
    # shellcheck disable=SC1090
    . "${env_file}"
    set +a
fi

echo "entrypoint: exec ${binary} (mode=${STORK_MODE})"
exec "${binary}" "$@"
