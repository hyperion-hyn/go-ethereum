package microstaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

var (
	testPub = BLSPublicKey_{Key: [48]byte{1}}

	nonZeroCommission = Commission_{
		Rate:              common.NewDecWithPrec(1, 1),
		RateForNextPeriod: common.NewDecWithPrec(1, 1),
		UpdateHeight:      common.Big1,
	}

	nonZeroMicrodelegation = Microdelegation_{
		DelegatorAddress: common.BigToAddress(common.Big1),
		Amount:           common.Big1,
		Reward:           common.Big2,
		PendingDelegation: PendingDelegation_{
			Amount:        common.Big1,
			UnlockedEpoch: common.NewDec(2),
		},
		Undelegation: Undelegation_{
			Amount: common.Big1,
			Epoch:  common.Big2,
		},
		Renewal: Renewal_{
			IsRenew:      true,
			UpdateHeight: big.NewInt(1),
		},
	}

	restakingReference = RestakingReference_{
		ValidatorAddress: common.BigToAddress(common.Big1),
	}
)

func TestMap3NodeWrapperBuilder(t *testing.T) {
	tests := []struct {
		map3Addr           common.Address
		operatorAddr       common.Address
		key                BLSPublicKey_
		commission         Commission_
		description        Description_
		creationHeight     *big.Int
		age                common.Dec
		status             Map3Status
		pendingEpoch       *big.Int
		activationEpoch    *big.Int
		releaseEpoch       common.Dec
		microdelegation    Microdelegation_
		restakingReference RestakingReference_
		accumulatedReward  *big.Int
	}{
		{
			map3Addr:     common.BigToAddress(common.Big1),
			operatorAddr: common.BigToAddress(common.Big1),
			key:          testPub,
			commission:   nonZeroCommission,
			description: Description_{
				Name:            "A",
				Identity:        "B",
				Website:         "C",
				SecurityContact: "D",
				Details:         "E",
			},
			creationHeight:     big.NewInt(1),
			age:                common.NewDec(1),
			status:             Active,
			pendingEpoch:       big.NewInt(1),
			activationEpoch:    big.NewInt(1),
			releaseEpoch:       common.NewDec(1),
			microdelegation:    nonZeroMicrodelegation,
			accumulatedReward:  big.NewInt(1),
			restakingReference: restakingReference,
		},
	}

	for i, test := range tests {
		n := NewMap3NodeWrapperBuilder().
			SetMap3Address(test.map3Addr).
			SetOperatorAddress(test.operatorAddr).
			AddNodeKey(test.key).
			SetCommission(test.commission).
			SetDescription(test.description).
			SetCreationHeight(test.creationHeight).
			SetAge(test.age).
			SetStatus(test.status).
			SetPendingEpoch(test.pendingEpoch).
			SetActivationEpoch(test.activationEpoch).
			SetReleaseEpoch(test.releaseEpoch).
			AddMicrodelegation(test.microdelegation).
			SetAccumulatedReward(test.accumulatedReward).
			SetRestakingReference(test.restakingReference).
			Build()

		exp := GetDefaultMap3NodeWrapper()
		exp.Map3Node.Map3Address = test.map3Addr
		exp.Map3Node.OperatorAddress = test.operatorAddr
		exp.Map3Node.NodeKeys.Keys = append(exp.Map3Node.NodeKeys.Keys, &test.key)
		exp.Map3Node.Commission = test.commission
		exp.Map3Node.Description = test.description
		exp.Map3Node.CreationHeight = test.creationHeight
		exp.Map3Node.Age = test.age
		exp.Map3Node.Status = uint8(test.status)
		exp.Map3Node.PendingEpoch = test.pendingEpoch
		exp.Map3Node.ActivationEpoch = test.activationEpoch
		exp.Map3Node.ReleaseEpoch = test.releaseEpoch
		exp.Microdelegations.Put(test.microdelegation.DelegatorAddress, test.microdelegation)
		exp.AccumulatedReward = test.accumulatedReward
		exp.RestakingReference = restakingReference
		exp.TotalDelegation = common.Big1
		exp.TotalPendingDelegation = common.Big1

		if err := CheckMap3NodeWrapperEqual(*n, exp); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}
