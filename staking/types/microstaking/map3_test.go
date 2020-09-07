package microstaking

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/bls"
	common2 "github.com/ethereum/go-ethereum/staking/types/common"
	"github.com/pkg/errors"
	"math/big"
	"reflect"
	"strings"
	"testing"
)

var (
	blsPubSigPairs  = makeBLSPubSigPairs(5)
	map3NodeAddr, _ = common.Bech32ToAddress("hyn1t2htvpfl862vnwdqnuekd9p4ulh3h6hdldamnd")
	operatorAddr, _ = common.Bech32ToAddress("hyn1pdv9lrdwl0rg5vglh4xtyrv3w123wsqket7zxy")

	oneMill = new(big.Int).Mul(big.NewInt(2000000), big.NewInt(1e18))

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

	validCommissionRates = Commission_{
		Rate:              halfRate,
		RateForNextPeriod: halfRate,
		UpdateHeight:      big.NewInt(10),
	}
)

func TestMap3Node_SanityCheck(t *testing.T) {
	tests := []struct {
		editMap3Node func(*Map3Node_)
		expErr       error
	}{
		{
			func(v *Map3Node_) {},
			nil,
		},
		{
			func(v *Map3Node_) { v.Description = invalidDescription },
			errors.New("exceed maximum name length"),
		},
		{
			func(v *Map3Node_) { v.NodeKeys.Keys = v.NodeKeys.Keys[:0] },
			errNeedAtLeastOneSlotKey,
		},
		{
			func(v *Map3Node_) {
				v.NodeKeys = NewEmptyBLSKeys()
				v.NodeKeys.Keys = append(v.NodeKeys.Keys, &blsPubSigPairs[0].pub, &blsPubSigPairs[1].pub)
			},
			ErrExcessiveBLSKeys,
		},
		{
			func(v *Map3Node_) { v.Commission.Rate = nilRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Map3Node_) { v.Commission.Rate = negativeRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Map3Node_) { v.Commission.Rate = invalidRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Map3Node_) { v.Commission.RateForNextPeriod = nilRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Map3Node_) { v.Commission.RateForNextPeriod = negativeRate },
			errInvalidCommissionRate,
		},
		{
			func(v *Map3Node_) { v.Commission.RateForNextPeriod = invalidRate },
			errInvalidCommissionRate,
		},
	}
	for i, test := range tests {
		v := makeMap3Node()
		test.editMap3Node(&v)
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

func makeMap3Node() Map3Node_ {
	c := validCommissionRates
	d := Description_{
		Name:     "Wayne",
		Identity: "wen",
		Website:  "harmony.one.wen",
		Details:  "best",
	}
	v := Map3Node_{
		Map3Address:     map3NodeAddr,
		OperatorAddress: operatorAddr,
		NodeKeys:        NewBLSKeysWithBLSKey(blsPubSigPairs[0].pub),
		Commission:      c,
		Description:     d,
		CreationHeight:  big.NewInt(12306),
		Age:             common.NewDecWithPrec(25, 1),
		Status:          uint8(Active),
		ActivationEpoch: big.NewInt(1),
		ReleaseEpoch:    common.NewDec(10),
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

func TestCreateMap3NodeFromNewMsg(t *testing.T) {
	tests := []struct {
		editCreateValidator func(*CreateMap3Node)
		expErr              error
	}{
		{
			editCreateValidator: func(cv *CreateMap3Node) {},
			expErr:              nil,
		},
		{
			editCreateValidator: func(cv *CreateMap3Node) { cv.NodeKeySig = blsPubSigPairs[2].sig },
			expErr:              errBLSKeysNotMatchSigs,
		},
	}
	for i, test := range tests {
		cn := makeCreateMap3Node()
		test.editCreateValidator(&cn)

		n, err := CreateMap3NodeFromNewMsg(&cn, map3NodeAddr, big.NewInt(10), big.NewInt(10))
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if err := assertMap3NodeAlignCreateMap3Node(n.Map3Node, cn); err != nil {
			t.Error(err)
		}
	}
}


//func TestUpdateMap3NodeFromEditMsg(t *testing.T) {
//	tests := []struct {
//		editMap3Node    EditMap3Node
//		editExpMap3Node func(*Map3Node_)
//		expErr          error
//	}{
//		{
//			editMap3Node:    EditMap3Node{Map3NodeAddress: map3NodeAddr},
//			editExpMap3Node: func(*Map3Node_) {},
//		},
//		{
//			// update Description.Name
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//				Description:      Description_{Name: "jacky@harmony.one"},
//			},
//			editExpMap3Node: func(v *Map3Node_) { v.Description.Name = "jacky@harmony.one" },
//		},
//		{
//			// Update CommissionRate
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//			},
//			editExpMap3Node: func(v *Map3Node_) { v.Commission.CommissionRates.Rate = halfRate },
//		},
//		{
//			// Update MaxTotalDelegation
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress:   map3NodeAddr,
//				MaxTotalDelegation: elevenK,
//			},
//			editExpMap3Node: func(v *Map3Node_) { v.MaxTotalDelegation = elevenK },
//		},
//		{
//			// Update MaxTotalDelegation to zero remain unchanged
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress:   map3NodeAddr,
//				MaxTotalDelegation: common.Big0,
//			},
//			editExpMap3Node: func(v *Map3Node_) {},
//		},
//		{
//			// Remove a bls pub key
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//				SlotKeyToRemove:  &blsPubSigPairs[0].pub,
//			},
//			editExpMap3Node: func(v *Map3Node_) { v.SlotPubKeys.Keys = nil },
//		},
//		{
//			// Add a bls pub key with signature
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//				SlotKeyToAdd:     &blsPubSigPairs[4].pub,
//				SlotKeyToAddSig:  &blsPubSigPairs[4].sig,
//			},
//			editExpMap3Node: func(v *Map3Node_) {
//				v.SlotPubKeys.Keys = append(v.SlotPubKeys.Keys, &blsPubSigPairs[4].pub)
//			},
//		},
//		{
//			// EditMap3Node having signature without pub will not be a update
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//				SlotKeyToAddSig:  &blsPubSigPairs[4].sig,
//			},
//			editExpMap3Node: func(v *Map3Node_) {},
//		},
//		{
//			// update status
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//				EPOSStatus:       Inactive,
//			},
//			editExpMap3Node: func(v *Map3Node_) { v.Status = Uint8(Inactive) },
//		},
//		{
//			// status to banned - not changed
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//				EPOSStatus:       Banned,
//			},
//			editExpMap3Node: func(v *Map3Node_) {},
//		},
//		{
//			// invalid address
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: common.BigToAddress(common.Big1),
//			},
//			expErr: errAddressNotMatch,
//		},
//		{
//			// invalid description
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//				Description:      invalidDescription,
//			},
//			expErr: errors.New("exceed maximum name length"),
//		},
//		{
//			// invalid removing bls key
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//				SlotKeyToRemove:  &blsPubSigPairs[4].pub,
//			},
//			expErr: errSlotKeyToRemoveNotFound,
//		},
//		{
//			// add pub not having valid signature
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//				SlotKeyToAdd:     &blsPubSigPairs[4].pub,
//				SlotKeyToAddSig:  &blsPubSigPairs[3].sig,
//			},
//			expErr: errBLSKeysNotMatchSigs,
//		},
//		{
//			// add pub key already exist in validator
//			editMap3Node: EditMap3Node{
//				Map3NodeAddress: map3NodeAddr,
//				SlotKeyToAdd:     &blsPubSigPairs[0].pub,
//				SlotKeyToAddSig:  &blsPubSigPairs[0].sig,
//			},
//			expErr: errSlotKeyToAddExists,
//		},
//	}
//	for i, test := range tests {
//		node := makeMap3Node()
//
//		err := UpdateMap3NodeFromEditMsg(&node, &test.editMap3Node)
//		if assErr := assertError(err, test.expErr); assErr != nil {
//			t.Errorf("Test %v: %v", i, assErr)
//		}
//		if (err != nil) || (test.expErr != nil) {
//			continue
//		}
//
//		expNode := makeMap3Node()
//		test.editExpMap3Node(&expNode)
//
//		if err := assertMap3NodeEqual(node, expNode); err != nil {
//			t.Errorf("Test %v: %v", i, err)
//		}
//	}
//}

// makeCreateMap3Node makes a structure of CreateValidator
func makeCreateMap3Node() CreateMap3Node {
	addr := operatorAddr
	desc := validDescription
	return CreateMap3Node{
		OperatorAddress: addr,
		Description:     desc,
		Commission:      halfRate,
		NodePubKey:      blsPubSigPairs[0].pub,
		NodeKeySig:      blsPubSigPairs[0].sig,
		Amount:          oneMill,
	}
}

func assertMap3NodeAlignCreateMap3Node(n Map3Node_, cn CreateMap3Node) error {
	if n.Map3Address != map3NodeAddr {
		return fmt.Errorf("map3 node address not equal")
	}
	if n.OperatorAddress != cn.OperatorAddress {
		return fmt.Errorf("operator address not equal")
	}
	if len(n.NodeKeys.Keys) != 1 {
		return fmt.Errorf("len(SlotPubKeys) not equal 1")
	}
	if !reflect.DeepEqual(*n.NodeKeys.Keys[0], cn.NodePubKey) {
		return fmt.Errorf("NodeKey not equal")
	}
	if !n.Commission.Rate.Equal(cn.Commission) {
		return fmt.Errorf("commissionRate not equal")
	}
	if !n.Commission.RateForNextPeriod.Equal(cn.Commission) {
		return fmt.Errorf("commissionRateForNextPeriod not equal")
	}
	if n.Commission.UpdateHeight.Cmp(n.CreationHeight) != 0 {
		return fmt.Errorf("validator's update height not equal to creation height")
	}
	if err := assertDescriptionEqual(n.Description, cn.Description); err != nil {
		return fmt.Errorf("description not expected: %v", err)
	}
	if n.CreationHeight.Cmp(big.NewInt(10)) != 0 {
		return fmt.Errorf("CreationHeight not equal")
	}
	if !n.Age.IsNil() {
		return fmt.Errorf("CreationHeight not nil")
	}
	if n.Status != uint8(Pending) {
		return fmt.Errorf("status not pending")
	}
	if n.ActivationEpoch != nil {
		return fmt.Errorf("ActivationEpoch not nil")
	}
	if !n.ReleaseEpoch.IsNil() {
		return fmt.Errorf("ReleaseEpoch not nil")
	}
	return nil
}