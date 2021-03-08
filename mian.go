package main

import (
	"github.com/boltdb/bolt"
	"pubChain/chain"
)

func main() {
	//生成blot.DB对象
	db, err := bolt.Open("pubchain",0600,nil)
	if err != nil {
		panic(err.Error())
	}

	blockChain := chain.NewBlockChain(db)
	blockChain.CreateGenesis([]byte("block0"))

}
