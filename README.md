# Ethereum Proof-of-Stake Devnet

This repository provides a docker-compose file to run a fully-functional, local development network for Ethereum with proof-of-stake enabled. This configuration uses [Prysm](https://github.com/prysmaticlabs/prysm) as a consensus client and [go-ethereum](https://github.com/ethereum/go-ethereum) for execution.

## Using

First, install Docker. Then, run:

```
git clone https://github.com/rauljordan/eth-pos-devnet && cd eth-pos-devnet
docker compose up -d
```

You will see the following:

```
$ docker compose up -d
[+] Running 7/7
 ⠿ Network eth-pos-devnet_default                          Created
 ⠿ Container eth-pos-devnet-geth-genesis-1                 Started
 ⠿ Container eth-pos-devnet-create-beacon-chain-genesis-1  Started
 ⠿ Container eth-pos-devnet-geth-account-1                 Started
 ⠿ Container eth-pos-devnet-geth-1                         Started
 ⠿ Container eth-pos-devnet-beacon-chain-1                 Started
 ⠿ Container eth-pos-devnet-validator-1                    Started
```
Next, you can inspect the logs of the different services launched and once the mining difficulty of go-ethereum reaches 50, proof-of-stake will be activated and the Prysm beacon chain will be driving consensus of blocks.

<img width="1728" alt="Screen Shot 2022-08-18 at 8 22 57 PM" src="https://user-images.githubusercontent.com/5572669/185518458-25a454a8-b70a-40a8-b3e6-d32770d16ca9.png">
