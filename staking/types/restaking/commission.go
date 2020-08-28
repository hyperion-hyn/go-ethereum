package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	hundredPercent = common.OneDec()
	zeroPercent    = common.ZeroDec()
)

// Copy deep copies the staking.CommissionRates
func (cr *CommissionRates_) Copy() CommissionRates_ {
	return CommissionRates_{
		Rate:          cr.Rate.Copy(),
		MaxRate:       cr.MaxRate.Copy(),
		MaxChangeRate: cr.MaxChangeRate.Copy(),
	}
}

func (c *Commission_) SanityCheck() error {
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


func (s *Storage_Commission_) Load() *Commission_ {
	s.CommissionRates().Rate().Value()
	s.CommissionRates().MaxChangeRate().Value()
	s.CommissionRates().MaxRate().Value()
	s.UpdateHeight().Value()
	return s.obj
}