#!/bin/sh
set -e

CONFIG_FILE="${CONFIG_FILE:-/etc/calendar/config.yaml}"

mkdir -p "$(dirname "$CONFIG_FILE")"

cat > "$CONFIG_FILE" <<EOF
logger:
  level: "${SCHEDULER_LOGGER_LEVEL:-info}"
storage:
  dbUrl: "${SCHEDULER_STORAGE_DBURL:-postgres://postgres:password@localhost:5432/calendar?sslmode=disable}"
rabbitmq:
  url: "${SCHEDULER_RABBITMQ_URL:-amqp://guest:guest@localhost:5672/}"
  queue_name: "${SCHEDULER_RABBITMQ_QUEUE:-calendar_notifications}"
scheduler:
  scan_interval: "${SCHEDULER_SCAN_INTERVAL:-5s}"
EOF

exec /opt/calendar/scheduler-app -config "$CONFIG_FILE"
