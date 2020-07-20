package committee

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/types"
	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/singleflight"
	"math/big"
	"time"
)

const (
	CommitteeSize = 88
)

var (
	// ErrValidNotInCommittee ..
	ErrValidNotInCommittee = errors.New("slot signer not this slot's subcommittee")
	// ErrCommitteeNil ..
	ErrCommitteeNil = errors.New("subcommittee is nil pointer")
)

// Slot represents node id (BLS address)
type Slot struct {
	EcdsaAddress common.Address     `json:"ecdsa-address"`
	BLSPublicKey types.BLSPublicKey `json:"bls-pubkey"`
	// nil means our node, 0 means not active, > 0 means staked node
	EffectiveStake *numeric.Dec `json:"effective-stake" rlp:"nil"`
}

func (s Slot) String() string {
	return fmt.Sprintf("%s:%s", s.EcdsaAddress.String(), s.EffectiveStake.String())
}

// SlotList is a list of Slot.
type SlotList []Slot

func (l SlotList) String() string {
	blsKeys := make([]string, len(l))
	for i, k := range l {
		blsKeys[i] = k.BLSPublicKey.Hex()
	}
	s, _ := json.Marshal(blsKeys)
	return string(s)
}

type Committee struct {
	Epoch *big.Int
	Slots SlotList
}

// BLSPublicKeys ..
func (c *Committee) BLSPublicKeys() ([]*bls.PublicKey, error) {
	if c == nil {
		return nil, ErrCommitteeNil
	}
	return lookupBLSPublicKeys(c)
}

func (c *Committee) Hash() (h common.Hash) {
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
	TotalEffectiveStaked numeric.Dec
}

// StakedValidators ..
func (c Committee) StakedValidators() *StakedSlots {
	countStakedValidator, countStakedBLSKey := 0, 0
	networkWideSlice, networkWideSet :=
		[]common.Address{}, map[common.Address]struct{}{}
	totalEffectiveStake := numeric.ZeroDec()

	for _, slot := range c.Slots {
		// an external validator,
		// non-nil EffectiveStake is how we known
		if addr := slot.EcdsaAddress; slot.EffectiveStake != nil {
			totalEffectiveStake = totalEffectiveStake.Add(*slot.EffectiveStake)
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

var (
	blsKeyCache singleflight.Group
)

func lookupBLSPublicKeys(
	c *Committee,
) ([]*bls.PublicKey, error) {
	key := c.Hash().Hex()
	results, err, _ := blsKeyCache.Do(
		key, func() (interface{}, error) {
			slice := make([]*bls.PublicKey, len(c.Slots))
			for j := range c.Slots {
				committerKey := &bls.PublicKey{}
				if err := c.Slots[j].BLSPublicKey.ToLibBLSPublicKey(
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