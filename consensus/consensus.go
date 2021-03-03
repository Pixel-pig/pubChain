package consensus

import (
	"math/big"
	"pubChain/chain"
	"pubChain/consensus/pow"
)

/**
 * 共识机制的标准化接口
 */
type Consensus interface {
	Run() ([32]byte, int64)
}

/**
 * pow共识算法
 */
func NewPoW(block chain.Block) Consensus {
	target := big.NewInt(1)
	target.Lsh(target, 255-pow.DIFFICULTY)
	return pow.Pow{Block: block, Target: target}
}
