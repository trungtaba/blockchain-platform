package trie

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/sotatek-dev/heta/types"
	"github.com/sotatek-dev/heta/utils"
)

// ONLY FOR TEST
type AccountType struct {
	Address types.Address
	Balance string
	Block   uint64
}

func MatchKey(a []byte, b []byte) int {
	var i, length = 0, len(a)
	if len(b) < length {
		length = len(b)
	}
	for ; i < length; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return i
}

// Input key: 16 bytes,
// Key after formating: 32 bytes,
// Last bytes: 16
func keybytesToHex(str []byte) []byte {
	l := len(str)*2 + 1
	var nibbles = make([]byte, l)
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	nibbles[l-1] = 16
	return nibbles
}

// ONLY FOR TEST
func GenNode(value int) *ShortNode {
	valueNod := GenNodeValue(value)
	short := GenShortNode(valueNod)
	return short
}

// ONLY FOR TEST
func GenNodeValue(value int) ValueNode {
	var b int
	utils.CloneValue(value, b)
	toAddress := types.Address{}
	addressBytes, err := hex.DecodeString("5b3e287767e13d88d11dfa919927c66492c30c41")
	if err != nil {
		panic(err)
	}
	copy(toAddress[20-len(addressBytes):], addressBytes)
	newAccc := &AccountType{
		Address: toAddress,
		Balance: "1000000" + string(value),
		Block:   uint64(b),
	}
	data := utils.Serialize(newAccc)
	node := ValueNode{}
	node = data
	return node
}

// ONLY FOR TEST
func GenShortNode(node ValueNode) *ShortNode {
	var hash types.Hash
	hash = sha256.Sum256(node)
	short := &ShortNode{hash[:], node, nodeFlag{}}
	return short
}
