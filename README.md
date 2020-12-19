# ansible-freeton

Roles of Ansible for install and monitor FreeTon node.

## System requirements

- Ubuntu 18 or newest

## Roles:

- **common** - preparing system and install dependencies
- **freeton** - build and setup FreeTon node
- **netdata** - real-time monitoring
- **prometheus-node-exporter** - exporter for hardware and OS metrics exposed, also this gives opportunity get _balance_ and _diff_ in freeton network

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
- Install nginx for close entry poins of monitoring systems
- Install and sync ntp server for avoid time shift

* System upgrade

## Installation

- Pull repository
- Add your host to `freeton` file
- Change role for installation (common should be always)
- Change nginx user/password for basic_auth in `vars/variables.yml`
- Add telegram bot token and group/chat id in `vars/variables.yml`
- Run ansible: `ansible-playbook freeton.yaml -i freeton --ask-sudo-pass`
- Ansible Build and setup node and save seed phrase `{{ install_path }}/ton-keys/seed_phrase.secret`
- Deploy wallet [instruction](https://docs.ton.dev/86757ecb2/v/0/p/94921e-multisignature-wallet-management-in-tonos-cli)
- install grafana [FreeTon Validator Dashboard](https://grafana.com/grafana/dashboards/13394)

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
