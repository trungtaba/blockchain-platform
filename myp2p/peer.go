package myp2p

import (
	"log"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

// ConnectToPeer ...
func ConnectToPeer(srv *p2p.Server, enodes []string) error {
	// Parsing the enode url
	for _, enode := range enodes {
		node, err := discover.ParseNode(enode)
		if err != nil {
			log.Printf("Failed to parse enode url with err: %v", err)
			return err
		}

		// Connecting to the peer
		srv.AddPeer(node)
	}

	return nil
}
