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
			func(n *Map3Node_) {},
			nil,
		},
		{
			func(n *Map3Node_) { n.Description = invalidDescription },
			errors.New("exceed maximum name length"),
		},
		{
			func(n *Map3Node_) { n.NodeKeys.Keys = n.NodeKeys.Keys[:0] },
			errNeedAtLeastOneSlotKey,
		},
		{
			func(n *Map3Node_) {
				n.NodeKeys = NewEmptyBLSKeys()
				n.NodeKeys.Keys = append(n.NodeKeys.Keys, &blsPubSigPairs[0].pub, &blsPubSigPairs[1].pub)
			},
			ErrExcessiveBLSKeys,
		},
		{
			func(n *Map3Node_) { n.Commission.Rate = nilRate },
			errInvalidCommissionRate,
		},
		{
			func(n *Map3Node_) { n.Commission.Rate = negativeRate },
			errInvalidCommissionRate,
		},
		{
			func(n *Map3Node_) { n.Commission.Rate = invalidRate },
			errInvalidCommissionRate,
		},
		{
			func(n *Map3Node_) { n.Commission.RateForNextPeriod = nilRate },
			errInvalidCommissionRate,
		},
		{
			func(n *Map3Node_) { n.Commission.RateForNextPeriod = negativeRate },
			errInvalidCommissionRate,
		},
		{
			func(n *Map3Node_) { n.Commission.RateForNextPeriod = invalidRate },
			errInvalidCommissionRate,
		},
	}
	for i, test := range tests {
		n := makeMap3Node()
		test.editMap3Node(&n)
		err := n.SanityCheck(MaxPubKeyAllowed)
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
	n := NewMap3NodeWrapperBuilder().
		SetMap3Address(map3NodeAddr).
		SetOperatorAddress(operatorAddr).
		AddNodeKey(blsPubSigPairs[0].pub).
		SetCommission(c).
		SetDescription(validDescription).
		SetCreationHeight(big.NewInt(12306)).
		SetAge(common.NewDecWithPrec(25, 1)).
		SetStatus(Active).
		SetActivationEpoch(big.NewInt(1)).
		SetReleaseEpoch(common.NewDec(10)).Build()
	return n.Map3Node
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
			editCreateValidator: func(cn *CreateMap3Node) {},
			expErr:              nil,
		},
		{
			editCreateValidator: func(cn *CreateMap3Node) { cn.NodeKeySig = blsPubSigPairs[2].sig },
			expErr:              errBLSKeysNotMatchSigs,
		},
	}
	for i, test := range tests {
		cn := makeCreateMap3Node()
		test.editCreateValidator(&cn)

		blockNum, epoch := big.NewInt(10), big.NewInt(10)
		n, err := CreateMap3NodeFromNewMsg(&cn, map3NodeAddr, blockNum, epoch)
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if err := assertMap3NodeAlignCreateMap3Node(n.Map3Node, cn, blockNum, epoch); err != nil {
			t.Error(err)
		}
	}
}

func TestUpdateMap3NodeFromEditMsg(t *testing.T) {
	tests := []struct {
		editMap3Node    EditMap3Node
		editExpMap3Node func(*Map3Node_)
		expErr          error
	}{
		{
			editMap3Node:    EditMap3Node{Map3NodeAddress: map3NodeAddr},
			editExpMap3Node: func(*Map3Node_) {},
		},
		{
			// update Description.Name
			editMap3Node: EditMap3Node{
				Map3NodeAddress: map3NodeAddr,
				Description:     Description_{Name: "jacky@harmony.one"},
			},
			editExpMap3Node: func(n *Map3Node_) { n.Description.Name = "jacky@harmony.one" },
		},
		{
			// Remove a bls pub key
			editMap3Node: EditMap3Node{
				Map3NodeAddress: map3NodeAddr,
				NodeKeyToRemove: &blsPubSigPairs[0].pub,
			},
			editExpMap3Node: func(n *Map3Node_) { n.NodeKeys = NewEmptyBLSKeys() },
		},
		{
			// Add a bls pub key with signature
			editMap3Node: EditMap3Node{
				Map3NodeAddress: map3NodeAddr,
				NodeKeyToAdd:    &blsPubSigPairs[4].pub,
				NodeKeyToAddSig: &blsPubSigPairs[4].sig,
			},
			editExpMap3Node: func(n *Map3Node_) {
				n.NodeKeys.Keys = append(n.NodeKeys.Keys, &blsPubSigPairs[4].pub)
			},
		},
		{
			// EditMap3Node having signature without pub will not be a update
			editMap3Node: EditMap3Node{
				Map3NodeAddress: map3NodeAddr,
				NodeKeyToAddSig: &blsPubSigPairs[4].sig,
			},
			editExpMap3Node: func(n *Map3Node_) {},
		},
		{
			// invalid address
			editMap3Node: EditMap3Node{
				Map3NodeAddress: common.BigToAddress(common.Big1),
			},
			expErr: errAddressNotMatch,
		},
		{
			// invalid description
			editMap3Node: EditMap3Node{
				Map3NodeAddress: map3NodeAddr,
				Description:     invalidDescription,
			},
			expErr: errors.New("exceed maximum name length"),
		},
		{
			// invalid removing bls key
			editMap3Node: EditMap3Node{
				Map3NodeAddress: map3NodeAddr,
				NodeKeyToRemove: &blsPubSigPairs[4].pub,
			},
			expErr: errNodeKeyToRemoveNotFound,
		},
		{
			// add pub not having valid signature
			editMap3Node: EditMap3Node{
				Map3NodeAddress: map3NodeAddr,
				NodeKeyToAdd:    &blsPubSigPairs[4].pub,
				NodeKeyToAddSig: &blsPubSigPairs[3].sig,
			},
			expErr: errBLSKeysNotMatchSigs,
		},
		{
			// add pub key already exist in validator
			editMap3Node: EditMap3Node{
				Map3NodeAddress: map3NodeAddr,
				NodeKeyToAdd:    &blsPubSigPairs[0].pub,
				NodeKeyToAddSig: &blsPubSigPairs[0].sig,
			},
			expErr: errNodeKeyToAddExists,
		},
	}
	for i, test := range tests {
		node := makeMap3Node()

		err := UpdateMap3NodeFromEditMsg(&node, &test.editMap3Node)
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
		if (err != nil) || (test.expErr != nil) {
			continue
		}

		expNode := makeMap3Node()
		test.editExpMap3Node(&expNode)

		if err := CheckMap3NodeEqual(node, expNode); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

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

func assertMap3NodeAlignCreateMap3Node(n Map3Node_, cn CreateMap3Node, blockNum *big.Int, epoch *big.Int) error {
	if n.Map3Address != map3NodeAddr {
		return fmt.Errorf("map3 node address not equal")
	}
	if n.OperatorAddress != cn.OperatorAddress {
		return fmt.Errorf("operator address not equal")
	}
	if len(n.NodeKeys.Keys) != 1 {
		return fmt.Errorf("len(NodeKeys) not equal 1")
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
	if n.CreationHeight.Cmp(blockNum) != 0 {
		return fmt.Errorf("CreationHeight not equal")
	}
	if !n.Age.IsZero() {
		return fmt.Errorf("node age not zero")
	}
	if n.Status != uint8(Pending) {
		return fmt.Errorf("status not pending")
	}
	if n.PendingEpoch.Cmp(epoch) != 0 {
		return fmt.Errorf("pending epoch not equal")
	}
	if n.ActivationEpoch.Cmp(common.Big0) != 0 {
		return fmt.Errorf("activation epoch not equal to 0")
	}
	if !n.ReleaseEpoch.IsZero() {
		return fmt.Errorf("release epoch not equal to 0")
	}
	return nil
}
