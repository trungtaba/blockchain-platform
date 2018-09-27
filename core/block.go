package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/sotatek-dev/heta/keystore"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sotatek-dev/heta/types"
	"github.com/sotatek-dev/heta/utils"
)

// Block ...
type Block struct {
	Timestamp     uint64
	Transactions  []*Transaction
	Height        uint64
	StateHash     types.Hash
	Hash          types.Hash
	Producer      types.Address
	PrevBlockHash types.Hash
}

// GetTimestamp is getter for Timestamp
func (b *Block) GetTimestamp() uint64 { return b.Timestamp }

// GetTransactions is getter for Transactions
func (b *Block) GetTransactions() []*Transaction { return b.Transactions }

// GetHeight is getter for Height
func (b *Block) GetHeight() uint64 { return b.Height }

// GetHash is getter for Hash
func (b *Block) GetHash() types.Hash { return b.Hash }

// GetProducer is getter for Transactions
func (b *Block) GetProducer() types.Address { return b.Producer }

// GetPrevBlockHash is getter for PrevBlockHash
func (b *Block) GetPrevBlockHash() types.Hash { return b.PrevBlockHash }

// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the block.
func NewBlock(txs []*Transaction, height uint64, prevBlockHash types.Hash) *Block {
	now := time.Now().Unix()
	nodeKey := os.Getenv("NODE_PRIVATE_KEY")
	producerKey := keystore.NewKeyFromPrivateKey(nodeKey)
	producer := producerKey.Address

	block := &Block{
		Timestamp:     uint64(now),
		Transactions:  txs,
		Height:        height,
		Hash:          [32]byte{},
		Producer:      producer,
		PrevBlockHash: prevBlockHash}

	data := bytes.Join(
		[][]byte{
			block.HashTransactions(),
			utils.IntToHex(int64(block.Timestamp)),
		},
		[]byte{},
	)

	var hashInt big.Int
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	block.Hash = hash

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction, producer *Account) *Block {
	genesisHeight := uint64(0)
	return NewBlock([]*Transaction{coinbase}, genesisHeight, [32]byte{})
}

// HashTransactions ...
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	data, err := rlp.EncodeToBytes(b)
	if err != nil {
		log.Panic(err)
	}
	return data
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	if len(d) == 0 {
		return &Block{}
	}
	block := new(Block)
	if err := rlp.Decode(bytes.NewReader(d), block); err != nil {
		log.Panic("Invalid block body RLP", err)
		return &Block{}
	}
	return block
}

// BlockRes ...
type BlockRes struct {
	Height          uint64   `json:"number"`
	Producer        string   `json:"producer"`
	BlockHash       string   `json:"hash"`
	PreviousHash    string   `json:"parentHash"`
	Timestamp       uint64   `json:"timestamp"`
	TransactionRoot string   `json:"transactionRoot"`
	Transaction     []*TxRes `json:"transactions"`
}

// GetBlock ...
func GetBlock(blockNumber uint64) BlockRes {
	blockRes := GetBlockByNumber(blockNumber)
	producer := blockRes.GetProducer()
	hash := blockRes.GetHash()
	prevHash := blockRes.GetPrevBlockHash()
	txs := blockRes.GetTransactions()
	txRess := []*TxRes{}
	for _, tx := range txs {
		txRes := tx.ToTxRes()
		txRess = append(txRess, &txRes)
	}

	return BlockRes{
		Height:       uint64(blockNumber),
		Producer:     producer.ToString(),
		BlockHash:    hex.EncodeToString(hash[:]),
		PreviousHash: hex.EncodeToString(prevHash[:]),
		Timestamp:    blockRes.GetTimestamp(),
		Transaction:  txRess,
	}
}

// GetBlock ByHash...
func GetBlockByHashReq(blockHash string) BlockRes {
	decodeBlockHash, err := hex.DecodeString(blockHash)
	if err != nil {
		log.Panic(err)
	}
	blockHashReq := types.Hash{}
	blockHashReq.SetBytes(decodeBlockHash)
	blockResHash := GetBlockByHash(blockHashReq)
	producer := blockResHash.GetProducer()
	prevHash := blockResHash.GetPrevBlockHash()
	txs := blockResHash.GetTransactions()
	txRess := []*TxRes{}
	for _, tx := range txs {
		txRes := tx.ToTxRes()
		txRess = append(txRess, &txRes)
	}

	return BlockRes{
		Height:       blockResHash.GetHeight(),
		Producer:     producer.ToString(),
		BlockHash:    hex.EncodeToString(blockHashReq[:]),
		PreviousHash: hex.EncodeToString(prevHash[:]),
		Timestamp:    blockResHash.GetTimestamp(),
		Transaction:  txRess,
	}
}

// BlockRequestFormat ...
type BlockRequestFormat struct {
	BlockNumber uint64 `json:"blockNumber"`
}

// HandleBlockInfoReq ...
func HandleBlockInfoReq(b []byte) BlockRes {
	reqJSON := BlockRequestFormat{}
	err := json.Unmarshal(b[:], &reqJSON)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	return GetBlock(reqJSON.BlockNumber)
}

// BlockRequestFormatByHash ...
type BlockRequestFormatByHash struct {
	BlockHash string `json:"blockHash"`
}

// HandleBlockInfoByHashReq ...
func HandleBlockInfoByHashReq(b []byte) BlockRes {
	reqJSON := BlockRequestFormatByHash{}
	err := json.Unmarshal(b[:], &reqJSON)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	return GetBlockByHashReq(reqJSON.BlockHash)
}
