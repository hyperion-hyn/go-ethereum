package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/pkg/errors"
	"math/big"
)

var (
	errDupMap3NodePubKey                    = errors.New("map3 node key exists")
	errInvalidMap3NodeOperator              = errors.New("invalid map3 node operator")
	errMicrodelegationNotExist              = errors.New("microdelegation does not exist")
	errInvalidNodeStateForDelegation        = errors.New("invalid map3 node status for delegation")
	errUnmicrodelegateNotAllowed            = errors.New("invalid map3 node status to unmicrodelegate")
	errInsufficientBalanceToUnmicrodelegate = errors.New("insufficient balance to unmicrodelegate")
	errMicrodelegationStillLocked           = errors.New("microdelegation still locked")
	errTerminateMap3NodeNotAllowed          = errors.New("not allow to terminate map3 node")

	errInvalidMap3NodeStatusToRestake = errors.New("invalid map3 node to restake")
	errMap3NodeAlreadyRestaking       = errors.New("map3 node already restaked.")
	errInvalidValidatorAddress        = errors.New("validator address not equal to the address of the validator map3 already restaked to")
)

type map3VerifierForRestaking struct {
}

func (m map3VerifierForRestaking) VerifyForCreatingValidator(stateDB vm.StateDB, msg *restaking.CreateValidator, signer common.Address) (participant, error) {
	node, err := stateDB.Map3NodeByAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	if !node.IsOperator(signer) {
		return nil, errInvalidSigner
	}

	if node.Map3Node().Status().Value() != uint8(microstaking.Active) {
		return nil, errInvalidMap3NodeStatusToRestake
	}

	if node.IsAlreadyRestaking() {
		return nil, errMap3NodeAlreadyRestaking
	}
	return map3NodeAsParticipant{node: node}, nil
}

func (m map3VerifierForRestaking) VerifyForEditingValidator(stateDB vm.StateDB, msg *restaking.EditValidator, signer common.Address) (participant, error) {
	node, err := stateDB.Map3NodeByAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	if !node.IsOperator(signer) {
		return nil, errInvalidSigner
	}

	if node.Map3Node().Status().Value() != uint8(microstaking.Active) {
		return nil, errInvalidMap3NodeStatusToRestake
	}

	if node.RestakingReference().ValidatorAddress().Value() != msg.ValidatorAddress {
		return nil, errInvalidValidatorAddress
	}
	return &map3NodeAsParticipant{node: node}, nil
}

func (m map3VerifierForRestaking) VerifyForRedelegating(stateDB vm.StateDB, msg *restaking.Redelegate, signer common.Address) (participant, error) {
	node, err := stateDB.Map3NodeByAddress(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	if !node.IsOperator(signer) {
		return nil, errInvalidSigner
	}

	if node.Map3Node().Status().Value() != uint8(microstaking.Active) {
		return nil, errInvalidMap3NodeStatusToRestake
	}

	if node.IsAlreadyRestaking() {
		return nil, errMap3NodeAlreadyRestaking
	}
	return map3NodeAsParticipant{node: node}, nil
}

func (m map3VerifierForRestaking) VerifyForUnredelegating(stateDB vm.StateDB, msg *restaking.Unredelegate, signer common.Address) (participant, error) {
	node, err := stateDB.Map3NodeByAddress(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	if !node.IsOperator(signer) {
		return nil, errInvalidSigner
	}

	if node.Map3Node().Status().Value() != uint8(microstaking.Active) {
		return nil, errInvalidMap3NodeStatusToRestake
	}

	if node.RestakingReference().ValidatorAddress().Value() != msg.ValidatorAddress {
		return nil, errInvalidValidatorAddress
	}
	return &map3NodeAsParticipant{node: node}, nil
}

func (m map3VerifierForRestaking) VerifyForCollectingReward(stateDB vm.StateDB, msg *restaking.CollectReward, signer common.Address) (participant, error) {
	node, err := stateDB.Map3NodeByAddress(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	if !node.IsOperator(signer) {
		return nil, errInvalidSigner
	}

	if node.Map3Node().Status().Value() != uint8(microstaking.Active) {
		return nil, errInvalidMap3NodeStatusToRestake
	}

	if node.RestakingReference().ValidatorAddress().Value() != msg.ValidatorAddress {
		return nil, errInvalidValidatorAddress
	}
	return &map3NodeAsParticipant{node: node}, nil
}

func checkMap3DuplicatedFields(state vm.StateDB, identity string, keys microstaking.BLSPublicKeys_) error {
	map3NodePool := state.Map3NodePool()
	if identity != "" {
		identitySet := map3NodePool.DescriptionIdentitySet()
		if identitySet.Get(identity).Value() {
			return errors.Wrapf(errDupIdentity, "duplicate identity %s", identity)
		}
	}
	if len(keys.Keys) != 0 {
		nodeKeySet := map3NodePool.NodeKeySet()
		for _, key := range keys.Keys {
			if nodeKeySet.Get(key.Hex()).Value() {
				return errors.Wrapf(errDupMap3NodePubKey, "duplicate public key %x", key.Hex())
			}
		}
	}
	return nil
}

func (verifier StakingVerifier) VerifyCreateMap3NodeMsg(stateDB vm.StateDB, chainContext ChainContext, epoch, blockNum *big.Int,
	msg *microstaking.CreateMap3Node, signer common.Address) (*microstaking.Map3NodeWrapper_, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}
	if chainContext == nil {
		return nil, errChainContextMissing
	}
	if epoch == nil {
		return nil, errEpochMissing
	}
	if blockNum == nil {
		return nil, errBlockNumMissing
	}
	if msg.Amount.Sign() == -1 {
		return nil, errNegativeAmount
	}
	if msg.OperatorAddress != signer {
		return nil, errInvalidSigner
	}

	if err := checkMap3DuplicatedFields(stateDB, msg.Description.Identity, microstaking.NewBLSKeysWithBLSKey(msg.NodePubKey)); err != nil {
		return nil, err
	}
	if !CanTransfer(stateDB, msg.OperatorAddress, msg.Amount) {
		return nil, errInsufficientBalanceForStake
	}

	_, minSelf, _ := network.LatestMap3StakingRequirement(blockNum, chainContext.Config())
	if minSelf.Cmp(msg.Amount) > 0 {
		return nil, ErrSelfDelegationTooSmall
	}

	// create map3 node
	map3Address := crypto.CreateAddress(signer, stateDB.GetNonce(signer))
	node, err := microstaking.CreateMap3NodeFromNewMsg(msg, map3Address, blockNum, epoch)
	if err != nil {
		return nil, err
	}
	if err := node.Map3Node.SanityCheck(microstaking.MaxPubKeyAllowed); err != nil {
		return nil, err
	}
	return node, nil
}

func (verifier StakingVerifier) VerifyEditMap3NodeMsg(stateDB vm.StateDB, epoch, blockNum *big.Int, msg *microstaking.EditMap3Node, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if epoch == nil {
		return errEpochMissing
	}
	if blockNum == nil {
		return errBlockNumMissing
	}
	if msg.OperatorAddress != signer {
		return errInvalidSigner
	}

	blsKeys := microstaking.NewEmptyBLSKeys()
	if msg.NodeKeyToAdd != nil {
		blsKeys.Keys = append(blsKeys.Keys, msg.NodeKeyToAdd)
	}
	if err := checkMap3DuplicatedFields(stateDB, msg.Description.Identity, blsKeys); err != nil {
		return err
	}
	wrapper, err := stateDB.Map3NodeByAddress(msg.Map3NodeAddress)
	if err != nil {
		return err
	}
	if !wrapper.IsOperator(msg.OperatorAddress) {
		return errInvalidMap3NodeOperator
	}

	node, err := wrapper.Map3Node().Load()
	if err != nil {
		return err
	}

	if err := microstaking.UpdateMap3NodeFromEditMsg(node, msg); err != nil {
		return err
	}
	if err := node.SanityCheck(microstaking.MaxPubKeyAllowed); err != nil {
		return err
	}
	return nil
}

func (verifier StakingVerifier) VerifyTerminateMap3NodeMsg(stateDB vm.StateDB, epoch *big.Int, msg *microstaking.TerminateMap3Node,
	signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if epoch == nil {
		return errEpochMissing
	}
	if msg.OperatorAddress != signer {
		return errInvalidSigner
	}
	node, err := stateDB.Map3NodeByAddress(msg.Map3NodeAddress)
	if err != nil {
		return err
	}
	if !node.IsOperator(msg.OperatorAddress) {
		return errInvalidMap3NodeOperator
	}

	if node.Map3Node().Status().Value() != uint8(microstaking.Pending) {
		return errTerminateMap3NodeNotAllowed
	}

	md, ok := node.Microdelegations().Get(signer)
	if !ok {
		return errMicrodelegationNotExist
	}

	if md.PendingDelegation().UnlockedEpoch().Value().GT(common.NewDecFromBigInt(epoch)) {
		return errMicrodelegationStillLocked
	}
	return nil
}

func (verifier StakingVerifier) VerifyMicrodelegateMsg(stateDB vm.StateDB, chainContext ChainContext, blockNum *big.Int, msg *microstaking.Microdelegate,
	signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if chainContext == nil {
		return errChainContextMissing
	}
	if blockNum == nil {
		return errBlockNumMissing
	}

	if signer != msg.DelegatorAddress {
		return errInvalidSigner
	}

	if msg.Amount.Sign() == -1 {
		return errNegativeAmount
	}

	wrapper, err := stateDB.Map3NodeByAddress(msg.Map3NodeAddress)
	if err != nil {
		return err
	}

	status := wrapper.Map3Node().Status().Value()
	if status != uint8(microstaking.Pending) {
		return errInvalidNodeStateForDelegation
	}

	// Check if there is enough liquid token to delegate
	if !CanTransfer(stateDB, msg.DelegatorAddress, msg.Amount) {
		return errInsufficientBalanceForStake
	}

	_, _, minDel := network.LatestMap3StakingRequirement(blockNum, chainContext.Config())
	if minDel.Cmp(msg.Amount) > 0 {
		return errDelegationTooSmall
	}
	return nil
}

func (verifier StakingVerifier) VerifyUnmicrodelegateMsg(stateDB vm.StateDB, chainContext ChainContext, blockNum *big.Int, epoch *big.Int,
	msg *microstaking.Unmicrodelegate, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if chainContext == nil {
		return errChainContextMissing
	}
	if blockNum == nil {
		return errBlockNumMissing
	}
	if epoch == nil {
		return errEpochMissing
	}
	if msg.Amount.Sign() == -1 {
		return errNegativeAmount
	}
	if msg.DelegatorAddress != signer {
		return errInvalidSigner
	}

	wrapper, err := stateDB.Map3NodeByAddress(msg.Map3NodeAddress)
	if err != nil {
		return err
	}

	// TODO(ATLAS): only pending status
	status := wrapper.Map3Node().Status().Value()
	if status != uint8(microstaking.Pending) {
		return errUnmicrodelegateNotAllowed
	}

	md, ok := wrapper.Microdelegations().Get(msg.DelegatorAddress)
	if !ok {
		return errMicrodelegationNotExist
	}

	p := md.PendingDelegation()
	if p.Amount().Value().Cmp(msg.Amount) < 0 {
		return errInsufficientBalanceToUnmicrodelegate
	}

	if p.UnlockedEpoch().Value().GT(common.NewDecFromBigInt(epoch)) {
		return errMicrodelegationStillLocked
	}

	if wrapper.IsOperator(msg.DelegatorAddress) {
		amt := big.NewInt(0).Sub(p.Amount().Value(), msg.Amount)
		total := amt.Add(amt, md.Amount().Value())

		_, minSelf, _ := network.LatestMap3StakingRequirement(blockNum, chainContext.Config())
		if minSelf.Cmp(total) > 0 {
			return ErrSelfDelegationTooSmall
		}
	}
	return nil
}

func (verifier StakingVerifier) VerifyCollectMicrodelRewardsMsg(stateDB vm.StateDB, msg *microstaking.CollectRewards, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if msg.DelegatorAddress != signer {
		return errInvalidSigner
	}
	map3NodePool := stateDB.Map3NodePool()
	delegationIndexMap := map3NodePool.DelegationIndexMapByDelegator().Get(msg.DelegatorAddress)
	if delegationIndexMap.Keys().Length() == 0 {
		return errNoRewardsToCollect
	}

	totalReward := big.NewInt(0)
	for i := 0; i < delegationIndexMap.Keys().Length(); i++ {
		nodeAddr := delegationIndexMap.Keys().Get(i).Value()
		node, err := stateDB.Map3NodeByAddress(nodeAddr)
		if err != nil {
			return err
		}
		if micro, ok := node.Microdelegations().Get(signer); ok {
			if micro.Reward().Value().Cmp(common.Big0) > 0 {
				totalReward.Add(totalReward, micro.Reward().Value())
			}
		} else {
			return errMicrodelegationNotExist
		}
	}

	if totalReward.Int64() == 0 {
		return errNoRewardsToCollect
	}
	return nil
}
