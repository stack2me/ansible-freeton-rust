#!/bin/bash

ELECTION_ID=$(/srv/main.ton.dev/ton/build/lite-client/lite-client -p "/home/ton/ton-keys/liteserver.pub" -a 127.0.0.1:3031 -rc "runmethod -1:3333333333333333333333333333333333333333333333333333333333333333 active_election_id" | grep -i 'result:' | awk '{print $3}')

TELEGRAM_BOT_TOKEN=
TELEGRAM_CHAT_ID=
KICK_START=$(cat /home/ton/ton-keys/kick_start.addr)
KICK_START_KEY="/home/ton/ton-keys/kick_start.key.json"
MULTISIG=$(cat /home/ton/ton-keys/$HOSTNAME.addr)
DEPOOL=$(cat /home/ton/ton-keys/depool.addr)
COUNT=0

# log function
function elog() {
    datestring=$(date +"%Y-%m-%d %H:%M:%S")
    echo -e "$datestring - $@"
}

# telegram message function
function notify() {
    curl -s \
        --data parse_mode=HTML \
        --data chat_id=${TELEGRAM_CHAT_ID} \
        --data text="<b>kicker</b>%0A$1" \
        --request POST https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage >/dev/null
}

# check open election
if
    [ $ELECTION_ID = 0 ]
then
    elog "No current elections"
    exit 1
else
    for FILE in "/var/tmp/*"; do
        if [ "$ELECTION_ID" = "$FILE" ]; then
            elog "Already submited"
            COUNT=1
        fi
    done
fi

if [ "$COUNT" == 0 ]; then
    elog "Sending Attempt to voting"
    notify "Sending Attempt to voting"
    touch "/var/tmp/$ELECTION_ID"
    elog "Kick depool"
    notify "Kick depool"
    cd /srv/main.ton.dev/ton/build/utils/ && ./tonos-cli call $KICK_START submitTransaction "{\"dest\":\"$DEPOOL\", \"value\":\"1000000000\", \"bounce\":\"true\", \"allBalance\":\"false\", \"payload\":\"te6ccgEBAQEABgAACCiAmCM=\"}" --abi /srv/main.ton.dev/configs/SafeMultisigWallet.abi.json --sign $KICK_START_KEY
    elog "sleep 1m"
    sleep 1m
    elog "Run validator_depool script"
    notify "Run validator_depool script"
    VALIDATOR=$(/srv/main.ton.dev/scripts/validator_depool.sh)
    elog $VALIDATOR
    sleep 30
    TR=$(cd /srv/main.ton.dev/ton/build/utils/ && ./tonos-cli run $MULTISIG getTransactions {} --abi /srv/main.ton.dev/configs/SafeMultisigWallet.abi.json | grep id | awk '{print $2}' | cut -d '"' -f2)
    notify "Transaction $TR"
fi
