package main

import (
	"encoding/hex"
	"fmt"
	"github.com/sotatek-dev/heta/core"
	"github.com/sotatek-dev/heta/crypto"
	"github.com/sotatek-dev/heta/types"
	"math/big"
	"testing"
)

/**
Private key 803e8fb594da8fd86177c115b0da6a3d0903f5c4a0d29948ededbbffe57e0051 (hex of big int D value)
Public key 04b0c5dd3ed31cfc5d3613ebd838c9f66ccdebfc19e515a88a75fed09a255beab35ed5c0ddf38a5fbb87327901cb94119125f9c79b293f1c7f3d22eb719e9579f8
D value 58006580981828425660811739450512947025097799432742941345161430916505614024785
*/
func TestVerifyingTransaction(t *testing.T) {
	d_va := big.Int{}
	d_va.SetString("58006580981828425660811739450512947025097799432742941345161430916505614024785", 10)
	stringHex, err := hex.DecodeString("04b0c5dd3ed31cfc5d3613ebd838c9f66ccdebfc19e515a88a75fed09a255beab35ed5c0ddf38a5fbb87327901cb94119125f9c79b293f1c7f3d22eb719e9579f8")

	pub, err := crypto.UnmarshalPubkey(stringHex)
	if err != nil {
		panic(err)
	}

	privateKeyECDSA, err := crypto.ToECDSA(d_va.Bytes())
	a := big.Int{}
	b, _ := hex.DecodeString("803e8fb594da8fd86177c115b0da6a3d0903f5c4a0d29948ededbbffe57e0051")
	a.SetBytes(b)
	if err != nil {
		panic(err)
	}

	// pushedTxData
	nonce := uint64(1)
	toAddress := types.Address{}
	toAddressData, _ := hex.DecodeString("5b3e287767e13d88d11dfa919927c66492c30c41")
	copy(toAddress[20-len(toAddressData):], toAddressData)
	value := big.Int{}
	value.SetString("100000", 10)
	payload := []byte("nothing")
	action := []byte("yuu") // yuu
	hash := types.Hash{}
	hashData, _ := hex.DecodeString("5b3e287767e13d88d11dfa919927c66492c30c41")
	copy(hash[32-len(hashData):], hashData)

	pushedTxData := core.Txdata{
		nonce,      // 1
		&toAddress, // yuu
		&value,     // "100000"
		payload,    // nothing
		action,     // yuu
		&hash,
	}

	sig, _ := pushedTxData.SignTransactionData(privateKeyECDSA)
	check := pushedTxData.VerifyTransactionData(pub, sig)

	fmt.Println("Pubkey obj                  : ", pub)
	fmt.Println(fmt.Sprintf("Private key (hex of D value): %x", privateKeyECDSA.D))
	fmt.Println("d value from private key    : ", a)
	fmt.Println("original d value            : ", d_va)
	fmt.Println("Data bytes                  : ", pushedTxData.Serialize())
	fmt.Println("Txid                        : ", hex.EncodeToString(pushedTxData.TxDataHash()))
	fmt.Println("Signature under hex encoding: ", hex.EncodeToString(sig[:]))
	fmt.Println("VERIFYING TRANSACTION RESULT: ", check)
}
