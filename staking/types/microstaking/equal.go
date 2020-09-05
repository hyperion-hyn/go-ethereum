package microstaking

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// CheckMap3NodeWrapperEqual checks the equality of microstaking.Map3NodeWrapper. If not equal, an
// error is returned. Note nil pointer is treated as zero in this compare function.
func CheckMap3NodeWrapperEqual(w1, w2 Map3NodeWrapper_) error {
	if err := checkMap3NodeWrapperEqual(w1, w2); err != nil {
		return fmt.Errorf("Map3NodeWrapper%v", err)
	}
	return nil
}

// CheckMap3NodeEqual checks the equality of validator. If not equal, an
// error is returned. Note nil pointer is treated as zero in this compare function.
func CheckMap3NodeEqual(v1, v2 Map3Node_) error {
	if err := checkMap3NodeEqual(v1, v2); err != nil {
		return fmt.Errorf("Map3Node%v", err)
	}
	return nil
}

func checkMap3NodeWrapperEqual(w1, w2 Map3NodeWrapper_) error {
	if err := checkMap3NodeEqual(w1.Map3Node, w2.Map3Node); err != nil {
		return fmt.Errorf(".Map3Node%v", err)
	}
	if err := checkMicrodelegationMapEqual(w1.Microdelegations, w2.Microdelegations); err != nil {
		return fmt.Errorf(".Microdelegations%v", err)
	}
	if w1.RestakingReference.ValidatorAddress != w2.RestakingReference.ValidatorAddress {
		return fmt.Errorf(".RestakingReference.ValidatorAddress not equal: %x / %x",
			w1.RestakingReference.ValidatorAddress, w2.RestakingReference.ValidatorAddress)
	}
	if err := checkBigIntEqual(w1.AccumulatedReward, w2.AccumulatedReward); err != nil {
		return fmt.Errorf(".AccumulatedReward %v", err)
	}
	if err := checkBigIntEqual(w1.TotalDelegation, w2.TotalDelegation); err != nil {
		return fmt.Errorf(".TotalDelegation %v", err)
	}
	if err := checkBigIntEqual(w1.TotalPendingDelegation, w2.TotalPendingDelegation); err != nil {
		return fmt.Errorf(".TotalPendingDelegation %v", err)
	}
	return nil
}

func checkMap3NodeEqual(v1, v2 Map3Node_) error {
	if v1.Map3Address != v2.Map3Address {
		return fmt.Errorf(".Map3Node address not equal: %x / %x", v1.Map3Address, v2.Map3Address)
	}
	if v1.OperatorAddress != v2.OperatorAddress {
		return fmt.Errorf(".Operator address not equal: %x / %x", v1.OperatorAddress, v2.OperatorAddress)
	}
	if err := checkPubKeysEqual(v1.NodeKeys, v2.NodeKeys); err != nil {
		return fmt.Errorf(".NodeKeys%v", err)
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
	if err := checkDecEqual(v1.Age, v2.Age); err != nil {
		return fmt.Errorf(".Age %v", err)
	}
	if v1.Status != v2.Status {
		return fmt.Errorf(".Status not equal: %v / %v", v1.Status, v2.Status)
	}
	if err := checkBigIntEqual(v1.ActivationEpoch, v2.ActivationEpoch); err != nil {
		return fmt.Errorf(".ActivationEpoch %v", err)
	}
	if err := checkDecEqual(v1.ReleaseEpoch, v2.ReleaseEpoch); err != nil {
		return fmt.Errorf(".ReleaseEpoch %v", err)
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
	if err := checkDecEqual(c1.Rate, c2.Rate); err != nil {
		return fmt.Errorf(".Rate %v", err)
	}
	if err := checkDecEqual(c1.RateForNextPeriod, c2.RateForNextPeriod); err != nil {
		return fmt.Errorf(".RateForNextPeriod %v", err)
	}
	if err := checkBigIntEqual(c1.UpdateHeight, c2.UpdateHeight); err != nil {
		return fmt.Errorf(".UpdateHeight %v", err)
	}
	return nil
}

func checkMicrodelegationMapEqual(ds1, ds2 MicrodelegationMap_) error {
	if len(ds1.Keys) != len(ds2.Keys) {
		return fmt.Errorf(".len not equal: %v / %v", len(ds1.Keys), len(ds2.Keys))
	}
	for _, key := range ds1.Keys {
		r1, _ := ds1.Get(*key)
		r2, _ := ds2.Get(*key)
		if err := checkMicrodelegationEqual(r1, r2); err != nil {
			return fmt.Errorf("[%v]%v", key, err)
		}
	}
	return nil
}

func checkMicrodelegationEqual(d1, d2 Microdelegation_) error {
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
	if err := checkPendingDelegationEqual(d1.PendingDelegation, d2.PendingDelegation); err != nil {
		return fmt.Errorf(".PendingDelegation%v", err)
	}
	if err := checkUndelegationEqual(d1.Undelegation, d2.Undelegation); err != nil {
		return fmt.Errorf(".Undelegation%v", err)
	}
	if err := checkRenewalEqual(d1.Renewal, d2.Renewal); err != nil {
		return fmt.Errorf(".Renewal%v", err)
	}
	return nil
}

func checkPendingDelegationEqual(pd1, pd2 PendingDelegation_) error {
	if err := checkBigIntEqual(pd1.Amount, pd2.Amount); err != nil {
		return fmt.Errorf(".Amount %v", err)
	}
	if err := checkDecEqual(pd1.UnlockedEpoch, pd2.UnlockedEpoch); err != nil {
		return fmt.Errorf(".UnlockedEpoch %v", err)
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

func checkRenewalEqual(r1, r2 Renewal_) error {
	if r1.IsRenew != r2.IsRenew {
		return fmt.Errorf(".IsRenew not equal: %v / %v", r1.IsRenew, r2.IsRenew)
	}
	if err := checkBigIntEqual(r1.UpdateHeight, r2.UpdateHeight); err != nil {
		return fmt.Errorf(".UpdateHeight %v", err)
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
