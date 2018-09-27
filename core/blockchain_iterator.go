package core

import (
	"github.com/sotatek-dev/heta/types"
	"github.com/syndtr/goleveldb/leveldb"
)

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHash types.Hash
	db          *leveldb.DB
}

// Next returns next block starting from the tip
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	db := i.db

	encodedBlock, _ := db.Get(blockKey(i.currentHash[:]), nil)
	block = DeserializeBlock(encodedBlock)

	i.currentHash = block.PrevBlockHash

	return block
}
