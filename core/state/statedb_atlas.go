package state

import (
    "math/big"
    "github.com/pkg/errors"


    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/log"
    "github.com/ethereum/go-ethereum/rlp"
    "github.com/ethereum/go-ethereum/staking/effective"
    stk "github.com/ethereum/go-ethereum/staking/types"
    "github.com/ethereum/go-ethereum/numeric"
    "github.com/ethereum/go-ethereum/staking"
)

var (
    errAddressNotPresent = errors.New("address not present in state")
)

// ValidatorWrapper retrieves the existing validator in the cache.
// The return value is a reference to the actual validator object in state.
// The modification on it will be committed to the state object when Finalize()
// is called.
func (s *StateDB) ValidatorWrapper(
	addr common.Address,
) (*stk.ValidatorWrapper, error) {
	// Read cache first
	cached, ok := s.stateValidators[addr]
	if ok {
		return cached, nil
	}

	val, err := s.ValidatorWrapperCopy(addr)
	if err != nil {
		return nil, err
	}
	// populate cache if the validator is not in it
	s.stateValidators[addr] = val
	return val, nil

}

// ValidatorWrapperCopy retrieves the existing validator as a copy from state object.
// Changes on the copy has to be explicitly commited with UpdateValidatorWrapper()
// to take effect.
func (s *StateDB) ValidatorWrapperCopy(
    addr common.Address,
) (*stk.ValidatorWrapper, error) {
    by := s.GetCode(addr)
    if len(by) == 0 {
        return nil, errAddressNotPresent
    }
    val := stk.ValidatorWrapper{}
    if err := rlp.DecodeBytes(by, &val); err != nil {
        return nil, errors.Wrapf(
            err,
            "could not decode for %s",
			common.MustAddressToBech32(addr),
		)
	}
	return &val, nil
}

const doNotEnforceMaxBLS = -1

// UpdateValidatorWrapper updates staking information of
// a given validator (including delegation info)
func (s *StateDB) UpdateValidatorWrapper(
	addr common.Address, val *stk.ValidatorWrapper,
) error {
	if err := val.SanityCheck(doNotEnforceMaxBLS); err != nil {
		return err
	}

	by, err := rlp.EncodeToBytes(val)
	if err != nil {
		return err
	}
	s.SetCode(addr, by)
	// update cache
	s.stateValidators[addr] = val
	return nil
}

// SetValidatorFlag checks whether it is a validator object
func (s *StateDB) SetValidatorFlag(addr common.Address) {
	s.SetState(addr, staking.IsValidatorKey, staking.IsValidator)
}

// UnsetValidatorFlag checks whether it is a validator object
func (s *StateDB) UnsetValidatorFlag(addr common.Address) {
	s.SetState(addr, staking.IsValidatorKey, common.Hash{})
}

// IsValidator checks whether it is a validator object
func (s *StateDB) IsValidator(addr common.Address) bool {
	so := s.getStateObject(addr)
	if so == nil {
		return false
	}
	return so.IsValidator(s.db)
}

var (
	zero = numeric.ZeroDec()
)

// AddReward distributes the reward to all the delegators based on stake percentage.
func (s *StateDB) AddReward(snapshot *stk.ValidatorWrapper, reward *big.Int, shareLookup map[common.Address]numeric.Dec) error {
	if reward.Cmp(common.Big0) == 0 {
	    log.Info("0 given as reward", "validator", []byte(snapshot.String()))
		return nil
	}

	curValidator, err := s.ValidatorWrapper(snapshot.Address)
	if err != nil {
		return errors.Wrapf(err, "failed to distribute rewards: validator does not exist")
	}

	if curValidator.Status == effective.Banned {
	    log.Info("cannot add reward to banned validator", "slashed-validator", []byte(curValidator.String()))
		return nil
	}

	rewardPool := big.NewInt(0).Set(reward)
	curValidator.BlockReward.Add(curValidator.BlockReward, reward)
	// Payout commission
	if r := snapshot.Validator.CommissionRates.Rate; r.GT(zero) {
		commissionInt := r.MulInt(reward).RoundInt()
		curValidator.Delegations[0].Reward.Add(
			curValidator.Delegations[0].Reward,
			commissionInt,
		)
		rewardPool.Sub(rewardPool, commissionInt)
	}

	// Payout each delegator's reward pro-rata
	totalRewardForDelegators := big.NewInt(0).Set(rewardPool)
	for i := range snapshot.Delegations {
		delegation := snapshot.Delegations[i]
		percentage, ok := shareLookup[delegation.DelegatorAddress]

		if !ok {
			return errors.Wrapf(err, "missing delegation shares for reward distribution")
		}

		rewardInt := percentage.MulInt(totalRewardForDelegators).RoundInt()
		curDelegation := curValidator.Delegations[i]
		curDelegation.Reward.Add(curDelegation.Reward, rewardInt)
		rewardPool.Sub(rewardPool, rewardInt)
	}

	// The last remaining bit belongs to the validator (remember the validator's self delegation is
	// always at index 0)
	if rewardPool.Cmp(common.Big0) > 0 {
		curValidator.Delegations[0].Reward.Add(curValidator.Delegations[0].Reward, rewardPool)
	}

	return nil
}

// ATLAS - END
