package consensus

import (
	"math/big"
)

/**
 * 共识机制的标准化接口
 */
type Consensus interface {
	Run() ([32]byte, int64)
}

/**
 * 获取block字段的接口
 */
type BlockIterface interface {
	GetVersion() int64
	GetPreHash() [32]byte
	GetTimeStamp() int64
	GetData() []byte
}

/**
 * pow共识算法
 */
func NewPoW(block BlockIterface) Consensus {
	target := big.NewInt(1)
	target.Lsh(target, 255-DIFFICULTY)
	return Pow{Getblock: block, Target: target}
}
