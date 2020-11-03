package microstaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	LockDurationInEpoch = 180
)

type LockDurationCalculator interface {
	Calculate(epoch, blockNum *big.Int) (activationEpoch *big.Int, releaseEpoch common.Dec)
}

type CalculatorForActivationAtEndOfEpoch struct {
}

func (c CalculatorForActivationAtEndOfEpoch) Calculate(epoch, blockNum *big.Int) (activationEpoch *big.Int, releaseEpoch common.Dec) {
	activationEpoch = new(big.Int).Set(epoch)
	releaseEpoch = common.NewDecFromBigInt(activationEpoch).Add(common.NewDec(LockDurationInEpoch))
	return activationEpoch, releaseEpoch
}

type Config interface {
	IsMicrostakingImprove(num *big.Int) bool
	IsLastBlock(blockNum uint64) bool
}

type CalculatorForActivationAtOnce struct {
	config Config
}

func (c CalculatorForActivationAtOnce) Calculate(epoch, blockNum *big.Int) (activationEpoch *big.Int, releaseEpoch common.Dec) {
	if c.config.IsLastBlock(blockNum.Uint64()) {
		activationEpoch = new(big.Int).Add(epoch, common.Big1)
		releaseEpoch = common.NewDecFromBigInt(activationEpoch).Add(common.NewDec(LockDurationInEpoch - 1))
	} else {
		activationEpoch = new(big.Int).Set(epoch)
		releaseEpoch = common.NewDecFromBigInt(activationEpoch).Add(common.NewDec(LockDurationInEpoch - 1))
	}
	return activationEpoch, releaseEpoch
}

type LockDurationCalculatorFactory struct {
}

func NewLockDurationCalculator(config Config, blockNum *big.Int) LockDurationCalculator {
	if config.IsMicrostakingImprove(blockNum) {
		return CalculatorForActivationAtOnce{config: config}
	}
	return CalculatorForActivationAtEndOfEpoch{}
}