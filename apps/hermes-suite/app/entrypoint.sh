#!/bin/bash
# Hermes Suite container entrypoint.
#
# Responsibilities:
#   * Optional UID/GID remap of the bundled `hermes` user (HERMES_UID/GID).
#   * chown of the data volume to that user.
#   * Drop root via gosu.
#   * Bootstrap default config files on the persistent volume on first run.
#   * Clean up stale gateway pid/lock files from a previous container.
#   * exec the CMD (supervisord by default).
#
# Modeled after the official hermes-agent entrypoint.

set -e

HERMES_HOME="${HERMES_HOME:-/opt/data}"
INSTALL_DIR="/opt/hermes"

# ---- Privilege handling (only when started as root) -----------------------
if [ "$(id -u)" = "0" ]; then
    if [ -n "${HERMES_UID:-}" ] && [ "${HERMES_UID}" != "$(id -u hermes)" ]; then
        echo "entrypoint: changing hermes UID to ${HERMES_UID}"
        usermod -u "${HERMES_UID}" hermes
    fi

    if [ -n "${HERMES_GID:-}" ] && [ "${HERMES_GID}" != "$(id -g hermes)" ]; then
        echo "entrypoint: changing hermes GID to ${HERMES_GID}"
        groupmod -o -g "${HERMES_GID}" hermes 2>/dev/null || true
    fi

    actual_uid="$(id -u hermes)"
    if [ "$(stat -c %u "${HERMES_HOME}" 2>/dev/null || echo "${actual_uid}")" != "${actual_uid}" ]; then
        echo "entrypoint: ${HERMES_HOME} not owned by ${actual_uid}, fixing"
        chown -R hermes:hermes "${HERMES_HOME}" 2>/dev/null \
            || echo "entrypoint: chown failed (rootless?) — continuing"
    fi

    echo "entrypoint: dropping root privileges"
    exec gosu hermes "$0" "$@"
fi

# ---- From here on we are the hermes user ----------------------------------
# shellcheck disable=SC1091
source "${INSTALL_DIR}/.venv/bin/activate"

mkdir -p "${HERMES_HOME}"/{cron,sessions,logs,hooks,memories,skills,skins,plans,workspace,home,webui,cache}

if [ ! -f "${HERMES_HOME}/.env" ] && [ -f "${INSTALL_DIR}/.env.example" ]; then
    cp "${INSTALL_DIR}/.env.example" "${HERMES_HOME}/.env"
    echo "entrypoint: created default .env — edit ${HERMES_HOME}/.env to add API keys"
fi

if [ ! -f "${HERMES_HOME}/config.yaml" ] && [ -f "${INSTALL_DIR}/cli-config.yaml.example" ]; then
    cp "${INSTALL_DIR}/cli-config.yaml.example" "${HERMES_HOME}/config.yaml"
    echo "entrypoint: created default config.yaml"
fi

if [ -f "${HERMES_HOME}/config.yaml" ]; then
    chmod 640 "${HERMES_HOME}/config.yaml" 2>/dev/null || true
fi

if [ ! -f "${HERMES_HOME}/SOUL.md" ] && [ -f "${INSTALL_DIR}/docker/SOUL.md" ]; then
    cp "${INSTALL_DIR}/docker/SOUL.md" "${HERMES_HOME}/SOUL.md"
fi

if [ -d "${INSTALL_DIR}/skills" ] && [ -f "${INSTALL_DIR}/tools/skills_sync.py" ]; then
    python3 "${INSTALL_DIR}/tools/skills_sync.py" 2>/dev/null || true
fi

# Stale lock files survive container restarts on the mounted volume.
for f in gateway.pid gateway.lock; do
    if [ -f "${HERMES_HOME}/${f}" ]; then
        echo "entrypoint: removing stale ${f}"
        rm -f "${HERMES_HOME}/${f}"
    fi
done

mkdir -p /var/log/supervisor /var/run/supervisor

echo "=========================================="
echo " Hermes Suite — all-in-one container"
echo "=========================================="
echo " Gateway:    http://0.0.0.0:8642"
echo " Dashboard:  http://0.0.0.0:9119"
echo " WebUI:      http://0.0.0.0:8787"
echo "=========================================="

exec "$@"
