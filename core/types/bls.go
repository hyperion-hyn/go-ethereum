package types

import (
    "bytes"
    "errors"
    "encoding/hex"
    "math/big"
    "fmt"
    "github.com/hyperion-hyn/bls/ffi/go/bls"
)

var (
    emptyBlsPubKey = BLSPublicKey{}
)

// PublicKeySizeInBytes ..
const (
    PublicKeySizeInBytes    = 48
    BlsSignatureSizeInBytes = 96
)


// BLSPublicKey defines the bls public key
type BLSPublicKey [PublicKeySizeInBytes]byte

// BLSSignature defines the bls signature
type BLSSignature [BlsSignatureSizeInBytes]byte



// Big ..
func (pk BLSPublicKey) Big() *big.Int {
    return new(big.Int).SetBytes(pk[:])
}

// IsEmpty returns whether the bls public key is empty 0 bytes
func (pk BLSPublicKey) IsEmpty() bool {
    return bytes.Compare(pk[:], emptyBlsPubKey[:]) == 0
}

// Hex returns the hex string of bls public key
func (pk BLSPublicKey) Hex() string {
    return hex.EncodeToString(pk[:])
}

// MarshalJSON ..
func (pk BLSPublicKey) MarshalJSON() ([]byte, error) {
    buf := bytes.Buffer{}
    buf.WriteString(`"`)
    buf.WriteString(pk.Hex())
    buf.WriteString(`"`)
    return buf.Bytes(), nil
}

// FromLibBLSPublicKey replaces the key contents with the given key,
func (pk *BLSPublicKey) FromLibBLSPublicKey(key *bls.PublicKey) error {
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
func (pk *BLSPublicKey) ToLibBLSPublicKey(key *bls.PublicKey) error {
    return key.Deserialize(pk[:])
}

// CompareBlsPublicKey compares two BLSPublicKey, lexicographically.
func CompareBlsPublicKey(k1, k2 BLSPublicKey) int {
    return bytes.Compare(k1[:], k2[:])
}

