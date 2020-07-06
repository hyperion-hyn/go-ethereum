package committee

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/types"
	"math/big"
)

const (
	CommitteeSize = 88
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

type CommitteeByEpoch map[uint64]Committee