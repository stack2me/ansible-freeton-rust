#!/bin/bash -eEx

echo "INFO: R-Node startup..."

echo "INFO: NETWORK_TYPE = ${NETWORK_TYPE}"

./ton_node --configs configs/ >>logs/ton.log

echo "INFO: R-Node startup... DONE"
