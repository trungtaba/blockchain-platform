package core

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"reflect"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/joho/godotenv"
	"github.com/sotatek-dev/heta/crypto"
	"github.com/sotatek-dev/heta/heta"
	"github.com/sotatek-dev/heta/types"
	"github.com/sotatek-dev/heta/utils"
)

const (
	TRANSFER            = "transfer"
	TRANSFER_EMPTY_DATA = "{}"
)

// Transaction ...
type Transaction struct {
	ID    []byte
	Actor *types.Address
	Data  Txdata
}

// Txdata ...
type Txdata struct {
	AccountNonce uint64         `json:"nonce"    gencodec:"required"`
	Recipient    *types.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int       `json:"value"    gencodec:"required"`
	Payload      []byte         `json:"input"    gencodec:"required"`
	Action       Action

	// This is only used when marshaling to JSON.
	Hash *types.Hash `json:"hash" rlp:"-"`
}

// TxRequestFormat ...
type TxRequestFormat struct {
	Pubkey    string `json:"pubkey"`
	Nonce     int    `json:"nonce"`
	Address   string `json:"address"`
	Payload   string `json:"payload"`
	Action    Action `json:"action"`
	Hash      string `json:"hash"`
	Signature string `json:"signature"`
}

// SignTxRequestFormat ...
type SignTxRequestFormat struct {
	PrivateKey string `json:"private_key"`
	Nonce      int    `json:"nonce"`
	Address    string `json:"address"`
	Amount     string `json:"amount"`
	Payload    string `json:"payload"`
	Action     Action `json:"action"`
	Hash       string `json:"hash"`
}

// TxRes ...
type TxRes struct {
	Txid             string `json:"hash"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      uint64 `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	Actor            string `json:"actor"`
	Nonce            string `json:"nonce"`
	Recipient        string
	Payload          string
	Action           Action `json:"action"`
	Amount           *big.Int
	Hash             string
}

// TransferFormat ...
type TransferFormat struct {
	Amount string `json:"amount"`
	To     string `json:"to"`
}

// TransactionRequestFormat ...
type TransactionRequestFormat struct {
	Txid string `json:"hash"`
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Panic(err)
	}
	return data
}

// Serialize ...
func (data Txdata) Serialize() []byte {
	d, err := rlp.EncodeToBytes(data)
	if err != nil {
		log.Panic(err)
	}
	return d
}

// Hash returns the Hash of the Transaction
func (tx *Transaction) Hash() []byte {
	var hash types.Hash

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// TxDataHash returns the Hash of the Transaction data
func (tx *Transaction) TxDataHash() []byte {
	var hash types.Hash
	hash = sha256.Sum256(tx.Data.Serialize())
	return hash[:]
}

// TxDataHash ...
func (Txdata *Txdata) TxDataHash() []byte {
	var hash types.Hash
	hash = sha256.Sum256(Txdata.Serialize())
	return hash[:]
}

// NewCoinbaseTX ...
func NewCoinbaseTX(actor *types.Address, data []byte) *Transaction {
	return NewTX(actor, uint64(0), actor, new(big.Int), data)
}

// NewTX ...
func NewTX(actor *types.Address, nonce uint64, to *types.Address, amount *big.Int, data []byte) *Transaction {
	d := Txdata{
		AccountNonce: nonce,
		Recipient:    to,
		Payload:      data,
		Amount:       new(big.Int),
		Action: Action{
			ActionName: TRANSFER,
			Data:       []byte(TRANSFER_EMPTY_DATA),
		},
	}
	tx := Transaction{nil, actor, d}
	tx.ID = tx.Hash()

	return &tx
}

// SignTransactionData ...
func (Txdata *Txdata) SignTransactionData(key *ecdsa.PrivateKey) (sig []byte, err error) {
	byteHash := Txdata.TxDataHash()
	return crypto.Sign(key, byteHash)
}

// VerifyTransactionData ...
func (Txdata *Txdata) VerifyTransactionData(pub *ecdsa.PublicKey, sign []byte) bool {
	pubkeyByte := crypto.FromECDSAPub(pub)
	byteHash := Txdata.TxDataHash()
	return crypto.VerifySignature(pubkeyByte, byteHash, sign)
}

// ToTxRes Convert transaction to TxRes format
func (tx *Transaction) ToTxRes() (txres TxRes) {
	dataHash := []byte{}
	if tx.Data.Hash != nil {
		dataHash = tx.Data.Hash[:]
	} else {
		dataHash = []byte{}
	}

	return TxRes{
		Txid:      hex.EncodeToString(tx.Hash()),
		Actor:     hex.EncodeToString(tx.Actor[:]),
		Recipient: hex.EncodeToString(tx.Data.Recipient[:]),
		Payload:   string(tx.Data.Payload),
		Action:    tx.Data.Action,
		Amount:    tx.Data.Amount,
		Hash:      hex.EncodeToString(dataHash),
	}
}

// GetTransaction ...
func GetTransaction(txid string) TxRes {
	// open connection
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
	if err != nil {
		log.Panic(err)
	}
	decodeTxid, err := hex.DecodeString(txid)
	if err != nil {
		log.Panic(err)
	}

	txHash := types.Hash{}
	txHash.SetBytes(decodeTxid)
	txLookup := txLookupKey(txHash)
	data, err := db.Get(txLookup, nil)
	txEntry := new(TxLookupEntry)
	utils.Deserialize(data, txEntry)

	block := GetBlockByHash(txEntry.BlockHash)
	txEntryInfo := block.GetTransactions()[txEntry.BlockIndex]
	return txEntryInfo.ToTxRes()
}

/** HandleTxInfoRequest ...
** Get transaction info
**/
func HandleTxInfoRequest(b []byte) TxRes {
	reqJSON := TransactionRequestFormat{}
	err := json.Unmarshal(b[:], &reqJSON)
	if err != nil {
		panic(err)
	}
	return GetTransaction(reqJSON.Txid)
}

/** Handle transaction requests ...
** This is only used for transfer transaction format
**/
func HandleTxRequest(b []byte) bool {
	tx := TxRequestFormat{}
	txData := ParseTransactionRequest(b, &tx)

	byteSignature, err := hex.DecodeString(tx.Signature)
	pubKeyByte, err := hex.DecodeString(tx.Pubkey)
	if err != nil {
		panic(err)
	}
	pubKey, err := crypto.UnmarshalPubkey(pubKeyByte)
	if err != nil {
		panic(err)
	}
	verify := txData.VerifyTransactionData(pubKey, byteSignature)
	return verify
}

/** Handle sign transaction requests ...
** This is only used for signing transaction format
**/
func HandleSignTxRequest(b []byte) string {
	tx := SignTxRequestFormat{}
	txData := ParseTransactionRequest(b, &tx)

	privateKey, err := crypto.HexToECDSA(tx.PrivateKey)
	if err != nil {
		panic(err)
	}
	signBytes, err := txData.SignTransactionData(privateKey)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(signBytes)
}

/** Parse transaction requests under byte format to struct value
** This is used for every structure that contains TxData Fields
**/
func ParseTransactionRequest(b []byte, tx interface{}) Txdata {
	err := json.Unmarshal(b[:], tx)
	if err != nil {
		panic(err)
	}

	action := reflect.Indirect(reflect.ValueOf(tx)).FieldByName("Action").Interface().(Action)
	nonce := uint64(reflect.Indirect(reflect.ValueOf(tx)).FieldByName("Nonce").Interface().(int))
	toAddress := types.Address{}
	addressBytes, _ := hex.DecodeString(reflect.Indirect(reflect.ValueOf(tx)).FieldByName("Address").Interface().(string))
	copy(toAddress[AddressLength-len(addressBytes):], addressBytes)
	payload := []byte(reflect.Indirect(reflect.ValueOf(tx)).FieldByName("Payload").Interface().(string))
	transfer := TransferFormat{}
	err = json.Unmarshal(action.Data, &transfer)
	if err != nil {
		panic(err)
	}

	amount := big.Int{}
	amount.SetString(transfer.Amount, 10)
	hash := types.Hash{}
	hashBytes, _ := hex.DecodeString(reflect.Indirect(reflect.ValueOf(tx)).FieldByName("Hash").Interface().(string))
	copy(hash[HashLength-len(hashBytes):], hashBytes)

	return Txdata{
		nonce,
		&toAddress,
		&amount,
		payload,
		action,
		&hash,
	}
}
