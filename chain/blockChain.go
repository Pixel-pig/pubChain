package chain

import "github.com/boltdb/bolt"

const BLOCKS  = "blocks"
const LASTHASH  = "lasthash"

/**
 * 定义区块链的结构体， 现在存在在内存中 Blocks []Block
 * @Author:朱健涛 time:2021/3/4
 * 持久化存储到介质中
 */
type BlockChain struct {
	//Blocks []Block
	//文件操作对象
	lastHash [32]byte
	DB *bolt.DB
}

/**
 * 创建区块链对象既实例化BlockChain
 */
func NewBlockChain(db *bolt.DB) BlockChain {
	return BlockChain{DB: db}
}

/**
 * 生成创世区块并添加到bolt文件中
 */
func (chain *BlockChain) CreateGenesis(genesisData []byte)  {
	genesis := CreateGenesisBloc(genesisData)
	genesisByte, _ := genesis.Serialize()

	db := chain.DB
	_ = db.Update(func(tx *bolt.Tx) error {
		//创建数据空间(实例化)
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil { //该空间不存在
			bucket, _ = tx.CreateBucket([]byte(BLOCKS))
		}
		err := bucket.Put(genesis.Hash[:], genesisByte) //存储区块
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(LASTHASH),genesis.Hash[:]) //记录最后一个区块hash
		if err != nil {
			return err
		}
		chain.lastHash = genesis.Hash
		return nil
	})
}

/**
 * 向链上追加区块
 */
func (chain *BlockChain) AddNewBlock(data []byte) {
	db := chain.DB
	_ = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))

		//创建一个新区块
		block := CreateBlock(chain.lastHash, data)
		blockByte, _ := block.Serialize()

		err := bucket.Put(block.Hash[:], blockByte)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(LASTHASH), block.Hash[:])
		if err != nil {
			return err
		}
		chain.lastHash = block.Hash
		return nil
	})

}
