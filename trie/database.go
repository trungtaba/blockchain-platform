package trie

import (
	"github.com/sotatek-dev/heta/types"
)

// CachedNode is all the information we know about a single cached node in the
// memory database write layer.
type CachedNode struct {
	node node   // Cached collapsed trie node, or raw rlp data
	size uint16 // Byte size of the useful cached data

	parents  uint16                // Number of live Nodes referencing this one
	children map[types.Hash]uint16 // External children referenced by this node

	flushCurrent types.Hash // Current node in the flush-list
	flushNext    types.Hash // Next node in the flush-list
}

// A database layer between disk and memory
type DataLayer struct {
	Nodes map[types.Hash]*CachedNode
}

// Push short node and full node into data layer
// That is saved in memory database
// TODO: node need to rlp encode before store down
func (dl *DataLayer) Push(node node) {
	cacheNode := CachedNode{
		node: node,
	}
	switch originNode := node.(type) {
	case *FullNode:
		cacheNode.flushCurrent.SetBytes(originNode.Flags.Hash)
	case *ShortNode:
		cacheNode.flushCurrent.SetBytes(originNode.Flags.Hash)
	}
	//fmt.Println("CACHE ~ ", hex.EncodeToString(cacheNode.flushCurrent[:]))
	dl.Nodes[cacheNode.flushCurrent] = &cacheNode
}

// Push short node and full node into data layer
// That is saved in memory database
// TODO: node need to rlp encode before store down
func (db *DataLayer) Get() {

}

// This func is used in Push()
// TODO: Implement this function to remove node from datalayer
func (db *DataLayer) Remove(flushHash types.Hash) {

}

// TODO: Implement this function to write entries into hard disk
func (db *DataLayer) WriteEntries() {

}

// TODO: Implement this function to get entries from hard disk
func (db *DataLayer) GetEntries() {

}
