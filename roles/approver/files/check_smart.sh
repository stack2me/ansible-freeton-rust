#! /bin/bash
MSIG="cat msig.addr"
DEPOOL="deroop.addr"
KEY="custodian_key_pair.json"

RES=$({{ ton_node_path }}/tools/tonos-cli --url $(cat /etc/ton/ton_network.txt) run $MSIG getTransactions {} --abi /etc/ton/ton-labs-contracts/solidity/safemultisig/SafeMultisigWallet.abi.json | grep id)
if [ -n "$RES" ]; then
    TR_ID=$(echo $RES | awk '{print $2}' | cut -d '"' -f2)
    DEST=$({{ ton_node_path }}/tools/tonos-cli --url $(cat /etc/ton/ton_network.txt) run $MSIG getTransactions {} --abi /etc/ton/ton-labs-contracts/solidity/safemultisig/SafeMultisigWallet.abi.json | grep dest | awk '{print $2}' | cut -d '"' -f2)
    if [ "$DEPOOL" = "$DEST" ]; then
        {{ ton_node_path }}/tools/tonos-cli --url $(cat /etc/ton/ton_network.txt) call $MSIG \
            confirmTransaction "{\"transactionId\":\"$TR_ID\"}" \
            --abi /etc/ton/ton-labs-contracts/solidity/safemultisig/SafeMultisigWallet.abi.json \
            --sign $KEY
        echo "Approved transaction: $TR_ID"
        TITLE="Transaction Approved"
    else
        echo "Transaction ALERT: $TR_ID"
        TITLE="Transaction ALERT"
    fi
    MESSAGE="transaction: $TR_ID"
fi

if [ -n "$TITLE" ]; then
    /usr/local/sbin/telegram_notifier.sh "<b>${TITLE}</b>%0A${MESSAGE}"
fi
