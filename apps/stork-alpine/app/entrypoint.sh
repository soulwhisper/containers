#!/bin/bash
set -e

STORK_MODE=${STORK_MODE:-server}

echo "Starting Stork in $STORK_MODE mode..."

if [ "$STORK_MODE" = "server" ]; then
    exec /usr/bin/supervisord -c /app/supervisor/server.conf
elif [ "$STORK_MODE" = "agent" ]; then
    exec /usr/bin/supervisord -c /app/supervisor/agent.conf
else
    echo "Unknown mode: $STORK_MODE. Use 'server' or 'agent'."
    exit 1
fi
