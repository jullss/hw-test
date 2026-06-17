#!/bin/sh
set -e

CONFIG_FILE="${CONFIG_FILE:-/etc/calendar/config.yaml}"

mkdir -p "$(dirname "$CONFIG_FILE")"

cat > "$CONFIG_FILE" <<EOF
logger:
  level: "${SENDER_LOGGER_LEVEL:-info}"
rabbitmq:
  url: "${SENDER_RABBITMQ_URL:-amqp://guest:guest@localhost:5672/}"
  queue_name: "${SENDER_RABBITMQ_QUEUE:-calendar_notifications}"
  sent_queue_name: "${SENDER_RABBITMQ_SENT_QUEUE:-calendar_notifications_sent}"
EOF

exec /opt/calendar/sender-app -config "$CONFIG_FILE"
