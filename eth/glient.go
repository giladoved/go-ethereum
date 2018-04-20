package eth

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/fraternal/ethapi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)


func now() string {
	return time.Now().Format("Jan 2 15:04:05 >")
}


func testA(eth *Ethereum) {

	pubServAPI := PublicEthereumAPI{eth}

	pubBlockChainApi := ethapi.NewPublicBlockChainAPI(eth.ApiBackend)

	blockNum := pubBlockChainApi.BlockNumber()

	ethAd, _ := pubServAPI.Etherbase()
	coinAd, _ := pubServAPI.Coinbase()
	pubAd := pubServAPI.Hashrate()

	fmt.Printf("%v\n%v\n%v\n%v\n", ethAd, coinAd, pubAd, blockNum)

	//am := accounts.NewManager(eth.ApiBackend

	//fmt.Printf("~~~~~~~~~~~~~~~~~~~\n\n\n\n")
	//fmt.Printf("%+v\n\n", backend.CurrentBlock())
	//fmt.Printf("damn")
	//
}


func countAccounts(eth *Ethereum) {
	var blockNum = 0
	var trans = 0
	var blocks = 0
	addresses := make(map[*common.Address]bool)

	fmt.Printf("%s Starting to count accounts...\n", now())

	var printEvery = 100000

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
				now(), blocks, blck, trans, len(addresses))
		}
	}

	fmt.Printf("Total blocks processed: %v Total transactions: %v\n", blocks, trans)

	var numAddresses = len(addresses)
	fmt.Printf("%s Total number of accounts: %v\n", now(), numAddresses)
	var positiveAddresses = 0

	publicBlockChainAPI := ethapi.NewPublicBlockChainAPI(eth.ApiBackend)

	for address := range addresses {
		if address == nil {
			fmt.Printf("??? Nil Address Found: %v\n", address)
			continue
		}

		balance, err := publicBlockChainAPI.GetBalance(nil, *address, rpc.LatestBlockNumber)
		if err != nil {
			fmt.Printf("%s Error in GetBalance(): %v\n", now(), err)
		} else if balance.Sign() > 0 {
			positiveAddresses += 1
		}
	}

	fmt.Printf("%s Number of accounts with positive balance: %v \n", now(), positiveAddresses)
	fmt.Printf("%s Exiting process.\n", now())

}


func countAccountsWithDBGetMetrics(eth *Ethereum) {
	var blockNum = 0
	var trans = 0
	var blocks = 0
	addresses := make(map[*common.Address]bool)

	var gets = int(eth.chainDb.MetricsDict()["user/gets"])
	var getsNow = 0

	var getBlockGets = 0
	var getTransactionGets = 0
	var getBalanceGets = 0

	var zeroTransactionBlocks = 0

	fmt.Printf("%s Starting to count accounts...\n", now())

	var printEvery = 100000

	for {
		blck := blockNum
		blockNum += 1

		if blck > 500000 {
			break
		}

		block := eth.blockchain.GetBlockByNumber(uint64(blck))
		if block == nil {
			break
		}

		getsNow = int(eth.chainDb.MetricsDict()["user/gets"])
		getBlockGets += getsNow - gets
		gets = getsNow

		trans = trans + block.Transactions().Len()
		if block.Transactions().Len() == 0 {
			zeroTransactionBlocks += 1
		}

		for i := 0; i < block.Transactions().Len(); i++ {
			transaction := block.Transactions()[i]

			getsNow = int(eth.chainDb.MetricsDict()["user/gets"])
			getTransactionGets += getsNow - gets
			gets = getsNow

			if transaction.Value().Sign() > 0 {
				addresses[transaction.To()] = true
			}
		}

		blocks += 1

		if blocks % printEvery == 0 {
			fmt.Printf("%s %v blocks processed. Number of transactions til block number %v: %v transactions. Number of accounts: %v\n",
				now(), blocks, blck, trans, len(addresses))
			fmt.Printf("%s Get metrics - getBlockGets: %v, getTransactionGets: %v\n",
				now(), getBlockGets, getTransactionGets)
		}
	}

	fmt.Printf("%s Total blocks processed: %v Total transactions: %v, Number of blocks with zero transactions: %v\n",
		now(), blocks, trans, zeroTransactionBlocks)

	var numAddresses = len(addresses)
	fmt.Printf("%s Total number of accounts: %v\n", now(), numAddresses)
	var positiveAddresses = 0

	publicBlockChainAPI := ethapi.NewPublicBlockChainAPI(eth.ApiBackend)

	for address := range addresses {
		if address == nil {
			fmt.Printf("??? Nil Address Found: %v\n", address)
			continue
		}

		balance, err := publicBlockChainAPI.GetBalance(nil, *address, rpc.LatestBlockNumber)
		if err != nil {
			fmt.Printf("%s Error in GetBalance(): %v\n", now(), err)
			continue
		}

		getsNow = int(eth.chainDb.MetricsDict()["user/gets"])
		getBalanceGets += getsNow - gets
		gets = getsNow

		if balance.Sign() > 0 {
			positiveAddresses += 1
		}
	}

	fmt.Printf("%s Get metrics - getBlockGets: %v, getTransactionGets: %v, getBalanceGets: %v\n",
		now(), getBlockGets, getTransactionGets, getBalanceGets)
	fmt.Printf("%s Number of accounts with positive balance: %v \n", now(), positiveAddresses)
	fmt.Printf("%s Exiting process.\n", now())

}







