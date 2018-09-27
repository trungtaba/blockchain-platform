package trie

import "fmt"

var indices = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "[17]"}

type node interface {
	fstring(string) string
	// cache and canUnload func
}

type (
	FullNode struct {
		Children [17]node // Actual trie node data to encode/decode (needs custom encoder)
		Flags    nodeFlag `rlp:"-"`
	}
	ShortNode struct {
		Key   []byte
		Val   node
		Flags nodeFlag `rlp:"-"`
	}
	HashNode  []byte
	ValueNode []byte
)

type nodeFlag struct {
	// un-implemented
	Hash HashNode
}

// implement node interface for 4 struct
func (n *FullNode) fstring(ind string) string {
	resp := fmt.Sprintf("[\n%s  ", ind)
	for i, node := range &n.Children {
		if node == nil {
			resp += fmt.Sprintf("%s: <nil> ", indices[i])
		} else {
			resp += fmt.Sprintf("%s: %v", indices[i], node.fstring(ind+"  "))
		}
	}
	return resp + fmt.Sprintf("\n%s] ", ind)
}
func (n *ShortNode) fstring(ind string) string {
	return fmt.Sprintf("{%x: %v} ", n.Key, n.Val.fstring(ind+"  "))
}
func (n HashNode) fstring(ind string) string {
	return fmt.Sprintf("<%x> ", []byte(n))
}
func (n ValueNode) fstring(ind string) string {
	return fmt.Sprintf("%x ", []byte(n))
}

// implement decode and encode node
