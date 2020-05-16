package state

import (
    "math/big"
    "github.com/pkg/errors"


    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/rlp"
    stk "github.com/ethereum/go-ethereum/staking/types"
    "github.com/ethereum/go-ethereum/numeric"
    "github.com/ethereum/go-ethereum/staking"
)

// ATLAS(yhx): this function is specific to atlas and should be moved to atlas package.

var (
    errAddressNotPresent = errors.New("address not present in state")
)

// ValidatorWrapper  ..
func (db *StateDB) ValidatorWrapper(
    addr common.Address,
) (*stk.ValidatorWrapper, error) {
    by := db.GetCode(addr)
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

// AddReward distributes the reward to all the delegators based on stake percentage.
func (s *StateDB) AddReward(snapshot *stk.ValidatorWrapper, reward *big.Int) error {
    rewardPool := big.NewInt(0).Set(reward)
	curValidator, err := s.ValidatorWrapper(snapshot.Address)
	if err != nil {
		return errors.Wrapf(err, "failed to distribute rewards: validator does not exist")
    }

	curValidator.BlockReward.Add(curValidator.BlockReward, reward)
    // Payout commission
    commissionInt := snapshot.Validator.CommissionRates.Rate.MulInt(reward).RoundInt()
    curValidator.Delegations[0].Reward.Add(curValidator.Delegations[0].Reward, commissionInt)
    rewardPool.Sub(rewardPool, commissionInt)
    totalRewardForDelegators := big.NewInt(0).Set(rewardPool)
    // Payout each delegator's reward pro-rata
    totalDelegationDec := numeric.NewDecFromBigInt(snapshot.TotalDelegation())
    for i := range snapshot.Delegations {
        delegation := snapshot.Delegations[i]
		// NOTE percentage = <this_delegator_amount>/<total_delegation>
        percentage := numeric.NewDecFromBigInt(delegation.Amount).Quo(totalDelegationDec)
        rewardInt := percentage.MulInt(totalRewardForDelegators).RoundInt()
        curDelegation := curValidator.Delegations[i]
        curDelegation.Reward.Add(curDelegation.Reward, rewardInt)
        rewardPool.Sub(rewardPool, rewardInt)
    }


    // The last remaining bit belongs to the validator (remember the validator's self delegation is always at index 0)
    if rewardPool.Cmp(common.Big0) > 0 {
        curValidator.Delegations[0].Reward.Add(curValidator.Delegations[0].Reward, rewardPool)
    }

    return s.UpdateValidatorWrapper(curValidator.Validator.Address, curValidator)
}

// ATLAS - END
