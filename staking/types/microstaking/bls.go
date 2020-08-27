package microstaking

import (
	"bytes"
	"encoding/hex"
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

// Storage_BLSPublicKeys_
func (s *Storage_BLSPublicKeys_) Length() int {
	return s.Keys().Length()
}

func (s *Storage_BLSPublicKeys_) Save(keys *BLSPublicKeys_) {
	length := len(keys.Keys)
	s.Keys().Resize(length)
	for i := 0; i < length; i++ {
		s.Keys().Get(i).Key().SetValue(keys.Keys[i].Key)
	}
}

func (s *Storage_BLSPublicKeys_) Get(index int) *BLSPublicKey_ {
	s.Keys().Get(index).Key().Value()
	return s.Keys().Get(index).obj
}

func (s *Storage_BLSPublicKeys_) Set(index int, key *BLSPublicKey_) {
	s.Keys().Get(index).Key().SetValue(key.Key)
}

func (s *Storage_BLSPublicKeys_) Remove(index int, keepOrder bool) {
	//remove current
	length := s.Length()
	lastOneStorage := s.Keys().Get(length - 1)
	//remove lastOne
	s.Keys().Get(length - 1).Key().SetValue([48]uint8{})
	//replace lastOne to index
	s.Keys().Get(index).Key().SetValue(lastOneStorage.Key().Value())
	//resize length
	s.Keys().Resize(length - 1)
}

func (s *Storage_BLSPublicKeys_) Push(key *BLSPublicKey_) {
	length := s.Length()

	//over length will auto resize , not resize again
	s.Keys().Get(length).Key().SetValue(key.Key)
}

func (s *Storage_BLSPublicKeys_) Pop() *BLSPublicKey_ {
	length := s.Length()

	blsPublicKeyTemp :=
		BLSPublicKey_{Key: s.Keys().Get(length - 1).Key().Value()}

	s.Keys().Get(length - 1).Key().SetValue([48]uint8{})
	s.Keys().Resize(length - 1)
	return &blsPublicKeyTemp
}

func (s *Storage_BLSPublicKeys_) Load() *BLSPublicKeys_ {
	length := s.Length()

	for i := 0; i < length; i++ {
		s.Keys().Get(i).Key().Value()
	}
	return s.obj
}
