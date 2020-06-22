package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Microdelegation represents the bond with tokens held by an account. It is
// owned by one delegator, and is associated with the voting power of one
// validator.
type Microdelegation struct {
	BaseDelegation
	PendingDelegations PendingDelegations
	IsAutoRenew        bool
}

// Microdelegations ..
type Microdelegations []Microdelegation

// String ..
func (d Microdelegations) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func (d Microdelegation) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

// PendingDelegation represents tokens during map3 node in pending state
type PendingDelegation struct {
	Amount *big.Int
	Epoch  *big.Int `json:"epoch"`
}

type PendingDelegations []PendingDelegation

//type MicrodelegationAugmenterEntries []MicrodelegationAugmenterEntry
//
//type MicrodelegationAugmenterEntry struct {
//	MicrodelegatorAddress common.Address
//	Amount                *big.Int
//}
//
//type MicrodelegationAugmenters []MicrodelegationAugmenter
//
//type MicrodelegationAugmenter struct {
//	MicrodelegationAugmenterEntries MicrodelegationAugmenterEntries
//	BlockHeight                     *big.Int
//	RewardSnapshot                  *big.Int
//}

// MicrodelegationIndexes is a slice of MicrodelegationIndex
type MicrodelegationIndexes []MicrodelegationIndex // TODO(ATLAS): need?

// MicrodelegationIndex stored the index of a delegation in the validator's delegation list
type MicrodelegationIndex struct {
	Map3NodeAddress common.Address
	Index           uint64
	BlockNum        *big.Int
}

// NewMicroDelegation creates a new microdelegation object
func NewMicroDelegation(delegatorAddr common.Address, amount *big.Int) Microdelegation {
	return Microdelegation{
		DelegatorAddress: delegatorAddr,
		Amount:           amount,
		Reward:           big.NewInt(0),
	}
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
