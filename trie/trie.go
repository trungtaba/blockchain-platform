package trie

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/sotatek-dev/heta/types"
	"github.com/sotatek-dev/heta/utils"
)

//	TODO: need to implement
type Trie struct {
	Root node
	dl   *DataLayer
}

func GenTrie(root node, dl *DataLayer) (*Trie, error) {
	trie := Trie{
		dl: dl,
	}
	trie.Root = root
	return &trie, nil
}

// TODO: Delete func now is not working
func (t *Trie) Push(key, value []byte) error {
	k := keybytesToHex(key)
	if len(value) != 0 {
		_, n, err := Insert(t.Root, nil, k, ValueNode(value))
		if err != nil {
			return err
		}
		t.Root = n
	} else {
		fmt.Println("Delete")
		_, n, err := t.Delete(t.Root, nil, k)
		if err != nil {
			return err
		}
		t.Root = n
	}
	return nil
}

// TODO: check value of hash after insert
// TODO: key of node need to HP encoded before insert into trie
func Insert(node node, prefix, key []byte, value node) (bool, node, error) {
	if len(key) == 0 {
		// if node is a value node
		if v, ok := node.(ValueNode); ok {
			return !bytes.Equal(v, value.(ValueNode)), value, nil
		}
		return true, value, nil
	}
	switch originNode := node.(type) {
	case *FullNode:
		dirty, newChildNode, err := Insert(originNode.Children[key[0]], append(prefix, key[0]), key[1:], value)
		if !dirty || err != nil {
			return false, originNode, err
		}
		originNode.Children[key[0]] = newChildNode
		return true, originNode, nil
	case nil:
		return true, &ShortNode{key, value, nodeFlag{}}, nil
	case *ShortNode:
		// If the whole key matches, keep this short node as is
		// and only update the value.
		matchLen := MatchKey(key, originNode.Key)
		if matchLen == len(originNode.Key) {
			// update value with key length == 0
			dirty, newChildNode, err := Insert(originNode.Val, append(prefix, key[:matchLen]...), key[matchLen:], value)
			if !dirty || err != nil {
				return false, newChildNode, err
			}
			return true, &ShortNode{originNode.Key, newChildNode, nodeFlag{}}, nil
		}
		// create branch for node
		branch := &FullNode{}
		var err error
		_, branch.Children[originNode.Key[matchLen]], err = Insert(nil, append(prefix, originNode.Key[:matchLen+1]...), originNode.Key[matchLen+1:], originNode.Val)
		if err != nil {
			return false, nil, err
		}
		_, branch.Children[key[matchLen]], err = Insert(nil, append(prefix, key[:matchLen+1]...), key[matchLen+1:], value)
		if err != nil {
			return false, nil, err
		}
		// Replace this ShortNode with the branch if it occurs at index 0.
		if matchLen == 0 {
			return true, branch, nil
		}
		// Otherwise, replace it with a short node leading up to the branch.
		return true, &ShortNode{key[:matchLen], branch, nodeFlag{}}, nil
	case HashNode:
		return true, originNode, nil

	default:
		panic(fmt.Sprintf("%T: invalid node: "))
	}
}

func (t *Trie) Find(paramNode node, key []byte) ValueNode {
	k := keybytesToHex(key)
	return Get(paramNode, k)
}

func Get(paramNode node, key []byte) ValueNode {
	switch originNode := paramNode.(type) {
	case ValueNode:
		return originNode
	case HashNode:
		return nil
	case *ShortNode:
		matchLen := MatchKey(key, originNode.Key)
		value := Get(originNode.Val, key[matchLen:])
		return value
	case *FullNode:
		var nodeGet ValueNode
		nodeGet = Get(originNode.Children[key[0]], key[1:])
		return nodeGet
	default: // origin node and short node
		return nil
	}
}

func (t *Trie) Delete(node node, prefix, key []byte) (bool, node, error) {
	return true, nil, nil
}

// TODO: implement hash function
func (t *Trie) Hash(node node) HashNode {
	switch originNode := node.(type) {
	case *FullNode:
		var allChildHash []HashNode
		for _, child := range originNode.Children {
			allChildHash = append(allChildHash, t.Hash(child))
		}
		fullHash := MultiNodeHash(allChildHash)
		t.CacheNode(fullHash[:], originNode)
		return fullHash
	case *ShortNode:
		var hash types.Hash
		hash = sha256.Sum256(utils.Serialize(originNode))
		t.CacheNode(hash[:], originNode)
		return hash[:]
	default: // value node and hash node
		return nil
	}
}

func MultiNodeHash(nodes []HashNode) HashNode {
	var hash types.Hash
	var longHash []byte
	for _, node := range nodes {
		longHash = append(longHash, node[:]...)
	}
	hash = sha256.Sum256(utils.Serialize(longHash))
	//fmt.Println(hex.EncodeToString(hash[:]))
	return hash[:]
}

// Add current hash into node and cache it
func (t *Trie) CacheNode(hash []byte, node node) {
	switch originNode := node.(type) {
	case *FullNode:
		originNode.Flags.Hash = hash
	case *ShortNode:
		originNode.Flags.Hash = hash
	}
	t.dl.Push(node)
}
