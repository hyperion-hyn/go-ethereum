package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/pkg/errors"
	"math/big"
	"sort"
)

//go:generate gencodec -type Delegation -field-override delegationMarshaling -out gen_delegation.go

const (
	// ValidatorLockPeriodInBlock is the number of epochs a undelegated token needs to be before it's released to the delegator's balance
	ValidatorLockPeriodInBlock = 5 * 1
	DelegatorLockPeriodInBlock = 5
)

var (
	errInsufficientBalance = errors.New("Insufficient balance to undelegate")
	errInvalidAmount       = errors.New("Invalid amount, must be positive")
)

// Delegation represents the bond with tokens held by an account. It is
// owned by one delegator, and is associated with the voting power of one
// validator.
type Delegation struct {
	DelegatorAddress common.Address `json:"delegator_address" gencodec:"required"`
	Amount           *big.Int       `json:"amount" gencodec:"required"`
	Reward           *big.Int       `json:"-"`
	Undelegations    []Undelegation `json:"-"`
}

type Delegations []Delegation

// Undelegation represents one undelegation entry
type Undelegation struct {
	Amount      *big.Int
	UnlockBlock *big.Int
}

// DelegationIndex stored the index of a delegation in the validator's delegation list
type DelegationIndex struct {
	ValidatorAddress common.Address
	Index            uint64
}

// field type overrides for gencodec
type delegationMarshaling struct {
	DelegatorAddress common.UnprefixedAddress
	Amount           *math.HexOrDecimal256
}

// NewDelegation creates a new delegation object
func NewDelegation(delegatorAddr common.Address, amount *big.Int) Delegation {
	return Delegation{
		DelegatorAddress: delegatorAddr,
		Amount:           amount,
	}
}

// Undelegate - append entry to the undelegation
func (d *Delegation) Undelegate(blockNum *big.Int, amt *big.Int, isValidator bool) error {
	if amt.Sign() <= 0 {
		return errInvalidAmount
	}
	if d.Amount.Cmp(amt) < 0 {
		return errInsufficientBalance
	}
	d.Amount.Sub(d.Amount, amt)
	unlockNum := big.NewInt(0).Set(blockNum).Add(blockNum, big.NewInt(DelegatorLockPeriodInBlock))
	if isValidator {
		unlockNum = big.NewInt(0).Set(blockNum).Add(blockNum, big.NewInt(ValidatorLockPeriodInBlock))
	}
	item := Undelegation{amt, unlockNum}
	d.Undelegations = append(d.Undelegations, item)

	// Always sort the undelegate by block num in increasing order
	sort.SliceStable(
		d.Undelegations,
		func(i, j int) bool { return d.Undelegations[i].UnlockBlock.Cmp(d.Undelegations[j].UnlockBlock) < 0 },
	)
	return nil
}

// TotalInUndelegation - return the total amount of token in undelegation (locking period)
func (d *Delegation) TotalInUndelegation() *big.Int {
	total := big.NewInt(0)
	for _, entry := range d.Undelegations {
		total.Add(total, entry.Amount)
	}
	return total
}

// RemoveUnlockedUndelegations removes all fully unlocked undelegations and returns the total sum
func (d *Delegation) RemoveUnlockedUndelegations(blockNum *big.Int) *big.Int {
	totalWithdraw := big.NewInt(0)
	count := 0
	for j := range d.Undelegations {
		if d.Undelegations[j].UnlockBlock.Cmp(blockNum) <= 0 {
			totalWithdraw.Add(totalWithdraw, d.Undelegations[j].Amount)
			count++
		} else {
			break
		}
	}
	d.Undelegations = d.Undelegations[count:]
	return totalWithdraw
}
