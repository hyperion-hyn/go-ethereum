package core

import (
    "math/big"

    "github.com/pkg/errors"

    "github.com/ethereum/go-ethereum/core/vm"
    "github.com/ethereum/go-ethereum/staking/network"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/log"
    "github.com/ethereum/go-ethereum/rlp"

    staking2 "github.com/ethereum/go-ethereum/staking"
    staking "github.com/ethereum/go-ethereum/staking/types"
)

var (
    errDupIdentity                 = errors.New("validator identity exists")
    errDupBlsKey                   = errors.New("BLS key exists")
)

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransitionEx(evm *vm.EVM, msg Message, gp *GasPool) *StateTransition {
    retval := NewStateTransition(evm, msg, gp)
    retval.bc = evm.ChainContext
    return retval
}

// ATLAS: ApplyStakingMessage computes the new state for staking message
func ApplyStakingMessage(evm *vm.EVM, msg Message, gp *GasPool) (uint64, error) {
    return NewStateTransitionEx(evm, msg, gp).StakingTransitionDb()
}

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

    switch msg.Type() {
    case types.StakeCreateVal:
        stkMsg := &staking.CreateValidator{}
        if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
            return 0, err
        }
        log.Info("staking", "type", msg.Type(), "gas", gas, "txn", stkMsg)
        if msg.From() != stkMsg.ValidatorAddress {
            return 0, errInvalidSigner
        }
        err = st.verifyAndApplyCreateValidatorTx(stkMsg, st.evm.BlockNumber)

    case types.StakeEditVal:
        stkMsg := &staking.EditValidator{}
        if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
            return 0, err
        }
        log.Info("staking", "type", msg.Type(), "gas", gas, "txn", stkMsg)
        if msg.From() != stkMsg.ValidatorAddress {
            return 0, errInvalidSigner
        }
        err = st.verifyAndApplyEditValidatorTx(stkMsg, st.evm.BlockNumber)

    case types.Delegate:
        stkMsg := &staking.Delegate{}
        if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
            return 0, err
        }
        log.Info("staking", "type", msg.Type(), "gas", gas, "txn", stkMsg)
        if msg.From() != stkMsg.DelegatorAddress {
            return 0, errInvalidSigner
        }
        err = st.verifyAndApplyDelegateTx(stkMsg)

    case types.Undelegate:
        stkMsg := &staking.Undelegate{}
        if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
            return 0, err
        }
        log.Info("staking", "type", msg.Type(), "gas", gas, "txn", stkMsg)
        if msg.From() != stkMsg.DelegatorAddress {
            return 0, errInvalidSigner
        }
        err = st.verifyAndApplyUndelegateTx(stkMsg)
    case types.CollectRewards:
        stkMsg := &staking.CollectRewards{}
        if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
            return 0, err
        }
        log.Info("staking", "type", msg.Type(), "gas", gas, "txn", stkMsg)
        if msg.From() != stkMsg.DelegatorAddress {
            return 0, errInvalidSigner
        }
        collectedRewards, err := st.verifyAndApplyCollectRewards(stkMsg)
        if err == nil {
            st.state.AddLog(&types.Log{
                Address:     stkMsg.DelegatorAddress,
                Topics:      []common.Hash{staking2.CollectRewardsTopic},
                Data:        collectedRewards.Bytes(),
                BlockNumber: st.evm.BlockNumber.Uint64(),
            })
        }
    default:
        return 0, staking.ErrInvalidStakingKind
    }
    st.refundGas()

	// Burn Txn Fees
	//txFee := new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice)
	//st.state.AddBalance(st.evm.Coinbase, txFee)

    return st.gasUsed(), err
}

func (st *StateTransition) verifyAndApplyCreateValidatorTx(
    createValidator *staking.CreateValidator, blockNum *big.Int,
) error {
    wrapper, err := VerifyAndCreateValidatorFromMsg(
		st.state, st.bc, st.evm.EpochNumber, blockNum, createValidator,
    )
    if err != nil {
        return err
    }
    if err := st.state.UpdateValidatorWrapper(wrapper.Address, wrapper); err != nil {
        return err
    }
	st.state.SetValidatorFlag(createValidator.ValidatorAddress)
	st.state.SubBalance(createValidator.ValidatorAddress, createValidator.Amount)
    return nil
}

func (st *StateTransition) verifyAndApplyEditValidatorTx(
    editValidator *staking.EditValidator, blockNum *big.Int,
) error {
    wrapper, err := VerifyAndEditValidatorFromMsg(
        st.state, st.bc, st.evm.EpochNumber, blockNum, editValidator,
    )
    if err != nil {
        return err
    }
    return st.state.UpdateValidatorWrapper(wrapper.Address, wrapper)
}

func (st *StateTransition) verifyAndApplyDelegateTx(delegate *staking.Delegate) error {
    wrapper, balanceToBeDeducted, err := VerifyAndDelegateFromMsg(st.state, delegate)
    if err != nil {
        return err
    }

    st.state.SubBalance(delegate.DelegatorAddress, balanceToBeDeducted)

    return st.state.UpdateValidatorWrapper(wrapper.Address, wrapper)
}

func (st *StateTransition) verifyAndApplyUndelegateTx(
    undelegate *staking.Undelegate,
) error {
    wrapper, err := VerifyAndUndelegateFromMsg(st.state, st.evm.EpochNumber, undelegate)
    if err != nil {
        return err
    }
    return st.state.UpdateValidatorWrapper(wrapper.Address, wrapper)
}

func (st *StateTransition) verifyAndApplyCollectRewards(collectRewards *staking.CollectRewards) (*big.Int, error) {
    if st.bc == nil {
        return network.NoReward, errors.New("[CollectRewards] No chain context provided")
    }
    // TODO(audit): make sure the delegation index is always consistent with onchain data
    delegations, err := st.bc.ReadDelegationsByDelegator(collectRewards.DelegatorAddress)
    if err != nil {
        return network.NoReward, err
    }
    updatedValidatorWrappers, totalRewards, err := VerifyAndCollectRewardsFromDelegation(
        st.state, delegations,
    )
    if err != nil {
        return network.NoReward, err
    }
    for _, wrapper := range updatedValidatorWrappers {
        if err := st.state.UpdateValidatorWrapper(wrapper.Address, wrapper); err != nil {
            return network.NoReward, err
        }
    }
    st.state.AddBalance(collectRewards.DelegatorAddress, totalRewards)
    return totalRewards, nil
}

// ATLAS
