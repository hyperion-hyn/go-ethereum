package microstaking

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	errInvalidCommissionRate = errors.New("commission rate should be a value ranging from 0.0 to 1.0")
)

var (
	hundredPercent = common.OneDec()
	zeroPercent    = common.ZeroDec()
)

func (c *Commission_) SanityCheck() error {
	if c.Rate.IsNil() || c.RateForNextPeriod.IsNil() {
		return errors.Wrap(errInvalidCommissionRate, "rate can not be nil")
	}

	if c.Rate.LT(zeroPercent) || c.Rate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "rate:%s", c.Rate.String(),
		)
	}
	if c.RateForNextPeriod.LT(zeroPercent) || c.RateForNextPeriod.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "rate for next period:%s", c.RateForNextPeriod.String(),
		)
	}
	return nil
}
