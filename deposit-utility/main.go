package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	contracts "github.com/prysmaticlabs/prysm/v4/contracts/deposit"
	"github.com/prysmaticlabs/prysm/v4/runtime/interop"
)

var amount32Eth = "32000000000000000000"

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
	numDeposits         = flag.Uint64("num-deposits", 1, "number of deposits to make")
	validatorStartIndex = flag.Uint64("validator-start-index", 64, "validator index for which to start making deposits")
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
	depositDataItems, depositDataRoots, err := interop.DepositDataFromKeys(privKeys, pubKeys)
	if err != nil {
		panic(err)
	}
	txOpts.Value = Amount32Eth()
	for i, depositData := range depositDataItems {
		tx, err := depositContract.Deposit(
			txOpts,
			depositData.PublicKey,
			depositData.WithdrawalCredentials,
			depositData.Signature,
			[32]byte(depositDataRoots[i]),
		)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Sent deposit #%d and tx hash %#x\n", i, tx.Hash())
	}
	fmt.Println("Successfully deposited all validators")
}
