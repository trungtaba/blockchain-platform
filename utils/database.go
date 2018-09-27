package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	dbFile      = ".data/chaindata-%s"
	dbStateFile = ".data/state-%s"
)

// CreateDB ...
func CreateDB(clientID string) (*leveldb.DB, error) {
	dbFile := fmt.Sprintf(dbFile, clientID)

	if !DbExists(dbFile) {
		log.Panicf("Cannot find blockchain db for clientID=%s", clientID)
		os.Exit(1)
	}

	return leveldb.OpenFile(dbFile, nil)
}

// CreateStateDB ...
func CreateStateDB(clientID string) (*leveldb.DB, error) {
	dbStateFile := fmt.Sprintf(dbStateFile, clientID)

	if !DbExists(dbStateFile) {
		log.Panicf("Cannot find state db for clientID=%s", clientID)
		os.Exit(1)
	}

	return leveldb.OpenFile(dbStateFile, nil)
}

// DbExists ...
func DbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
