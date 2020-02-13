package staking

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
)

//go:generate gencodec -type Delegation -field-override delegationMarshaling -out gen_delegation.go

// Delegation represents the bond with tokens held by an account. It is
// owned by one delegator, and is associated with the voting power of one
// validator.
type Delegation struct {
	DelegatorAddress common.Address `json:"delegator_address" gencodec:"required"`
	Amount           *big.Int       `json:"amount" gencodec:"required"`
	Reward           *big.Int       `json:"reward"`
	Undelegations    []Undelegation `json:"undelegations"`
	// TODO more fields
}

type Delegations []Delegation

type delegationMarshaling struct {
	DelegatorAddress common.UnprefixedAddress
	Amount	*math.HexOrDecimal256
}

// Undelegation represents one undelegation entry
type Undelegation struct {
	Amount *big.Int
	Epoch  *big.Int
	// TODO more fields
}

// NewDelegation creates a new delegation object
func NewDelegation(delegatorAddr common.Address, amount *big.Int) Delegation {
	return Delegation{
		DelegatorAddress: delegatorAddr,
		Amount:           amount,
	}
}
