package core

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/sotatek-dev/heta/heta"
	"github.com/sotatek-dev/heta/types"
	"github.com/sotatek-dev/heta/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Voter ...
type Voter struct {
	Address     types.Address
	VoteBalance *big.Int
	Timestamp   uint64
}

// Candidate ...
type Candidate struct {
	Timestamp uint64
	Address   string
	Vote      *big.Int
	Voters    []Voter
}

// headerHashKey = headerPrefix + num (uint64 big endian) + headerHashSuffix
func candidateKey(address string) []byte {
	return append(candidatePrefix, []byte(address)...)
}

// UpdateCandidate ...
func UpdateCandidate(address string, candidate Candidate) {
	db := GetDB(nil)
	if err := db.Put(candidateKey(address), utils.Serialize(candidate), nil); err != nil {
		log.Panic("Failed to update candidate info", "err", err)
	}
}

// GetCandidateNumber ...
func GetCandidateNumber() uint64 {
	db := GetDB(nil)
	iter := db.NewIterator(util.BytesPrefix(candidatePrefix), nil)
	number := uint64(0)
	for iter.Next() {
		number++
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Panic("GetCandidateNumber", "err", err)
	}
	return number
}

// GetCandidates ...
func GetCandidates() []*Candidate {
	db := GetDB(&opt.Options{ReadOnly: true})
	iter := db.NewIterator(util.BytesPrefix(candidatePrefix), nil)
	var candidates []*Candidate
	for iter.Next() {
		data := iter.Value()
		candidate := new(Candidate)
		utils.Deserialize(data, candidate)
		candidates = append(candidates, candidate)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Panic("GetCandidateNumber", "err", err)
	}
	return candidates
}

// GetDB ...
func GetDB(options *opt.Options) *leveldb.DB {
	clientID := os.Getenv("LOCAL_CLIENT_ID")
	if clientID == "" {
		fmt.Printf("LOCAL_CLIENT_ID is not set in environment variables!")
		os.Exit(1)
	}

	dbFile := fmt.Sprintf(dbFile, clientID)

	if !dbExists(dbFile) {
		fmt.Println("Init blockchain first.")
		os.Exit(1)
	}

	db := heta.FullNode.ChainDB
	return db
}
