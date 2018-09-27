package myp2p

import (
	"math/big"
)

// VoteContent ...
type VoteContent struct {
	Timestamp uint64
	Address   string
	Vote      *big.Int
	PublicKey []byte
	IsUnvote  bool
}

// VoteMsg ...
type VoteMsg struct {
	VoteContent
	Signature []byte
}
