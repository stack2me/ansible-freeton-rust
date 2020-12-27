#!/bin/bash
TELEGRAM_CHAT_ID=$(head -n 1 /etc/ton/telegram.secret)
TELEGRAM_BOT_TOKEN=$(tail -n 1 /etc/ton/telegram.secret)
HOST=$(hostname)

curl -s \
    --data parse_mode=HTML \
    --data chat_id=${TELEGRAM_CHAT_ID} \
    --data text="<b>host: $HOST</b>%0A$1" \
    --request POST https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage >/dev/null
