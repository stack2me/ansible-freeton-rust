# ansible-freeton

Roles of Ansible for install and monitor FreeTon node.

## System requirements

- Ubuntu 18 or newest

## Roles:

- **common** - preparing system and install dependencies
- **rust​** - install and upgrade rust
- **freeton** - build and setup FreeTon node
- **netdata** - real-time monitoring
- **prometheus-node-exporter** - exporter for hardware and OS metrics exposed, also this gives opportunity get _balance_ and _diff_ in freeton network
- **notifier**​ - telegram notifications about node events
- **promtail**​ - log collector for Loki

## Functional

- Freeton Install

  - Creating user and group
  - Cronjob for validator script
  - All logs in one folder /var/log/...
  - Systemd for control status of node and restart in fail case
  - Logrotate for archive logs

- Node Monitoring
  - Install netdata for realtime status <host>/netdata
  - install prometheus-node-exporter for collect metrics
    - collecting data about node status(node diff, wallet balance, total validators, if your node became validator, open elections)
- Notifications about node events like an open election, approve/reject transactions etc...
- Script automatically sending stake
- Script control transactions (confirm/reject) with notifications to telegram
- Install nginx for close entry poins of monitoring systems
- Install and sync ntp server for avoid time shift

* System upgrade

## Installation 1.1

- Pull repository
- Add your host to `freeton` file
- Change role for installation (common should be always)
- Change nginx user/password for basic_auth in `vars/variables.yml`
- Add telegram bot token and group/chat id in `vars/variables.yml`
- Run ansible: `ansible-playbook freeton.yaml -i freeton --ask-sudo-pass`
- Deploy wallet [instruction](https://docs.ton.dev/86757ecb2/v/0/p/94921e-multisignature-wallet-management-in-tonos-cli)
- install grafana [FreeTon Validator Dashboard](https://grafana.com/grafana/dashboards/13394)

## Installation 1.2

### Setup validators keys

- Generate Multisig wallet [SafeMultisig](https://github.com/tonlabs/ton-labs-contracts/tree/master/solidity/safemultisig).
- - Multisig address should be situated `/etc/ton/configs/${VALIDATOR_NAME}.addr`
- - Multisig keys should be situated `/etc/ton/configs/keys/msig.keys.json`
- Generate Kicker wallet [Kicker](https://github.com/tonlabs/ton-labs-contracts/tree/master/solidity/safemultisig).
- - Kicker address should be situated `/etc/ton/configs/kick_start.addr`
- - Kicker keys should be situated `/etc/ton/configs/keys/kick_start.key.json`
- Generate Depool wallet [Depoolv3](https://docs.ton.dev/86757ecb2/p/04040b-run-depool-v3).
- - Kicker address should be situated `/etc/configs/ton/depool.addr`
- - Kicker keys should be situated `/etc/configs/ton/keys/depool.key.json`

## Custom metrics in prometheus-node-exporter

- **ton_node_diff** - seconds until synchronization will complete
- **ton_node_balance** - current wallet balance
- **ton_total_validators** - number of validators
- **ton_election_num** - election numbers
- **ton_elections** - election status (0 - closed 1 - open)
- **ton_aggregateBlockSignatures** - number of signed blocks by node
- **ton_getTransactionsCount** - numbers of transaction
- **ton_getAccountsCount** - total accounts in net.ton.dev network
- **ton_getAccountsTotalBalance** - total balance of all accounts
- **ton_aggregateBlocks** - blocks by current validators

## Alerts

Before install grafana template (grafana\*freeton\*node\*alerts.json) please replace "\_NODE-IP:8080" on IP address and port of your server.

## Example Dashboard based on prometheus-node-exporter

![Alt text](images/dashboard.png?raw=true "FreeTon dashboard")

![Alt text](images/dashboard2.png?raw=true "FreeTon dashboard part2")

## Todo List:

- Alerts for log
- Dashboard for logs
- Distribution by binaries
- Security improvements
- Alerts based on prometheus
