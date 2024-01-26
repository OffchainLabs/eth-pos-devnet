package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	contracts "github.com/prysmaticlabs/prysm/v4/contracts/deposit"
	"github.com/prysmaticlabs/prysm/v4/runtime/interop"
)

var (
	depositContractAddrFlag = flag.String(
		"deposit-contract-address", "0x4242424242424242424242424242424242424242", "deposit contract address",
	)
	rpcEndpoint = flag.String(
		"rpc-endpoint", "http://localhost:8545", "Ethereum JSON-RPC endpoint",
	)
	ethereumPrivKey = flag.String(
		"priv-key", "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622", "Ethereum private key to send the deposit tx from",
	)
	numDeposits            = flag.Uint64("num-deposits", 1, "number of deposits to make")
	validatorStartIndex    = flag.Uint64("validator-start-index", 64, "validator index for which to start making deposits")
	checkStatus            = flag.Bool("check-validator-status", false, "check the status of a validator index")
	validatorIndex         = flag.Uint64("validator-index", 64, "the validator index to check the status for")
	beaconChainAPIEndpoint = flag.String("beacon-api-endpoint", "http://localhost:3500", "beacon API endpoint")
	amount32Eth            = "32000000000000000000"
)

func Amount32Eth() *big.Int {
	amount, _ := new(big.Int).SetString(amount32Eth, 10)
	return amount
}

func main() {
	flag.Parse()
	ctx := context.Background()
	contractAddr := common.BytesToAddress(hexutil.MustDecode(*depositContractAddrFlag))
	rpcClient, err := rpc.Dial(*rpcEndpoint)
	if err != nil {
		panic(err)
	}
	client := ethclient.NewClient(rpcClient)
	depositContract, err := contracts.NewDepositContract(contractAddr, client)
	if err != nil {
		panic(err)
	}
	if *checkStatus {
		endpoint := fmt.Sprintf("%s:/eth/v1/beacon/states/head/validators/%d", *beaconChainAPIEndpoint, *validatorIndex)
		response, err := http.Get(endpoint)
		if err != nil {
			panic(err)
		}
		resp := make(map[string]any)
		if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", resp)
		return
	}
	validatorPrivateKey, err := crypto.HexToECDSA(*ethereumPrivKey)
	if err != nil {
		panic(err)
	}
	l1ChainId, err := client.ChainID(ctx)
	if err != nil {
		panic(err)
	}
	txOpts, err := bind.NewKeyedTransactorWithChainID(validatorPrivateKey, l1ChainId)
	if err != nil {
		panic(err)
	}
	balance, err := client.BalanceAt(ctx, txOpts.From, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Account balance", balance.String())
	privKeys, pubKeys, err := interop.DeterministicallyGenerateKeys(*validatorStartIndex, *numDeposits)
	if err != nil {
		panic(err)
	}
	depositDataItems, _, err := interop.DepositDataFromKeys(privKeys, pubKeys)
	if err != nil {
		panic(err)
	}
	txOpts.Value = Amount32Eth()
	domain, err := signing.ComputeDomain(
		params.BeaconConfig().DomainDeposit,
		nil, /*forkVersion*/
		nil, /*genesisValidatorsRoot*/
	)
	if err != nil {
		panic(err)
	}
	for i, depositData := range depositDataItems {
		if err := contracts.VerifyDepositSignature(depositData, domain); err != nil {
			panic("deposit failed to verify")
		}
		root, err := depositData.HashTreeRoot()
		if err != nil {
			panic(err)
		}
		tx, err := depositContract.Deposit(
			txOpts,
			depositData.PublicKey,
			depositData.WithdrawalCredentials,
			depositData.Signature,
			root,
		)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Sent deposit #%d and tx hash %#x\n", i, tx.Hash())
	}
	fmt.Println("Successfully deposited all validators")
}
