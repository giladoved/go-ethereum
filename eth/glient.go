package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/fraternal/ethapi"
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

