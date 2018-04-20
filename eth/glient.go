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


// place the following on the first line of the function to time how long a function takes:
// defer trackTimeOfFunction(time.Now(), "~~ Count Accounts ~~")
func trackTimeOfFunction(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n\n", name, elapsed)
}


func countAccounts(eth *Ethereum) {
	var blockNum = 0
	var trans = 0
	var blocks = 0
	addresses := make(map[*common.Address]bool)

	fmt.Printf("%s Starting to count accounts...\n", now())

	var printEvery = 100000

	var startTimeTransactions = time.Now()
	var startTimeBlocks = time.Now()
	var startTimeAccounts = time.Now()
	var elapsedTimeTransactions = time.Since(time.Now())
	var elapsedTimeBlocks = time.Since(time.Now())
	var elapsedTimeAccounts = time.Since(time.Now())

	startTimeBlocks = time.Now()
	for {
		blck := blockNum
		blockNum += 1

		//if (blck > 200000) {
		//    break
		//}

		block := eth.blockchain.GetBlockByNumber(uint64(blck))
		if block == nil {
			break
		}

		trans = trans + block.Transactions().Len()

		startTimeTransactions = time.Now()
		for i := 0; i < block.Transactions().Len(); i++ {
			transaction := block.Transactions()[i]
			if transaction.Value().Sign() > 0 {
				addresses[transaction.To()] = true
			}
		}
		elapsedTimeTransactions += time.Since(startTimeTransactions)

		blocks += 1

		if blocks % printEvery == 0 {
			fmt.Printf("%s Time to process blocks: %s. Time to process transactions: %s\n", now(), elapsedTimeBlocks, elapsedTimeTransactions)
			fmt.Printf("%s %v blocks processed. Number of transactions til block number %v: %v transactions. Number of accounts: %v\n",
				now(), blocks, blck, trans, len(addresses))
		}
	}
	elapsedTimeBlocks = time.Since(startTimeBlocks)
	fmt.Printf("Total time to process blocks: %s Total time to process transactions: %s\n", elapsedTimeBlocks, elapsedTimeTransactions)
	fmt.Printf("Total blocks processed: %v Total transactions: %v\n", blocks, trans)

	var numAddresses = len(addresses)
	fmt.Printf("%s Total number of accounts: %v\n", now(), numAddresses)
	var positiveAddresses = 0

	publicBlockChainAPI := ethapi.NewPublicBlockChainAPI(eth.ApiBackend)

	startTimeAccounts = time.Now()
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
	elapsedTimeAccounts = time.Since(startTimeAccounts)

	fmt.Printf("%s Total time to process accounts: %s\n", time.Now().String(), elapsedTimeAccounts)
	fmt.Printf("%s Number of accounts with positive balance: %v \n", time.Now().String(), positiveAddresses)
	fmt.Printf("%s Exiting process.\n", time.Now().String())
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

	var startTimeTransactions = time.Now()
	var startTimeBlocks = time.Now()
	var startTimeAccounts = time.Now()
	var elapsedTimeTransactions = time.Since(time.Now())
	var elapsedTimeBlocks = time.Since(time.Now())
	var elapsedTimeAccounts = time.Since(time.Now())

	var printEvery = 100000

	startTimeBlocks = time.Now()
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

		startTimeTransactions = time.Now()
		for i := 0; i < block.Transactions().Len(); i++ {
			transaction := block.Transactions()[i]

			getsNow = int(eth.chainDb.MetricsDict()["user/gets"])
			getTransactionGets += getsNow - gets
			gets = getsNow

			if transaction.Value().Sign() > 0 {
				addresses[transaction.To()] = true
			}
		}
		elapsedTimeTransactions += time.Since(startTimeTransactions)

		blocks += 1

		if blocks % printEvery == 0 {
			fmt.Printf("%s Time to process blocks: %s. Time to process transactions: %s\n", now(), elapsedTimeBlocks, elapsedTimeTransactions)
			fmt.Printf("%s %v blocks processed. Number of transactions til block number %v: %v transactions. Number of accounts: %v\n",
				now(), blocks, blck, trans, len(addresses))
			fmt.Printf("%s Get metrics - getBlockGets: %v, getTransactionGets: %v\n",
				now(), getBlockGets, getTransactionGets)
		}
	}
	elapsedTimeBlocks = time.Since(startTimeBlocks)
	fmt.Printf("Total time to process blocks: %s Total time to process transactions: %s\n", elapsedTimeBlocks, elapsedTimeTransactions)
	fmt.Printf("%s Total blocks processed: %v Total transactions: %v, Number of blocks with zero transactions: %v\n",
		now(), blocks, trans, zeroTransactionBlocks)

	var numAddresses = len(addresses)
	fmt.Printf("%s Total number of accounts: %v\n", now(), numAddresses)
	var positiveAddresses = 0

	publicBlockChainAPI := ethapi.NewPublicBlockChainAPI(eth.ApiBackend)

	startTimeAccounts = time.Now()
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
	elapsedTimeAccounts = time.Since(startTimeAccounts)

	fmt.Printf("%s Total time to process accounts: %s\n", time.Now().String(), elapsedTimeAccounts)
	fmt.Printf("%s Get metrics - getBlockGets: %v, getTransactionGets: %v, getBalanceGets: %v\n",
		now(), getBlockGets, getTransactionGets, getBalanceGets)
	fmt.Printf("%s Number of accounts with positive balance: %v \n", now(), positiveAddresses)
	fmt.Printf("%s Exiting process.\n", now())

}


func countTransactions(eth *Ethereum) {
	var blockNum = 0
	var trans = 0
	var blocks = 0

	fmt.Printf("%s Starting to count transactions...\n", now())

	var startTimeTransactions = time.Now()
	var startTimeBlocks = time.Now()
	var elapsedTimeTransactions = time.Since(time.Now())
	var elapsedTimeBlocks = time.Since(time.Now())

	var printEvery = 10000

	startTimeBlocks = time.Now()
	for {
		blck := blockNum
		blockNum += 1

		block := eth.blockchain.GetBlockByNumber(uint64(blck))
		if block == nil {
			break
		}

		startTimeTransactions = time.Now()
		trans = trans + block.Transactions().Len()
		elapsedTimeTransactions = time.Since(startTimeTransactions)

		blocks += 1

		if blocks % printEvery == 0 {
			fmt.Printf("%s %v blocks processed. Number of transactions til block number %v: %v transactions.\n",
				now(), blocks, blck, trans)
		}
	}
	elapsedTimeBlocks = time.Since(startTimeBlocks)
	fmt.Printf("Total time to process blocks: %s Total time to process transactions: %s\n", elapsedTimeBlocks, elapsedTimeTransactions)
	fmt.Printf("Total blocks processed: %v Total transactions: %v\n", blocks, trans)
	fmt.Printf("%s Exiting process.\n", now())

}




