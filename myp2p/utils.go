package myp2p

import (
	"net"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/netutil"
	log "github.com/inconshreveable/log15"
)

// NewBootNodeServer ...
func NewBootNodeServer(listenAddr string, netrestrict string) *discover.Table {

	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		log.Error(err.Error())
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Error(err.Error())
	}

	realaddr := conn.LocalAddr().(*net.UDPAddr)

	var restrictList *netutil.Netlist
	if netrestrict != "" {
		restrictList, err = netutil.ParseNetlist(netrestrict)
		if err != nil {
			log.Error(err.Error())
		}
	}
	nodeKeyHex := "19f7e3cc1f65994aa61eea90a4d91e5dbca9c027cc378a60e4101d5b14a8530c"
	nodeKey, err := crypto.HexToECDSA(nodeKeyHex)
	cfg := discover.Config{
		PrivateKey:   nodeKey,
		AnnounceAddr: realaddr,
		NetRestrict:  restrictList,
	}
	node, err := discover.ListenUDP(conn, cfg)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("UDP listener up", "self", node.Self())
	return node
}
