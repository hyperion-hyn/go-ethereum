package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/pkg/errors"
	"math/big"
)

var (
	errValidatorNotExist         = errors.New("validator does not exist")
	errValidatorSnapshotNotExist = errors.New("validator snapshot does not exist")
	errRedelegationNotExist      = errors.New("redelegation does not exist")

	validatorStorageAddress = common.HexToAddress("0x69270f88069d56dc62bd62b0b9f2b302a2b820a8")
)

func (s *StateDB) ValidatorPool() *restaking.Storage_ValidatorPool_ {
	// singleton
	s.validatorOnce.Do(func() {
		var g restaking.Global_t
		globalSt := restaking.New(&g, s, validatorStorageAddress, common.Big0)
		s.validatorPool = globalSt.ValidatorPool()
	})
	return s.validatorPool
}

func (s *StateDB) ValidatorByAddress(validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	validator, ok := s.ValidatorPool().Validators().Get(validatorAddress)
	if !ok {
		return nil, errValidatorNotExist
	}
	return validator, nil
}

func (s *StateDB) ValidatorSnapshotByAddress(validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	validator, ok := s.ValidatorPool().ValidatorSnapshots().Get(validatorAddress)
	if !ok {
		return nil, errValidatorSnapshotNotExist
	}
	return validator, nil
}

func (s *StateDB) ValidatorList() []common.Address {
	return s.ValidatorPool().Validators().AllKeys()
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

	if curValidator.Validator().Status().Value() == uint8(restaking.Banned) {
		log.Info("cannot add reward to banned validator", "validator", valAddr)
		return nil
	}
	curValidator.AddBlockReward(reward)

	// Payout each delegator's (non-operator) reward
	rewardPool := big.NewInt(0).Set(reward)
	commissionRate := snapshot.Validator().Commission().CommissionRates().Rate().Value()
	noncommisionRate := common.OneDec().Sub(commissionRate)
	if noncommisionRate.GT(common.ZeroDec()) {
		for _, delegatorAddress := range snapshot.Redelegations().AllKeys() {
			if snapshot.IsOperator(delegatorAddress) {
				continue
			}

			percentage, ok := shareLookup[delegatorAddress]
			if !ok {
				return errors.Wrapf(err, "missing delegation shares for reward distribution")
			}
			if percentage.IsNil() || percentage.IsZero() {
				continue
			}
			rewardInt := percentage.MulInt(rewardPool).Mul(noncommisionRate).RoundInt()

			curDelegation, ok := curValidator.Redelegations().Get(delegatorAddress)
			if !ok {
				return errors.Wrap(errRedelegationNotExist, "missing delegation for reward distribution")
			}
			curDelegation.AddReward(rewardInt)
			rewardPool.Sub(rewardPool, rewardInt)
		}
	}

	// Payout each operator's reward
	totalDelegationByOperator := snapshot.TotalDelegationByOperator().Value()
	if totalDelegationByOperator.Sign() == 0 {
		return errors.New("missing total delegation of operator")
	}

	rewardForOperators := big.NewInt(0).Set(rewardPool)
	emptyAddress := common.Address{}
	largestOperator := emptyAddress
	largestAmount := common.Big0
	for _, operator := range snapshot.Validator().OperatorAddresses().AllKeys() {
		redelegationSnapshot, ok := snapshot.Redelegations().Get(operator)
		if !ok {
			continue
		}
		amtSnapshot := redelegationSnapshot.Amount().Value()
		if amtSnapshot.Sign() == 0 {
			continue
		}
		if amtSnapshot.Cmp(largestAmount) > 0 {
			largestAmount = amtSnapshot
			largestOperator = operator
		}
		percentage := common.NewDecFromBigInt(amtSnapshot).QuoInt(totalDelegationByOperator)
		rewardInt := percentage.MulInt(rewardForOperators).RoundInt()

		curDelegation, ok := curValidator.Redelegations().Get(operator)
		if !ok {
			return errors.Wrap(errRedelegationNotExist, "missing delegation of operator for reward distribution")
		}
		curDelegation.AddReward(rewardInt)
		rewardPool.Sub(rewardPool, rewardInt)
	}

	// The last remaining bit belongs to the operator with largest delegation
	if rewardPool.Cmp(common.Big0) > 0 && largestOperator != emptyAddress {
		redelegation, _ := curValidator.Redelegations().Get(largestOperator)
		redelegation.AddReward(rewardPool)
	}
	return nil
}

/**
 * IncreaseValidatorNonceIfZero avoids account state of validators would be delete if its nonce and balance are zero
 */
func (s *StateDB) IncreaseValidatorNonceIfZero() {
	if s.GetNonce(validatorStorageAddress) == 0 {
		s.SetNonce(validatorStorageAddress, 1)
	}
}
