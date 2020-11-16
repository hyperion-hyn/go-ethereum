package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/staking/burning"
	"github.com/ethereum/go-ethereum/staking/network"
	"math/big"
)

const (
	preBurningAmount = 1 // TODO(ATLAS): calculate
)

func CheckAndPreburnToken(chain ChainContext, stateDB vm.StateDB, blockNum *big.Int) {
	if blockNum.Cmp(big.NewInt(int64(chain.Config().Atlas.HYNBurningBlock))) == 0 {
		amount := big.NewInt(preBurningAmount)
		stateDB.SubBalance(foundationAddress, amount)

		// write off-chain record
		receipt := burning.Receipt{
			InternalAmount: amount,
			ExternalAmount: common.Big0,
			BlockNum:       blockNum,
		}
		receipt.DoHash()
		db := chain.Database()
		rawdb.WriteTokenBurningReceipt(db, receipt)
	}
}

func CanBurnAtEndOfEach30Epochs(chain ChainContext, blockNum, epoch *big.Int) bool {
	if !chain.Config().Atlas.IsHYNBurning(blockNum) {
		return false
	}
	return epoch.Uint64()%30 == 0
}

func BurnTokenByEach30Epochs(chain ChainContext, stateDB vm.StateDB, blockNum *big.Int) error {
	// internal burning
	numOfScalingCycle := network.NumOfScalingCycle(blockNum, chain.Config())
	require, _, _ := network.LatestMicrostakingRequirement(blockNum, chain.Config())
	numOfMap3Node := stateDB.Map3NodePool().Map3Nodes().Keys().Length()
	amount, err := burning.CalculateInternalBurningAmount(numOfMap3Node, numOfScalingCycle, require)
	if err != nil {
		return err
	}
	stateDB.SubBalance(foundationAddress, amount)

	// write off-chain burning receipt
	receipt := burning.Receipt{
		InternalAmount: amount,
		ExternalAmount: common.Big0,
		BlockNum:       blockNum,
	}
	receipt.DoHash()
	db := chain.Database()
	rawdb.WriteTokenBurningReceipt(db, receipt)
	return nil
}
