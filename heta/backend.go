package heta

import (
<<<<<<< HEAD
	"log"
	"os"

=======
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/log"
>>>>>>> c838dea44f9ab4826534580d6b8aaecba2d35422
	"github.com/sotatek-dev/heta/trie"
	"github.com/sotatek-dev/heta/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

// Heta ...
type Heta struct {
	ChainDB   *leveldb.DB
	StateDB   *leveldb.DB
	StateTrie *trie.Trie
}

var (
	// FullNode ...
	FullNode Heta
	once     sync.Once
)

<<<<<<< HEAD
// Heta ...
type Heta struct {
	ChainDB   *leveldb.DB
	StateDB   *leveldb.DB
	StateTrie *trie.Trie
}

// New create a fullnode instance
func New() {
	clientID := os.Getenv("LOCAL_CLIENT_ID")
	chainDB, err := utils.CreateDB(clientID)
	if err != nil {
		log.Panic("Cannot open blockchain database")
	}
	stateDB, err := utils.CreateStateDB(clientID)
	if err != nil {
		log.Panic("Cannot open state database")
	}

	// make db global
	FullNode = Heta{
		ChainDB: chainDB,
		StateDB: stateDB,
	}
=======
// New create a fullnode instance
func New() Heta {
	once.Do(func() {
		clientID := os.Getenv("LOCAL_CLIENT_ID")
		chainDB, err := utils.CreateDB(clientID)
		if err != nil {
			log.Error(err.Error())
		}
		stateDB, err := utils.CreateStateDB(clientID)
		if err != nil {
			log.Error(err.Error())
		}

		// make db global
		FullNode = Heta{
			ChainDB: chainDB,
			StateDB: stateDB,
		}
		log.Debug("Create full node successfully", "client id", clientID)
	})
	return FullNode
}

// GetFullNode return the heta full node
func GetFullNode() Heta {
	return FullNode
>>>>>>> c838dea44f9ab4826534580d6b8aaecba2d35422
}
