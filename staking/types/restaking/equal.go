package restaking

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// CheckValidatorWrapperEqual checks the equality of staking.ValidatorWrapper. If not equal, an
// error is returned. Note nil pointer is treated as zero in this compare function.
func CheckValidatorWrapperEqual(w1, w2 ValidatorWrapper_) error {
	if err := checkValidatorWrapperEqual(w1, w2); err != nil {
		return fmt.Errorf("ValidatorWrapper%v", err)
	}
	return nil
}

// CheckValidatorEqual checks the equality of validator. If not equal, an
// error is returned. Note nil pointer is treated as zero in this compare function.
func CheckValidatorEqual(v1, v2 Validator_) error {
	if err := checkValidatorEqual(v1, v2); err != nil {
		return fmt.Errorf("Validator%v", err)
	}
	return nil
}

func checkValidatorWrapperEqual(w1, w2 ValidatorWrapper_) error {
	if err := checkValidatorEqual(w1.Validator, w2.Validator); err != nil {
		return fmt.Errorf(".Validator%v", err)
	}
	if err := checkRedelegationMapEqual(w1.Redelegations, w2.Redelegations); err != nil {
		return fmt.Errorf(".Redelegations%v", err)
	}
	if err := checkBigIntEqual(w1.Counters.NumBlocksToSign, w2.Counters.NumBlocksToSign); err != nil {
		return fmt.Errorf(".Counters.NumBlocksToSign %v", err)
	}
	if err := checkBigIntEqual(w1.Counters.NumBlocksSigned, w2.Counters.NumBlocksSigned); err != nil {
		return fmt.Errorf(".Counters.NumBlocksSigned %v", err)
	}
	if err := checkBigIntEqual(w1.BlockReward, w2.BlockReward); err != nil {
		return fmt.Errorf(".BlockReward %v", err)
	}
	if err := checkBigIntEqual(w1.TotalDelegation, w2.TotalDelegation); err != nil {
		return fmt.Errorf(".TotalDelegation %v", err)
	}
	if err := checkBigIntEqual(w1.TotalDelegationFromOperators, w2.TotalDelegationFromOperators); err != nil {
		return fmt.Errorf(".TotalDelegationFromOperators %v", err)
	}
	return nil
}

func checkValidatorEqual(v1, v2 Validator_) error {
	if v1.ValidatorAddress != v2.ValidatorAddress {
		return fmt.Errorf(".Validator address not equal: %x / %x", v1.ValidatorAddress, v2.ValidatorAddress)
	}
	if err := checkAddressSetEqual(v1.OperatorAddresses, v2.OperatorAddresses); err != nil {
		return fmt.Errorf(".Operator addresses not equal %v", err)
	}
	if err := checkPubKeysEqual(v1.SlotPubKeys, v2.SlotPubKeys); err != nil {
		return fmt.Errorf(".SlotPubKeys%v", err)
	}
	if err := checkBigIntEqual(v1.LastEpochInCommittee, v2.LastEpochInCommittee); err != nil {
		return fmt.Errorf(".LastEpochInCommittee %v", err)
	}
	if err := checkBigIntEqual(v1.MaxTotalDelegation, v2.MaxTotalDelegation); err != nil {
		return fmt.Errorf(".MaxTotalDelegation %v", err)
	}
	if v1.Status != v2.Status {
		return fmt.Errorf(".Status not equal: %v / %v", v1.Status, v2.Status)
	}
	if err := checkCommissionEqual(v1.Commission, v2.Commission); err != nil {
		return fmt.Errorf(".Commission%v", err)
	}
	if err := checkDescriptionEqual(v1.Description, v2.Description); err != nil {
		return fmt.Errorf(".Description%v", err)
	}
	if err := checkBigIntEqual(v1.CreationHeight, v2.CreationHeight); err != nil {
		return fmt.Errorf(".CreationHeight %v", err)
	}
	return nil
}

func checkAddressSetEqual(a1, a2 IterableAddressSet_) error {
	if len(a1.Keys) != len(a2.Keys) {
		return fmt.Errorf(".len of keys not equal: %v / %v", len(a1.Keys), len(a2.Keys))
	}
	if len(a1.Map) != len(a2.Map) {
		return fmt.Errorf(".len of map not equal: %v / %v", len(a1.Map), len(a2.Map))
	}
	for i := range a1.Keys {
		if *(a1.Keys[i]) != *(a2.Keys[i]) {
			return fmt.Errorf("[%v] not equal in array: %x / %x", i, a1.Keys[i], a2.Keys[i])
		}
		k := *(a1.Keys[i])
		if *(a1.Map[k]) != *(a2.Map[k]) {
			return fmt.Errorf("[%v] not equal in map: %x / %x", k, a1.Map[k], a2.Map[k])
		}
	}
	return nil
}

func checkPubKeysEqual(pubs1, pubs2 BLSPublicKeys_) error {
	if len(pubs1.Keys) != len(pubs2.Keys) {
		return fmt.Errorf(".len not equal: %v / %v", len(pubs1.Keys), len(pubs2.Keys))
	}
	for i := range pubs1.Keys {
		if pubs1.Keys[i].Key != pubs2.Keys[i].Key {
			return fmt.Errorf("[%v] not equal: %x / %x", i, pubs1.Keys[i], pubs2.Keys[i])
		}
	}
	return nil
}

func checkDescriptionEqual(d1, d2 Description_) error {
	if d1.Name != d2.Name {
		return fmt.Errorf(".Name not equal: %v / %v", d1.Name, d2.Name)
	}
	if d1.Identity != d2.Identity {
		return fmt.Errorf(".Identity not equal: %v / %v", d1.Identity, d2.Identity)
	}
	if d1.Website != d2.Website {
		return fmt.Errorf(".Website not equal: %v / %v", d1.Website, d2.Website)
	}
	if d1.Details != d2.Details {
		return fmt.Errorf(".Details not equal: %v / %v", d1.Details, d2.Details)
	}
	if d1.SecurityContact != d2.SecurityContact {
		return fmt.Errorf(".SecurityContact not equal: %v / %v", d1.SecurityContact, d2.SecurityContact)
	}
	return nil
}

func checkCommissionEqual(c1, c2 Commission_) error {
	if err := checkCommissionRateEqual(c1.CommissionRates, c2.CommissionRates); err != nil {
		return fmt.Errorf(".CommissionRate%v", err)
	}
	if err := checkBigIntEqual(c1.UpdateHeight, c2.UpdateHeight); err != nil {
		return fmt.Errorf(".UpdateHeight %v", err)
	}
	return nil
}

func checkCommissionRateEqual(cr1, cr2 CommissionRates_) error {
	if err := checkDecEqual(cr1.Rate, cr2.Rate); err != nil {
		return fmt.Errorf(".Rate %v", err)
	}
	if err := checkDecEqual(cr1.MaxChangeRate, cr2.MaxChangeRate); err != nil {
		return fmt.Errorf(".MaxChangeRate %v", err)
	}
	if err := checkDecEqual(cr1.MaxRate, cr2.MaxRate); err != nil {
		return fmt.Errorf(".MaxRate %v", err)
	}
	return nil
}

func checkRedelegationMapEqual(ds1, ds2 IterableRedelegationMap_) error {
	if len(ds1.Keys) != len(ds2.Keys) {
		return fmt.Errorf(".len not equal: %v / %v", len(ds1.Keys), len(ds2.Keys))
	}
	for _, key := range ds1.Keys {
		r1, _ := ds1.Get(*key)
		r2, _ := ds2.Get(*key)
		if err := checkRedelegationEqual(r1, r2); err != nil {
			return fmt.Errorf("[%v]%v", key, err)
		}
	}
	return nil
}

func checkRedelegationEqual(d1, d2 Redelegation_) error {
	if d1.DelegatorAddress != d2.DelegatorAddress {
		return fmt.Errorf(".DelegatorAddress not equal: %x / %x",
			d1.DelegatorAddress, d2.DelegatorAddress)
	}
	if err := checkBigIntEqual(d1.Amount, d2.Amount); err != nil {
		return fmt.Errorf(".Amount %v", err)
	}
	if err := checkBigIntEqual(d1.Reward, d2.Reward); err != nil {
		return fmt.Errorf(".Reward %v", err)
	}
	if err := checkUndelegationEqual(d1.Undelegation, d2.Undelegation); err != nil {
		return fmt.Errorf(".Undelegation%v", err)
	}
	return nil
}

func checkUndelegationEqual(ud1, ud2 Undelegation_) error {
	if err := checkBigIntEqual(ud1.Amount, ud2.Amount); err != nil {
		return fmt.Errorf(".Amount %v", err)
	}
	if err := checkBigIntEqual(ud1.Epoch, ud2.Epoch); err != nil {
		return fmt.Errorf(".Epoch %v", err)
	}
	return nil
}

func checkDecEqual(d1, d2 common.Dec) error {
	if d1.IsNil() {
		d1 = common.ZeroDec()
	}
	if d2.IsNil() {
		d2 = common.ZeroDec()
	}
	if !d1.Equal(d2) {
		return fmt.Errorf("not equal: %v / %v", d1, d2)
	}
	return nil
}

func checkBigIntEqual(i1, i2 *big.Int) error {
	if i1 == nil {
		i1 = big.NewInt(0)
	}
	if i2 == nil {
		i2 = big.NewInt(0)
	}
	if i1.Cmp(i2) != 0 {
		return fmt.Errorf("not equal: %v / %v", i1, i2)
	}
	return nil
}
