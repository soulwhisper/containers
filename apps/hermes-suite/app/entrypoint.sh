#!/bin/bash
# Hermes Suite container entrypoint.
#
# Responsibilities:
#   * Optional UID/GID remap of the bundled `hermes` user (HERMES_UID/GID).
#   * chown of the data volume to that user.
#   * Bootstrap default config files on the persistent volume on first run.
#   * Clean up stale gateway pid/lock files from a previous container.
#   * exec the CMD (supervisord as root by default).
#
# Modeled after the official hermes-agent entrypoint.

set -e

HERMES_HOME="${HERMES_HOME:-/opt/data}"
WORKSPACE="${HERMES_WEBUI_DEFAULT_WORKSPACE:-/workspace}"
INSTALL_DIR="/opt/hermes"

# ---- Root path: privileged setup, then drop to hermes ----------------------
if [ "$(id -u)" = "0" ]; then
    if [ -n "${HERMES_UID:-}" ] && [ "${HERMES_UID}" != "$(id -u hermes)" ]; then
        echo "entrypoint: changing hermes UID to ${HERMES_UID}"
        usermod -u "${HERMES_UID}" hermes
    fi

    if [ -n "${HERMES_GID:-}" ] && [ "${HERMES_GID}" != "$(id -g hermes)" ]; then
        echo "entrypoint: changing hermes GID to ${HERMES_GID}"
        groupmod -o -g "${HERMES_GID}" hermes 2>/dev/null || true
    fi

    target_uid="$(id -u hermes)"
    for d in "${HERMES_HOME}" "${WORKSPACE}"; do
        [ -d "$d" ] || continue
        cur_uid="$(stat -c %u "$d" 2>/dev/null || echo "${target_uid}")"
        if [ "${cur_uid}" != "${target_uid}" ]; then
            echo "entrypoint: chown ${d} -> hermes:hermes"
            chown -R hermes:hermes "${d}" 2>/dev/null \
                || echo "entrypoint: chown ${d} failed (rootless?) — continuing" >&2
        fi
    done

    exec gosu hermes "$0" "$@"
fi

# ---- From here: running as a non-root user (hermes by default) -------------
mkdir -p "${HERMES_HOME}"/{cron,sessions,logs,hooks,memories,skills,skins,plans,workspace,home,webui,cache}

# Seed defaults idempotently.
if [ -f "${INSTALL_DIR}/.env.example" ] && [ ! -f "${HERMES_HOME}/.env" ]; then
    cp "${INSTALL_DIR}/.env.example" "${HERMES_HOME}/.env"
    echo "entrypoint: created ${HERMES_HOME}/.env — edit it to add API keys"
fi
if [ -f "${INSTALL_DIR}/cli-config.yaml.example" ] && [ ! -f "${HERMES_HOME}/config.yaml" ]; then
    cp "${INSTALL_DIR}/cli-config.yaml.example" "${HERMES_HOME}/config.yaml"
    echo "entrypoint: created default config.yaml"
fi
[ -f "${HERMES_HOME}/config.yaml" ] && chmod 600 "${HERMES_HOME}/config.yaml" 2>/dev/null || true
if [ -f "${INSTALL_DIR}/docker/SOUL.md" ] && [ ! -f "${HERMES_HOME}/SOUL.md" ]; then
    cp "${INSTALL_DIR}/docker/SOUL.md" "${HERMES_HOME}/SOUL.md"
    echo "entrypoint: created default SOUL.md"
fi

# Bundled skills sync (best-effort).
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

# ---- Service selection -----------------------------------------------------

SERVICES="${HERMES_SERVICES:-gateway,webui}"
CONF_D="${HERMES_HOME}/.supervisor/conf.d"

mkdir -p "${CONF_D}"
find "${CONF_D}" -maxdepth 1 -name '*.conf' -delete 2>/dev/null || true

normalized=""
old_ifs="$IFS"
IFS=','
for svc in ${SERVICES}; do
    svc="$(echo "$svc" | tr -d '[:space:]')"
    [ -z "$svc" ] && continue
    snippet="/app/programs/${svc}.conf"
    if [ ! -f "$snippet" ]; then
        echo "entrypoint: unknown service '${svc}' (no $snippet)" >&2
        echo "entrypoint: available services: gateway, dashboard, webui" >&2
        exit 1
    fi
    ln -sf "$snippet" "${CONF_D}/${svc}.conf"
    normalized="${normalized}${normalized:+,}${svc}"
done
IFS="$old_ifs"

if [ -z "$normalized" ]; then
    echo "entrypoint: HERMES_SERVICES is empty" >&2
    exit 1
fi

echo "=========================================="
echo " Hermes Suite — services: ${normalized}"
case ",${normalized}," in *,gateway,*)   echo " - gateway:   http://0.0.0.0:8642" ;; esac
case ",${normalized}," in *,dashboard,*) echo " - dashboard: http://0.0.0.0:9119" ;; esac
case ",${normalized}," in *,webui,*)     echo " - webui:     http://0.0.0.0:8787" ;; esac
echo "=========================================="

# ---- Single-service fast path: skip supervisord, exec PID 1 ----------------
case "$normalized" in
    gateway)
        exec "${INSTALL_DIR}/.venv/bin/hermes" gateway run
        ;;
    dashboard)
        exec "${INSTALL_DIR}/.venv/bin/hermes" dashboard \
            --host 0.0.0.0 --port 9119 --insecure --no-open
        ;;
    webui)
        cd "${WEBUI_DIR}"
        exec "${INSTALL_DIR}/.venv/bin/python" "${WEBUI_DIR}/server.py"
        ;;
esac

# ---- Multi-service: hand off to CMD (supervisord) --------------------------
exec "$@"
