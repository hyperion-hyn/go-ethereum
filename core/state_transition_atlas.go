package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/network"
	staking "github.com/ethereum/go-ethereum/staking/types"
	"github.com/pkg/errors"
	"math/big"
)

var (
	ErrInvalidStakingKind = errors.New("bad staking kind")
)

// ApplyStakingMessage computes the new state for staking message
func ApplyStakingMessage(evm *vm.EVM, msg Message, gp *GasPool) (uint64, error) {
	return NewStateTransition(evm, msg, gp).StakingTransitionDb()
}

// StakingTransitionDb will transition the state by applying the staking message and
// returning the result including the used gas. It returns an error if failed.
// It is used for staking transaction only
func (st *StateTransition) StakingTransitionDb() (usedGas uint64, err error) {
	if err = st.preCheck(); err != nil {
		return 0, err
	}
	msg := st.msg

	sender := vm.AccountRef(msg.From())
	homestead := st.evm.ChainConfig().IsHomestead(st.evm.BlockNumber)
	istanbul := st.evm.ChainConfig().IsIstanbul(st.evm.BlockNumber)

	// Pay intrinsic gas
	// TODO(ATLAS): gas?
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
	case types.StakeCreateNode:
		stkMsg := &staking.CreateMap3Node{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyCreateMap3NodeTx(stkMsg, msg.From())
	case types.StakeEditNode:
		stkMsg := &staking.EditMap3Node{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyEditMap3NodeTx(stkMsg, msg.From())
	case types.Microdelegate:
		stkMsg := &staking.Microdelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyMicrodelegateTx(stkMsg, msg.From())
	case types.Unmicrodelegate:
		stkMsg := &staking.Unmicrodelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyUnmicrodelegateTx(stkMsg, msg.From())
	case types.CollectMicrodelRewards:
		stkMsg := &staking.CollectMicrodelegationRewards{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		_, err := st.verifyAndApplyCollectMicrodelRewardsTx(stkMsg, msg.From())
		if err != nil {
			return 0, err
		}
		// TODO: Add log for reward ?
	case types.StakeCreateVal:
		stkMsg := &staking.CreateValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		if msg.From() != stkMsg.InitiatorAddress {
			return 0, errInvalidSigner
		}
		err = st.verifyAndApplyCreateValidatorTx(stkMsg, st.evm.BlockNumber)
	case types.StakeEditVal:
		stkMsg := &staking.EditValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		if msg.From() != stkMsg. {
			return 0, errInvalidSigner
		}
		err = st.verifyAndApplyEditValidatorTx(stkMsg, st.evm.BlockNumber)
	case types.Redelegate:
		stkMsg := &staking.Redelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		//utils.Logger().Info().Msgf("[DEBUG STAKING] staking type: %s, gas: %d, txn: %+v", msg.Type(), gas, stkMsg)
		//if msg.From() != stkMsg.MicrodelegatorAddress {
		//	return 0, errInvalidSigner
		//}
		err = st.verifyAndApplyDelegateTx(stkMsg)
	case types.Unredelegate:
		stkMsg := &staking.Unredelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		//utils.Logger().Info().Msgf("[DEBUG STAKING] staking type: %s, gas: %d, txn: %+v", msg.Type(), gas, stkMsg)
		//if msg.From() != stkMsg.MicrodelegatorAddress {
		//	return 0, errInvalidSigner
		//}
		err = st.verifyAndApplyUndelegateTx(stkMsg)
	case types.CollectRedelRewards:
		stkMsg := &staking.CollectRedelegationRewards{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		//utils.Logger().Info().Msgf("[DEBUG STAKING] staking type: %s, gas: %d, txn: %+v", msg.Type(), gas, stkMsg)
		//if msg.From() != stkMsg.MicrodelegatorAddress {
		//	return 0, errInvalidSigner
		//}
		collectedRewards, tempErr := st.verifyAndApplyCollectRewards(stkMsg)
		err = tempErr
		// TODO: Add log for reward ?
	default:
		return 0, ErrInvalidStakingKind
	}
	st.refundGas()

	// TODO(ATLAS) Txn Fees
	//txFee := new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice)
	//st.state.AddBalance(st.evm.Coinbase, txFee)

	return st.gasUsed(), err
}

func (st *StateTransition) verifyAndApplyCreateMap3NodeTx(createMap3Node *staking.CreateMap3Node, signer common.Address) error {
	minTotal, minSelf, _ := staking.CalcMinTotalNodeStake(st.evm.BlockNumber, st.evm.ChainConfig())
	wrapper, err := VerifyAndCreateMap3NodeFromMsg(
		st.state, st.evm.EpochNumber, st.evm.BlockNumber, createMap3Node, signer, minSelf,
	)
	if err != nil {
		return err
	}
	st.state.Map3NodePool().GetNodes().Put(wrapper.Map3Node.NodeAddress, wrapper)
	st.state.SubBalance(signer, createMap3Node.Amount)

	// TODO(ATLAS): check for activating the new node
	// check minTotal <= curren total && 20%

	return nil
}

func (st *StateTransition) verifyAndApplyEditMap3NodeTx(editMap3Node *staking.EditMap3Node, signer common.Address) error {
	err := VerifyAndEditMap3NodeFromMsg(
		st.state, st.evm.EpochNumber, st.evm.BlockNumber, editMap3Node, signer,
	)
	return err
}

func (st *StateTransition) verifyAndApplyMicrodelegateTx(microdelegate *staking.Microdelegate, signer common.Address) error {
	minTotal, _, minDel := staking.CalcMinTotalNodeStake(st.evm.BlockNumber, st.evm.ChainConfig())
	balanceToBeDeducted, err := VerifyAndMicrodelegateFromMsg(st.state, st.evm.EpochNumber, microdelegate, minDel, signer)
	if err != nil {
		return err
	}
	st.state.SubBalance(signer, balanceToBeDeducted)

	// TODO(ATLAS): check for activating the new node
	// check minTotal <= curren total && 20%

	return nil
}

func (st *StateTransition) verifyAndApplyUnmicrodelegateTx(unmicrodelegate *staking.Unmicrodelegate, signer common.Address) error {
	balanceToBeAdded, err := VerifyAndUnmicrodelegateFromMsg(st.state, st.evm.EpochNumber, unmicrodelegate, signer)
	if err != nil {
		return err
	}
	st.state.AddBalance(signer, balanceToBeAdded)
	return nil
}

func (st *StateTransition) verifyAndApplyCollectMicrodelRewardsTx(
	collectRewards *staking.CollectMicrodelegationRewards, signer common.Address,
) (*big.Int, error) {
	totalRewards, err := VerifyAndCollectMicrodelRewardsFromDelegation(st.state, collectRewards, signer)
	if err != nil {
		return network.NoReward, err
	}
	return totalRewards, nil
}

func (st *StateTransition) verifyAndApplyCreateValidatorTx(
	createValidator *staking.CreateValidator, blockNum *big.Int,
) error {
	wrapper, err := VerifyAndCreateValidatorFromMsg(
		st.state, st.evm.EpochNumber, blockNum, createValidator,
	)
	if err != nil {
		return err
	}
	if err := st.state.UpdateValidator(wrapper); err != nil {
		return err
	}
	// TODO(ATLAS): update node state
	return nil
}

func (st *StateTransition) verifyAndApplyEditValidatorTx(
	editValidator *EditValidator, blockNum *big.Int,
) error {
	wrapper, err := VerifyAndEditValidatorFromMsg(
		st.state, st.bc, st.evm.EpochNumber, blockNum, editValidator,
	)
	if err != nil {
		return err
	}
	return st.state.UpdateValidatorWrapper(wrapper.Address, wrapper)
}

func (st *StateTransition) verifyAndApplyDelegateTx(delegate *Redelegate) error {
	wrapper, balanceToBeDeducted, err := VerifyAndDelegateFromMsg(st.state, delegate)
	if err != nil {
		return err
	}

	st.state.SubBalance(delegate.DelegatorAddress, balanceToBeDeducted)

	return st.state.UpdateValidatorWrapper(wrapper.Address, wrapper)
}

func (st *StateTransition) verifyAndApplyUndelegateTx(
	undelegate *Unredelegate,
) error {
	wrapper, err := VerifyAndUndelegateFromMsg(st.state, st.evm.EpochNumber, undelegate)
	if err != nil {
		return err
	}
	return st.state.UpdateValidatorWrapper(wrapper.Address, wrapper)
}

func (st *StateTransition) verifyAndApplyCollectRewards(collectRewards *staking.CollectRedelegationRewards) (*big.Int, error) {
	if st.bc == nil {
		return network.NoReward, errors.New("[CollectRedelegationRewards] No chain context provided")
	}
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
