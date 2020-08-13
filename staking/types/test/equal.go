package staketest

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

// CheckValidatorEqual checks the equality of validator. If not equal, an
// error is returned. Note nil pointer is treated as zero in this compare function.
func CheckValidatorEqual(v1, v2 *restaking.Validator_) error {
	if err := checkValidatorEqual(v1, v2); err != nil {
		return fmt.Errorf("validator%v", err)
	}
	return nil
}

func checkValidatorEqual(v1, v2 *restaking.Validator_) error {
	if v1.ValidatorAddress != v2.ValidatorAddress {
		return fmt.Errorf(".Validator address not equal: %x / %x", v1.ValidatorAddress, v2.ValidatorAddress)
	}
	if err := checkAddressesEqual(v1.OperatorAddresses.Keys, v2.OperatorAddresses.Keys); err != nil {
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

func checkAddressesEqual(a1, a2 []*common.Address) error {
	if len(a1) != len(a2) {
		return fmt.Errorf(".len not equal: %v / %v", len(a1), len(a2))
	}
	for i := range a1 {
		if *a1[i] != *a2[i] {
			return fmt.Errorf("[%v] not equal: %x / %x", i, a1[i], a2[i])
		}
	}
	return nil
}

func checkPubKeysEqual(pubs1, pubs2 restaking.BLSPublicKeys_) error {
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

func checkDescriptionEqual(d1, d2 restaking.Description_) error {
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

func checkCommissionEqual(c1, c2 restaking.Commission_) error {
	if err := checkCommissionRateEqual(c1.CommissionRates, c2.CommissionRates); err != nil {
		return fmt.Errorf(".CommissionRate%v", err)
	}
	if err := checkBigIntEqual(c1.UpdateHeight, c2.UpdateHeight); err != nil {
		return fmt.Errorf(".UpdateHeight %v", err)
	}
	return nil
}

func checkCommissionRateEqual(cr1, cr2 restaking.CommissionRates_) error {
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
