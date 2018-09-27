package core

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sotatek-dev/heta/trie"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"math/big"
	"os"

	"github.com/sotatek-dev/heta/crypto"
	"github.com/sotatek-dev/heta/heta"
	"github.com/sotatek-dev/heta/keystore"
	"github.com/sotatek-dev/heta/types"
	"github.com/sotatek-dev/heta/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	// HashLength ...
	HashLength = 32

	// AddressLength ...
	AddressLength = 20
)

// Account ...
type Account struct {
	Key *keystore.Key
}

// AccountLookupEntry is a positional metadata to help looking up the state of an account
type AccountLookupEntry struct {
	Balance big.Int
	Stake   big.Int
	Vote    big.Int
	Index   uint64
}

// BalanceResFormat ...
type BalanceResFormat struct {
	Balance string `json:"balance"`
	Address string `json:"address"`
	Index   uint64 `json:"index"`
}

// AccountKeyFormat ...
type AccountKeyFormat struct {
	Address    string `json:"address"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

// AccountRequestFormat ...
type AccountRequestFormat struct {
	Address string `json:"address"`
}

// State
type State struct {
	Address types.Address
	Balance big.Int
}

// NewAccount ...
func NewAccount() *Account {
	key := keystore.NewKey()
	account := Account{key}
	return &account
}

// GetStateDB ...
func GetStateDB(options *opt.Options) *leveldb.DB {
	clientID := os.Getenv("LOCAL_CLIENT_ID")
	if clientID == "" {
		fmt.Printf("LOCAL_CLIENT_ID is not set in environment variables!")
		os.Exit(1)
	}

	dbStateFile := fmt.Sprintf(dbStateFile, clientID)

	if !dbExists(dbStateFile) {
		fmt.Println("Init blockchain first.")
		os.Exit(1)
	}

	db := heta.FullNode.StateDB

	return db
}

// NewAccountFromPrivateKey ...
func NewAccountFromPrivateKey(privateKey string) *Account {
	key := keystore.NewKeyFromPrivateKey(privateKey)
	account := Account{key}
	return &account
}

// ValidateAddress check if address if valid
func ValidateAddress(address string) bool {
	return keystore.ValidateAddress(address)
}

// AccountLookupKey = accountLookupPrefix + hash
func AccountLookupKey(address types.Address) []byte {
	return append(accountLookupPrefix, address[:]...)
}

// WriteAccountLookupEntries find address and push it into database if it is not existed
// else update amount of that address
func WriteAccountLookupEntries(db *leveldb.DB, block *Block) {
	for _, tx := range block.GetTransactions() {
		// subtract actor balance
		SubBalance(db, tx.Actor, tx.Data.Amount)
		// add receiver balance
		AddBalance(db, tx.Data.Recipient, tx.Data.Amount)
	}
}

// AddBalance ...
func AddBalance(db *leveldb.DB, account *types.Address, amountModify *big.Int) {
	accountEntry := new(AccountLookupEntry)
	accountBytes, err := db.Get(AccountLookupKey(*account), nil)
	utils.Deserialize(accountBytes, accountEntry)
	newAccountBalance := new(big.Int).Add(&accountEntry.Balance, amountModify)
	newAccountEntry := AccountLookupEntry{
		Balance: *newAccountBalance,
		Stake:   accountEntry.Stake,
		Vote:    accountEntry.Vote,
		Index:   uint64(1),
	}
	newAccountData := utils.Serialize(newAccountEntry)
	err = db.Put(AccountLookupKey(*account), newAccountData, nil)
	if err != nil {
		log.Panic("Failed to add address balance", " err ", err)
	}
	leaf := GenStateLeaf(account[:], &newAccountEntry)
	heta.FullNode.StateTrie.Push(account[:], leaf)
}

// InitBalance ...
func InitBalance(db *leveldb.DB, account *types.Address, amountModify *big.Int) {
	accountEntry := new(AccountLookupEntry)
	accountBytes, err := db.Get(AccountLookupKey(*account), nil)
	utils.Deserialize(accountBytes, accountEntry)
	newAccountBalance := new(big.Int).Add(&accountEntry.Balance, amountModify)
	newAccountEntry := AccountLookupEntry{
		Balance: *newAccountBalance,
		Stake:   accountEntry.Stake,
		Vote:    accountEntry.Vote,
		Index:   uint64(1),
	}
	newAccountData := utils.Serialize(newAccountEntry)
	err = db.Put(AccountLookupKey(*account), newAccountData, nil)
	if err != nil {
		log.Panic("Failed to add address balance", " err ", err)
	}
}

// SubBalance ...
func SubBalance(db *leveldb.DB, account *types.Address, amountModify *big.Int) {
	accountEntry := new(AccountLookupEntry)
	accountBytes, err := db.Get(AccountLookupKey(*account), nil)
	utils.Deserialize(accountBytes, accountEntry)
	newAccountBalance := new(big.Int).Sub(&accountEntry.Balance, amountModify)
	newAccountEntry := AccountLookupEntry{
		Balance: *newAccountBalance,
		Stake:   accountEntry.Stake,
		Vote:    accountEntry.Vote,
		Index:   uint64(1),
	}
	newAccountData := utils.Serialize(newAccountEntry)
	err = db.Put(AccountLookupKey(*account), newAccountData, nil)
	if err != nil {
		log.Panic("Failed to subtract address balance", "err", err)
	}
	leaf := GenStateLeaf(account[:], &newAccountEntry)
	heta.FullNode.StateTrie.Push(account[:], leaf)
}

// WriteGenesisAccountLookupEntries ...
func WriteGenesisAccountLookupEntries(db *leveldb.DB, block *Block) {
	genesisAmount := new(big.Int)
	genesisAmount.SetString("10000000000000000000000000", 10)
	fmt.Println("INIT MONEY")
	for _, tx := range block.GetTransactions() {
		var account types.Address
		account.SetBytes(tx.Actor[:])
		InitBalance(db, &account, genesisAmount)
	}
}

// GetAccount ...
func GetAccount(addr types.Address) BalanceResFormat {
	accountBytes := heta.FullNode.StateTrie.Find(heta.FullNode.StateTrie.Root, addr[:])
	accountData := new(State)
	utils.Deserialize(accountBytes, accountData)
	accountRes := BalanceResFormat{
		Address: hex.EncodeToString(addr[:]),
		Balance: accountData.Balance.String(),
		Index:   0,
	}
	return accountRes
}

// HandleBalanceReq ...
func HandleBalanceReq(b []byte) BalanceResFormat {
	reqJSON := AccountRequestFormat{}
	_ = json.Unmarshal(b[:], &reqJSON)

	accountBytes, err := hex.DecodeString(reqJSON.Address)
	if err != nil {
		panic(err)
	}
	account := new(types.Address)
	account.SetBytes(accountBytes)
	return GetAccount(*account)
}

// HandleAddressRegReq ...
func HandleAddressRegReq(b []byte) *AccountKeyFormat {
	account := NewAccount()
	address := hex.EncodeToString(account.Key.Address[:])
	privateKey := hex.EncodeToString(crypto.FromECDSA(account.Key.PrivateKey))

	pubkey := account.Key.PrivateKey.PublicKey
	ecdsaPubkey := ecdsa.PublicKey{Curve: pubkey.Curve, X: pubkey.X, Y: pubkey.Y}
	publicKey := hex.EncodeToString(crypto.FromECDSAPub(&ecdsaPubkey))

	accountKey := AccountKeyFormat{
		Address:    address,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	return &accountKey
}

// GenStateTrie ...
func GenStateTrie() *trie.Trie {
	m := make(map[types.Hash]*trie.CachedNode)
	dl := trie.DataLayer{
		Nodes: m,
	}
	stateTrie, err := trie.GenTrie(nil, &dl) // nil root
	statedb := GetStateDB(nil)

	iter := statedb.NewIterator(util.BytesPrefix(accountLookupPrefix), nil)
	for iter.Next() {
		data := iter.Value()
		account := new(AccountLookupEntry)
		utils.Deserialize(data, account)
		value := GenStateLeaf(iter.Key()[len(accountLookupPrefix):], account)
		stateTrie.Push(iter.Key()[len(accountLookupPrefix):], value)
	}
	hash := stateTrie.Hash(stateTrie.Root)
	fmt.Println("HASH: ", hex.EncodeToString(hash))
	iter.Release()
	err = iter.Error()
	if err != nil {
		log.Panic("IterateAccountDb", "err", err)
	}
	return stateTrie
}

// GenStateTrie ...
func GenStateLeaf(address []byte, account *AccountLookupEntry) []byte {
	cloneAcc := new(AccountLookupEntry)
	utils.CloneValue(account, cloneAcc)
	toAddress := types.Address{}
	copy(toAddress[20-len(address):], address)
	newStateValue := &State{
		Address: toAddress,
		Balance: cloneAcc.Balance,
	}
	data := utils.Serialize(newStateValue)
	return data
}
