package eth

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/fraternal/ethapi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)


func testA(eth *Ethereum) {

	pubServAPI := PublicEthereumAPI{eth}

	pubBlockChainApi := ethapi.NewPublicBlockChainAPI(eth.ApiBackend)

	blockNum := pubBlockChainApi.BlockNumber()

	ethAd, _ := pubServAPI.Etherbase()
	coinAd, _ := pubServAPI.Coinbase()
	pubAd := pubServAPI.Hashrate()

	fmt.Printf("%v\n%v\n%v\n%v\n", ethAd, coinAd, pubAd, blockNum)

	//fmt.Printf("~~~~~~~~~~~~~~~~~~~\n\n\n\n")
	//fmt.Printf("%+v\n\n", backend.CurrentBlock())
	//fmt.Printf("damn")
	//
}


func testB(eth *Ethereum) {
	//am := accounts.NewManager(eth.ApiBackend)
}

func countAccounts(eth *Ethereum) {
	blockNum := 0
	trans := 0
	blocks := 0
	addresses := make(map[*common.Address]bool)

	fmt.Printf("%s Starting to count accounts...\n", time.Now().String())

	printEvery := 100000

	for {
		blck := blockNum
		blockNum += 1

		// if (blck > 200000) {
		//     break
		// }

		block := eth.blockchain.GetBlockByNumber(uint64(blck))
		if block == nil {
			break
		}

		trans = trans + block.Transactions().Len()

		for i := 0; i < block.Transactions().Len(); i++ {
			transaction := block.Transactions()[i]
			if transaction.Value().Sign() > 0 {
				addresses[transaction.To()] = true
			}
		}

		blocks += 1

		if blocks % printEvery == 0 {
			fmt.Printf("%s %v blocks processed. Number of transactions til block number %v: %v transactions. Number of accounts: %v\n",
				time.Now().String(), blocks, blck, trans, len(addresses))
		}
	}

	fmt.Printf("Total blocks processed: %v Total transactions: %v\n", blocks, trans)

	//let positiveAddresses = []

	var numAddresses = len(addresses)
	fmt.Printf("%s Total number of accounts: %v\n", time.Now().String(), numAddresses)
	var positiveAddresses = 0

	publicBlockChainAPI := ethapi.NewPublicBlockChainAPI(eth.ApiBackend)

	for address := range addresses {
		if address == nil {
			fmt.Println(address)
			continue
		}

		balance, err := publicBlockChainAPI.GetBalance(nil, *address, rpc.LatestBlockNumber)
		if err != nil {
			fmt.Printf("%s Error in GetBalance(): %v\n", time.Now().String(), err)
		} else if balance.Sign() > 0 {
			positiveAddresses += 1
		}
	}

	fmt.Printf("%s Number of accounts with positive balance: %v \n", time.Now().String(), positiveAddresses)
	fmt.Printf("%s Exiting process.\n", time.Now().String())

}

