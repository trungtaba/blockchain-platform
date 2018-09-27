package trie

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/sotatek-dev/heta/types"
	"github.com/sotatek-dev/heta/utils"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var dbFile = "../.data/chaindata-%s"

// key: 9302e99a53668155839a575711629a8720da83b4b1a74962fd203a9d20cf8da7
func TestInsert(t *testing.T) {
	m := make(map[types.Hash]*CachedNode)
	dl := DataLayer{
		Nodes: m,
	}
	trie, err := GenTrie(nil, &dl) // nil Root
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1000; i++ {
		short := GenNode(i)
		trie.Push(short.Key, short.Val.(ValueNode))
	}
	key, _ := hex.DecodeString("9302e99a53668155839a575711629a8720da83b4b1a74962fd203a9d20cf8da7")
	findingResult := trie.Find(trie.Root, key)
	fmt.Println("Finding result ", findingResult)
}

// Override into trie with key
// Key: 9302e99a53668155839a575711629a8720da83b4b1a74962fd203a9d20cf8da7
func TestOverride(t *testing.T) {
	dl := new(DataLayer)
	hexStringKey := "9302e99a53668155839a575711629a8720da83b4b1a74962fd203a9d20cf8da7"
	newValue := GenNode(1000000)
	trie, err := GenTrie(nil, dl) // nil Root
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
		short := GenNode(i)
		trie.Push(short.Key, short.Val.(ValueNode))
	}
	key, _ := hex.DecodeString(hexStringKey)
	findingResult := trie.Find(trie.Root, key)
	fmt.Println("Finding result 1: ", findingResult)

	trie.Push(key, newValue.Val.(ValueNode))
	newFindingResult := trie.Find(trie.Root, key)
	fmt.Println("Finding result 2: ", newFindingResult)
}

// Use Gen node value
// Problem: memory leak
func TestMultiInsertion(t *testing.T) {
	m := make(map[types.Hash]*CachedNode)
	dl := DataLayer{
		Nodes: m,
	}
	trie, err := GenTrie(nil, &dl) // nil Root
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10000000; i++ {
		short := GenNode(i)
		trie.Push(short.Key, short.Val.(ValueNode))
	}
	hash := trie.Hash(trie.Root)
	fmt.Println("HASH: ", hex.EncodeToString(hash))
}

// 100.000.000 (key, values), about 7 minutes running
// Then hashing all Nodes is generated
func TestBigGeneration(t *testing.T) {
	m := make(map[types.Hash]*CachedNode)
	dl := DataLayer{
		Nodes: m,
	}
	trie, err := GenTrie(nil, &dl) // nil Root
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100000000; i++ {
		intVal, err := strconv.Atoi("0")
		if err != nil {
			panic(err)
		}
		toAddress := types.Address{}
		addressBytes, err := hex.DecodeString("5b3e287767e13d88d11dfa919927c66492c30c41")
		if err != nil {
			panic(err)
		}
		copy(toAddress[20-len(addressBytes):], addressBytes)
		newAccc := &AccountType{
			Address: toAddress,
			Balance: "1000000" + string(i),
			Block:   uint64(intVal),
		}
		data := utils.Serialize(newAccc)
		node := ValueNode{}
		node = data
		var hash types.Hash
		hash = sha256.Sum256(data)
		shortNode := &ShortNode{hash[:], node, nodeFlag{}}
		trie.Push(shortNode.Key, node)
	}
	hash := trie.Hash(trie.Root)
	fmt.Println("HASH: ", hex.EncodeToString(hash))
}

// key: 9302e99a53668155839a575711629a8720da83b4b1a74962fd203a9d20cf8da7
func TestNodeHash(t *testing.T) {
	m := make(map[types.Hash]*CachedNode)
	dl := DataLayer{
		Nodes: m,
	}
	trie, err := GenTrie(nil, &dl) // nil Root
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1000; i++ {
		short := GenNode(i)
		trie.Push(short.Key, short.Val.(ValueNode))
	}
	key, _ := hex.DecodeString("9302e99a53668155839a575711629a8720da83b4b1a74962fd203a9d20cf8da7")
	findingResult := trie.Find(trie.Root, key)

	stateRoot := trie.Hash(trie.Root)
	hashRoot := trie.Root.(*FullNode).Flags.Hash
	fmt.Println("Finding result ", findingResult)
	fmt.Println("Hash result ", hex.EncodeToString(stateRoot))
	fmt.Println("Root hash   ", hex.EncodeToString(hashRoot))
	assert.Equal(t, stateRoot, hashRoot)
}

func TestFindKey(t *testing.T) {

}
