package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
<<<<<<< HEAD
	"log"
	"reflect"
	"runtime"

=======
>>>>>>> c838dea44f9ab4826534580d6b8aaecba2d35422
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sotatek-dev/heta/types"
	"log"
	"reflect"
	"runtime"
)

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) types.Hash {
	var h types.Hash
	h.SetBytes(b)
	return h
}

// Serialize serializes the block
func Serialize(d interface{}) []byte {
	data, err := rlp.EncodeToBytes(d)
	if err != nil {
		log.Panic(err)
	}
	return data
}

// Hash hash a struct
func Hash(d interface{}) types.Hash {
	data, err := rlp.EncodeToBytes(d)
	if err != nil {
		log.Panic(err)
	}
	return BytesToHash(data)
}

// Deserialize deserializes a block
func Deserialize(d []byte, object interface{}) {
	if len(d) == 0 {
		return
	}
	if err := rlp.Decode(bytes.NewReader(d), object); err != nil {
		log.Panic("Invalid block body RLP", err)
		return
	}
}

// EncodeUInt64 ...
func EncodeUInt64(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// DecodeUInt64 ...
func DecodeUInt64(d []byte) uint64 {
	number := binary.BigEndian.Uint64(d)
	return number
}

// CloneValue function
// Source: pointer
// Destin: allocated struct value
func CloneValue(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		destin = x.Interface()
	}
}

// MakeName creates a node name that follows the ethereum convention
// for such names. It adds the operation system name and Go runtime version
// the name.
func MakeName(name, version string) string {
	return fmt.Sprintf("%s/v%s/%s/%s", name, version, runtime.GOOS, runtime.Version())
}

// CloneValue ...
// CloneValue function
// Source: pointer
// Destin: allocated struct value
func CloneValue(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		destin = x.Interface()
	}
}
