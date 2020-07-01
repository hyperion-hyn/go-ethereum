package types

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"sort"
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
	Undelegations    Undelegations
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
func NewRedelegation(delegatorAddr common.Address,
	amount *big.Int) Redelegation {
	return Redelegation{
		DelegatorAddress: delegatorAddr,
		Amount:           amount,
		Reward:           big.NewInt(0),
	}
}

// Unredelegate - append entry to the undelegation
func (d *Redelegation) Unredelegate(epoch *big.Int, amt *big.Int) error {
	if amt.Sign() <= 0 {
		return errInvalidAmount
	}
	if d.Amount.Cmp(amt) < 0 {
		return errInsufficientBalance
	}
	d.Amount.Sub(d.Amount, amt)

	exist := false
	for _, entry := range d.Undelegations {
		if entry.Epoch.Cmp(epoch) == 0 {
			exist = true
			entry.Amount.Add(entry.Amount, amt)
			return nil
		}
	}

	if !exist {
		item := Undelegation{amt, epoch}
		d.Undelegations = append(d.Undelegations, item)

		// Always sort the undelegate by epoch in increasing order
		sort.SliceStable(
			d.Undelegations,
			func(i, j int) bool { return d.Undelegations[i].Epoch.Cmp(d.Undelegations[j].Epoch) < 0 },
		)
	}

	return nil
}

// TotalInUndelegation - return the total amount of token in undelegation (locking period)
func (d *Redelegation) TotalInUndelegation() *big.Int {
	total := big.NewInt(0)
	for i := range d.Undelegations {
		total.Add(total, d.Undelegations[i].Amount)
	}
	return total
}

// DeleteEntry - delete an entry from the undelegation
// Opimize it
func (d *Redelegation) DeleteEntry(epoch *big.Int) {
	entries := []Undelegation{}
	for i, entry := range d.Undelegations {
		if entry.Epoch.Cmp(epoch) == 0 {
			entries = append(d.Undelegations[:i], d.Undelegations[i+1:]...)
		}
	}
	if entries != nil {
		d.Undelegations = entries
	}
}

// RemoveUnlockedUndelegations removes all fully unlocked
// undelegations and returns the total sum
func (d *Redelegation) RemoveUnlockedUndelegations(
	curEpoch, lastEpochInCommittee *big.Int, lockPeriod int,
) *big.Int {
	totalWithdraw := big.NewInt(0)
	count := 0
	for j := range d.Undelegations {
		if big.NewInt(0).Sub(curEpoch, d.Undelegations[j].Epoch).Int64() >= int64(lockPeriod) ||
			big.NewInt(0).Sub(curEpoch, lastEpochInCommittee).Int64() >= int64(lockPeriod) {
			// need to wait at least 7 epochs to withdraw; or the validator has been out of committee for 7 epochs
			totalWithdraw.Add(totalWithdraw, d.Undelegations[j].Amount)
			count++
		} else {
			break
		}
	}
	d.Undelegations = d.Undelegations[count:]
	return totalWithdraw
}
