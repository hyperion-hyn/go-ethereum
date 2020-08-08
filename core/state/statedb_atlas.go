package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/staking/effective"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/pkg/errors"
	"math/big"
)

var (
	errValidatorNotExist    = errors.New("staking validator does not exist")
	errRedelegationNotExist = errors.New("redelegation does not exist")

	ValidatorStorageAddress = common.Address{}
)

func (s *StateDB) ValidatorPool() *restaking.Storage_ValidatorPool_ {
	if s.validatorPool == nil {
		var g restaking.Global_t
		globalSt := restaking.New(&g, s, ValidatorStorageAddress, common.Big0)
		s.validatorPool = globalSt.ValidatorPool()
	}
	return s.validatorPool
}

func (s *StateDB) ValidatorByAddress(validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	validator, ok := s.ValidatorPool().Validators().Get(validatorAddress)
	if !ok {
		return nil, errValidatorNotExist
	}
	return validator, nil
}

func (s *StateDB) ValidatorList() []common.Address {
	return s.validatorPool.Validators().AllKeys()
}

// AddRedelegationReward distributes the reward to all the delegators based on stake percentage.
func (s *StateDB) AddRedelegationReward(snapshot *restaking.Storage_ValidatorWrapper_, reward *big.Int,
	shareLookup map[common.Address]common.Dec) error {
	valAddr := snapshot.Validator().ValidatorAddress().Value()
	if reward.Cmp(common.Big0) == 0 {
		log.Info("0 given as reward", "validator", valAddr)
		return nil
	}

	curValidator, err := s.ValidatorByAddress(valAddr)
	if err != nil {
		return errors.Wrapf(err, "failed to distribute rewards: validator does not exist")
	}

	if curValidator.Validator().Status().Value().Uint64() == uint64(effective.Banned) {
		log.Info("cannot add reward to banned validator", "validator", valAddr)
		return nil
	}
	curValidator.AddBlockReward(reward)

	// Payout each delegator's (non-operator) reward
	rewardPool := big.NewInt(0).Set(reward)
	commissionRate := snapshot.Validator().Commission().CommissionRates().Rate().Value()
	noncommisionRate := common.OneDec().Sub(commissionRate)
	for _, delegatorAddress := range snapshot.Redelegations().AllKeys() {
		if snapshot.Validator().OperatorAddresses().Set().Get(delegatorAddress).Value() {
			continue
		}

		percentage, ok := shareLookup[delegatorAddress]
		if !ok {
			return errors.Wrapf(err, "missing delegation shares for reward distribution")
		}
		rewardInt := percentage.MulInt(rewardPool).Mul(noncommisionRate).RoundInt()

		curDelegation, ok := curValidator.Redelegations().Get(delegatorAddress)
		if !ok {
			return errRedelegationNotExist
		}
		curDelegation.AddAmount(rewardInt)
		rewardPool.Sub(rewardPool, rewardInt)
	}

	// Payout each operator's reward
	rewardForOperators := big.NewInt(0).Set(rewardPool)
	var largestOperator common.Address
	largestAmount := common.Big0
	for _, operator := range snapshot.Validator().OperatorAddresses().AllKeys() {
		redelegation, ok := snapshot.Redelegations().Get(operator)
		if !ok {
			continue
		}
		if redelegation.Amount().Value().Cmp(largestAmount) > 0 {
			largestAmount = big.NewInt(0).Set(redelegation.Amount().Value())
			largestOperator = operator
		}
		percentage := common.NewDecFromBigInt(redelegation.Amount().Value()).QuoInt(snapshot.TotalDelegationByOperator().Value())
		rewardInt := percentage.MulInt(rewardForOperators).RoundInt()

		curDelegation, ok := curValidator.Redelegations().Get(operator)
		if !ok {
			return errRedelegationNotExist
		}
		curDelegation.AddReward(rewardInt)
		rewardPool.Sub(rewardPool, rewardInt)
	}

	// The last remaining bit belongs to the operator with largest delegation
	if rewardForOperators.Cmp(common.Big0) > 0 {
		redelegation, _ := snapshot.Redelegations().Get(largestOperator)
		redelegation.AddReward(rewardForOperators)
	}
	return nil
}
