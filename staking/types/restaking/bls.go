package restaking

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/pkg/errors"
	"math/big"
)

var (
	emptyBLSPubKey = BLSPublicKey_{}
)

func NewEmptyBLSKeys() BLSPublicKeys_ {
	return BLSPublicKeys_{Keys: make([]*BLSPublicKey_, 0)}
}

func NewBLSKeysWithBLSKey(key BLSPublicKey_) BLSPublicKeys_ {
	keys := NewEmptyBLSKeys()
	keys.Keys = append(keys.Keys, &key)
	return keys
}

// Big ..
func (pk BLSPublicKey_) Big() *big.Int {
	return new(big.Int).SetBytes(pk.Key[:])
}

// IsEmpty returns whether the bls public key is empty 0 bytes
func (pk BLSPublicKey_) IsEmpty() bool {
	return bytes.Equal(pk.Key[:], emptyBLSPubKey.Key[:])
}

// Hex returns the hex string of bls public key
func (pk BLSPublicKey_) Hex() string {
	return hex.EncodeToString(pk.Key[:])
}

// ToLibBLSPublicKey copies the key contents into the given key.
func (pk *BLSPublicKey_) ToLibBLSPublicKey(key *bls.PublicKey) error {
	return key.Deserialize(pk.Key[:])
}

// FromLibBLSPublicKey replaces the key contents with the given key,
func (pk *BLSPublicKey_) FromLibBLSPublicKey(key *bls.PublicKey) error {
	bs := key.Serialize()
	if len(bs) != len(pk.Key) {
		return errors.Errorf(
			"key size (BLS) size mismatch, expected %d have %d", len(pk.Key), len(bs),
		)
	}
	copy(pk.Key[:], bs)
	return nil
}

func (pk *BLSPublicKey_) UnmarshalJSON(data []byte) error {
	var h string
	if err := json.Unmarshal(data, &h); err != nil {
		return err
	}
	key := bls.PublicKey{}
	if err := key.DeserializeHexStr(h); err != nil {
		return err
	}
	return pk.FromLibBLSPublicKey(&key)
}

func (pk *BLSPublicKey_) MarshalJSON() ([]byte, error) {
	return json.Marshal(pk.Hex())
}

// Storage_BLSPublicKey_
func (s *Storage_BLSPublicKey_) Equal(key_ *BLSPublicKey_) bool {
	return s.Key().Value() == key_.Key
}

// Storage_BLSPublicKeys_
func (s *Storage_BLSPublicKeys_) Length() int {
	return s.Keys().Length()
}

func (s *Storage_BLSPublicKeys_) Get(index int) *Storage_BLSPublicKey_ {
	if index < 0 || index >= s.Length() {
		panic("out of range")
	}
	return s.Keys().Get(index)
}

func (s *Storage_BLSPublicKeys_) Set(index int, key *BLSPublicKey_) {
	s.Keys().Get(index).Clear()
	s.Keys().Get(index).Save(key)
}

func (s *Storage_BLSPublicKeys_) Remove(index int) {
	//remove current
	length := s.Length()
	if index < 0 || index >= length {
		panic("out of range")
	}

	//replace lastOne to index if length !=1
	if length > 1 {
		lastOne := s.Keys().Get(length - 1).load()
		s.Set(index, lastOne)
	}
	s.Keys().Get(length - 1).Clear() //remove lastOne
	s.Keys().Resize(length - 1)	//resize length
}

func (s *Storage_BLSPublicKeys_) Push(key *BLSPublicKey_) {
	length := s.Length()
	s.Set(length, key)
}
