#!/bin/ash
set -eu

cd /usr/src/app

if [ -f tmp/pids/server.pid ]; then
  rm tmp/pids/server.pid
fi

echo "Preparing database..."
bundle exec rails db:prepare:with_data

echo "Setting database file ownership (SQLite3 only)..."
bundle exec rake db:chown

echo "Cleaning up old cache files..."
bundle exec rake tmp:cache:clear

if [ ! -d "$PLUGINS_PATH" ]; then
  echo "Creating plugin directory..."
  mkdir -p "$PLUGINS_PATH"
fi
if [ -w "$PLUGINS_PATH" ]; then
  chown "$PUID:$PGID" "$PLUGINS_PATH"
fi

echo "Launching application..."
export RAILS_PORT="$PORT"
export RAILS_LOG_TO_STDOUT=true
exec "$@"
