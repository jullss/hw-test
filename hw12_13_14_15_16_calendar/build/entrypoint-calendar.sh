#!/bin/sh
set -e

CONFIG_FILE="${CONFIG_FILE:-/etc/calendar/config.yaml}"

mkdir -p "$(dirname "$CONFIG_FILE")"

cat > "$CONFIG_FILE" <<EOF
logger:
  level: "${CALENDAR_LOGGER_LEVEL:-info}"

listen:
  host: "${CALENDAR_LISTEN_HOST:-0.0.0.0}"
  port: "${CALENDAR_LISTEN_PORT:-8888}"
  grpc_port: "${CALENDAR_LISTEN_GRPC_PORT:-50051}"

storage:
  type: "${CALENDAR_STORAGE_TYPE:-sql}"
  dbUrl: "${CALENDAR_STORAGE_DBURL:-postgres://postgres:password@localhost:5432/calendar?sslmode=disable}"
EOF

exec /opt/calendar/calendar-app -config "$CONFIG_FILE"
