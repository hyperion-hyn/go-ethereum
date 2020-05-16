package types

import (
    "bytes"
    "errors"
    "encoding/hex"
    "math/big"
    "fmt"
    "github.com/harmony-one/bls/ffi/go/bls"
)

var (
    emptyBlsPubKey = BlsPublicKey{}
)

// PublicKeySizeInBytes ..
const (
    PublicKeySizeInBytes    = 48
    BlsSignatureSizeInBytes = 96
)


// BlsPublicKey defines the bls public key
type BlsPublicKey [PublicKeySizeInBytes]byte

// BlsSignature defines the bls signature
type BlsSignature [BlsSignatureSizeInBytes]byte



// Big ..
func (pk BlsPublicKey) Big() *big.Int {
    return new(big.Int).SetBytes(pk[:])
}

// IsEmpty returns whether the bls public key is empty 0 bytes
func (pk BlsPublicKey) IsEmpty() bool {
    return bytes.Compare(pk[:], emptyBlsPubKey[:]) == 0
}

// Hex returns the hex string of bls public key
func (pk BlsPublicKey) Hex() string {
    return hex.EncodeToString(pk[:])
}

// MarshalJSON ..
func (pk BlsPublicKey) MarshalJSON() ([]byte, error) {
    buf := bytes.Buffer{}
    buf.WriteString(`"`)
    buf.WriteString(pk.Hex())
    buf.WriteString(`"`)
    return buf.Bytes(), nil
}

// FromLibBLSPublicKey replaces the key contents with the given key,
func (pk *BlsPublicKey) FromLibBLSPublicKey(key *bls.PublicKey) error {
    bytes := key.Serialize()
    if len(bytes) != len(pk) {
        return errors.New( fmt.Sprint("BLS public key size mismatch",
            "expected", len(pk),
            "actual", len(bytes)))
    }
    copy(pk[:], bytes)
    return nil
}

// ToLibBLSPublicKey copies the key contents into the given key.
func (pk *BlsPublicKey) ToLibBLSPublicKey(key *bls.PublicKey) error {
    return key.Deserialize(pk[:])
}

// CompareBlsPublicKey compares two BlsPublicKey, lexicographically.
func CompareBlsPublicKey(k1, k2 BlsPublicKey) int {
    return bytes.Compare(k1[:], k2[:])
}

