package restaking

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/staking/effective"
	"github.com/pkg/errors"
	"math/big"
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
		editValidator_ func(*Validator_)
		expErr         error
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
			func(v *Validator_) { v.Commission.CommissionRates.Rate = negativeRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Validator_) { v.Commission.CommissionRates.Rate = invalidRate },
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
				v.SlotPubKeys = BLSPublicKeys_{Keys: []*BLSPublicKey_{
					&blsPubSigPairs[0].pub, &blsPubSigPairs[0].pub,
				}}
			},
			errDuplicateSlotKeys,
		},
	}
	for i, test := range tests {
		v := makeValidValidator()
		test.editValidator_(&v)
		err := v.SanityCheck(MaxPubKeyAllowed)
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
	}
}

func TestUpdateDescription(t *testing.T) {
	tests := []struct {
		raw    Description_
		update Description_
		expect Description_
		expErr error
	}{
		{
			raw: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "wenDetails",
			},
			update: Description_{
				Name:            "Jacky",
				Identity:        "jw",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of Jacky",
			},
			expect: Description_{
				Name:            "Jacky",
				Identity:        "jw",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of Jacky",
			},
		},
		{
			raw: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "wenDetails",
			},
			update: Description_{},
			expect: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "wenDetails",
			},
		},
		{
			raw: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "wenDetails",
			},
			update: Description_{
				Details: "new details",
			},
			expect: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "new details",
			},
		},
		{
			raw: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "wenDetails",
			},
			update: Description_{
				Website: "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongwebsite",
			},
			expErr: errors.New("exceed Maximum Length website"),
		},
	}
	for i, test := range tests {
		err := test.raw.UpdateFrom(&test.update)
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if err := assertDescriptionEqual(test.raw, test.expect); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

func TestDescription_EnsureLength(t *testing.T) {
	tests := []struct {
		desc   Description_
		expErr error
	}{
		{
			desc: Description_{
				Name:            "Jacky Wang",
				Identity:        "jacky@harmony.one",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of jacky",
			},
			expErr: nil,
		},
		{
			desc:   Description_{},
			expErr: nil,
		},
		{
			desc: Description_{
				Name:            "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongname",
				Identity:        "jacky@harmony.one",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of jacky",
			},
			expErr: errors.New("exceed maximum name length"),
		},
		{
			desc: Description_{
				Name:            "Jacky Wang",
				Identity:        "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongidentity",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of jacky",
			},
			expErr: errors.New("exceed Maximum Length identity"),
		},
		{
			desc: Description_{
				Name:            "Jacky Wang",
				Identity:        "jacky@harmony.one",
				Website:         "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongwebsite",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of jacky",
			},
			expErr: errors.New("exceed Maximum Length website"),
		},
		{
			desc: Description_{
				Name:            "Jacky Wang",
				Identity:        "jacky@harmony.one",
				Website:         "harmony.one/jacky",
				SecurityContact: "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongcontact",
				Details:         "Details of jacky",
			},
			expErr: errors.New("exceed Maximum Length"),
		},
		{
			desc: Description_{
				Name:            "Jacky Wang",
				Identity:        "jacky@harmony.one",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongdetail",
			},
			expErr: errors.New("exceed Maximum Length for details"),
		},
	}
	for i, test := range tests {
		err := test.desc.EnsureLength()
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
	}
}

type blsPubSigPair struct {
	pub BLSPublicKey_
	sig BLSSignature
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
	msgHash := crypto.Keccak256([]byte(BLSVerificationStr))
	blsSig := blsPriv.SignHash(msgHash)

	pub := BLSPublicKey_{}
	copy(pub.Key[:], blsPub.Serialize())

	var sig BLSSignature
	copy(sig[:], blsSig.Serialize())

	return blsPubSigPair{pub, sig}
}

// makeValidValidator makes a valid Validator data structure
func makeValidValidator() Validator_ {
	cr := validCommissionRates
	c := Commission_{cr, big.NewInt(300)}
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
		SlotPubKeys: BLSPublicKeys_{Keys: []*BLSPublicKey_{&blsPubSigPairs[0].pub}},
		LastEpochInCommittee: big.NewInt(20),
		MaxTotalDelegation:   twelveK,
		Status:               big.NewInt(int64(effective.Active)),
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

// compare two descriptions' items
func assertDescriptionEqual(d1, d2 Description_) error {
	if d1.Name != d2.Name {
		return fmt.Errorf("name not equal: [%v] / [%v]", d1.Name, d2.Name)
	}
	if d1.Identity != d2.Identity {
		return fmt.Errorf("identity not equal: [%v] / [%v]", d1.Identity, d2.Identity)
	}
	if d1.Website != d2.Website {
		return fmt.Errorf("website not equal: [%v] / [%v]", d1.Website, d2.Website)
	}
	if d1.SecurityContact != d2.SecurityContact {
		return fmt.Errorf("security contact not equal: [%v] / [%v]", d1.SecurityContact, d2.SecurityContact)
	}
	if d1.Details != d2.Details {
		return fmt.Errorf("details not equal: [%v] / [%v]", d1.Details, d2.Details)
	}
	return nil
}
