package main

import (
	"fmt"
	"pubChain/chain"
)

func main() {
	genesisBloc := chain.CreateGenesisBloc([]byte("hello word!"))
	fmt.Println("创世区块", genesisBloc)
	block1 := chain.CreateBloc(genesisBloc.Hash, []byte("hello word!"))
	fmt.Println("第一个区块", block1)
}
