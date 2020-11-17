package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
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
		actualAmount := burnTokenFromFoundationAccount(stateDB, foundationAddress, preBurningAmount)
		log.Info("genesis internal burning", "block", blockNum, "expect", preBurningAmount, "actual", actualAmount)

		// write off-chain record
		receipt := burning.Receipt{
			InternalAmount: actualAmount,
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
	expectedAmount, err := burning.CalculateInternalBurningAmount(numOfActiveMap3Node, numOfScalingCycle, require)
	if err != nil {
		return err
	}

	actualAmount := burnTokenFromFoundationAccount(stateDB, foundationAddress, expectedAmount)
	log.Info("internal burning", "block", blockNum, "expect", expectedAmount, "actual", actualAmount)
	if actualAmount.Sign() == 0 {
		return nil
	}

	// write off-chain burning receipt
	receipt := burning.Receipt{
		InternalAmount: actualAmount,
		ExternalAmount: common.Big0,
		BlockNum:       blockNum,
	}
	receipt.DoHash()
	db := chain.Database()
	rawdb.WriteTokenBurningReceipt(db, receipt)
	return nil
}

func burnTokenFromFoundationAccount(stateDB vm.StateDB, foundationAddress common.Address, amount *big.Int) *big.Int {
	balance := stateDB.GetBalance(foundationAddress)
	if balance.Cmp(amount) < 0 {
		amount = balance
	}
	stateDB.SubBalance(foundationAddress, amount)
	return amount
}