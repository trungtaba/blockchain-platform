package core

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/sotatek-dev/heta/heta"
	"github.com/sotatek-dev/heta/types"
	"github.com/sotatek-dev/heta/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	dbFile                = ".data/chaindata-%s"
	dbStateFile           = ".data/state-%s"
	blockPrefix           = []byte("b")
	blockNumberHashPrefix = []byte("bnh")
	latestBlockPrefix     = []byte("l")
	txLookupPrefix        = []byte("txl")
	accountLookupPrefix   = []byte("acc")
	candidatePrefix       = []byte("c")
)

const PRODUCER_REWARD uint64 = 10

// BlockChain zzz
type BlockChain struct {
	tip types.Hash
	db  *leveldb.DB
}

// TxLookupEntry is a positional metadata to help looking up the data content of
// a transaction or receipt given only its hash.
type TxLookupEntry struct {
	BlockHash  types.Hash
	BlockIndex uint64
	Index      uint64
}

// encodeBlockNumber encodes a block number as big endian uint64
func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// headerHashKey = headerPrefix + num (uint64 big endian) + headerHashSuffix
func blockKey(hash []byte) []byte {
	return append(blockPrefix, hash...)
}

func blockNumberHashKey(number uint64) []byte {
	return append(blockNumberHashPrefix, encodeBlockNumber(number)...)
}

func latestBlockKey() []byte {
	return latestBlockPrefix
}

// txLookupKey = txLookupPrefix + hash
func txLookupKey(hash types.Hash) []byte {
	return append(txLookupPrefix, hash[:]...)
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address, clientID string) *BlockChain {
	dbFile := fmt.Sprintf(dbFile, clientID)
	if dbExists(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip types.Hash

	privateKey := os.Getenv("NODE_PRIVATE_KEY")
	account := NewAccountFromPrivateKey(privateKey)
	cbtx := NewCoinbaseTX(&account.Key.Address, []byte("Random coinbase tx."))
	genesis := NewGenesisBlock(cbtx, account)

	db, err := leveldb.OpenFile(dbFile, nil)
	if err != nil {
		log.Error(err.Error())
	}
	defer db.Close()

	stateDbFile := fmt.Sprintf(dbStateFile, clientID)
	if dbExists(stateDbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	stateDb, err := leveldb.OpenFile(stateDbFile, nil)
	if err != nil {
<<<<<<< HEAD
		log.Panic(err)
=======
		log.Error(err.Error())
>>>>>>> c838dea44f9ab4826534580d6b8aaecba2d35422
	}
	defer stateDb.Close()

	err = db.Put(blockKey(genesis.Hash[:]), genesis.Serialize(), nil)
	if err != nil {
		log.Error(err.Error())
	}

	err = db.Put(blockNumberHashKey(genesis.GetHeight()), genesis.Hash[:], nil)
	if err != nil {
		log.Error(err.Error())
	}

	err = db.Put(latestBlockKey(), genesis.Hash[:], nil)
	if err != nil {
		log.Error(err.Error())
	}

	WriteTxLookupEntries(db, genesis)
	WriteGenesisAccountLookupEntries(stateDb, genesis)

	tip = genesis.Hash

	bc := BlockChain{tip, db}
	return &bc
}

// NewBlockChain creates a new Blockchain with genesis Block
func NewBlockChain(clientID string) *BlockChain {
	dbFile := fmt.Sprintf(dbFile, clientID)
	if dbExists(dbFile) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip types.Hash
	db := heta.FullNode.ChainDB

	latestHash, _ := db.Get(latestBlockKey(), nil)

	copy(tip[:], latestHash[:HashLength])

	bc := BlockChain{tip, db}

	return &bc
}

// AddBlock saves the block into the blockchain
func (bc *BlockChain) AddBlock(block *Block) bool {
	db := bc.db
	blockInDb, _ := db.Get(blockKey(block.Hash[:]), nil)
	if blockInDb != nil {
		return false
	}

	blockData := block.Serialize()
	err := db.Put(blockKey(block.Hash[:]), blockData, nil)
	if err != nil {
		log.Error(err.Error())
	}
	err = db.Put(blockNumberHashKey(block.GetHeight()), block.Hash[:], nil)
	if err != nil {
		log.Error(err.Error())
	}

<<<<<<< HEAD
	fmt.Println("BLOCK STATE HASH: ", hex.EncodeToString(block.StateHash[:]))
=======
>>>>>>> c838dea44f9ab4826534580d6b8aaecba2d35422
	WriteTxLookupEntries(db, block)
	statedb := GetStateDB(nil)
	WriteAccountLookupEntries(statedb, block)
	ProducerReward(statedb, &block.Producer)

	pullBlockHash := hex.EncodeToString(block.StateHash[:])
	calculatedHash := hex.EncodeToString(heta.FullNode.StateTrie.Hash(heta.FullNode.StateTrie.Root))
<<<<<<< HEAD
	fmt.Println("CALCULATED STATE HASH: ", calculatedHash)
=======
	log.Debug("Caculate state hash", "block state hash", pullBlockHash, "caculate state hash", calculatedHash)
>>>>>>> c838dea44f9ab4826534580d6b8aaecba2d35422
	if pullBlockHash != calculatedHash {
		panic(fmt.Errorf("Error (1801) when downloading block, replaying blocks required "))
	}

	lastHash, _ := db.Get(latestBlockKey(), nil)
	lastBlockData, _ := db.Get(blockKey(lastHash), nil)
	lastBlock := DeserializeBlock(lastBlockData)

	if block.Height > lastBlock.Height {
		err = db.Put(latestBlockKey(), block.Hash[:], nil)
		if err != nil {
			log.Error(err.Error())
		}
		bc.tip = block.Hash
	}
	return true
}

// Iterator ...
func (bc *BlockChain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// GetDB TODO: move this method to another place
func (bc *BlockChain) GetDB() *leveldb.DB {
	return bc.db
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// MineBlock mines a new block with the provided transactions
func (bc *BlockChain) MineBlock(transactions []*Transaction) *Block {
	var lastHash types.Hash
	var lastHeight uint64
	var stateHash types.Hash
	statedb := GetStateDB(nil)

	db := bc.db
	data, _ := db.Get(latestBlockKey(), nil)
	copy(lastHash[:], data[:HashLength])

	blockData, _ := db.Get(blockKey(lastHash[:]), nil)
	block := DeserializeBlock(blockData)
	lastHeight = block.Height
	newBlock := NewBlock(transactions, lastHeight+1, lastHash)

	ProducerReward(statedb, &newBlock.Producer)
	WriteTxLookupEntries(db, newBlock)
	WriteAccountLookupEntries(statedb, newBlock)

	rootHash := heta.FullNode.StateTrie.Hash(heta.FullNode.StateTrie.Root)
	copy(stateHash[:], rootHash[:HashLength])
	newBlock.StateHash = stateHash

	err := db.Put(blockKey(newBlock.Hash[:]), newBlock.Serialize(), nil)
	if err != nil {
		log.Error(err.Error())
	}

	err = db.Put(blockNumberHashKey(newBlock.GetHeight()), newBlock.Hash[:], nil)
	if err != nil {
		log.Error(err.Error())
	}

	err = db.Put(latestBlockKey(), newBlock.Hash[:], nil)
	if err != nil {
		log.Error(err.Error())
	}

	bc.tip = newBlock.Hash

	return newBlock
}

// InfoResponse TODO: move this struct to another place
type InfoResponse struct {
	ChainID               string `json:"chainID"`
	ServerVersion         string `json:"version"`
	LastIrreversibleBlock uint64 `json:"lastIrreversibleBlock"`
	LastBlock             uint64 `json:"lastBlock"`
	LastBlockHash         string `json:"lastBlockHash"`
	LastBlockTime         uint64 `json:"lastBlockTime"`
	LastBlockProducer     string `json:"lastBlockProducer"`
}

// PeersResponse TODO: move this struct to another place
type PeersResponse struct {
	ChainID       string `json:"addr"`
	ServerVersion string `json:"version"`
}

// TotalStakedResponse TODO: move this struct to another place
type TotalStakedResponse struct {
	User  string `json:"totalUsers"`
	Coins string `json:"totalCoins"`
}

// ListProducersResponse TODO: move this struct to another place
type ListProducersResponse struct {
	AddRess string `json:"address"`
	Votes   string `json:"votes"`
}

// HandleChainInfoReq ...
func HandleChainInfoReq(b []byte) InfoResponse {
	return GetChainInfo()
}

// HandlePeersReq ...
func HandlePeersReq(b []byte) PeersResponse {
	return GetPeersInfo()
}

// HandleTotalStaked ...
func HandleTotalStaked(b []byte) TotalStakedResponse {
	return GetTotalStaked()
}

// HandleLitsProducers ...
func HandleListProducers(b []byte) ListProducersResponse {
	return GetListProducers()
}

// GetChainInfo TODO: move this method to another place
func GetChainInfo() InfoResponse {
	var lastHash types.Hash
	var lastHeight uint64

	db := heta.FullNode.ChainDB

	data, err := db.Get(latestBlockKey(), nil)
	if err != nil {
		return InfoResponse{
			"1",
			"1.00.00",
			0,
			0,
			"",
			uint64(time.Now().UnixNano() / int64(time.Millisecond)),
			"sotatek",
		}
	}
	copy(lastHash[:], data[:HashLength])

	blockData, _ := db.Get(blockKey(lastHash[:]), nil)
	block := DeserializeBlock(blockData)

	lastHeight = block.GetHeight()
	lastBlockTime := uint64(block.GetTimestamp())

	lastProducer := block.GetProducer()

	return InfoResponse{
		"1",
		"1.00.00",
		0,
		lastHeight,
		hex.EncodeToString(lastHash[:]),
		lastBlockTime,
		lastProducer.ToString(),
	}
}

// GetChainInfo TODO: move this method to another place
func GetPeersInfo() PeersResponse {
	return PeersResponse{
		"123.13.252.47:50276",
		"1.00.00",
	}
}

// GetTotalStaked TODO: move this method to another place
func GetTotalStaked() TotalStakedResponse {
	return TotalStakedResponse{
		"544",
		"100200300",
	}
}

// GetChainInfo TODO: move this method to another place
func GetListProducers() ListProducersResponse {
	return ListProducersResponse{
		"0xecf73def92e8935a0758716d3cda19e11697fe04",
		"12356",
	}
}

// GetBlockByNumber ...
func GetBlockByNumber(number uint64) *Block {
	db := GetDB(&opt.Options{ReadOnly: true})
	data, _ := db.Get(blockNumberHashKey(number), nil)
	hash := types.Hash{}
	if len(data) != 0 {
		hash = utils.BytesToHash(data)
	}

	data, _ = db.Get(blockKey(hash[:]), nil)
	return DeserializeBlock(data)
}

// GetBlockByHash ...
func GetBlockByHash(hash types.Hash) *Block {
	db := GetDB(&opt.Options{ReadOnly: true})
	data, _ := db.Get(blockKey(hash[:]), nil)
	return DeserializeBlock(data)
}

// WriteTxLookupEntries stores a positional metadata for every transaction from
// a block, enabling hash based transaction and receipt lookups.
func WriteTxLookupEntries(db *leveldb.DB, block *Block) {
	for i, tx := range block.GetTransactions() {
		entry := TxLookupEntry{
			BlockHash:  block.GetHash(),
			BlockIndex: block.GetHeight(),
			Index:      uint64(i),
		}
		data := utils.Serialize(entry)

		var txHash types.Hash
		txHash.SetBytes(tx.Hash())

		if err := db.Put(txLookupKey(txHash), data, nil); err != nil {
			log.Error("Failed to store transaction lookup entry", "err", err.Error())
		}
	}
}

// ProducerReward ...
func ProducerReward(db *leveldb.DB, producer *types.Address) {
	reward := big.NewInt(int64(PRODUCER_REWARD))
	AddBalance(db, producer, reward)
}
