package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/effective"
	staking "github.com/ethereum/go-ethereum/staking/types"
	"github.com/pkg/errors"
	"math/big"
)

var (
	ErrMap3NodeNotExist  = errors.New("staking map3 node does not exist")
	ErrValidatorNotExist = errors.New("staking validator does not exist")
)

func (s *StateDB) Map3NodePool() *staking.Map3NodePoolStorage {
	return s.map3NodePool
}

func (s *StateDB) ValidatorPool() *staking.ValidatorPoolStorage {
	return s.validatorPool
}

func (s *StateDB) Map3NodeByAddress(nodeAddress common.Address) (*staking.Map3NodeWrapperStorage, error) {
	node, ok := s.Map3NodePool().GetNodes().Get(nodeAddress)
	if !ok {
		return nil, ErrMap3NodeNotExist
	}
	return node, nil
}

func (s *StateDB) ValidatorByAddress(validatorAddress common.Address) (*staking.ValidatorWrapperStorage, error) {
	validator, ok := s.ValidatorPool().GetValidators().Get(validatorAddress)
	if !ok {
		return nil, ErrValidatorNotExist
	}
	return validator, nil
}

func (s *StateDB) ValidatorList() []common.Address {
	return s.validatorPool.GetValidators().Keys()
}

// AddRedelegationReward distributes the reward to all the delegators based on stake percentage.
func (s *StateDB) AddRedelegationReward(snapshot *staking.ValidatorWrapperStorage, reward *big.Int,
	shareLookup map[common.Address]numeric.Dec) error {
	valAddr := snapshot.GetValidator().GetValidatorAddress()
	if reward.Cmp(common.Big0) == 0 {
		log.Info("0 given as reward", "validator", valAddr)
		return nil
	}

	curValidator, err := s.ValidatorByAddress(valAddr)
	if err != nil {
		return errors.Wrapf(err, "failed to distribute rewards: validator does not exist")
	}

	if curValidator.GetValidator().GetStatus() == effective.Banned {
		log.Info("cannot add reward to banned validator", "validator", valAddr)
		return nil
	}
	curValidator.SetBlockReward(big.NewInt(0).Add(curValidator.GetBlockReward(), reward))

	// Payout each delegator's (non-initiator) reward
	rewardPool := big.NewInt(0).Set(reward)
	commissionRate := snapshot.GetValidator().GetCommission().GetCommissionRates().GetRate()
	noncommisionRate := numeric.OneDec().Sub(*commissionRate)
	for _, delegatorAddress := range snapshot.GetRedelegations().Keys() {
		if snapshot.GetValidator().GetInitiatorAddresses().Contain(delegatorAddress) {
			continue
		}

		percentage, ok := shareLookup[delegatorAddress]
		if !ok {
			return errors.Wrapf(err, "missing delegation shares for reward distribution")
		}
		rewardInt := percentage.MulInt(rewardPool).Mul(noncommisionRate).RoundInt()

		curDelegation, ok := curValidator.GetRedelegations().Get(delegatorAddress)
		if !ok {
			return core.ErrRedelegationNotExist
		}
		curDelegation.SetReward(big.NewInt(0).Add(curDelegation.GetReward(), rewardInt))
		rewardPool.Sub(rewardPool, rewardInt)
	}

	// Payout each initiator's reward
	rewardForInitiators := big.NewInt(0).Set(rewardPool)
	var largestInitiator common.Address
	largestAmount := common.Big0
	for _, initiator := range snapshot.GetValidator().GetInitiatorAddresses().Keys() {
		redelegation, ok := snapshot.GetRedelegations().Get(initiator)
		if !ok {
			continue
		}
		if redelegation.GetAmount().Cmp(largestAmount) > 0 {
			largestAmount = big.NewInt(0).Set(redelegation.GetAmount())
			largestInitiator = initiator
		}
		percentage := numeric.NewDecFromBigInt(redelegation.GetAmount()).QuoInt(snapshot.GetTotalDelegationByInitiator())
		rewardInt := percentage.MulInt(rewardForInitiators).RoundInt()

		curDelegation, ok := curValidator.GetRedelegations().Get(initiator)
		if !ok {
			return core.ErrRedelegationNotExist
		}
		curDelegation.SetReward(big.NewInt(0).Add(curDelegation.GetReward(), rewardInt))
		rewardPool.Sub(rewardPool, rewardInt)
	}

	// The last remaining bit belongs to the initiator with largest delegation
	if rewardForInitiators.Cmp(common.Big0) > 0 {
		redelegation, _ := snapshot.GetRedelegations().Get(largestInitiator)
		redelegation.SetReward(big.NewInt(0).Add(redelegation.GetReward(), rewardForInitiators))
	}
	return nil
}

// AddRedelegationReward distributes the reward to all the delegators based on stake percentage.
func (s *StateDB) AddMicrodelegationReward(snapshot *staking.Map3NodeWrapperStorage, reward *big.Int,
	shareLookup map[common.Address]numeric.Dec) error {
	return nil
}