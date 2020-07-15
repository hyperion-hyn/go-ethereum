package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/numeric"
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

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransitionWithChainContext(evm *vm.EVM, msg Message, gp *GasPool, bc ChainContext) *StateTransition {
	return &StateTransition{
		gp:       gp,
		evm:      evm,
		msg:      msg,
		gasPrice: msg.GasPrice(),
		value:    msg.Value(),
		data:     msg.Data(),
		state:    evm.StateDB,
		bc:       bc,
	}
}

// ApplyStakingMessage computes the new state for staking message
func ApplyStakingMessage(evm *vm.EVM, msg Message, gp *GasPool, bc ChainContext) (uint64, error) {
	return NewStateTransitionWithChainContext(evm, msg, gp, bc).StakingTransitionDb()
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
	case types.StopNodeStake:
		stkMsg := &staking.StopMap3Node{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyStopMap3NodeTx(stkMsg, msg.From())
	case types.ResumeNodeStake:
		stkMsg := &staking.ResumeMap3Node{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyResumeMap3NodeTx(stkMsg, msg.From())
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
	case types.DivideNodeStake:
		stkMsg := &staking.DivideMap3NodeStake{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyDivideNodeStakeTx(stkMsg, msg.From())
	case types.RenewNodeStake:
		stkMsg := &staking.RenewMap3NodeStake{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyRenewNodeStakeTx(stkMsg, msg.From())
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
	epoch, blockNum := st.evm.EpochNumber, st.evm.BlockNumber
	minTotal, minSelf, _ := staking.CalcMinTotalNodeStake(blockNum, st.evm.ChainConfig())
	wrapper, err := VerifyCreateMap3NodeMsg(st.state, epoch, blockNum, msg, signer, minSelf)
	if err != nil {
		return err
	}
	st.state.SubBalance(signer, msg.Amount)

	// save the new map3 node into pool
	nodePool := st.state.Map3NodePool()
	nodePool.GetNodes().Put(wrapper.Map3Node.NodeAddress, wrapper)
	keySet := nodePool.GetNodeKeySet()
	for _, key := range wrapper.Map3Node.NodeKeys {
		keySet.Put(key.Hex())
	}
	nodePool.GetDescriptionIdentitySet().Put(msg.Description.Identity)
	addNodeAddressToAddressSet(nodePool.GetNodeAddressSetByDelegator(), signer, wrapper.Map3Node.NodeAddress)

	wrapperSt, _ := nodePool.GetNodes().Get(wrapper.Map3Node.NodeAddress)
	if CanActivateMap3Node(wrapperSt, minTotal, minSelf) {
		if err := ActivateMap3Node(wrapperSt, st.evm.EpochNumber); err != nil {
			return err
		}
	}
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

func (st *StateTransition) verifyAndApplyStopMap3NodeTx(msg *staking.StopMap3Node, signer common.Address) error {
	if err := VerifyStopMap3NodeMsg(st.state, st.evm.EpochNumber, msg, signer); err != nil {
		return err
	}

	nodePool := st.state.Map3NodePool()
	wrapper, _ := nodePool.GetNodes().Get(msg.Map3NodeAddress)
	for _, delegator := range wrapper.GetMicrodelegations().Keys() {
		delegation, ok := wrapper.GetMicrodelegations().Get(delegator)
		if !ok {
			return ErrMicrodelegationNotExist
		}
		amt := delegation.GetPendingDelegation().GetAmount()
		st.state.AddBalance(delegator, amt)
	}

	// TODO: remove all delegation and index
	return nil
}

func (st *StateTransition) verifyAndApplyResumeMap3NodeTx(msg *staking.ResumeMap3Node, signer common.Address) error {
	minTotal, minSelf, _ := staking.CalcMinTotalNodeStake(st.evm.BlockNumber, st.evm.ChainConfig())
	if err := VerifyResumeMap3NodeMsg(st.state, msg, minSelf, signer); err != nil {
		return err
	}
	return nil
}

func (st *StateTransition) verifyAndApplyMicrodelegateTx(msg *staking.Microdelegate, signer common.Address) error {
	minTotal, minSelf, minDel := staking.CalcMinTotalNodeStake(st.evm.BlockNumber, st.evm.ChainConfig())
	if err := VerifyMicrodelegateMsg(st.state, msg, minDel, signer); err != nil {
		return err
	}

	pool := st.state.Map3NodePool()
	wrapper, _ := pool.GetNodes().Get(msg.Map3NodeAddress)
	status := wrapper.GetNodeState().GetStatus()
	unlockedEpoch := numeric.NewDecFromBigInt(st.evm.EpochNumber).Add(numeric.NewDec(staking.PendingDelegationLockPeriodInEpoch))
	if microdelegation, ok := wrapper.GetMicrodelegations().Get(msg.DelegatorAddress); ok {
		if status == staking.Active {
			microdelegation.SetAmount(big.NewInt(0).Add(microdelegation.GetAmount(), msg.Amount))
		} else { // Pending or Dividing
			pd := microdelegation.GetPendingDelegation()
			if pd == nil {
				microdelegation.SetPendingDelegation(&staking.PendingDelegation{
					Amount:        msg.Amount,
					UnlockedEpoch: unlockedEpoch,
				})
			} else {
				// TODO: weighted average
			}
		}
	} else {
		m := staking.NewMicrodelegation(
			msg.DelegatorAddress, msg.Amount,
			unlockedEpoch,
			status != staking.Active,
		)
		wrapper.GetMicrodelegations().Put(msg.DelegatorAddress, &m)
		addNodeAddressToAddressSet(pool.GetNodeAddressSetByDelegator(), signer, msg.Map3NodeAddress)
	}

	if status == staking.Active {
		wrapper.SetTotalDelegation(big.NewInt(0).Add(wrapper.GetTotalDelegation(), msg.Amount))
		// TODO: add redelegation
	} else {
		wrapper.SetTotalPendingDelegation(big.NewInt(0).Add(wrapper.GetTotalPendingDelegation(), msg.Amount))
	}
	st.state.SubBalance(signer, msg.Amount)

	if CanActivateMap3Node(wrapper, minTotal, minSelf) {
		if wrapper.GetRedelegationReference() != common.Address0 {
			// TODO: add redelegation
		}
		if err := ActivateMap3Node(wrapper, st.evm.EpochNumber); err != nil {
			return err
		}
	}
	return nil
}

func (st *StateTransition) verifyAndApplyUnmicrodelegateTx(msg *staking.Unmicrodelegate, signer common.Address) error {
	_, minSelf, _ := staking.CalcMinTotalNodeStake(st.evm.BlockNumber, st.evm.ChainConfig())
	if err := VerifyUnmicrodelegateMsg(st.state, st.evm.EpochNumber, msg, minSelf, signer); err != nil {
		return err
	}
	pool := st.state.Map3NodePool()
	wrapper, _ := pool.GetNodes().Get(msg.Map3NodeAddress)
	md, _ := wrapper.GetMicrodelegations().Get(msg.DelegatorAddress)
	amt := md.GetPendingDelegation().GetAmount()

	if amt.Cmp(msg.Amount) > 0 {
		md.GetPendingDelegation().SetAmount(big.NewInt(0).Sub(amt, msg.Amount))
	} else { // amt == msg.Amount
		md, _ := wrapper.GetMicrodelegations().Get(msg.DelegatorAddress)
		if md.GetAmount().Cmp(common.Big0) == 0 {
			wrapper.GetMicrodelegations().Remove(msg.DelegatorAddress)
			nodeAddrSet, _ := pool.GetNodeAddressSetByDelegator().Get(msg.DelegatorAddress)
			nodeAddrSet.Remove(msg.Map3NodeAddress)
		} else {
			md.SetPendingDelegation(nil)
		}
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

func (st *StateTransition) verifyAndApplyDivideNodeStakeTx(msg *staking.DivideMap3NodeStake, signer common.Address) error {
	epoch, blockNum := st.evm.EpochNumber, st.evm.BlockNumber
	minTotal, minSelf, _ := staking.CalcMinTotalNodeStake(blockNum, st.evm.ChainConfig())
	if err, err := VerifyDivideNodeStakeMsg(st.state, epoch, blockNum, msg, signer, minTotal, minSelf); err != nil {
		return err
	}

	// TODO: save new node, update delegation index, redelegation

	return nil
}

func (st *StateTransition) verifyAndApplyRenewNodeStakeTx(msg *staking.RenewMap3NodeStake, signer common.Address) error {
	if err := VerifyRenewNodeStakeMsg(st.state, st.bc, st.evm.EpochNumber, st.evm.BlockNumber, msg, signer); err != nil {
		return err
	}

	pool := st.state.Map3NodePool()
	wrapper, _ := pool.GetNodes().Get(msg.Map3NodeAddress)
	md, _ := wrapper.GetMicrodelegations().Get(msg.DelegatorAddress)
	md.SetRenewal(&staking.Renewal{
		IsRenew:      msg.IsRenew,
		UpdateHeight: st.evm.BlockNumber,
	})

	if !msg.CommissionRate.IsNil() {
		wrapper.GetMap3Node().GetCommission().GetCommissionRates().SetRate(&msg.CommissionRate)
		wrapper.GetMap3Node().GetCommission().SetUpdateHeight(st.evm.BlockNumber)
	}
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
	node.SetRedelegationReference(wrapper.Validator.ValidatorAddress)
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
	if initiatorNode.GetMap3Node().GetInitiatorAddress() == signer { // same node initiator
		wrapper.GetValidator().GetInitiatorAddresses().Put(msg.DelegatorAddress)
	}

	node, _ := nodePool.GetNodes().Get(msg.DelegatorAddress)
	node.SetRedelegationReference(wrapper.GetValidator().GetValidatorAddress())
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

	// TODO: collect reward and remove delegation if remaining amount == 0

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
	// piecewise calculation

	return reward, nil
}

func updateDescription(
	curDecs *staking.DescriptionStorage, newDesc *staking.Description, identitySet *staking.DescriptionIdentitySetStorage,
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

// TODO: state machine
func CanActivateMap3Node(wrapper *staking.Map3NodeWrapperStorage, minTotal, minSelf *big.Int) bool {
	if wrapper.GetNodeState().GetStatus() == staking.Active {
		return false
	}

	if big.NewInt(0).Add(wrapper.GetTotalPendingDelegation(), wrapper.GetTotalDelegation()).Cmp(minTotal) >= 0 {
		initiator := wrapper.GetMap3Node().GetInitiatorAddress()
		m, ok := wrapper.GetMicrodelegations().Get(initiator)
		if !ok {
			return false
		}

		if big.NewInt(0).Add(m.GetAmount(), m.GetPendingDelegation().GetAmount()).Cmp(minSelf) >= 0 {
			return true
		}
	}
	return false
}

func ActivateMap3Node(wrapper *staking.Map3NodeWrapperStorage, epoch *big.Int) error {
	// change pending delegation
	for _, delegator := range wrapper.GetMicrodelegations().Keys() {
		delegation, ok := wrapper.GetMicrodelegations().Get(delegator)
		if !ok {
			return ErrMicrodelegationNotExist
		}
		pd := delegation.GetPendingDelegation().GetAmount()
		delegation.SetAmount(big.NewInt(0).Add(delegation.GetAmount(), pd))
		delegation.SetPendingDelegation(nil)
	}
	totalPending := wrapper.GetTotalPendingDelegation()
	totalDel := wrapper.GetTotalDelegation()
	wrapper.SetTotalDelegation(big.NewInt(0).Add(totalDel, totalPending))
	wrapper.SetTotalPendingDelegation(common.Big0)

	// update state
	nodeState := wrapper.GetNodeState()
	status := nodeState.GetStatus()
	nodeState.SetStatus(staking.Active)
	nodeState.SetActivationEpoch(epoch)
	if status == staking.Dividing {
		releaseEpoch := nodeState.GetReleaseEpoch()
		timeLeft := releaseEpoch.Sub(numeric.NewDecFromBigInt(epoch))
		if timeLeft.IsNegative() {
			timeLeft = numeric.ZeroDec()
		}
		avgTime := WeightedAverageTime(numeric.NewDecFromBigInt(totalDel), timeLeft,
			numeric.NewDecFromBigInt(totalPending), staking.Map3NodeLockPeriodInEpoch)
		nodeState.SetReleaseEpoch(&avgTime)
	} else { // Pending
		time := numeric.OneDec().Mul(staking.Map3NodeLockPeriodInEpoch)
		nodeState.SetReleaseEpoch(&time)
	}
	return nil
}

func WeightedAverageTime(weight1, time1, weight2, time2 numeric.Dec) numeric.Dec {
	p1 := weight1.Mul(time1)
	p2 := weight2.Mul(time2)
	result := p1.Add(p2)
	return result.Quo(weight1.Add(weight2))
}

func addNodeAddressToAddressSet(nodeAddressSetByDelegator *staking.Map3NodeAddressSetByDelegatorStorage, delegator, nodeAddr common.Address) {
	if nodeAddrSet, ok := nodeAddressSetByDelegator.Get(delegator); ok {
		nodeAddrSet.Put(nodeAddr)
	} else {
		nodeAddressSetByDelegator.Put(delegator, &staking.AddressSet{
			nodeAddr: struct{}{},
		})
	}
}

func LookupMicrodelegationShares(wrapper *staking.Map3NodeWrapperStorage) (map[common.Address]numeric.Dec, error) {
	result := map[common.Address]numeric.Dec{}

	totalDelegationDec := numeric.NewDecFromBigInt(wrapper.GetTotalDelegation())
	if totalDelegationDec.IsZero() {
		log.Info("zero total delegation during AddReward delegation payout",
			"validator-snapshot", wrapper.GetMap3Node().GetNodeAddress().Hex())
		return result, nil
	}

	for _, key := range wrapper.GetMicrodelegations().Keys() {
		delegation, ok := wrapper.GetMicrodelegations().Get(key)
		if !ok {
			return nil, ErrMicrodelegationNotExist
		}
		percentage := numeric.NewDecFromBigInt(delegation.GetAmount()).Quo(totalDelegationDec)
		result[key] = percentage
	}
	return result, nil
}
