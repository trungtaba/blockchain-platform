package keystore

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/pborman/uuid"
	"github.com/sotatek-dev/heta/crypto"
	"github.com/sotatek-dev/heta/types"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

// Key means a keypair, private key and its address
type Key struct {
	ID         uuid.UUID
	Address    types.Address
	PrivateKey *ecdsa.PrivateKey
}

// NewKey creates a new key
func NewKey() *Key {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	id := uuid.NewRandom()
	key := &Key{
		ID:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}

// NewKeyFromPrivateKey creates a new key from private key string
func NewKeyFromPrivateKey(privateKey string) *Key {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Panic(err)
	}

	id := uuid.NewRandom()
	key := &Key{
		ID:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// ValidateAddress check if address if valid
func ValidateAddress(address string) bool {
	// pubKeyHash := utils.Base58Decode([]byte(address))
	// actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	// version := pubKeyHash[0]
	// pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	// targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	// return bytes.Compare(actualChecksum, targetChecksum) == 0
	return true
}
