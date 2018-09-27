package types

import (
	"encoding/hex"
)

const (
	// AddressLength is the desired length of an address
	AddressLength = 20
)

// Address is a 160-bit hex number actually
type Address [AddressLength]byte

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// ToString Convert an address from hex number to string
func (a *Address) ToString() string {
	return hex.EncodeToString(a[:])
}
