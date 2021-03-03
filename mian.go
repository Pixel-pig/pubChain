package main

import (
	"fmt"
	"pubChain/chain"
)

func main() {
	genesisBloc := chain.CreateGenesisBloc([]byte("hello word!"))
	fmt.Println("创世区块", genesisBloc)
}
