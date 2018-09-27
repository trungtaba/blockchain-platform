package types

import (
	"bytes"
)

const (
	// HashLength is the desired length of a hash
	HashLength = 32
)

// Hash is a 256-bit hex number actually
type Hash [HashLength]byte

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// IsZeroHash checks whether a hash is all-zero
func IsZeroHash(hash Hash) bool {
	zeroHash := make([]byte, 32)
	return bytes.Equal(hash[:], zeroHash)
}
