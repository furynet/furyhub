---
order: 3
---

# Local Testnet

For testing or developing purpose, you may want to setup a local testnet.

## Single Node Testnet

**Requirements:**

- [Install grid](../get-started/install.md)

:::tip
We use the default [home directory](intro.md#home-directory) for all the following examples
:::

### grid init

Initialize the genesis.json file that will help you to bootstrap the network

```bash
grid init testing --chain-id=testing
```

### create a key

Create a key to hold your validator account

```bash
grid keys add MyValidator
```

### grid add-genesis-account

Add that key into the genesis.app_state.accounts array in the genesis file

:::tip
this command lets you set the number of coins. Make sure this account has some ufury which is the only staking coin on GRIDhub
:::

```bash
grid add-genesis-account $(grid keys show MyValidator --address) 150000000ugrid
```

### grid gentx

Generate the transaction that creates your validator. The gentxs are stored in `~/.grid/config/gentx/`

```bash
grid gentx MyValidator 100000000ugrid --chain-id=testing 
```

### grid collect-gentxs

Add the generated staking transactions to the genesis file

```bash
grid collect-gentxs
```

### grid start

Change the default token denom to `ufury`

```bash
sed -i 's/stake/ufury/g' $HOME/.grid/config/genesis.json
```

Now it‘s ready to start `grid`

```bash
grid start
```

### grid unsafe-reset-all

You can use this command to reset your node, including the local blockchain database, address book file, and resets priv_validator.json to the genesis state.

This is useful when your local blockchain database somehow breaks and you are not able to sync or participate in the consensus.

```bash
grid unsafe-reset-all
```

### grid tendermint

Query the unique node id which can be used in p2p connection, e.g. the `seeds` and `persistent_peers` in the [config.toml](intro.md#cnofig-toml) are formatted as `<node-id>@ip:26656`.

The node id is stored in the [node_key.json](intro.md#node_key-json).

```bash
grid tendermint show-node-id
```

Query the [Tendermint Pubkey](../concepts/validator-faq.md#tendermint-key) which is used to [identify your validator](../cli-client/stake/create-validator.md), and the corresponding private key will be used to sign the Pre-vote/Pre-commit in the consensus.

The [Tendermint Key](../concepts/validator-faq.md#tendermint-key) is stored in the [priv_validator.json](intro.md#priv_validator-json) which is [required to be backed up](../concepts/validator-faq.md#how-to-backup-the-validator) once you become a validator.

```bash
grid tendermint show-validator
```

Query the bech32 prefixed validator address

```bash
grid tendermint show-address
```

### grid export

Please refer to [Export Blockchain State](export.md)

## Multiple Nodes Testnet

**Requirements:**

- [Install grid](../get-started/install.md)
- [Install jq](https://stedolan.github.io/jq/download/)
- [Install docker](https://docs.docker.com/engine/installation/)
- [Install docker-compose](https://docs.docker.com/compose/install/)

### Build and Init

```bash
# Work from the gridiron repo
cd [your-gridiron-repo]

# Build the linux binary in ./build
make build-linux

# Quick init a 4-node testnet configs
make testnet-init
```

The `make testnet-init` generates config files for a 4-node testnet in the `./build/nodecluster` directory by calling the `grid testnet` command:

```bash
$ tree -L 3 build/nodecluster/
build/nodecluster/
├── gentxs
│   ├── node0.json
│   ├── node1.json
│   ├── node2.json
│   └── node3.json
├── node0
│   ├── grid
│   │   ├── config
│   │   └── data
│   └── gridcli
│       ├── key_seed.json
│       └── keys
├── node1
│   ├── grid
│   │   ├── config
│   │   └── data
│   └── gridcli
│       └── key_seed.json
├── node2
│   ├── grid
│   │   ├── config
│   │   └── data
│   └── gridcli
│       └── key_seed.json
└── node3
    ├── grid
    │   ├── config
    │   └── data
    └── gridcli
        └── key_seed.json
```

### Start

```bash
make testnet-start
```

This command creates a 4-node network using the ubuntu:16.04 docker image. The ports for each node are found in this table:

| Node      | P2P Port | RPC Port |
| --------- | -------- | -------- |
| gridnode0 | 26656    | 26657    |
| gridnode1 | 26659    | 26660    |
| gridnode2 | 26661    | 26662    |
| gridnode3 | 26663    | 26664    |

To update the binary, just rebuild it and restart the nodes:

```bash
make build-linux testnet-start
```

### Stop

To stop all the running nodes:

```bash
make testnet-stop
```

### Clean

To stop all the running nodes and delete all the files in the `build/` directory:

```bash
make testnet-clean
```
