package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/sotatek-dev/heta/types"
	"golang.org/x/crypto/sha3"
	"math/big"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

// Keccak256 ...
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return secp256k1.S256()
}

// UnmarshalPubkey converts bytes to a secp256k1 public key.
func UnmarshalPubkey(pub []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(S256(), pub)
	if x == nil {
		return nil, fmt.Errorf("error when unmarshal public key")
	}
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}, nil
}

// FromECDSAPub ...
func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}

// FromECDSA exports a private key into a binary dump.
func FromECDSA(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}
	return math.PaddedBigBytes(priv.D, priv.Params().BitSize/8)
}

// ToECDSA creates a private key with the given D value.
func ToECDSA(d []byte) (*ecdsa.PrivateKey, error) {
	return toECDSA(d, true)
}

func ToECDSAUnsafe(d []byte) *ecdsa.PrivateKey {
	priv, _ := toECDSA(d, false)
	return priv
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = S256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
}

// PubkeyToAddress ...
func PubkeyToAddress(p ecdsa.PublicKey) types.Address {
	pubBytes := FromECDSAPub(&p)
	return types.BytesToAddress(Keccak256(pubBytes[1:])[12:])
}

// Sign data
func Sign(key *ecdsa.PrivateKey, data []byte) (sig []byte, err error) {
	if len(data) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(data))
	}
	seckey := math.PaddedBigBytes(key.D, key.Params().BitSize/8)
	defer zeroBytes(seckey)
	return secp256k1.Sign(data, seckey)
}

// Verify data
func VerifySignature(pub, hash, signature []byte) bool {
	// pass 64 bytes first of signature
	return secp256k1.VerifySignature(pub, hash, signature[:64])
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(hash, sig []byte) ([]byte, error) {
	return secp256k1.RecoverPubkey(hash, sig)
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	s, err := Ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}

	x, y := elliptic.Unmarshal(S256(), s)
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}, nil
}

// HexToECDSA parses a secp256k1 private key.
func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return ToECDSA(b)
}
