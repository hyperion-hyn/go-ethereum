package core

import (
    "bytes"
    "math/big"

    "github.com/pkg/errors"

    "github.com/ethereum/go-ethereum/core/vm"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/rlp"
    "github.com/ethereum/go-ethereum/log"

    staking "github.com/ethereum/go-ethereum/staking/types"
)

// ATLAS
// StakingTransitionDb will transition the state by applying the staking message and
// returning the result including the used gas. It returns an error if failed.
// It is used for staking transaction only
func (st *StateTransition) StakingTransitionDb() (usedGas uint64, err error) {
    if err = st.preCheck(); err != nil {
        return
    }
    msg := st.msg
    sender := vm.AccountRef(msg.From())
    homestead := st.evm.ChainConfig().IsHomestead(st.evm.BlockNumber)
    istanbul := st.evm.ChainConfig().IsIstanbul(st.evm.BlockNumber)

    // Pay intrinsic gas
    // TODO: propose staking-specific formula for staking transaction
    gas, err := IntrinsicGas(st.data, false, homestead, istanbul)
    if err != nil {
        return 0, err
    }
    if err = st.useGas(gas); err != nil {
        return 0, err
    }

    // Increment the nonce for the next transaction
    st.state.SetNonce(msg.From(), st.state.GetNonce(sender.Address())+1)

    // ATLAS(yhx): should we verify st.evm.BlockNumber and st.evm.BlockNumber?
    switch msg.Type() {
    case types.StakeCreateVal:
        stkMsg := &staking.CreateValidator{}
        if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
            return 0, err
        }
        log.Info("Staking Message", "Type", msg.Type(), "Message", stkMsg)
        if msg.From() != stkMsg.ValidatorAddress {
            return 0, errInvalidSigner
        }
        err = st.applyCreateValidatorTx(stkMsg, st.evm.BlockNumber)

    case types.StakeEditVal:
        stkMsg := &staking.EditValidator{}
        if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
            return 0, err
        }
        log.Info("Staking Message", "Type", msg.Type(), "Message", stkMsg)
        if msg.From() != stkMsg.ValidatorAddress {
            return 0, errInvalidSigner
        }
        err = st.applyEditValidatorTx(stkMsg, st.evm.BlockNumber)

    case types.Delegate:
        stkMsg := &staking.Delegate{}
        if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
            return 0, err
        }
        log.Info("Staking Message", "Type", msg.Type(), "Message", stkMsg)
        if msg.From() != stkMsg.DelegatorAddress {
            return 0, errInvalidSigner
        }
        err = st.applyDelegateTx(stkMsg)

    case types.Undelegate:
        stkMsg := &staking.Undelegate{}
        if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
            return 0, err
        }
        log.Info("Staking Message", "Type", msg.Type(), "Message", stkMsg)
        if msg.From() != stkMsg.DelegatorAddress {
            return 0, errInvalidSigner
        }
        err = st.applyUndelegateTx(stkMsg)
    case types.CollectRewards:
        stkMsg := &staking.CollectRewards{}
        if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
            return 0, err
        }
        log.Info("Staking Message", "Type", msg.Type(), "Message", stkMsg)
        if msg.From() != stkMsg.DelegatorAddress {
            return 0, errInvalidSigner
        }
        err = st.applyCollectRewards(stkMsg)
    default:
        return 0, staking.ErrInvalidStakingKind
    }
    st.refundGas()

    txFee := new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice)
    st.state.AddBalance(st.evm.Coinbase, txFee)

    return st.gasUsed(), err
}

func (st *StateTransition) applyCreateValidatorTx(createValidator *staking.CreateValidator, blockNum *big.Int) error {
    if createValidator.Amount.Sign() == -1 {
        return errNegativeAmount
    }

    if val := createValidator.ValidatorAddress; st.state.IsValidator(val) {
        return errors.Wrapf(errValidatorExist, common.MustAddressToBech32(val))
    }

    if !CanTransfer(st.state, createValidator.ValidatorAddress, createValidator.Amount) {
        return errInsufficientBalanceForStake
    }

    v, err := staking.CreateValidatorFromNewMsg(createValidator, blockNum)
    if err != nil {
        return err
    }

    delegations := []staking.Delegation{
        staking.NewDelegation(v.Address, createValidator.Amount),
    }
    wrapper := staking.ValidatorWrapper{*v, delegations}

    if err := st.state.UpdateStakingInfo(v.Address, &wrapper); err != nil {
        return err
    }

    st.state.SetValidatorFlag(v.Address)

    st.state.SubBalance(v.Address, createValidator.Amount)
    return nil
}

func (st *StateTransition) applyEditValidatorTx(editValidator *staking.EditValidator, blockNum *big.Int) error {
    if !st.state.IsValidator(editValidator.ValidatorAddress) {
        return errValidatorNotExist
    }

    wrapper := st.state.GetStakingInfo(editValidator.ValidatorAddress)
    if wrapper == nil {
        return errValidatorNotExist
    }

    if err := staking.UpdateValidatorFromEditMsg(&wrapper.Validator, editValidator); err != nil {
        return err
    }
    newRate := wrapper.Validator.Rate

    // TODO: make sure we are reading from the correct snapshot
    snapshotValidator, err := st.bc.ReadValidatorSnapshot(wrapper.Address)
    if err != nil {
        return err
    }
    rateAtBeginningOfEpoch := snapshotValidator.Rate

    if rateAtBeginningOfEpoch.IsNil() || (!newRate.IsNil() && !rateAtBeginningOfEpoch.Equal(newRate)) {
        wrapper.Validator.UpdateHeight = blockNum
    }

    if newRate.Sub(rateAtBeginningOfEpoch).Abs().GT(wrapper.Validator.MaxChangeRate) {
        return errCommissionRateChangeTooFast
    }

    if newRate.GT(wrapper.Validator.MaxRate) {
        return errCommissionRateChangeTooHigh
    }

    if err := st.state.UpdateStakingInfo(wrapper.Address, wrapper); err != nil {
        return err
    }
    return nil
}

func (st *StateTransition) applyDelegateTx(delegate *staking.Delegate) error {
    if delegate.Amount.Sign() == -1 {
        return errNegativeAmount
    }

    if !st.state.IsValidator(delegate.ValidatorAddress) {
        return errValidatorNotExist
    }
    wrapper := st.state.GetStakingInfo(delegate.ValidatorAddress)
    if wrapper == nil {
        return errValidatorNotExist
    }

    stateDB := st.state
    delegatorExist := false
    for i := range wrapper.Delegations {
        delegation := &wrapper.Delegations[i]
        if bytes.Equal(delegation.DelegatorAddress.Bytes(), delegate.DelegatorAddress.Bytes()) {
            delegatorExist = true
            totalInUndelegation := delegation.TotalInUndelegation()
            balance := stateDB.GetBalance(delegate.DelegatorAddress)
            // If the sum of normal balance and the total amount of tokens in undelegation is greater than the amount to delegate
            if big.NewInt(0).Add(totalInUndelegation, balance).Cmp(delegate.Amount) >= 0 {
                // Firstly use the tokens in undelegation to delegate (redelegate)
                delegateBalance := big.NewInt(0).Set(delegate.Amount)
                // Use the latest undelegated token first as it has the longest remaining locking time.
                i := len(delegation.Undelegations) - 1
                for ; i >= 0; i-- {
                    if delegation.Undelegations[i].Amount.Cmp(delegateBalance) <= 0 {
                        delegateBalance.Sub(delegateBalance, delegation.Undelegations[i].Amount)
                    } else {
                        delegation.Undelegations[i].Amount.Sub(delegation.Undelegations[i].Amount, delegateBalance)
                        delegateBalance = big.NewInt(0)
                        break
                    }
                }

                delegation.Undelegations = delegation.Undelegations[:i+1]
                delegation.Amount.Add(delegation.Amount, delegate.Amount)
                err := stateDB.UpdateStakingInfo(wrapper.Validator.Address, wrapper)

                // Secondly, if all locked token are used, try use the balance.
                if err == nil && delegateBalance.Cmp(big.NewInt(0)) > 0 {
                    stateDB.SubBalance(delegate.DelegatorAddress, delegateBalance)
                }
                return err
            }
            return errors.Wrapf(
                errInsufficientBalanceForStake,
                "total-delegated %s own-current-balance %s amount-to-delegate %s",
                totalInUndelegation.String(),
                balance.String(),
                delegate.Amount.String(),
            )
        }
    }

    if !delegatorExist {
        if CanTransfer(stateDB, delegate.DelegatorAddress, delegate.Amount) {
            newDelegator := staking.NewDelegation(delegate.DelegatorAddress, delegate.Amount)
            wrapper.Delegations = append(wrapper.Delegations, newDelegator)

            if err := stateDB.UpdateStakingInfo(wrapper.Validator.Address, wrapper); err == nil {
                stateDB.SubBalance(delegate.DelegatorAddress, delegate.Amount)
            } else {
                return err
            }
        }
    }

    return nil
}

func (st *StateTransition) applyUndelegateTx(undelegate *staking.Undelegate) error {
    if undelegate.Amount.Sign() == -1 {
        return errNegativeAmount
    }

    if !st.state.IsValidator(undelegate.ValidatorAddress) {
        return errValidatorNotExist
    }
    wrapper := st.state.GetStakingInfo(undelegate.ValidatorAddress)
    if wrapper == nil {
        return errValidatorNotExist
    }

    stateDB := st.state
    delegatorExist := false
    for i := range wrapper.Delegations {
        delegation := &wrapper.Delegations[i]
        if bytes.Equal(delegation.DelegatorAddress.Bytes(), undelegate.DelegatorAddress.Bytes()) {
            delegatorExist = true

            err := delegation.Undelegate(st.evm.EpochNumber, undelegate.Amount)
            if err != nil {
                return err
            }
            err = stateDB.UpdateStakingInfo(wrapper.Validator.Address, wrapper)
            return err
        }
    }
    if !delegatorExist {
        return errNoDelegationToUndelegate
    }
    return nil
}

func (st *StateTransition) applyCollectRewards(collectRewards *staking.CollectRewards) error {
    if st.bc == nil {
        return errors.New("[CollectRewards] No chain context provided")
    }
    chainContext := st.bc
    delegations, err := chainContext.ReadDelegationsByDelegator(collectRewards.DelegatorAddress)

    if err != nil {
        return err
    }

    totalRewards := big.NewInt(0)
    for i := range delegations {
        wrapper := st.state.GetStakingInfo(delegations[i].ValidatorAddress)
        if wrapper == nil {
            return errValidatorNotExist
        }

        if uint64(len(wrapper.Delegations)) > delegations[i].Index {
            delegation := &wrapper.Delegations[delegations[i].Index]
            if delegation.Reward.Cmp(big.NewInt(0)) > 0 {
                totalRewards.Add(totalRewards, delegation.Reward)
            }

            delegation.Reward.SetUint64(0)
        }

        err = st.state.UpdateStakingInfo(wrapper.Validator.Address, wrapper)
        if err != nil {
            return err
        }
    }
    if totalRewards.Int64() == 0 {
        return errNoRewardsToCollect
    }
    st.state.AddBalance(collectRewards.DelegatorAddress, totalRewards)
    return nil
}

// ATLAS