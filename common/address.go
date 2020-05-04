package common

import (
    "bytes"
    "encoding/hex"
    "math/big"

    "github.com/pkg/errors"
    "golang.org/x/crypto/sha3"

    "github.com/ethereum/go-ethereum/common/bech32"
)

var (
    emptyAddress = Address{}
)

// IsBech32Address verifies whether a string can represent a valid bech32-encoded
// Harmony address or not.
func IsBech32Address(s string) bool {
    hrp, bytes, err := bech32.DecodeAndConvert(s)
    if err != nil || (hrp != "hyn" && hrp != "thyn") || len(bytes) != AddressLength {
        return false
    }
    return true
}

// IsEmpty gets whether the address contains all 0 bytes
func (a Address) IsEmpty() bool {
    return bytes.Compare(a[:], emptyAddress[:]) == 0
}

// Big converts an address to a big integer.
func (a Address) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Bech32 returns an bip0173-compliant string representation of the address.
func (a Address) Bech32() string {
    unchecksummed := hex.EncodeToString(a[:])
    sha := sha3.NewLegacyKeccak256()
    sha.Write([]byte(unchecksummed))
    hash := sha.Sum(nil)

    result := []byte(unchecksummed)
    for i := 0; i < len(result); i++ {
        hashByte := hash[i/2]
        if i%2 == 0 {
            hashByte = hashByte >> 4
        } else {
            hashByte &= 0xf
        }
        if result[i] > '9' && hashByte > 7 {
            result[i] -= 32
        }
    }
    return "0x" + string(result)
}

// TODO ek – the following functions use Ethereum addresses until we have a
//  proper abstraction set in place.

// ParseBech32Addr decodes the given bech32 address and populates the given
// human-readable-part string and address with the decoded result.
func ParseBech32Addr(b32 string, hrp *string, addr *Address) error {
    h, b, err := bech32.DecodeAndConvert(b32)
    if err != nil {
        return errors.Wrapf(err, "cannot decode %#v as bech32 address", b32)
    }
    if len(b) != AddressLength {
        return errors.Errorf("decoded bech32 %#v has invalid length %d",
            b32, len(b))
    }
    *hrp = h
    addr.SetBytes(b)
    return nil
}

// BuildBech32Addr encodes the given human-readable-part string and address
// into a bech32 address.
func BuildBech32Addr(hrp string, addr Address) (string, error) {
    return bech32.ConvertAndEncode(hrp, addr.Bytes())
}

// MustBuildBech32Addr encodes the given human-readable-part string and
// address into a bech32 address.  It panics on error.
func MustBuildBech32Addr(hrp string, addr Address) string {
    b32, err := BuildBech32Addr(hrp, addr)
    if err != nil {
        panic(err)
    }
    return b32
}

// Bech32AddressHRP is the human-readable part of the Harmony address used by
// this process.
var Bech32AddressHRP = "hyn"

// Bech32ToAddress decodes the given bech32 address.
func Bech32ToAddress(b32 string) (addr Address, err error) {
    var hrp string
    err = ParseBech32Addr(b32, &hrp, &addr)
    if err == nil && hrp != Bech32AddressHRP {
        err = errors.Errorf("%#v is not a %#v address", b32, Bech32AddressHRP)
    }
    return
}

// MustBech32ToAddress decodes the given bech32 address.  It panics on error.
func MustBech32ToAddress(b32 string) Address {
    addr, err := Bech32ToAddress(b32)
    if err != nil {
        panic(err)
    }
    return addr
}

// AddressToBech32 encodes the given address into bech32 format.
func AddressToBech32(addr Address) (string, error) {
    return BuildBech32Addr(Bech32AddressHRP, addr)
}

// MustAddressToBech32 encodes the given address into bech32 format.
// It panics on error.
func MustAddressToBech32(addr Address) string {
    b32, err := BuildBech32Addr(Bech32AddressHRP, addr)
    if err != nil {
        panic(err)
    }
    return b32
}

// ParseAddr parses the given address, either as bech32 or as hex.
// The result can be 0x00..00 if the passing param is not a correct address.
func ParseAddr(s string) Address {
    if addr, err := Bech32ToAddress(s); err == nil {
        return addr
    }
    // The result can be 0x00...00 if the passing param is not a correct address.
    return HexToAddress(s)
}
