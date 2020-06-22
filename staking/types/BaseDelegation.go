package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type BaseDelegation struct {
	DelegatorAddress common.Address
	Amount           *big.Int
	Reward           *big.Int
	Undelegations    Undelegations
}

// Undelegation represents one undelegation entry
type Undelegation struct {
	Amount *big.Int `json:"amount"`
	Epoch  *big.Int `json:"epoch"`
}

// Undelegations ..
type Undelegations []Undelegation

// String ..
func (u Undelegations) String() string {
	s, _ := json.Marshal(u)
	return string(s)
}
