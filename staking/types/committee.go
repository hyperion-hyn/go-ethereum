package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/numeric"
	"math/big"
)

const (
	CommitteeSize = 88
)

// Slot represents node id (BLS address)
type Slot struct {
	Address        common.Address `json:"address"`
	EffectiveStake numeric.Dec    `json:"effective-stake"`
}

func (s Slot) String() string {
	return fmt.Sprintf("%s:%s", s.Address.String(), s.EffectiveStake.String())
}

// SlotList is a list of Slot.
type SlotList []Slot

type Committee struct {
	BlockNum *big.Int
	Slots    SlotList `json:"subcommittee"`
}
