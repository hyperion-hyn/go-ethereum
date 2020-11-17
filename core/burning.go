package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/staking/burning"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"math/big"
)

var (
	preBurningAmount = new(big.Int).Mul(big.NewInt(510000000), big.NewInt(params.Ether)) // (100 * 57 *6 +55 * 50 * 6 + 100 * 1 *3) * 10^18
)

func CheckAndPreburnToken(chain ChainContext, stateDB vm.StateDB, blockNum *big.Int) {
	if blockNum.Cmp(big.NewInt(int64(chain.Config().Atlas.HYNBurningBlock))) == 0 {
		stateDB.SubBalance(foundationAddress, preBurningAmount)

		// write off-chain record
		receipt := burning.Receipt{
			InternalAmount: preBurningAmount,
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
	numOfActiveMap3Node := stateDB.Map3NodePool().SizeByMap3Status(microstaking.Active)
	amount, err := burning.CalculateInternalBurningAmount(numOfActiveMap3Node, numOfScalingCycle, require)
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
