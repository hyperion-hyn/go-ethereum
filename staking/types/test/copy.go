package staketest

import (
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

// CopyValidatorWrapper deep copies staking.ValidatorWrapper
func CopyValidatorWrapper(w restaking.ValidatorWrapper_) restaking.ValidatorWrapper_ {
	cp := restaking.ValidatorWrapper_{
		Validator:     CopyValidator(w.Validator),
		Redelegations: CopyDelegationMap(w.Redelegations),
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

// CopyValidator deep copies restaking.Validator
func CopyValidator(v restaking.Validator_) restaking.Validator_ {
	cp := restaking.Validator_{
		ValidatorAddress: v.ValidatorAddress,
		Status:           v.Status,
		Commission:       CopyCommission(v.Commission),
		Description:      v.Description,
	}
	if v.SlotPubKeys.Keys != nil {
		cp.SlotPubKeys = restaking.BLSPublicKeys_{Keys:make([]*restaking.BLSPublicKey_, len(v.SlotPubKeys.Keys))}
		copy(cp.SlotPubKeys.Keys, v.SlotPubKeys.Keys)
	}
	if v.LastEpochInCommittee != nil {
		cp.LastEpochInCommittee = new(big.Int).Set(v.LastEpochInCommittee)
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
func CopyCommission(c restaking.Commission_) restaking.Commission_ {
	cp := restaking.Commission_{
		CommissionRates: c.CommissionRates.Copy(),
	}
	if c.UpdateHeight != nil {
		cp.UpdateHeight = new(big.Int).Set(c.UpdateHeight)
	}
	return cp
}

// CopyDelegations deeps copy restaking.Delegations
func CopyDelegationMap(ds restaking.RedelegationMap_) restaking.RedelegationMap_ {
	cp := restaking.NewRelegationMap()
	if ds.Keys == nil || len(ds.Keys) == 0 {
		return cp
	}
	for _, key := range ds.Keys {
		d := ds.Get(*key)
		cd := CopyDelegation(*d)
		cp.Put(*key, &cd)
	}
	return cp
}

// CopyDelegation copies restaking.Redelegation_
func CopyDelegation(d restaking.Redelegation_) restaking.Redelegation_ {
	cp := restaking.Redelegation_{
		DelegatorAddress: d.DelegatorAddress,
		Undelegation:     CopyUndelegation(d.Undelegation),
	}
	if d.Amount != nil {
		cp.Amount = new(big.Int).Set(d.Amount)
	}
	if d.Reward != nil {
		cp.Reward = new(big.Int).Set(d.Reward)
	}
	return cp
}

// CopyUndelegation deep copies restaking.Undelegation
func CopyUndelegation(ud restaking.Undelegation_) restaking.Undelegation_ {
	cp := restaking.Undelegation_{}
	if ud.Amount != nil {
		cp.Amount = new(big.Int).Set(ud.Amount)
	}
	if ud.Epoch != nil {
		cp.Epoch = new(big.Int).Set(ud.Epoch)
	}
	return cp
}
