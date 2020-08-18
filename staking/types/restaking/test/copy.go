package staketest

import (
	"github.com/ethereum/go-ethereum/common"
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
	if w.TotalDelegation != nil {
		cp.TotalDelegation = new(big.Int).Set(w.TotalDelegation)
	}
	if w.TotalDelegationByOperator != nil {
		cp.TotalDelegationByOperator = new(big.Int).Set(w.TotalDelegationByOperator)
	}
	return cp
}

// CopyValidator deep copies restaking.Validator
func CopyValidator(v restaking.Validator_) restaking.Validator_ {
	cp := restaking.Validator_{
		ValidatorAddress:  v.ValidatorAddress,
		OperatorAddresses: CopyAddressSet(v.OperatorAddresses),
		Status:            v.Status,
		Commission:        CopyCommission(v.Commission),
		Description:       v.Description,
	}
	cp.SlotPubKeys = CopySlotPubKeys(v.SlotPubKeys)
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

// CopyAddressSet deep copy the AddressSet
func CopyAddressSet(s restaking.AddressSet_) restaking.AddressSet_ {
	cp := restaking.AddressSet_{}
	if s.Keys != nil {
		cp.Keys = []*common.Address{}
		for _, key := range s.Keys {
			k := *key
			cp.Keys = append(cp.Keys, &k)
		}
	}

	if s.Set != nil {
		cp.Set = make(map[common.Address]*bool)
		for addr, bo := range s.Set {
			b := *bo
			cp.Set[addr] = &b
		}
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

func CopySlotPubKeys(blsKeys restaking.BLSPublicKeys_) restaking.BLSPublicKeys_ {
	cp := restaking.BLSPublicKeys_{}
	if blsKeys.Keys != nil {
		cp.Keys = make([]*restaking.BLSPublicKey_, 0)
		for i := 0; i < len(blsKeys.Keys); i++ {
			c := CopySlotPubKey(*blsKeys.Keys[i])
			cp.Keys = append(cp.Keys, &c)
		}
	}
	return cp
}

func CopySlotPubKey(blsKey restaking.BLSPublicKey_) restaking.BLSPublicKey_ {
	cp := restaking.BLSPublicKey_{Key: blsKey.Key}
	return cp
}

// CopyDelegations deeps copy restaking.Delegations
func CopyDelegationMap(ds restaking.RedelegationMap_) restaking.RedelegationMap_ {
	if ds.Keys == nil {
		return restaking.RedelegationMap_{}
	}
	cp := restaking.NewRedelegationMap()
	if len(ds.Keys) == 0 {
		return cp
	}
	for _, key := range ds.Keys {
		d, _ := ds.Get(*key)
		cd := CopyRedelegation(d)
		cp.Put(*key, cd)
	}
	return cp
}

// CopyRedelegation copies restaking.Redelegation_
func CopyRedelegation(d restaking.Redelegation_) restaking.Redelegation_ {
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
