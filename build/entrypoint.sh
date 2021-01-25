#!/bin/bash -eEx

echo "INFO: R-Node startup..."

echo "INFO: NETWORK_TYPE = ${NETWORK_TYPE}"

exec /home/ton/ton_node --configs /home/ton/configs/ >>logs/ton.log 2>>logs/ton.err.log

echo "INFO: R-Node startup... DONE"
