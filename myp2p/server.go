package myp2p

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

// NewServer Creating new p2p server
func NewServer(name string, port uint64, protocols []p2p.Protocol, enodes []*discover.Node, isRandom bool) (*p2p.Server, error) {
	pkey, err := crypto.GenerateKey()
	if !isRandom {
		nodeKeyHex := "19f7e3cc1f65994aa61eea90a4d91e5dbca9c027cc378a60e4101d5b14a8530c"
		pkey, err = crypto.HexToECDSA(nodeKeyHex)
	}
	if err != nil {
		log.Printf("Generate private key failed with err: %v", err)
		return nil, err
	}

	cfg := p2p.Config{
		PrivateKey:      pkey,
		Name:            name,
		MaxPeers:        10,
		Protocols:       protocols,
		EnableMsgEvents: true,
		BootstrapNodes:  enodes,
	}

	if port > 0 {
		cfg.ListenAddr = fmt.Sprintf(":%d", port)
	}
	srv := &p2p.Server{
		Config: cfg,
	}

	err = srv.Start()
	if err != nil {
		log.Printf("Start server failed with err: %v", err)
		return nil, err
	}

	return srv, nil
}
