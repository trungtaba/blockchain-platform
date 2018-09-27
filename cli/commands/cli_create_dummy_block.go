package cli

import (
	"fmt"
	"math/big"

	log "github.com/inconshreveable/log15"
	"github.com/sotatek-dev/heta/core"
)

// CreateDummyBlock ...
func (cli *CLI) CreateDummyBlock(clientID string, numOfBlocks int) {
	bc := core.NewBlockChain(clientID)
	log.Debug("Created dummy block")
	txs := createDummyTransactions(numOfBlocks)
	bc.MineBlock(txs)
}

// CreateBlockWithTransactions ...
func (cli *CLI) CreateBlockWithTransactions(clientID string, txs []*core.Transaction) {
	bc := core.NewBlockChain(clientID)
	msg := fmt.Sprintf("Created block with %d transactions", len(txs))
	log.Debug(msg)
	bc.MineBlock(txs)
}

func createDummyTransactions(numOfTransaction int) []*core.Transaction {
	msg := fmt.Sprintf("Will create %d dummy transactions.", numOfTransaction)
	log.Debug(msg)

	var txs []*core.Transaction

	for i := 0; i < numOfTransaction; i++ {
		account := core.NewAccount()
		tx := core.NewTX(&account.Key.Address, uint64(0), &account.Key.Address, new(big.Int), []byte{})
		txs = append(txs, tx)
	}

	return txs
}
