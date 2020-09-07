package microstaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Map3NodeWrapperBuilder struct {
	wrapper Map3NodeWrapper_
}

func (b *Map3NodeWrapperBuilder) SetMap3Address(map3Address common.Address) *Map3NodeWrapperBuilder {
	b.wrapper.Map3Node.Map3Address = map3Address
	return b
}

func (b *Map3NodeWrapperBuilder) SetOperatorAddress(operator common.Address) *Map3NodeWrapperBuilder {
	b.wrapper.Map3Node.OperatorAddress = operator
	return b
}

func (b *Map3NodeWrapperBuilder) AddNodeKey(key BLSPublicKey_) *Map3NodeWrapperBuilder {
	b.wrapper.Map3Node.NodeKeys.Keys = append(b.wrapper.Map3Node.NodeKeys.Keys, &key)
	return b
}

func (b *Map3NodeWrapperBuilder) AddNodeKeys(keys BLSPublicKeys_) *Map3NodeWrapperBuilder {
	for _, k := range keys.Keys {
		b.AddNodeKey(*k)
	}
	return b
}

func (b *Map3NodeWrapperBuilder) SetCommission(commission Commission_) *Map3NodeWrapperBuilder {
	c, _ := commission.Copy()
	b.wrapper.Map3Node.Commission = *c
	return b
}

func (b *Map3NodeWrapperBuilder) SetDescription(description Description_) *Map3NodeWrapperBuilder {
	b.wrapper.Map3Node.Description = description
	return b
}

func (b *Map3NodeWrapperBuilder) SetCreationHeight(creationHeight *big.Int) *Map3NodeWrapperBuilder {
	b.wrapper.Map3Node.CreationHeight = big.NewInt(0).Set(creationHeight)
	return b
}

func (b *Map3NodeWrapperBuilder) SetAge(age common.Dec) *Map3NodeWrapperBuilder {
	b.wrapper.Map3Node.Age = age
	return b
}

func (b *Map3NodeWrapperBuilder) SetStatus(status Map3Status) *Map3NodeWrapperBuilder {
	b.wrapper.Map3Node.Status = uint8(status)
	return b
}

func (b *Map3NodeWrapperBuilder) SetActivationEpoch(activationEpoch *big.Int) *Map3NodeWrapperBuilder {
	b.wrapper.Map3Node.ActivationEpoch = big.NewInt(0).Set(activationEpoch)
	return b
}

func (b *Map3NodeWrapperBuilder) SetPendingEpoch(pendingEpoch *big.Int) *Map3NodeWrapperBuilder {
	b.wrapper.Map3Node.PendingEpoch = big.NewInt(0).Set(pendingEpoch)
	return b
}

func (b *Map3NodeWrapperBuilder) SetReleaseEpoch(releaseEpoch common.Dec) *Map3NodeWrapperBuilder {
	b.wrapper.Map3Node.ReleaseEpoch = releaseEpoch
	return b
}

func (b *Map3NodeWrapperBuilder) AddMicrodelegation(microdelegation Microdelegation_) *Map3NodeWrapperBuilder {
	b.wrapper.Microdelegations.Put(microdelegation.DelegatorAddress, microdelegation)
	if microdelegation.Amount.Cmp(common.Big0) > 0 {
		b.wrapper.TotalDelegation.Add(b.wrapper.TotalDelegation, microdelegation.Amount)
	}
	if microdelegation.PendingDelegation.Amount.Cmp(common.Big0) > 0 {
		b.wrapper.TotalPendingDelegation.Add(b.wrapper.TotalPendingDelegation, microdelegation.PendingDelegation.Amount)
	}
	return b
}

func (b *Map3NodeWrapperBuilder) SetAccumulatedReward(reward *big.Int) *Map3NodeWrapperBuilder {
	b.wrapper.AccumulatedReward = big.NewInt(0).Set(reward)
	return b
}

func (b *Map3NodeWrapperBuilder) SetRestakingReference(restakingReference RestakingReference_) *Map3NodeWrapperBuilder {
	b.wrapper.RestakingReference = restakingReference
	return b
}

func (b *Map3NodeWrapperBuilder) Build() *Map3NodeWrapper_ {
	return &b.wrapper
}

func NewMap3NodeWrapperBuilder() *Map3NodeWrapperBuilder {
	return &Map3NodeWrapperBuilder{
		wrapper: Map3NodeWrapper_{
			Map3Node: Map3Node_{
				Map3Address:     common.Address{},
				OperatorAddress: common.Address{},
				NodeKeys:        NewEmptyBLSKeys(),
				Commission: Commission_{
					Rate:              common.NewDec(0),
					RateForNextPeriod: common.NewDec(0),
					UpdateHeight:      big.NewInt(0),
				},
				Description:     Description_{},
				CreationHeight:  big.NewInt(0),
				Age:             common.NewDec(0),
				Status:          uint8(Pending),
				PendingEpoch:    big.NewInt(0),
				ActivationEpoch: big.NewInt(0),
				ReleaseEpoch:    common.NewDec(0),
			},
			Microdelegations: NewMicrodelegationMap(),
			RestakingReference: RestakingReference_{
				ValidatorAddress: common.Address{},
			},
			AccumulatedReward:      big.NewInt(0),
			TotalDelegation:        big.NewInt(0),
			TotalPendingDelegation: big.NewInt(0),
		},
	}
}
