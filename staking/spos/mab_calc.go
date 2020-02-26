package spos

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/numeric"
	"math/big"
)

const (
	MaxUpdatedInterval = 10 // N
)

type ValidatorMAB struct {
	ValidatorAddress common.Address
	DelegationMABs   []DelegationMAB
	Records          []mabRecord
}

type DelegationMAB struct {
	DelegatorAddress common.Address
	Records          []mabRecord
}

type mabRecord struct {
	Amount   *big.Int
	BlockNum *big.Int
	MAB      numeric.Dec
}

func (v ValidatorMAB) Calc(blockNum *big.Int) numeric.Dec {
	// genesis
	if blockNum.Cmp(big.NewInt(0)) == 0 && len(v.Records) > 0 && v.Records[0].BlockNum.Cmp(big.NewInt(0)) == 0 {
		return v.Records[0].MAB
	}

	for i := len(v.Records) - 1; i >= 0; i-- {
		record := v.Records[i]
		if record.BlockNum.Cmp(blockNum) == 0 {
			return record.MAB
		}
		if record.BlockNum.Cmp(blockNum) > 0 {
			continue
		}
		return calcMAB(big.NewInt(0).Set(blockNum).Sub(blockNum, record.BlockNum), record.Amount, record.MAB, record.Amount)
	}
	return numeric.NewDec(0)
}

// UpdateMAB should be called only when validator's delegation changes
func (v *ValidatorMAB) UpdateMAB(blockNum, amount *big.Int) {
	v.Records = updateMABRecords(v.Records, blockNum, amount)
}

func (d DelegationMAB) Calc(blockNum *big.Int) numeric.Dec {
	// genesis
	if blockNum.Cmp(big.NewInt(0)) == 0 && len(d.Records) > 0 && d.Records[0].BlockNum.Cmp(big.NewInt(0)) == 0 {
		return d.Records[0].MAB
	}

	for i := len(d.Records) - 1; i >= 0; i-- {
		record := d.Records[i]
		if record.BlockNum.Cmp(blockNum) == 0 {
			return record.MAB
		}
		if record.BlockNum.Cmp(blockNum) > 0 {
			continue
		}
		return calcMAB(big.NewInt(0).Set(blockNum).Sub(blockNum, record.BlockNum), record.Amount, record.MAB, record.Amount)
	}
	return numeric.NewDec(0)
}

func (d *DelegationMAB) UpdateMAB(blockNum, amount *big.Int) {
	d.Records = updateMABRecords(d.Records, blockNum, amount)
}

func updateMABRecords(records []mabRecord, blockNum, amount *big.Int) []mabRecord {
	// Genesis, Amount == MAB
	if blockNum.Cmp(big.NewInt(0)) == 0 {
		records = append(records, mabRecord{
			Amount:   amount,
			BlockNum: blockNum,
			MAB:      numeric.NewDecFromBigInt(amount),
		})
		return records
	}

	if len(records) == 0 {
		records = append(records, mabRecord{
			Amount:   amount,
			BlockNum: blockNum,
			MAB:      numeric.NewDec(0),
		})
		return records
	}

	lastBlock := records[len(records)-1].BlockNum
	lastBalance := records[len(records)-1].Amount
	lastMAB := records[len(records)-1].MAB
	if lastBlock.Cmp(blockNum) >= 0 {
		log.Error("Block number should be greater than last block number",
			"BlockNum", blockNum, "lastBlockNum", lastBlock)
		return records
	}
	records = append(records, mabRecord{
		Amount:   amount,
		BlockNum: blockNum,
		MAB:      calcMAB(big.NewInt(0).Set(blockNum).Sub(blockNum, lastBlock), lastBalance, lastMAB, amount),
	})
	return records
}

func calcMAB(blockNumDiff, lastBalance *big.Int, lastMAB numeric.Dec, balance *big.Int) numeric.Dec {
	alpha := numeric.NewDecFromInt(blockNumDiff).QuoInt64(MaxUpdatedInterval)
	if alpha.GT(numeric.NewDec(1)) {
		alpha = numeric.NewDec(1)
	}

	mab := alpha.MulInt(lastBalance).Add(numeric.NewDec(1).Sub(alpha).Mul(lastMAB))
	if numeric.NewDecFromInt(balance).GT(mab) {
		return mab
	}
	return numeric.NewDecFromInt(balance)
}
