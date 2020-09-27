package microstaking

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	errCommissionRateNil = errors.New("commission rate cannot be nil")
)

var (
	tenPercent    = common.NewDecWithPrec(1, 1) // 10%
	twentyPercent = common.NewDecWithPrec(2, 1) // 20%
)

func (c *Commission_) SanityCheck(minCommissionRate, maxCommissionRate common.Dec) error {
	if c.Rate.IsNil() || c.RateForNextPeriod.IsNil() {
		return errCommissionRateNil
	}

	if c.Rate.LT(minCommissionRate) || c.Rate.GT(maxCommissionRate) {
		return errors.Errorf("commission rate should be a value ranging from %v to %v, rate:%s", minCommissionRate, maxCommissionRate, c.Rate.String())
	}
	if c.RateForNextPeriod.LT(minCommissionRate) || c.RateForNextPeriod.GT(maxCommissionRate) {
		return errors.Errorf("commission rate should be a value ranging from %v to %v, rate:%s", minCommissionRate, maxCommissionRate, c.RateForNextPeriod.String())
	}
	return nil
}
