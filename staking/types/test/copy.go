package staketest

import (
	"github.com/ethereum/go-ethereum/staking/types"
	"math/big"
)

// CopyValidatorWrapper deep copies types.ValidatorWrapper
func CopyValidatorWrapper(w types.ValidatorWrapper) types.ValidatorWrapper {
	cp := types.ValidatorWrapper{
		Validator:     CopyValidator(w.Validator),
		Redelegations: CopyDelegations(w.Redelegations),
	}
	if w.Counters.NumBlocksSigned != nil {
		cp.Counters.NumBlocksSigned = new(big.Int).Set(w.Counters.NumBlocksSigned)
	}
	if w.Counters.NumBlocksToSign != nil {
		cp.Counters.NumBlocksToSign = new(big.Int).Set(w.Counters.NumBlocksToSign)
	}
	if w.BlockReward != nil {
		cp.BlockReward = new(big.Int).Set(w.BlockReward)
	}
	return cp
}

// CopyValidator deep copies types.Validator
func CopyValidator(v types.Validator) types.Validator {
	cp := types.Validator{
		Address:     v.Address,
		Status:      v.Status,
		Commission:  CopyCommission(v.Commission),
		Description: v.Description,
	}
	if v.SlotPubKeys != nil {
		cp.SlotPubKeys = make([]types.BLSPublicKey, len(v.SlotPubKeys))
		copy(cp.SlotPubKeys, v.SlotPubKeys)
	}
	if v.LastEpochInCommittee != nil {
		cp.LastEpochInCommittee = new(big.Int).Set(v.LastEpochInCommittee)
	}
	if v.MinSelfDelegation != nil {
		cp.MinSelfDelegation = new(big.Int).Set(v.MinSelfDelegation)
	}
	if v.MaxTotalDelegation != nil {
		cp.MaxTotalDelegation = new(big.Int).Set(v.MaxTotalDelegation)
	}
	if v.CreationHeight != nil {
		cp.CreationHeight = new(big.Int).Set(v.CreationHeight)
	}
	return cp
}

// CopyCommission deep copy the Commission
func CopyCommission(c types.Commission) types.Commission {
	cp := types.Commission{
		CommissionRates: c.CommissionRates.Copy(),
	}
	if c.UpdateHeight != nil {
		cp.UpdateHeight = new(big.Int).Set(c.UpdateHeight)
	}
	return cp
}

// CopyDelegations deeps copy types.Redelegations
func CopyDelegations(ds types.Redelegations) types.Redelegations {
	if ds == nil {
		return nil
	}
	cp := make(types.Redelegations, 0, len(ds))
	for _, d := range ds {
		cp = append(cp, CopyDelegation(d))
	}
	return cp
}

// CopyDelegation copies types.Redelegation
func CopyDelegation(d types.Redelegation) types.Redelegation {
	cp := types.Redelegation{
		DelegatorAddress: d.DelegatorAddress,
		Undelegations:    CopyUndelegations(d.Undelegations),
	}
	if d.Amount != nil {
		cp.Amount = new(big.Int).Set(d.Amount)
	}
	if d.Reward != nil {
		cp.Reward = new(big.Int).Set(d.Reward)
	}
	return cp
}

// CopyUndelegations deep copies types.Undelegations
func CopyUndelegations(uds types.Undelegations) types.Undelegations {
	if uds == nil {
		return nil
	}
	cp := make(types.Undelegations, 0, len(uds))
	for _, ud := range uds {
		cp = append(cp, CopyUndelegation(ud))
	}
	return cp
}

// CopyUndelegation deep copies types.Undelegation
func CopyUndelegation(ud types.Undelegation) types.Undelegation {
	cp := types.Undelegation{}
	if ud.Amount != nil {
		cp.Amount = new(big.Int).Set(ud.Amount)
	}
	if ud.Epoch != nil {
		cp.Epoch = new(big.Int).Set(ud.Epoch)
	}
	return cp
}
