package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/bls/ffi/go/bls"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/singleflight"
	"time"
)

// BLSPublicKeys ..
func (c *Committee_) BLSPublicKeys() ([]*bls.PublicKey, error) {
	if c == nil {
		return nil, ErrCommitteeNil
	}
	return lookupBLSPublicKeys(c)
}

func (c *Committee_) Hash() (h common.Hash) {
	hw := sha3.New256()
	rlp.Encode(hw, c)
	hw.Sum(h[:0])
	return h
}

// StakedSlots gives overview of members
// in a subcommittee (aka a shard)
type StakedSlots struct {
	CountStakedValidator int
	CountStakedBLSKey    int
	Addrs                []common.Address
	LookupSet            map[common.Address]struct{}
	TotalEffectiveStaked common.Dec
}

// StakedValidators ..
func (c *Committee_) StakedValidators() *StakedSlots {
	countStakedValidator, countStakedBLSKey := 0, 0
	networkWideSlice, networkWideSet :=
		[]common.Address{}, map[common.Address]struct{}{}
	totalEffectiveStake := common.ZeroDec()

	for _, slot := range c.Slots.Entrys {
		// an external validator,
		// non-nil EffectiveStake is how we known
		if addr := slot.EcdsaAddress; slot.EffectiveStake.IsPositive() {
			totalEffectiveStake = totalEffectiveStake.Add(slot.EffectiveStake)
			countStakedBLSKey++
			if _, seen := networkWideSet[addr]; !seen {
				countStakedValidator++
				networkWideSet[addr] = struct{}{}
				networkWideSlice = append(networkWideSlice, addr)
			}
		}
	}

	return &StakedSlots{
		CountStakedValidator: countStakedValidator,
		CountStakedBLSKey:    countStakedBLSKey,
		Addrs:                networkWideSlice,
		LookupSet:            networkWideSet,
		TotalEffectiveStaked: totalEffectiveStake,
	}
}

// Storage_Slots_
func (s *Storage_Slots_) Length() int {
	return s.Entrys().Length()
}

func (s *Storage_Slots_) Get(index int) *Storage_Slot_ {
	if index < 0 || index >= s.Length() {
		panic("out of range")
	}
	return s.Entrys().Get(index)
}

func (s *Storage_Slots_) Set(index int, key *Slot_) {
	s.Entrys().Get(index).Clear()
	s.Entrys().Get(index).Save(key)
}

func (s *Storage_Slots_) Remove(index int) {
	// remove from index
	length := s.Length()
	if index < 0 || index >= length {
		panic("out of range")
	}

	//set lastEntity to index
	if length > 1 {
		lastOne := s.Entrys().Get(length - 1).load()
		s.Set(index, lastOne)
	}
	s.Entrys().Get(length - 1).Clear() //remove lastOne
	s.Entrys().Resize(length - 1)	//resize length
}

func (s *Storage_Slots_) Push(slot *Slot_) {
	length := s.Length()
	s.Set(length, slot)
}


var (
	blsKeyCache singleflight.Group
)

func lookupBLSPublicKeys(c *Committee_) ([]*bls.PublicKey, error) {
	key := c.Hash().Hex()
	results, err, _ := blsKeyCache.Do(
		key, func() (interface{}, error) {
			slice := make([]*bls.PublicKey, len(c.Slots.Entrys))
			for j := range c.Slots.Entrys {
				committerKey := &bls.PublicKey{}
				if err := c.Slots.Entrys[j].BLSPublicKey.ToLibBLSPublicKey(
					committerKey,
				); err != nil {
					return nil, err
				}
				slice[j] = committerKey
			}
			// Only made once
			go func() {
				time.Sleep(25 * time.Minute)
				blsKeyCache.Forget(key)
			}()
			return slice, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return results.([]*bls.PublicKey), nil
}
