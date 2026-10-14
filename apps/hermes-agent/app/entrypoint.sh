#!/bin/bash
# hermes-agent container entrypoint.
#
#   * Optional UID/GID remap of the bundled `hermes` user (HERMES_UID/GID).
#   * chown of the data volume to that user.
#   * Bootstrap default config files on the persistent volume on first run.
#   * Mirror the baked caveman skills onto the volume (exact sync per skill).
#   * Clean up stale gateway pid/lock files from a previous container.
#   * exec the CMD (supervisord as `hermes` by default).
#
# Modeled after the upstream hermes-agent entrypoint; the image's own
# s6-overlay entrypoint is not used (it needs root + a writable /run).

set -e

HERMES_HOME="${HERMES_HOME:-/opt/data}"
INSTALL_DIR="/opt/hermes"
BUNDLED="/bundled"

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
    for d in "${HERMES_HOME}" "/workspace"; do
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
mkdir -p "${HERMES_HOME}"/{cron,sessions,logs,hooks,memories,skills,skins,plans,workspace,home,cache}

# Seed defaults idempotently (standalone/docker UX; in k8s an initContainer
# seeds the authoritative config onto the volume before this runs).
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

# Upstream bundled-skills sync (best-effort; file location differs by version)
if [ -f "${INSTALL_DIR}/tools/skills_sync.py" ]; then
    python3 "${INSTALL_DIR}/tools/skills_sync.py" 2>/dev/null || true
fi

# Baked caveman skills: exact mirror per skill dir (stale copies of a previous
# build must not linger on the PVC)
if [ -d "${BUNDLED}/skills" ]; then
    mkdir -p "${HERMES_HOME}/skills/productivity"
    for d in "${BUNDLED}"/skills/*/; do
        name="$(basename "$d")"
        rm -rf "${HERMES_HOME}/skills/productivity/$name"
        cp -r "$d" "${HERMES_HOME}/skills/productivity/$name"
    done
fi

# Stale lock files survive container restarts on the mounted volume.
for f in gateway.pid gateway.lock; do
    if [ -f "${HERMES_HOME}/${f}" ]; then
        echo "entrypoint: removing stale ${f}"
        rm -f "${HERMES_HOME}/${f}"
    fi
done

# ---- Service selection -----------------------------------------------------

SERVICES="${HERMES_SERVICES:-gateway,dashboard}"
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
        echo "entrypoint: available services: gateway, dashboard" >&2
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
echo " hermes-agent — services: ${normalized}"
case ",${normalized}," in *,gateway,*)   echo " - gateway:   http://0.0.0.0:8642" ;; esac
case ",${normalized}," in *,dashboard,*) echo " - dashboard: http://0.0.0.0:9119" ;; esac
echo "=========================================="

# ---- Single-service fast path: skip supervisord, exec PID 1 ----------------
case "$normalized" in
    gateway)
        exec "${INSTALL_DIR}/.venv/bin/hermes" gateway run
        ;;
    dashboard)
        cd "${HERMES_HOME}"
        exec "${INSTALL_DIR}/.venv/bin/hermes" dashboard \
            --host 0.0.0.0 --port 9119 --no-open
        ;;
esac

# ---- Multi-service: hand off to CMD (supervisord) --------------------------
exec "$@"
