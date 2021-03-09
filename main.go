package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"pubChain/chain"
)

func main() {
	//生成blot.DB对象
	db, err := bolt.Open("pubchain",0600,nil)
	if err != nil {
		panic(err.Error())
	}

	blockChain, _ := chain.NewBlockChain(db)

	for blockChain.HasNext() {
		block := blockChain.Next()
		fmt.Println(block)
		fmt.Println(block.Hash)
		fmt.Println("--------")
	}

}
