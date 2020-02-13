package staking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Delegation represents the bond with tokens held by an account. It is
// owned by one delegator, and is associated with the voting power of one
// validator.
type Delegation struct {
	DelegatorAddress common.Address `json:"delegator_address"`
	Amount           *big.Int       `json:"amount"`
	Reward           *big.Int       `json:"reward"`
	Undelegations    []Undelegation `json:"undelegations"`
}

// Undelegation represents one undelegation entry
type Undelegation struct {
	Amount *big.Int
	Epoch  *big.Int
}

type Delegations []Delegation

// NewDelegation creates a new delegation object
func NewDelegation(delegatorAddr common.Address, amount *big.Int) Delegation {
	return Delegation{
		DelegatorAddress: delegatorAddr,
		Amount:           amount,
	}
}

func (delegations Delegations) Amount() *big.Int {
	amount := big.NewInt(0)
	for _, d := range delegations {
		amount = amount.Add(amount, d.Amount)
	}
	return amount
}