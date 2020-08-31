package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/bls/ffi/go/bls"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/singleflight"
	"math/big"
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

func (s *Storage_Slots_) Load() []*Slot_ {
	slotsLength := s.Length()
	for i := 0; i < slotsLength; i++ {
		s.Get(i).Load()
	}
	return s.obj.Entrys
}

func (s *Storage_Slots_) Get(index int) *Storage_Slot_ {
	return s.Entrys().Get(index)
}

func (s *Storage_Slots_) Set(index int, key *Slot_) {
	s.Entrys().Get(index).Save(key)
}

func (s *Storage_Slots_) Remove(index int, keepOrder bool) {
	// remove from index
	oldEntriesLength := s.Entrys().Length()

	//set lastEntity to index
	lastEntry := s.Entrys().Get(oldEntriesLength - 1)
	if oldEntriesLength > 1 {
		s.Entrys().Get(index).Save(&Slot_{
			EcdsaAddress: lastEntry.EcdsaAddress().Value(),
			BLSPublicKey: BLSPublicKey_{
				Key: lastEntry.BLSPublicKey().Key().Value(),
			},
			EffectiveStake: lastEntry.EffectiveStake().Value(),
		})
	}

	//set lastEntity to zero
	lastEntry.Clean()

	//resize slice
	s.Entrys().Resize(oldEntriesLength - 1)
}

func (s *Storage_Slots_) Push(slot *Slot_) {
	entityLength := s.Entrys().Length()
	s.Entrys().Get(entityLength).Save(slot)
}

func (s *Storage_Slots_) Pop() *Storage_Slot_ {
	entityLength := s.Entrys().Length()
	storageSlot := s.Entrys().Get(entityLength - 1)
	s.Remove(entityLength-1, false)
	return storageSlot
}

func (s *Storage_Slots_) UpdateSlots(slots Slots_) {
	// remove old
	length := s.Length()
	for i := 0; i < length; i++ {
		s.Get(i).Clean()
	}
	//set new
	newSlotsLength := len(slots.Entrys)
	s.Entrys().Resize(newSlotsLength)
	for i := 0; i < newSlotsLength; i++ {
		s.Get(i).Save(slots.Entrys[i])
	}
}

func (s *Storage_Slot_) Clean() {
	s.EffectiveStake().SetValue(common.NewDec(int64(0)))
	s.EcdsaAddress().SetValue(common.BigToAddress(big.NewInt(0)))
	s.BLSPublicKey().Key().SetValue([48]uint8{})
}

func (s *Storage_Slot_) Save(key *Slot_) {
	s.BLSPublicKey().Key().SetValue(key.BLSPublicKey.Key)
	s.EcdsaAddress().SetValue(key.EcdsaAddress)
	if !key.EffectiveStake.IsNil() {
		s.EffectiveStake().SetValue(key.EffectiveStake)
	}
}

func (s *Storage_Slot_) Load() *Slot_ {
	s.BLSPublicKey().Key().Value()
	s.EcdsaAddress().Value()
	s.EffectiveStake().Value()
	return s.obj
}

// Storage_Committee_
func (s *Storage_Committee_) Load() *Committee_ {
	s.Epoch().Value()
	s.Slots().Load()
	return s.obj
}

var (
	blsKeyCache singleflight.Group
)

func lookupBLSPublicKeys(
	c *Committee_,
) ([]*bls.PublicKey, error) {
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
