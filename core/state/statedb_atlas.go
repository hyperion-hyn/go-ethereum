package state

import (
    "math/big"
    "fmt"
    "errors"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/rlp"
    stk "github.com/ethereum/go-ethereum/staking/types"
    "github.com/ethereum/go-ethereum/numeric"
    "github.com/ethereum/go-ethereum/staking"
)

// ATLAS(yhx): this function is specific to atlas and should be moved to atlas package.
// GetStakingInfo returns staking information of a given validator (including delegation info)
func (s *StateDB) GetStakingInfo(addr common.Address) *stk.ValidatorWrapper {
    //ATLAS(yhx): MaxCodeSize = 24576, maybe code is not suitable to store staking information.
    by := s.GetCode(addr)
    if len(by) == 0 {
        return nil
    }
    val := stk.ValidatorWrapper{}
    err := rlp.DecodeBytes(by, &val)
    if err != nil {
        fmt.Printf("GetStakingInfo unable to decode: %v\n", err)
        return nil
    }
    return &val
}

// ATLAS(yhx): this function is specific to atlas and should be moved to atlas package.
// UpdateStakingInfo update staking information of a given validator (including delegation info)
func (s *StateDB) UpdateStakingInfo(addr common.Address, val *stk.ValidatorWrapper) error {
    // TODO: check ValidatorWrapper's compliance
    if err := val.SanityCheck(); err != nil {
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

    curValidator := s.GetStakingInfo(snapshot.Validator.Address)
    if curValidator == nil {
        return errors.New("failed to distribute rewards: validator does not exist")
    }

    // Payout commission
    commissionInt := snapshot.Validator.CommissionRates.Rate.MulInt(reward).RoundInt()

    curValidator.Delegations[0].Reward.Add(curValidator.Delegations[0].Reward, commissionInt)
    rewardPool.Sub(rewardPool, commissionInt)

    totalRewardForDelegators := big.NewInt(0).Set(rewardPool)

    // Payout each delegator's reward pro-rata
    totalDelegationDec := numeric.NewDecFromBigInt(snapshot.TotalDelegation())
    for i := range snapshot.Delegations {
        delegation := snapshot.Delegations[i]
        percentage := numeric.NewDecFromBigInt(delegation.Amount).Quo(totalDelegationDec) // percentage = <this_delegator_amount>/<total_delegation>
        rewardInt := percentage.MulInt(totalRewardForDelegators).RoundInt()

        curDelegation := curValidator.Delegations[i]
        curDelegation.Reward.Add(curDelegation.Reward, rewardInt)
        rewardPool.Sub(rewardPool, rewardInt)
    }

    // The last remaining bit belongs to the validator (remember the validator's self delegation is always at index 0)
    if rewardPool.Cmp(big.NewInt(0)) > 0 {
        curValidator.Delegations[0].Reward.Add(curValidator.Delegations[0].Reward, rewardPool)
    }

    return s.UpdateStakingInfo(curValidator.Validator.Address, curValidator)
}

// ATLAS - END
