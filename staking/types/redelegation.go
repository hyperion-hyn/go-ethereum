package types

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	errInsufficientBalance = errors.New("insufficient balance to undelegate")
	errInvalidAmount       = errors.New("invalid amount, must be positive")
)

const (
	// LockPeriodInEpoch is the number of epochs a undelegated token needs to be before it's released to the delegator's balance
	LockPeriodInEpoch = 7
)

// Redelegation represents the bond with tokens held by an account. It is
// owned by one delegator, and is associated with the voting power of one
// validator.
type Redelegation struct {
	DelegatorAddress common.Address
	Amount           *big.Int
	Reward           *big.Int
	Undelegation     Undelegation
}

// Redelegations ..
type Redelegations map[common.Address]Redelegation

func (d Redelegation) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

type RedelegationReference struct {
	ValidatorAddress        common.Address
	ReleasedTotalDelegation *big.Int // for a portion of released amount
}

// NewDelegation creates a new delegation object
func NewRedelegation(delegatorAddr common.Address, amount *big.Int) Redelegation {
	return Redelegation{
		DelegatorAddress: delegatorAddr,
		Amount:           amount,
		Reward:           big.NewInt(0),
	}
}
