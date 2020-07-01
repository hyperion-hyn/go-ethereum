package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	PendingDelegationLockPeriodInEpoch = 7
)

// Microdelegation represents the bond with tokens held by an account. It is
// owned by one delegator, and is associated with the voting power of one
// validator.
type Microdelegation struct {
	DelegatorAddress   common.Address
	Amount             *big.Int
	Reward             *big.Int
	Undelegations      Undelegations
	PendingDelegations PendingDelegations
	AutoRenew          bool
}

// Microdelegations ..
type Microdelegations map[common.Address]Microdelegation

// String ..
func (d Microdelegations) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func (d Microdelegation) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

// Undelegation represents one undelegation entry
type Undelegation struct {
	Amount *big.Int `json:"amount"`
	Epoch  *big.Int `json:"epoch"`
}

// Undelegations ..
type Undelegations []Undelegation

// PendingDelegation represents tokens during map3 node in pending state
type PendingDelegation struct {
	Amount *big.Int
	Epoch  *big.Int
}

type PendingDelegations []PendingDelegation

// NewMicrodelegation creates a new microdelegation object
func NewMicrodelegation(
	delegator common.Address, amount, epoch *big.Int, autoRenew, pending bool,
) Microdelegation {
	d := Microdelegation{
		DelegatorAddress: delegator,
		Amount:           big.NewInt(0),
		Reward:           big.NewInt(0),
		AutoRenew:        autoRenew,
	}
	if pending {
		d.PendingDelegations = append(d.PendingDelegations, PendingDelegation{
			Amount: amount,
			Epoch:  epoch,
		})
	} else {
		d.Amount = amount
	}
	return d
}

// Unmicrodelegate - append entry to the unmicrodelegation
func (d *Microdelegation) Unmicrodelegate(epoch *big.Int, amt *big.Int) error {
	if amt.Sign() <= 0 {
		return errInvalidAmount
	}
	if d.Amount.Cmp(amt) < 0 {
		return errInsufficientBalance
	}
	d.Amount.Sub(d.Amount, amt)

	//exist := false
	//for _, entry := range d.Undelegations {
	//	if entry.Epoch.Cmp(epoch) == 0 {
	//		exist = true
	//		entry.Amount.Add(entry.Amount, amt)
	//		return nil
	//	}
	//}
	//
	//if !exist {
	//	item := Undelegation{amt, epoch}
	//	d.Undelegations = append(d.Undelegations, item)
	//
	//	// Always sort the undelegate by epoch in increasing order
	//	sort.SliceStable(
	//		d.Undelegations,
	//		func(i, j int) bool { return d.Undelegations[i].Epoch.Cmp(d.Undelegations[j].Epoch) < 0 },
	//	)
	//}

	return nil
}

// DeleteEntry - delete an entry from the undelegation
// Opimize it
func (d *Microdelegation) DeleteEntry(epoch *big.Int) {
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
