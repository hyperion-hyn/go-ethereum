package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	hundredPercent = common.OneDec()
	zeroPercent    = common.ZeroDec()
)

var (
	errCommissionRateTooLarge = errors.New("commission rate and change rate can not be larger than max commission rate")
	errInvalidCommissionRate  = errors.New("commission rate, change rate and max rate should be a value ranging from 0.0 to 1.0")
)

func (c *Commission_) SanityCheck() error {
	if c.CommissionRates.Rate.IsNil() || c.CommissionRates.MaxChangeRate.IsNil() || c.CommissionRates.MaxRate.IsNil() {
		return errors.Wrap(errInvalidCommissionRate, "rate can not be nil")
	}

	if c.CommissionRates.Rate.LT(zeroPercent) || c.CommissionRates.Rate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "rate:%s", c.CommissionRates.Rate.String(),
		)
	}

	if c.CommissionRates.MaxRate.LT(zeroPercent) || c.CommissionRates.MaxRate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "max rate:%s", c.CommissionRates.MaxRate.String(),
		)
	}

	if c.CommissionRates.MaxChangeRate.LT(zeroPercent) ||
		c.CommissionRates.MaxChangeRate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "max change rate:%s", c.CommissionRates.MaxChangeRate.String(),
		)
	}

	if c.CommissionRates.Rate.GT(c.CommissionRates.MaxRate) {
		return errors.Wrapf(
			errCommissionRateTooLarge,
			"rate:%s max rate:%s", c.CommissionRates.Rate.String(),
			c.CommissionRates.MaxRate.String(),
		)
	}

	if c.CommissionRates.MaxChangeRate.GT(c.CommissionRates.MaxRate) {
		return errors.Wrapf(
			errCommissionRateTooLarge,
			"rate:%s max change rate:%s", c.CommissionRates.Rate.String(),
			c.CommissionRates.MaxChangeRate.String(),
		)
	}
	return nil
}
