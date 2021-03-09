package chain

import (
	"bytes"
	"errors"
	"github.com/boltdb/bolt"
)

const BLOCKS = "blocks"
const LASTHASH = "lasthash"

/**
 * 定义区块链的结构体， 现在存在在内存中 Blocks []Block
 * @Author:朱健涛 time:2021/3/4
 * 持久化存储到介质中
 */
type BlockChain struct {
	//Blocks []Block
	//文件操作对象
	DB                *bolt.DB
	LastBlock         Block    //最新的区块
	IteratorBlockHash [32]byte //迭代器迭代到的区块
}

/**
 * 创建区块链对象既实例化BlockChain
 */
func NewBlockChain(db *bolt.DB) (BlockChain, error) {
	var lastBlock Block
	var err error
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("bolt数据库中bucket不存在")
			return err
		}

		//那到最后一个区块
		lastHash := bucket.Get([]byte(LASTHASH))
		BolckByte := bucket.Get(lastHash)
		lastBlock, err = Deserialize(BolckByte)
		if err != nil {
			err = errors.New("区块反序列化错误")
			return err
		}
		return nil
	})
	return BlockChain{DB: db, LastBlock: lastBlock, IteratorBlockHash: lastBlock.Hash}, err
}

/**
 * 生成创世区块并添加到bolt文件中
 */
func (chain *BlockChain) CreateGenesis(genesisData []byte) {
	genesis := CreateGenesisBloc(genesisData)
	genesisByte, _ := genesis.Serialize()

	db := chain.DB
	_ = db.Update(func(tx *bolt.Tx) error {
		//创建数据空间(实例化)（获取数据空间）
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil { //该空间不存在，第一次创建数据空间
			bucket, _ = tx.CreateBucket([]byte(BLOCKS))
			chain.LastBlock = genesis
			chain.IteratorBlockHash = chain.LastBlock.Hash
		}
		if bucket != nil {
			lastHash := bucket.Get([]byte(LASTHASH))
			if len(lastHash) == 0 { //空间中没有数据
				//存储区块
				_ = bucket.Put(genesis.Hash[:], genesisByte)
				//记录最后一个区块hash(更新lastHash)
				_ = bucket.Put([]byte(LASTHASH), genesis.Hash[:])
			}
		}

		return nil
	})
}

/**
 * 向链上追加区块
 */
func (chain *BlockChain) AddNewBlock(data []byte) error {
	db := chain.DB
	lastBlock := chain.LastBlock
	var err error

	/* 生成一个新区块,并序列化, 更新链上的最后得一个区块 */
	NewBlock := CreateBlock(lastBlock.Hash, data)
	blockByte, _ := NewBlock.Serialize()

	/* 向bolt文件中添加新区块*/
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("bolt数据库中bucket不存在")
			return err
		}
		//存储区块
		_ = bucket.Put(NewBlock.Hash[:], blockByte)
		//记录最后一个区块hash(更新lastHash)
		_ = bucket.Put([]byte(LASTHASH), NewBlock.Hash[:])
		chain.LastBlock = NewBlock
		chain.IteratorBlockHash = chain.LastBlock.Hash
		return nil
	})
	if err != nil {
		return err
	}

	return err
}

/**
 * 获取最新的区块
 */
func (chain *BlockChain) GetLastBlock() Block {
	return chain.LastBlock
}

/**
 * 获取bolt文件中的所有区块
 */
func (chain *BlockChain) GetAllBlocks() ([]Block, error) {
	db := chain.DB
	blocks := make([]Block, 0)
	genesisHash := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var err error
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("bolt数据库中bucket不存在")
			return err
		}

		/* 先获取lastBlock */
		lastHash := chain.LastBlock.Hash
		/* 中间变量用于循环 */
		var currentBlockhash []byte
		currentBlockhash = lastHash[:] //将最后一个区块hash设置为当前区块的hash
		/* 循环添加数据区块 */
		for {
			currentBlockBytes := bucket.Get(currentBlockhash)
			currentBlock, err := Deserialize(currentBlockBytes)
			if err != nil {
				err = errors.New("区块反序列化错误")
				break
			}
			blocks = append(blocks, currentBlock)
			currentBlockhash = currentBlock.PreHash[:]

			if bytes.Compare(currentBlockhash, genesisHash[:]) == 0 {
				break
			}
		}
		return err
	})
	return blocks, err
}

/**
 * 迭代器 (实现迭代器接口的方法)
 */
/* 用于判断是否还有区块 */
func (chain *BlockChain) HasNext() bool {
	//查看当前区块是否是创世区块
	db := chain.DB
	var hasNext bool
	genesisHash := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	_ = db.View(func(tx *bolt.Tx) error {
		currentBlockHash := chain.IteratorBlockHash
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败")
		}
		currentBlockByte := bucket.Get(currentBlockHash[:])
		currentBlock, _ := Deserialize(currentBlockByte)
		if bytes.Compare(currentBlock.Hash[:], genesisHash[:] ) == 0 {
			hasNext = false
		}else {
			hasNext = true
		}
		return nil
	})
	return hasNext
}

/* 去除下一个区块(既将得到的区块更新迭代blockhash) */
func (chain *BlockChain) Next() Block {
	db := chain.DB
	var currentBlock Block
	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败")
		}

		currentBlockByte := bucket.Get(chain.IteratorBlockHash[:])
		currentBlock, _ = Deserialize(currentBlockByte)
		//下一个区块的hash(当前区块以获取，需获取下一个区块)
		chain.IteratorBlockHash = currentBlock.PreHash

		return nil
	})
	return currentBlock
}
