package eth

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/fraternal/ethapi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)


func now() string {
	return time.Now().Format("Jan 2 15:04:05")
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
	//fmt.Printf("%+v\n\n", backend.CurrentBlock())
}


// place the following on the first line of the function to time how long a function takes:
// defer trackTimeOfFunction(time.Now(), "~~ Count Accounts ~~")
func trackTimeOfFunction(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n\n", name, elapsed)
}


func countAccountsWithDBGetMetrics(eth *Ethereum) {
	var blockNum = 0
	var totalTransactions = 0
	var zeroTransactionBlocks = 0
	addressTransactionMap := make(map[*common.Address]bool)

	// DB metrics
	var gets = int(eth.chainDb.MetricsDict()["user/gets"])
	var getsNow = 0
	var getBlockGets = 0
	var getTransactionGets = 0
	var getBalanceGets = 0

	// Timers
	var startTimeTransactions = time.Now()
	var startTimeBlocks = time.Now()
	var startTimeAccounts = time.Now()
	var elapsedTimeTransactions = time.Since(time.Now())
	var elapsedTimeBlocks = time.Since(time.Now())
	var elapsedTimeAccounts = time.Since(time.Now())

	var printEvery = 100000

	fmt.Printf("\n%s \n\t Starting to process blocks and transactions...\n", now())

	startTimeBlocks = time.Now()

	for {
		// Increment the block number
		currentBlockNumber := blockNum
		blockNum++

		// Call ETH blockchain API to get the current block
		block := eth.blockchain.GetBlockByNumber(uint64(currentBlockNumber))
		if block == nil {
			break
		}

		// Update DB metrics for the current block
		getsNow = int(eth.chainDb.MetricsDict()["user/gets"])
		getBlockGets += getsNow - gets
		gets = getsNow

		// Get transactions in the current block
		totalTransactions = totalTransactions + block.Transactions().Len()
		if block.Transactions().Len() == 0 {
			zeroTransactionBlocks += 1
		}

		startTimeTransactions = time.Now()

		// Loop through transactions in the current block
		for i := 0; i < block.Transactions().Len(); i++ {
			transaction := block.Transactions()[i]

			// Update DB metrics for the current transaction (doesn't work because block is already in memory)
			getsNow = int(eth.chainDb.MetricsDict()["user/gets"])
			getTransactionGets += getsNow - gets
			gets = getsNow

			// Add positive value transaction receiver address to address/transaction map
			if transaction.Value().Sign() > 0 {
				addressTransactionMap[transaction.To()] = true
			}
		}

		elapsedTimeTransactions += time.Since(startTimeTransactions)
		elapsedTimeBlocks = time.Since(startTimeBlocks)

		if (blockNum + 1) % printEvery == 0 {
			fmt.Printf("\n%s \n", now())
			fmt.Printf("\t ~~~~~~~~~ %v blocks processed ~~~~~~~~ \n", blockNum + 1)
			fmt.Printf("\t Number of transactions \t %v \n", totalTransactions)
			fmt.Printf("\t Number of accounts \t\t %v \n", len(addressTransactionMap))
			fmt.Printf("\t Transaction process time \t %s \n", elapsedTimeTransactions)
			fmt.Printf("\t Block process time \t\t %s \n", elapsedTimeBlocks)
			fmt.Printf("\t DB Metrics \n\t\t Block Gets \t\t %v \n\t\t Transaction Gets \t %v \n", getBlockGets, getTransactionGets)
		}
	}

	elapsedTimeBlocks = time.Since(startTimeBlocks)
	var numAddresses = len(addressTransactionMap)

	fmt.Printf("\n%s \n", now())
	fmt.Printf("\t ~~~~~~~~~ Finished blocks and transactions ~~~~~~~~ \n")
	fmt.Printf("\t Total blocks \t\t\t\t %v \n", blockNum + 1)
	fmt.Printf("\t Total transactions \t\t\t %v \n", totalTransactions)
	fmt.Printf("\t Total accounts \t\t\t %v \n", numAddresses)
	fmt.Printf("\t Blocks with zero transactions \t\t %v \n", zeroTransactionBlocks)
	fmt.Printf("\t Total block process time \t\t %s \n", elapsedTimeBlocks)
	fmt.Printf("\t Total transaction process time \t %s \n", elapsedTimeTransactions)
	fmt.Printf("\t ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ \n")

	publicBlockChainAPI := ethapi.NewPublicBlockChainAPI(eth.ApiBackend)

	var positiveAddresses = 0
	var addressIndex = 0
	var printEveryAddress = numAddresses / 100

	fmt.Printf("\n%s \n\t Starting to process accounts...\n\n", now())

	startTimeAccounts = time.Now()

	// Loop through accounts in the address/transaction map
	for address := range addressTransactionMap {
		addressIndex++

		if address == nil {
			fmt.Printf("%s > Nil Address Found: %v\n", now(), address)
			continue
		}

		// Call ETH public blockchain API to get the current account's balance
		balance, err := publicBlockChainAPI.GetBalance(nil, *address, rpc.LatestBlockNumber)
		if err != nil {
			fmt.Printf("%s > Error in GetBalance(): %v\n", now(), err)
			continue
		}

		// Update DB metrics for the current account
		getsNow = int(eth.chainDb.MetricsDict()["user/gets"])
		getBalanceGets += getsNow - gets
		gets = getsNow

		if balance.Sign() > 0 {
			positiveAddresses += 1
		}

		// Print account traversal progress
		if addressIndex % printEveryAddress == 0 {
			fmt.Printf("%s > Processed %v/%v accounts \n", now(), addressIndex, len(addressTransactionMap))
		}
	}

	elapsedTimeAccounts = time.Since(startTimeAccounts)

	fmt.Printf("\n%s \n", now())
	fmt.Printf("\t ~~~~~~~~~~~~~~~~~ Finished accounts ~~~~~~~~~~~~~~~~ \n")
	fmt.Printf("\t Total account process time \t\t %s\n", elapsedTimeAccounts)
	fmt.Printf("\t Accounts with positive balance \t %v \n", positiveAddresses)
	fmt.Printf("\t DB Metrics \n\t\t Block Gets \t\t\t %v \n\t\t Transaction Gets \t\t %v \n\t\t Balance Gets \t\t\t %v \n",
		getBlockGets, getTransactionGets, getBalanceGets)
	fmt.Printf("\t ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ \n")

	fmt.Printf("\n%s \n\t Exiting process.\n\n", now())

}







