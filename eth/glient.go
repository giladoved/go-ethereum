package eth

import (
	"fmt",
	""
)

func testA(eth *Ethereum) {

	//fmt.Printf("%v, %v", i, v)
	pubServAPI := PublicEthereumAPI{eth}
	pubBlockChainApi := NewPublicBlockChainAPI(eth)
	BlockNumber(eth)
	ethAd, _ := pubServAPI.Etherbase()
	coinAd, _ := pubServAPI.Coinbase()
	pubAd := pubServAPI.Hashrate()
	fmt.Printf("%v\n%v\n%v", ethAd, coinAd, pubAd)

	//fmt.Printf("~~~~~~~~~~~~~~~~~~~\n\n\n\n")
	//fmt.Printf("%+v\n\n", backend.CurrentBlock())
	//fmt.Printf("damn")
	//
}

//func testB(somePublicAPI *) {
//
//
//
//}
