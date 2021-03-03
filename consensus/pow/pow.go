package pow

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"pubChain/chain"
	"pubChain/utils"
)

const DIFFICULTY = 10

/**
 * Pow算法的结构体
 */
type Pow struct {
	Block  chain.Block
	Target *big.Int
}

/**
 * @author:朱健涛 时间：2021/3/2
 * 初步实现PoW算法功能
 */
func (pow Pow) Run() ([32]byte, int64) {
	var nonce int64 = 0
	for {
		hash := SetNowHash(pow.Block, nonce)
		target := pow.Target
		result := bytes.Compare(hash[:], target.Bytes())
		if result == -1 {
			return hash, nonce
		}
		nonce++
	}
}

/**
 * 计算当区块的hash值
 */
func SetNowHash(block chain.Block, nonce int64) [32]byte{
	versionByte := utils.Int2byte(block.Version)
	timeStampByte := utils.Int2byte(block.TimeStamp)
	nonceByte := utils.Int2byte(nonce)
	blockByte := bytes.Join([][]byte{
		versionByte,
		timeStampByte,
		nonceByte,
		block.PreHash[:],
		block.Data[:],
	}, []byte{})
	return sha256.Sum256(blockByte)
}