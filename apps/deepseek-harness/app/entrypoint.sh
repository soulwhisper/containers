#!/bin/sh
# deepseek-harness appliance entrypoint.
#   Profile:  $1 or $DSH_PROFILE (web | tavern), default web.
#   Seed:     plugins baked in /opt/dsh-seed are linked/copied into $DSH_HOME
#             on first boot and re-synced (code only) when the bundled
#             versions change. User data (profile-data/) is never touched.
set -eu

PROFILE="${1:-${DSH_PROFILE:-web}}"
case "$PROFILE" in
  web|tavern) ;;
  *) echo "unknown profile '$PROFILE' (expected: web | tavern)" >&2; exit 1 ;;
esac

SEED=/opt/dsh-seed
STAMP="$DSH_HOME/.bundled-version"

# The base image carries no /opt/dsh-seed (plugin-free basement); children
# bake one. Seeding is a no-op then.
if [ -d "$SEED" ]; then
  if [ ! -f "$STAMP" ]; then
    if [ -z "$(ls -A "$DSH_HOME" 2>/dev/null)" ]; then
      echo "entrypoint: fresh DSH_HOME, seeding baked profiles" >&2
      mkdir -p "$DSH_HOME"
      cp -a "$SEED/." "$DSH_HOME/"
      cp "$SEED/BUNDLED_VERSIONS" "$STAMP"
    else
      echo "entrypoint: DSH_HOME has data but no version stamp; syncing plugin code only" >&2
      for p in web tavern; do
        [ -d "$SEED/profiles/$p" ] || continue
        mkdir -p "$DSH_HOME/profiles/$p"
        rsync -a --delete "$SEED/profiles/$p/node_modules/" "$DSH_HOME/profiles/$p/node_modules/"
        rsync -a --ignore-existing --exclude node_modules "$SEED/profiles/$p/" "$DSH_HOME/profiles/$p/" || true
      done
      cp "$SEED/BUNDLED_VERSIONS" "$STAMP"
    fi
  elif ! cmp -s "$SEED/BUNDLED_VERSIONS" "$STAMP"; then
    echo "entrypoint: bundled versions changed, re-syncing plugin code" >&2
    diff "$SEED/BUNDLED_VERSIONS" "$STAMP" >&2 || true
    for p in web tavern; do
      [ -d "$SEED/profiles/$p" ] || continue
      mkdir -p "$DSH_HOME/profiles/$p"
      rsync -a --delete "$SEED/profiles/$p/node_modules/" "$DSH_HOME/profiles/$p/node_modules/"
      rsync -a --ignore-existing --exclude node_modules "$SEED/profiles/$p/" "$DSH_HOME/profiles/$p/" || true
    done
    cp "$SEED/BUNDLED_VERSIONS" "$STAMP"
  fi
fi

# web relay: dsh listens on 127.0.0.1:DSH_PORT only; the relay binds the
# container's concrete address so the Service/LAN can reach it (design
# borrowed from runzhliu/deepseek-harness-docker, MIT).
node /app/dsh-web-relay.cjs &
relay_pid=$!

stop() {
  kill -TERM "$relay_pid" 2>/dev/null || true
}
trap stop INT TERM HUP

echo "entrypoint: starting dsh with profile '$PROFILE' (relay listens on ${DSH_PORT:-3080})" >&2
dsh --profile "$PROFILE" --no-open &
dsh_pid=$!

# forward the first exiting child; if dsh dies, stop the relay too
wait "$dsh_pid"
status=$?
stop
wait "$relay_pid" 2>/dev/null || true
exit "$status"
