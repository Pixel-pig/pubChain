package consensus

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"pubChain/utils"
)

const DIFFICULTY = 10

/**
 * Pow算法的结构体
 */
type Pow struct {
	Getblock BlockIterface
	Target   *big.Int
}

/**
 * @author:朱健涛 时间：2021/3/2
 * 初步实现PoW算法功能
 */
func (pow Pow) Run() ([32]byte, int64) {
	var nonce int64 = 0
	 hashBig := new(big.Int)
	for {
		hash := SetNowHash(pow.Getblock, nonce)
		hashBig = hashBig.SetBytes(hash[:])
		target := pow.Target
		result := hashBig.Cmp(target)
		if result == -1 {
			return hash, nonce
		}
		nonce++
	}
}

/**
 * 计算当区块的hash值
 */
func SetNowHash(getblock BlockIterface, nonce int64) [32]byte{
	versionByte := utils.Int2byte(getblock.GetVersion())
	timeStampByte := utils.Int2byte(getblock.GetTimeStamp())
	preHashByte := getblock.GetPreHash()
	nonceByte := utils.Int2byte(nonce)
	blockByte := bytes.Join([][]byte{
		versionByte,
		timeStampByte,
		nonceByte,
		preHashByte[:],
		getblock.GetData(),
	}, []byte{})
	return sha256.Sum256(blockByte)
}