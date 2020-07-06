package types

import (
	"github.com/ethereum/go-ethereum/core/storage"
	"github.com/ethereum/go-ethereum/numeric"
	"math/big"
)

type CommissionStorage struct {
	storage *storage.Storage
}

func (c *CommissionStorage) ToCommission() *Commission {
	return nil
}

func (c *CommissionStorage) GetCommissionRates() *CommissionRatesStorage {
	return nil
}

func (c *CommissionStorage) SetCommissionRates(rates *CommissionRates) {
}

func (c *CommissionStorage) GetUpdateHeight() *big.Int {
	return nil
}

func (c *CommissionStorage) SetUpdateHeight(updateHeight *big.Int) {
}

type CommissionRatesStorage struct {
	storage *storage.Storage
}

func (c *CommissionRatesStorage) ToCommissionRates() *CommissionRates {
	return nil
}

func (c *CommissionRatesStorage) GetRate() *numeric.Dec {
	return nil
}

func (c *CommissionRatesStorage) SetRate(rate *numeric.Dec) {
}

func (c *CommissionRatesStorage) GetMaxRate() *numeric.Dec {
	return nil
}

func (c *CommissionRatesStorage) SetMaxRate(maxRate *numeric.Dec) {
}

func (c *CommissionRatesStorage) GetMaxChangeRate() *numeric.Dec {
	return nil
}

func (c *CommissionRatesStorage) SetMaxChangeRate(maxChangeRate *numeric.Dec) {
}
