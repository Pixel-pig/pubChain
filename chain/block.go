package chain

import (
	"bytes"
	"encoding/gob"
	"pubChain/consensus"
	"time"
)

const VERSION = 0x01

/**
 * 区块结构
 */
type Block struct {
	Version int64
	PreHash [32]byte
	Hash    [32]byte //当前区块的hash
	//默克尔根
	TimeStamp int64
	//难度值
	Nonce int64
	Data  []byte
}

/**
 * 区块序列化
 * 将区块序列化为byte类型
 */
func (block *Block) Serialize() ([]byte, error){
	buff := new(bytes.Buffer)
	encoder := gob.NewEncoder(buff)
	err := encoder.Encode(&block)
	return buff.Bytes(), err
}

/**
 * 区块的反序列化
 */
func Deserialize(data []byte) (Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	return block, err
}

/**
 * 创建一个新区块
 */
func CreateBlock(preHash [32]byte, data []byte) Block {
	block := Block{}
	block.Version = VERSION
	block.PreHash = preHash
	block.TimeStamp = time.Now().Unix()
	block.Data = data

	//利用共识算法得到一个当前区块的Nonce和hash
	//这里使用接口实现方便之后的扩展
	cons := consensus.NewPoW(block)
	block.Hash, block.Nonce = cons.Run()
	return block
}

/**
 * 创建创世区块
 */
func CreateGenesisBloc(data []byte) Block {
	genesis := Block{}
	genesis.Version = VERSION
	genesis.PreHash = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	genesis.TimeStamp = time.Now().Unix()
	genesis.Data = data

	//利用共识算法得到一个当前区块的Nonce和hash
	//这里使用接口实现方便之后的扩展
	cons := consensus.NewPoW(genesis)
	genesis.Hash, genesis.Nonce  = cons.Run()
	return genesis
}

/**
 * 获取block字段的get，set方法
 */
func (block Block) GetVersion() int64{
	return block.Version
}

func (block Block) GetPreHash() [32]byte{
	return block.PreHash
}

func (block Block) GetTimeStamp() int64{
	return block.TimeStamp
}

func (block Block) GetData() []byte{
	return block.Data
}

