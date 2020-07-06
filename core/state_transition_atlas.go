package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/effective"
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
	case types.SplitNode:
		stkMsg := &staking.SplitNode{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplySplitNodeTx(stkMsg, msg.From())
	case types.StakeCreateVal:
		stkMsg := &staking.CreateValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyCreateValidatorTx(stkMsg, msg.From())
	case types.StakeEditVal:
		stkMsg := &staking.EditValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyEditValidatorTx(stkMsg, msg.From())
	case types.Redelegate:
		stkMsg := &staking.Redelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyRedelegateTx(stkMsg, msg.From())
	case types.Unredelegate:
		stkMsg := &staking.Unredelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyUnredelegateTx(stkMsg, msg.From())
	case types.CollectRedelRewards:
		stkMsg := &staking.CollectRedelegationRewards{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		collectedRewards, err := st.verifyAndApplyCollectRedelRewards(stkMsg, msg.From())
		if err != nil {
			return 0, err
		}
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

func (st *StateTransition) verifyAndApplyCreateMap3NodeTx(msg *staking.CreateMap3Node, signer common.Address) error {
	minTotal, minSelf, _ := staking.CalcMinTotalNodeStake(st.evm.BlockNumber, st.evm.ChainConfig())
	wrapper, err := VerifyCreateMap3NodeMsg(
		st.state, st.evm.EpochNumber, st.evm.BlockNumber, msg, signer, minSelf,
	)
	if err != nil {
		return err
	}
	nodePool := st.state.Map3NodePool()
	nodePool.GetNodes().Put(wrapper.Map3Node.NodeAddress, wrapper)
	keySet := nodePool.GetNodeKeySet()
	for _, key := range wrapper.Map3Node.NodeKeys {
		keySet.Put(key.Hex())
	}
	nodePool.GetDescriptionIdentitySet().Put(msg.Description.Identity)
	if nodeAddrSet, ok := nodePool.GetNodeAddressSetByDelegator().Get(signer); ok {
		nodeAddrSet.Put(wrapper.Map3Node.NodeAddress)
	} else {
		nodePool.GetNodeAddressSetByDelegator().Put(signer, &staking.AddressSet{
			wrapper.Map3Node.NodeAddress: struct{}{},
		})
	}
	st.state.SubBalance(signer, msg.Amount)

	// TODO(ATLAS): check for activating the new node
	// check minTotal <= curren total && 20%

	return nil
}

func (st *StateTransition) verifyAndApplyEditMap3NodeTx(msg *staking.EditMap3Node, signer common.Address) error {
	if err := VerifyEditMap3NodeMsg(st.state, st.evm.EpochNumber, st.evm.BlockNumber, msg, signer); err != nil {
		return err
	}

	pool := st.state.Map3NodePool()
	wrapper, _ := pool.GetNodes().Get(msg.Map3NodeAddress)
	node := wrapper.GetMap3Node()

	updateDescription(node.GetDescription(), msg.Description, pool.GetDescriptionIdentitySet())

	if !msg.CommissionRate.IsNil() {
		node.GetCommission().GetCommissionRates().SetRate(&msg.CommissionRate)
		node.GetCommission().SetUpdateHeight(st.evm.BlockNumber)
	}

	if msg.NodeKeyToRemove != nil {
		for i := 0; i < node.GetNodeKeys().Len(); i++ {
			if msg.NodeKeyToRemove.Hex() == node.GetNodeKeys().Get(i).Hex() {
				node.GetNodeKeys().Remove(i, false)
				pool.GetNodeKeySet().Remove(msg.NodeKeyToRemove.Hex())
				break
			}
		}
	}

	if msg.NodeKeyToAdd != nil {
		node.GetNodeKeys().Push(msg.NodeKeyToAdd)
		pool.GetNodeKeySet().Put(msg.NodeKeyToAdd.Hex())
	}
	return nil
}

func (st *StateTransition) verifyAndApplyMicrodelegateTx(msg *staking.Microdelegate, signer common.Address) error {
	minTotal, _, minDel := staking.CalcMinTotalNodeStake(st.evm.BlockNumber, st.evm.ChainConfig())
	if err := VerifyMicrodelegateMsg(st.state, msg, minDel, signer); err != nil {
		return err
	}

	pool := st.state.Map3NodePool()
	wrapper, _ := pool.GetNodes().Get(msg.Map3NodeAddress)

	status := wrapper.GetNodeState().GetStatus()
	if status == staking.Active {
		// TODO(ATLAS): collect reward from validator as initiator
	}

	if microdelegation, ok := wrapper.GetMicrodelegations().Get(msg.DelegatorAddress); ok {
		if status == staking.Pending {
			pd := microdelegation.GetPendingDelegation()
			if pd == nil {
				microdelegation.SetPendingDelegation(&staking.PendingDelegation{
					Amount:        msg.Amount,
					UnlockedEpoch: big.NewInt(0).Add(st.evm.EpochNumber, big.NewInt(staking.PendingDelegationLockPeriodInEpoch)),
				})
			} else {
				// TODO: weighted average
			}
		} else {	// Active
			microdelegation.SetAmount(big.NewInt(0).Add(microdelegation.GetAmount(), msg.Amount))
		}
		microdelegation.SetAutoRenew(msg.AutoRenew)
	} else {
		m := staking.NewMicrodelegation(
			msg.DelegatorAddress, msg.Amount, st.evm.EpochNumber,
			msg.AutoRenew, status == staking.Pending,
		)
		wrapper.GetMicrodelegations().Put(msg.DelegatorAddress, &m)
	}

	if status == staking.Active {
		wrapper.SetTotalPendingDelegation(big.NewInt(0).Add(wrapper.GetTotalPendingDelegation(), msg.Amount))
	} else {
		wrapper.SetTotalDelegation(big.NewInt(0).Add(wrapper.GetTotalDelegation(), msg.Amount))
	}

	st.state.SubBalance(signer, msg.Amount)

	// TODO(ATLAS): check for activating the new node
	// check minTotal <= curren total && 20%

	return nil
}

func (st *StateTransition) verifyAndApplyUnmicrodelegateTx(msg *staking.Unmicrodelegate, signer common.Address) error {
	if err := VerifyUnmicrodelegateMsg(st.state, st.evm.EpochNumber, msg, signer); err != nil {
		return err
	}
	pool := st.state.Map3NodePool()
	wrapper, _ := pool.GetNodes().Get(msg.Map3NodeAddress)
	md, _ := wrapper.GetMicrodelegations().Get(msg.DelegatorAddress)
	amt := md.GetPendingDelegation().GetAmount()
	if amt.Cmp(msg.Amount) > 0 {
		md.GetPendingDelegation().SetAmount(big.NewInt(0).Sub(amt, msg.Amount))
	} else {	// amt == msg.Amount
		wrapper.GetMicrodelegations().Remove(msg.DelegatorAddress)
		nodeAddrSet, _ := pool.GetNodeAddressSetByDelegator().Get(msg.DelegatorAddress)
		nodeAddrSet.Remove(msg.Map3NodeAddress)
	}
	wrapper.SetTotalPendingDelegation(big.NewInt(0).Sub(wrapper.GetTotalPendingDelegation(), msg.Amount))
	st.state.AddBalance(signer, msg.Amount)
	return nil
}

func (st *StateTransition) verifyAndApplyCollectMicrodelRewardsTx(
	msg *staking.CollectMicrodelegationRewards, signer common.Address,
) (*big.Int, error) {
	if err := VerifyCollectMicrodelRewardsDelegation(st.state, msg, signer); err != nil {
		return network.NoReward, err
	}

	map3NodePool := st.state.Map3NodePool()
	nodeAddressSet, _ := map3NodePool.GetNodeAddressSetByDelegator().Get(signer)
	totalRewards := common.Big0
	for _, nodeAddr := range nodeAddressSet.Keys() {
		node, _ := map3NodePool.GetNodes().Get(nodeAddr)
		micro, _ := node.GetMicrodelegations().Get(signer)
		if micro.GetReward().Cmp(common.Big0) > 0 {
			totalRewards.Add(totalRewards, micro.GetReward())
			micro.SetReward(common.Big0)
		}
	}
	return totalRewards, nil
}

func (st *StateTransition) verifyAndApplySplitNodeTx(msg *staking.SplitNode, signer common.Address) error {
	if err := VerifySplitNodeMsg(st.state, st.evm.BlockNumber, msg, signer); err != nil {
		return err
	}

	// TODO: split

	return nil
}

func (st *StateTransition) verifyAndApplyCreateValidatorTx(msg *staking.CreateValidator, signer common.Address) error {
	wrapper, err := VerifyCreateValidatorMsg(st.state, st.evm.BlockNumber, msg, signer)
	if err != nil {
		return err
	}
	validatorPool := st.state.ValidatorPool()
	validatorPool.GetValidators().Put(wrapper.Validator.ValidatorAddress, wrapper)

	keySet := validatorPool.GetSlotKeySet()
	for _, key := range wrapper.Validator.SlotPubKeys {
		keySet.Put(key.Hex())
	}
	validatorPool.GetDescriptionIdentitySet().Put(msg.Description.Identity)

	node, _ := st.state.Map3NodePool().GetNodes().Get(msg.InitiatorAddress)
	node.SetRedelegationReference(&staking.RedelegationReference{
		ValidatorAddress: wrapper.Validator.ValidatorAddress,
	})
	return nil
}

func (st *StateTransition) verifyAndApplyEditValidatorTx(msg *staking.EditValidator, signer common.Address) error {
	if err := VerifyEditValidatorMsg(st.state, st.evm.EpochNumber, st.evm.BlockNumber, msg, signer); err != nil {
		return err
	}

	pool := st.state.ValidatorPool()
	wrapper, _ := pool.GetValidators().Get(msg.ValidatorAddress)
	validator := wrapper.GetValidator()

	updateDescription(validator.GetDescription(), msg.Description, pool.GetDescriptionIdentitySet())

	if !msg.CommissionRate.IsNil() {
		validator.GetCommission().GetCommissionRates().SetRate(&msg.CommissionRate)
		validator.GetCommission().SetUpdateHeight(st.evm.BlockNumber)
	}

	if msg.SlotKeyToRemove != nil {
		for i := 0; i < validator.GetSlotPubKeys().Len(); i++ {
			if msg.SlotKeyToRemove == validator.GetSlotPubKeys().Get(i) {
				validator.GetSlotPubKeys().Remove(i, false)
				pool.GetSlotKeySet().Remove(msg.SlotKeyToRemove.Hex())
				break
			}
		}
	}

	if msg.SlotKeyToAdd != nil {
		validator.GetSlotPubKeys().Push(msg.SlotKeyToAdd)
		pool.GetSlotKeySet().Put(msg.SlotKeyToAdd.Hex())
	}

	if msg.EPOSStatus == effective.Active || msg.EPOSStatus == effective.Inactive {
		validator.SetStatus(msg.EPOSStatus)
	}
	return nil
}

func (st *StateTransition) verifyAndApplyRedelegateTx(msg *staking.Redelegate, signer common.Address) error {
	err := VerifyRedelegateMsg(st.state, msg, signer)
	if err != nil {
		return err
	}

	validatorPool := st.state.ValidatorPool()
	wrapper, _ := validatorPool.GetValidators().Get(msg.ValidatorAddress)

	if redelegation, ok := wrapper.GetRedelegations().Get(msg.DelegatorAddress); ok {
		amt := redelegation.GetAmount()
		redelegation.SetAmount(big.NewInt(0).Add(amt, msg.Amount))
	} else {
		m := staking.NewRedelegation(msg.DelegatorAddress, msg.Amount)
		wrapper.GetRedelegations().Put(msg.DelegatorAddress, &m)
	}
	wrapper.SetTotalDelegation(big.NewInt(0).Add(wrapper.GetTotalDelegation(), msg.Amount))

	nodePool := st.state.Map3NodePool()
	initiatorNode, _ := nodePool.GetNodes().Get(wrapper.GetValidator().GetInitiatorAddresses().Keys()[0])
	if initiatorNode.GetMap3Node().GetInitiatorAddress() == signer {	// same node initiator
		wrapper.GetValidator().GetInitiatorAddresses().Put(msg.DelegatorAddress)
	}

	node, _ := nodePool.GetNodes().Get(msg.DelegatorAddress)
	node.SetRedelegationReference(&staking.RedelegationReference{
		ValidatorAddress: wrapper.GetValidator().GetValidatorAddress(),
	})
	return nil
}

func (st *StateTransition) verifyAndApplyUnredelegateTx(msg *staking.Unredelegate, signer common.Address) error {
	if err := VerifyUnredelegateMsg(st.state, st.evm.EpochNumber, msg, signer); err != nil {
		return err
	}

	validator, _ := st.state.ValidatorPool().GetValidators().Get(msg.ValidatorAddress)
	redelegation, _ := validator.GetRedelegations().Get(msg.DelegatorAddress)
	amt := redelegation.GetAmount()
	redelegation.SetUndelegation(&staking.Undelegation{
		Amount: amt,
		Epoch:  st.evm.EpochNumber,
	})
	redelegation.SetAmount(common.Big0)
	validator.SetTotalDelegation(big.NewInt(0).Sub(validator.GetTotalDelegation(), amt))
	return nil
}

func (st *StateTransition) verifyAndApplyCollectRedelRewards(msg *staking.CollectRedelegationRewards, signer common.Address) (*big.Int, error) {
	if err := VerifyCollectRedelRewardsMsg(st.state, msg, signer); err != nil {
		return network.NoReward, err
	}

	validator, _ := st.state.ValidatorPool().GetValidators().Get(msg.ValidatorAddress)
	redelegation, _ := validator.GetRedelegations().Get(msg.DelegatorAddress)
	reward := redelegation.GetReward()

	node, _ := st.state.Map3NodePool().GetNodes().Get(msg.DelegatorAddress)
	// TODO: lookup share of microdelegation

	return reward, nil
}


func updateDescription(
	curDecs *staking.DescriptionStorage, newDesc *staking.Description,identitySet *staking.DescriptionIdentitySetStorage,
) {
	if newDesc.Name != "" {
		curDecs.SetName(newDesc.Name)
	}
	if newDesc.Identity != "" {
		identitySet.Remove(curDecs.GetIdentity())
		curDecs.SetIdentity(newDesc.Identity)
		identitySet.Put(newDesc.Identity)
	}
	if newDesc.Website != "" {
		curDecs.SetWebsite(newDesc.Website)
	}
	if newDesc.SecurityContact != "" {
		curDecs.SetSecurityContact(newDesc.SecurityContact)
	}
	if newDesc.Details != "" {
		curDecs.SetDetails(newDesc.Details)
	}
}