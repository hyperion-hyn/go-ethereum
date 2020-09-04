package restaking

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common2 "github.com/ethereum/go-ethereum/staking/types/common"
	"github.com/harmony-one/harmony/crypto/bls"
	"github.com/pkg/errors"
	"math/big"
	"reflect"
	"strings"
	"testing"
)

var (
	blsPubSigPairs   = makeBLSPubSigPairs(5)
	validatorAddr, _ = common.Bech32ToAddress("hyn1t2htvpfl862vnwdqnuekd9p4ulh3h6hdldamnd")
	operatorAddr, _  = common.Bech32ToAddress("hyn1pdv9lrdwl0rg5vglh4xtyrv3w123wsqket7zxy")
)

var (
	zeroDec = common.ZeroDec()
	oneDec  = common.OneDec()
	nineK   = new(big.Int).Mul(big.NewInt(9000), big.NewInt(1e18))
	tenK    = new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18))
	elevenK = new(big.Int).Mul(big.NewInt(11000), big.NewInt(1e18))
	twelveK = new(big.Int).Mul(big.NewInt(12000), big.NewInt(1e18))
	twentyK = new(big.Int).Mul(big.NewInt(20000), big.NewInt(1e18))

	nilRate      = common.Dec{}
	negativeRate = common.NewDec(-1)
	zeroRate     = common.ZeroDec()
	halfRate     = common.NewDecWithPrec(5, 1)
	oneRate      = common.NewDec(1)
	invalidRate  = common.NewDec(2)
)

var (
	validDescription = Description_{
		Name:            "Jacky Wang",
		Identity:        "jacky@harmony.one",
		Website:         "harmony.one/jacky",
		SecurityContact: "jacky@harmony.one",
		Details:         "Details of jacky",
	}

	invalidDescription = Description_{
		Name:            "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongname",
		Identity:        "jacky@harmony.one",
		Website:         "harmony.one/jacky",
		SecurityContact: "jacky@harmony.one",
		Details:         "Details of jacky",
	}

	validCommissionRates = CommissionRates_{
		Rate:          zeroRate,
		MaxRate:       zeroRate,
		MaxChangeRate: zeroRate,
	}
)

func TestValidator_SanityCheck(t *testing.T) {
	tests := []struct {
		editValidator func(*Validator_)
		expErr        error
	}{
		{
			func(v *Validator_) {},
			nil,
		},
		{
			func(v *Validator_) { v.Description = invalidDescription },
			errors.New("exceed maximum name length"),
		},
		{
			func(v *Validator_) { v.SlotPubKeys.Keys = v.SlotPubKeys.Keys[:0] },
			errNeedAtLeastOneSlotKey,
		},
		{
			func(v *Validator_) { v.MaxTotalDelegation = nil },
			errNilMaxTotalDelegation,
		},
		{
			func(v *Validator_) { v.Commission.CommissionRates.Rate = nilRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Validator_) { v.Commission.CommissionRates.Rate = negativeRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Validator_) { v.Commission.CommissionRates.Rate = invalidRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Validator_) { v.Commission.CommissionRates.MaxRate = nilRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Validator_) { v.Commission.CommissionRates.MaxRate = negativeRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Validator_) { v.Commission.CommissionRates.MaxRate = invalidRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Validator_) { v.Commission.CommissionRates.MaxChangeRate = nilRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Validator_) { v.Commission.CommissionRates.MaxChangeRate = negativeRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Validator_) { v.Commission.CommissionRates.MaxChangeRate = invalidRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Validator_) {
				v.Commission.CommissionRates.Rate, v.Commission.CommissionRates.MaxRate = oneRate, halfRate
			},
			errCommissionRateTooLarge,
		},
		{
			func(v *Validator_) {
				v.Commission.CommissionRates.MaxChangeRate, v.Commission.CommissionRates.MaxRate = oneRate, halfRate
			},
			errCommissionRateTooLarge,
		},
		{
			func(v *Validator_) {
				v.SlotPubKeys = NewEmptyBLSKeys()
				v.SlotPubKeys.Keys = append(v.SlotPubKeys.Keys, &blsPubSigPairs[0].pub, &blsPubSigPairs[1].pub)
			},
			ErrExcessiveBLSKeys,
		},
	}
	for i, test := range tests {
		v := makeValidValidator()
		test.editValidator(&v)
		err := v.SanityCheck(MaxPubKeyAllowed)
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
	}
}

type blsPubSigPair struct {
	pub BLSPublicKey_
	sig common2.BLSSignature
}

func makeBLSPubSigPairs(size int) []blsPubSigPair {
	pairs := make([]blsPubSigPair, 0, size)
	for i := 0; i != size; i++ {
		pairs = append(pairs, makeBLSPubSigPair())
	}
	return pairs
}

func makeBLSPubSigPair() blsPubSigPair {
	blsPriv := bls.RandPrivateKey()
	blsPub := blsPriv.GetPublicKey()
	msgHash := crypto.Keccak256([]byte(common2.BLSVerificationStr))
	blsSig := blsPriv.SignHash(msgHash)

	pub := BLSPublicKey_{}
	copy(pub.Key[:], blsPub.Serialize())

	var sig common2.BLSSignature
	copy(sig[:], blsSig.Serialize())

	return blsPubSigPair{pub, sig}
}

// makeValidValidator makes a valid Validator data structure
func makeValidValidator() Validator_ {
	cr := validCommissionRates
	c := Commission_{CommissionRates: cr, UpdateHeight: big.NewInt(300)}
	d := Description_{
		Name:     "Wayne",
		Identity: "wen",
		Website:  "harmony.one.wen",
		Details:  "best",
	}
	v := Validator_{
		ValidatorAddress: validatorAddr,
		OperatorAddresses: AddressSet_{
			Keys: []*Address{&operatorAddr},
			Set:  map[Address]*Bool{operatorAddr: func() *bool { t := true; return &t }()},
		},
		SlotPubKeys:          NewBLSKeysWithBLSKey(blsPubSigPairs[0].pub),
		LastEpochInCommittee: big.NewInt(20),
		MaxTotalDelegation:   twelveK,
		Status:               uint8(Active),
		Commission:           c,
		Description:          d,
		CreationHeight:       big.NewInt(12306),
	}
	return v
}

func assertError(gotErr, expErr error) error {
	if (gotErr == nil) != (expErr == nil) {
		return fmt.Errorf("error unexpected [%v] / [%v]", gotErr, expErr)
	}
	if gotErr == nil {
		return nil
	}
	if !strings.Contains(gotErr.Error(), expErr.Error()) {
		return fmt.Errorf("error unexpected [%v] / [%v]", gotErr, expErr)
	}
	return nil
}

func TestCreateValidatorFromNewMsg(t *testing.T) {
	tests := []struct {
		editCreateValidator func(*CreateValidator)
		expErr              error
	}{
		{
			editCreateValidator: func(cv *CreateValidator) {},
			expErr:              nil,
		},
		{
			editCreateValidator: func(cv *CreateValidator) { cv.SlotKeySig = blsPubSigPairs[2].sig },
			expErr:              errBLSKeysNotMatchSigs,
		},
	}
	for i, test := range tests {
		cv := makeCreateValidator()
		test.editCreateValidator(&cv)

		v, err := CreateValidatorFromNewMsg(&cv, validatorAddr, common.Big1)
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if err := assertValidatorAlignCreateValidator(*v, cv); err != nil {
			t.Error(err)
		}
	}
}

func TestUpdateValidatorFromEditMsg(t *testing.T) {
	tests := []struct {
		editValidator    EditValidator
		editExpValidator func(*Validator_)
		expErr           error
	}{
		{
			editValidator:    EditValidator{ValidatorAddress: validatorAddr},
			editExpValidator: func(*Validator_) {},
		},
		{
			// update Description.Name
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				Description:      Description_{Name: "jacky@harmony.one"},
			},
			editExpValidator: func(v *Validator_) { v.Description.Name = "jacky@harmony.one" },
		},
		{
			// Update CommissionRate
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				CommissionRate:   &halfRate,
			},
			editExpValidator: func(v *Validator_) { v.Commission.CommissionRates.Rate = halfRate },
		},
		{
			// Update MaxTotalDelegation
			editValidator: EditValidator{
				ValidatorAddress:   validatorAddr,
				MaxTotalDelegation: elevenK,
			},
			editExpValidator: func(v *Validator_) { v.MaxTotalDelegation = elevenK },
		},
		{
			// Update MaxTotalDelegation to zero remain unchanged
			editValidator: EditValidator{
				ValidatorAddress:   validatorAddr,
				MaxTotalDelegation: common.Big0,
			},
			editExpValidator: func(v *Validator_) {},
		},
		{
			// Remove a bls pub key
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				SlotKeyToRemove:  &blsPubSigPairs[0].pub,
			},
			editExpValidator: func(v *Validator_) { v.SlotPubKeys.Keys = nil },
		},
		{
			// Add a bls pub key with signature
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				SlotKeyToAdd:     &blsPubSigPairs[4].pub,
				SlotKeyToAddSig:  &blsPubSigPairs[4].sig,
			},
			editExpValidator: func(v *Validator_) {
				v.SlotPubKeys.Keys = append(v.SlotPubKeys.Keys, &blsPubSigPairs[4].pub)
			},
		},
		{
			// EditValidator having signature without pub will not be a update
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				SlotKeyToAddSig:  &blsPubSigPairs[4].sig,
			},
			editExpValidator: func(v *Validator_) {},
		},
		{
			// update status
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				EPOSStatus:       Inactive,
			},
			editExpValidator: func(v *Validator_) { v.Status = Uint8(Inactive) },
		},
		{
			// status to banned - not changed
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				EPOSStatus:       Banned,
			},
			editExpValidator: func(v *Validator_) {},
		},
		{
			// invalid address
			editValidator: EditValidator{
				ValidatorAddress: common.BigToAddress(common.Big1),
			},
			expErr: errAddressNotMatch,
		},
		{
			// invalid description
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				Description:      invalidDescription,
			},
			expErr: errors.New("exceed maximum name length"),
		},
		{
			// invalid removing bls key
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				SlotKeyToRemove:  &blsPubSigPairs[4].pub,
			},
			expErr: errSlotKeyToRemoveNotFound,
		},
		{
			// add pub not having valid signature
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				SlotKeyToAdd:     &blsPubSigPairs[4].pub,
				SlotKeyToAddSig:  &blsPubSigPairs[3].sig,
			},
			expErr: errBLSKeysNotMatchSigs,
		},
		{
			// add pub key already exist in validator
			editValidator: EditValidator{
				ValidatorAddress: validatorAddr,
				SlotKeyToAdd:     &blsPubSigPairs[0].pub,
				SlotKeyToAddSig:  &blsPubSigPairs[0].sig,
			},
			expErr: errSlotKeyToAddExists,
		},
	}
	for i, test := range tests {
		val := makeValidValidator()

		err := UpdateValidatorFromEditMsg(&val, &test.editValidator)
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
		if (err != nil) || (test.expErr != nil) {
			continue
		}

		expVal := makeValidValidator()
		test.editExpValidator(&expVal)

		if err := assertValidatorEqual(val, expVal); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

func getPubsFromPairs(pairs []blsPubSigPair, indexes []int) BLSPublicKeys_ {
	pubs := NewEmptyBLSKeys()
	for _, index := range indexes {
		pubs.Keys = append(pubs.Keys, &pairs[index].pub)
	}
	return pubs
}

func getSigsFromPairs(pairs []blsPubSigPair, indexes []int) []common2.BLSSignature {
	sigs := make([]common2.BLSSignature, 0, len(indexes))
	for _, index := range indexes {
		sigs = append(sigs, pairs[index].sig)
	}
	return sigs
}

// makeCreateValidator makes a structure of CreateValidator
func makeCreateValidator() CreateValidator {
	addr := operatorAddr
	desc := validDescription
	cr := validCommissionRates
	return CreateValidator{
		OperatorAddress:    addr,
		Description:        desc,
		CommissionRates:    cr,
		MaxTotalDelegation: twelveK,
		SlotPubKey:         blsPubSigPairs[0].pub,
		SlotKeySig:         blsPubSigPairs[0].sig,
	}
}

func assertValidatorEqual(v1, v2 Validator_) error {
	if v1.ValidatorAddress != v2.ValidatorAddress {
		return fmt.Errorf("validator address not equal: %v / %v", v1.ValidatorAddress, v2.ValidatorAddress)
	}
	if !reflect.DeepEqual(v1.OperatorAddresses, v2.OperatorAddresses) {
		return fmt.Errorf("operator addresses not equal")
	}
	if len(v1.SlotPubKeys.Keys) != len(v2.SlotPubKeys.Keys) {
		return fmt.Errorf("len(SlotPubKeys) not equal: %v / %v", len(v1.SlotPubKeys.Keys), len(v2.SlotPubKeys.Keys))
	}
	for i := range v1.SlotPubKeys.Keys {
		pk1, pk2 := v1.SlotPubKeys.Keys[i], v2.SlotPubKeys.Keys[i]
		if *pk1 != *pk2 {
			return fmt.Errorf("SlotPubKeys[%v] not equal: %s / %s", i, pk1.Hex(), pk2.Hex())
		}
	}
	if v1.LastEpochInCommittee.Cmp(v2.LastEpochInCommittee) != 0 {
		return fmt.Errorf("LastEpochInCommittee not equal: %v / %v", v1.LastEpochInCommittee, v2.LastEpochInCommittee)
	}
	if v1.MaxTotalDelegation.Cmp(v2.MaxTotalDelegation) != 0 {
		return fmt.Errorf("MaxTotalDelegation not equal: %v / %v", v1.MaxTotalDelegation, v2.MaxTotalDelegation)
	}
	if v1.Status != v2.Status {
		return fmt.Errorf("status not equal: %v / %v", v1.Status, v2.Status)
	}
	if err := assertCommissionRatesEqual(v1.Commission.CommissionRates, v2.Commission.CommissionRates); err != nil {
		return fmt.Errorf("validator.Commission: %v", err)
	}
	if err := assertDescriptionEqual(v1.Description, v2.Description); err != nil {
		return fmt.Errorf("validator.Description: %v", err)
	}
	if v1.CreationHeight.Cmp(v2.CreationHeight) != 0 {
		return fmt.Errorf("CreationHeight not equal: %v / %v", v1.CreationHeight, v2.CreationHeight)
	}
	return nil
}

func assertValidatorAlignCreateValidator(v Validator_, cv CreateValidator) error {
	if v.ValidatorAddress != validatorAddr {
		return fmt.Errorf("validator address not equal")
	}
	if _, ok := v.OperatorAddresses.Set[cv.OperatorAddress]; !ok {
		return fmt.Errorf("operator address not equal")
	}
	if len(v.SlotPubKeys.Keys) != 1 {
		return fmt.Errorf("len(SlotPubKeys) not equal 1")
	}
	if !reflect.DeepEqual(*v.SlotPubKeys.Keys[0], cv.SlotPubKey) {
		return fmt.Errorf("SlotPubKey not equal")
	}
	if v.LastEpochInCommittee.Cmp(new(big.Int)) != 0 {
		return fmt.Errorf("LastEpochInCommittee not zero")
	}
	if v.MaxTotalDelegation.Cmp(cv.MaxTotalDelegation) != 0 {
		return fmt.Errorf("MaxTotalDelegation not equal")
	}
	if v.Status != uint8(Active) {
		return fmt.Errorf("status not active")
	}
	if err := assertCommissionRatesEqual(v.Commission.CommissionRates, cv.CommissionRates); err != nil {
		return fmt.Errorf("commissionRate not expected: %v", err)
	}
	if v.Commission.UpdateHeight.Cmp(v.CreationHeight) != 0 {
		return fmt.Errorf("validator's update height not equal to creation height")
	}
	if err := assertDescriptionEqual(v.Description, cv.Description); err != nil {
		return fmt.Errorf("description not expected: %v", err)
	}
	return nil
}

func assertCommissionRatesEqual(c1, c2 CommissionRates_) error {
	if !c1.Rate.Equal(c2.Rate) {
		return fmt.Errorf("rate not equal: %v / %v", c1.Rate, c2.Rate)
	}
	if !c1.MaxRate.Equal(c2.MaxRate) {
		return fmt.Errorf("max rate not equal: %v / %v", c1.MaxRate, c2.MaxRate)
	}
	if !c1.MaxChangeRate.Equal(c2.MaxChangeRate) {
		return fmt.Errorf("max change rate not equal: %v / %v", c1.MaxChangeRate, c2.MaxChangeRate)
	}
	return nil
}
