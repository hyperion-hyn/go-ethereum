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
	errMap3NodeNotExist                     = errors.New("map3 node does not exist")
	errDupMap3NodeIdentity                  = errors.New("map3 node identity exists")
	errDupMap3NodePubKey                    = errors.New("map3 node key exists")
	errInvalidMap3NodeOperator              = errors.New("invalid map3 node operator")
	errMicrodelegationNotExist              = errors.New("microdelegation does not exist")
	errInvalidNodeStatusForDelegation       = errors.New("invalid map3 node status for delegation")
	errUnmicrodelegateNotAllowed            = errors.New("invalid map3 node status to unmicrodelegate")
	errInsufficientBalanceToUnmicrodelegate = errors.New("insufficient balance to unmicrodelegate")
	errMicrodelegationStillLocked           = errors.New("microdelegation still locked")
	errTerminateMap3NodeNotAllowed          = errors.New("not allow to terminate map3 node")
	errSelfDelegationTooSmall               = errors.New("self delegation amount too small")
	errEditTerminatedMap3NodeNotAllowed     = errors.New("not allow to edit terminated map3 node")
	errMap3NodeRenewalNotAllowed            = errors.New("not allow to renew map3 node")
	errChangeRenewalDecisionNotAllowed      = errors.New("not allow to change renewal decision")
	errCommissionUpdateNotAllow             = errors.New("not allow to update commission by non-operator")
	errMap3NodeNotRenewalAnyMore            = errors.New("map3 node not renewal any more")

	errInvalidMap3NodeStatusToRestake = errors.New("invalid map3 node status to restake")
	errMap3NodeAlreadyRestaking       = errors.New("map3 node already restaked")
	errInvalidValidatorAddress        = errors.New("validator address not equal to the address of the validator map3 already restaked to")
)

type map3VerifierForRestaking struct {
	chainContext ChainContext
}

func (m map3VerifierForRestaking) VerifyForCreatingValidator(stateDB vm.StateDB, msg *restaking.CreateValidator, signer common.Address) (participant, error) {
	node, err := stateDB.Map3NodeByAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	if !node.IsOperator(signer) {
		return nil, errInvalidSigner
	}

	if !node.Map3Node().AtStatus(microstaking.Active) {
		return nil, errInvalidMap3NodeStatusToRestake
	}

	if node.IsRestaking() {
		return nil, errMap3NodeAlreadyRestaking
	}
	return map3NodeAsParticipant{
		stateDB: stateDB,
		chain:   m.chainContext,
		node:    node,
	}, nil
}

func (m map3VerifierForRestaking) VerifyForEditingValidator(stateDB vm.StateDB, msg *restaking.EditValidator, signer common.Address) (participant, error) {
	node, err := stateDB.Map3NodeByAddress(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	if !node.IsOperator(signer) {
		return nil, errInvalidSigner
	}

	if !node.Map3Node().AtStatus(microstaking.Active) {
		return nil, errInvalidMap3NodeStatusToRestake
	}

	if node.RestakingReference().ValidatorAddress().Value() != msg.ValidatorAddress {
		return nil, errInvalidValidatorAddress
	}
	return map3NodeAsParticipant{
		stateDB: stateDB,
		chain:   m.chainContext,
		node:    node,
	}, nil
}

func (m map3VerifierForRestaking) VerifyForRedelegating(stateDB vm.StateDB, msg *restaking.Redelegate, signer common.Address) (participant, error) {
	node, err := stateDB.Map3NodeByAddress(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	if !node.IsOperator(signer) {
		return nil, errInvalidSigner
	}

	if !node.Map3Node().AtStatus(microstaking.Active) {
		return nil, errInvalidMap3NodeStatusToRestake
	}

	if node.IsRestaking() {
		return nil, errMap3NodeAlreadyRestaking
	}
	return map3NodeAsParticipant{
		stateDB: stateDB,
		chain:   m.chainContext,
		node:    node,
	}, nil
}

func (m map3VerifierForRestaking) VerifyForUnredelegating(stateDB vm.StateDB, msg *restaking.Unredelegate, signer common.Address) (participant, error) {
	node, err := stateDB.Map3NodeByAddress(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	if !node.IsOperator(signer) {
		return nil, errInvalidSigner
	}

	if !node.Map3Node().AtStatus(microstaking.Active) {
		return nil, errInvalidMap3NodeStatusToRestake
	}

	if node.RestakingReference().ValidatorAddress().Value() != msg.ValidatorAddress {
		return nil, errInvalidValidatorAddress
	}
	return map3NodeAsParticipant{
		stateDB: stateDB,
		chain:   m.chainContext,
		node:    node,
	}, nil
}

func (m map3VerifierForRestaking) VerifyForCollectingReward(stateDB vm.StateDB, msg *restaking.CollectReward, signer common.Address) (participant, error) {
	node, err := stateDB.Map3NodeByAddress(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	if !node.IsOperator(signer) {
		return nil, errInvalidSigner
	}

	if !node.Map3Node().AtStatus(microstaking.Active) {
		return nil, errInvalidMap3NodeStatusToRestake
	}

	if node.RestakingReference().ValidatorAddress().Value() != msg.ValidatorAddress {
		return nil, errInvalidValidatorAddress
	}
	return map3NodeAsParticipant{
		stateDB: stateDB,
		chain:   m.chainContext,
		node:    node,
	}, nil
}

func checkMap3DuplicatedFields(state vm.StateDB, identity string, keys microstaking.BLSPublicKeys_) error {
	map3NodePool := state.Map3NodePool()
	if identity != "" {
		identitySet := map3NodePool.DescriptionIdentitySet()
		if identitySet.Get(identity).Value() {
			return errors.Wrapf(errDupMap3NodeIdentity, "duplicate identity %s", identity)
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

	requireTotal, requireSelf, _ := network.LatestMicrostakingRequirement(blockNum, chainContext.Config())
	if requireSelf.Cmp(msg.Amount) > 0 {
		return nil, errSelfDelegationTooSmall
	}
	var percent *common.Dec
	if !chainContext.Config().Atlas.IsMicrostakingImprove(blockNum) {
		p := common.NewDecFromInt(msg.Amount).QuoInt(requireTotal)
		percent = &p
	}

	// create map3 node
	map3Address := crypto.CreateAddress(signer, stateDB.GetNonce(signer))
	node, err := microstaking.CreateMap3NodeFromNewMsg(msg, map3Address, blockNum, epoch)
	if err != nil {
		return nil, err
	}
	if err := node.Map3Node.SanityCheck(microstaking.MaxPubKeyAllowed, percent); err != nil {
		return nil, err
	}
	return node, nil
}

func (verifier StakingVerifier) VerifyEditMap3NodeMsg(stateDB vm.StateDB, chainContext ChainContext, epoch, blockNum *big.Int,
	msg *microstaking.EditMap3Node, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if chainContext == nil {
		return errChainContextMissing
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

	if node.Status == uint8(microstaking.Terminated) {
		return errEditTerminatedMap3NodeNotAllowed
	}

	if err := microstaking.UpdateMap3NodeFromEditMsg(node, msg); err != nil {
		return err
	}
	if err := node.SanityCheck(microstaking.MaxPubKeyAllowed, nil); err != nil {
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

	if !node.Map3Node().AtStatus(microstaking.Pending) {
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

	node, err := stateDB.Map3NodeByAddress(msg.Map3NodeAddress)
	if err != nil {
		return err
	}

	if !node.Map3Node().AtStatus(microstaking.Pending) {
		return errInvalidNodeStatusForDelegation
	}

	// Check if there is enough liquid token to delegate
	if !CanTransfer(stateDB, msg.DelegatorAddress, msg.Amount) {
		return errInsufficientBalanceForStake
	}

	_, _, requireDel := network.LatestMicrostakingRequirement(blockNum, chainContext.Config())
	if requireDel.Cmp(msg.Amount) > 0 {
		return errDelegationTooSmall
	}
	return nil
}

func (verifier StakingVerifier) VerifyUnmicrodelegateMsg(stateDB vm.StateDB, chainContext ChainContext, blockNum, epoch *big.Int,
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

	node, err := stateDB.Map3NodeByAddress(msg.Map3NodeAddress)
	if err != nil {
		return err
	}

	// TODO(ATLAS): only pending status
	if !node.Map3Node().AtStatus(microstaking.Pending) {
		return errUnmicrodelegateNotAllowed
	}

	md, ok := node.Microdelegations().Get(msg.DelegatorAddress)
	if !ok {
		return errMicrodelegationNotExist
	}

	p := md.PendingDelegation()
	if p.Amount().Value().Cmp(msg.Amount) < 0 {
		return errInsufficientBalanceToUnmicrodelegate
	}

	if p.UnlockedEpoch().Value().GTE(common.NewDecFromBigInt(epoch)) {
		return errMicrodelegationStillLocked
	}

	if node.IsOperator(msg.DelegatorAddress) {
		amt := big.NewInt(0).Sub(p.Amount().Value(), msg.Amount)
		self := amt.Add(amt, md.Amount().Value())

		requireTotal, requireSelf, _ := network.LatestMicrostakingRequirement(blockNum, chainContext.Config())
		if chainContext.Config().Atlas.IsMicrostakingImprove(blockNum) {
			if requireSelf.Cmp(self) > 0 {
				return errSelfDelegationTooSmall
			}
		} else {
			percent := common.NewDecFromInt(self).QuoInt(requireTotal)
			commissionRate := node.Map3Node().Commission().Rate().Value()
			if requireSelf.Cmp(self) > 0 || percent.LT(commissionRate) {
				return errSelfDelegationTooSmall
			}
		}
	}
	return nil
}

func (verifier StakingVerifier) VerifyCollectMicrostakingRewardsMsg(stateDB vm.StateDB, msg *microstaking.CollectRewards, signer common.Address) error {
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
			if micro.Reward().Value().Sign() > 0 {
				totalReward.Add(totalReward, micro.Reward().Value())
			}
		} else {
			return errMicrodelegationNotExist
		}
	}

	if totalReward.Sign() == 0 {
		return errNoRewardsToCollect
	}
	return nil
}

func (verifier StakingVerifier) VerifyRenewMap3NodeMsg(stateDB vm.StateDB, chainContext ChainContext, blockNum, epoch *big.Int,
	msg *microstaking.RenewMap3Node, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if chainContext == nil {
		return errChainContextMissing
	}
	if epoch == nil {
		return errEpochMissing
	}
	if blockNum == nil {
		return errBlockNumMissing
	}
	if msg.DelegatorAddress != signer {
		return errInvalidSigner
	}

	node, err := stateDB.Map3NodeByAddress(msg.Map3NodeAddress)
	if err != nil {
		return err
	}
	if !node.Map3Node().AtStatus(microstaking.Active) {
		return errMap3NodeRenewalNotAllowed
	}

	md, ok := node.Microdelegations().Get(msg.DelegatorAddress)
	if !ok {
		return errMicrodelegationNotExist
	}

	if !md.Renewal().AtStatus(microstaking.Undecided) {
		return errChangeRenewalDecisionNotAllowed
	}

	curEpoch := common.NewDecFromBigInt(epoch)
	releaseEpoch := node.Map3Node().ReleaseEpoch().Value()
	if node.IsOperator(msg.DelegatorAddress) {
		// within the penultimate 7 epochs
		intervalFrom := releaseEpoch.Sub(common.NewDec(2*microstaking.RenewalTimeWindowInEpoch - 1))
		intervalTo := releaseEpoch.Sub(common.NewDec(7))
		if !curEpoch.BTE(intervalFrom, intervalTo) {
			return errMap3NodeRenewalNotAllowed
		}

		if msg.IsRenew {
			var percent *common.Dec
			if !chainContext.Config().Atlas.IsMicrostakingImprove(blockNum) {
				// self delegation proportion
				requireTotal, _, _ := network.LatestMicrostakingRequirement(blockNum, chainContext.Config())
				p := common.NewDecFromInt(md.Amount().Value()).QuoInt(requireTotal)
				percent = &p
			}

			node, err := node.Map3Node().Load()
			if err != nil {
				return err
			}

			if msg.NewCommissionRate != nil && !msg.NewCommissionRate.IsNil() {
				node.Commission.RateForNextPeriod = *msg.NewCommissionRate
			}
			return node.SanityCheck(microstaking.MaxPubKeyAllowed, percent)
		}
	} else {
		if msg.NewCommissionRate != nil {
			return errCommissionUpdateNotAllow
		}

		mdByOperator, ok := node.Microdelegations().Get(node.Map3Node().OperatorAddress().Value())
		if !ok {
			return errMicrodelegationNotExist
		}

		if mdByOperator.Renewal().AtStatus(microstaking.NotRenewed) {
			return errMap3NodeNotRenewalAnyMore
		}

		// the last 7 epoch
		//If the operator has already decided to renew, the participant can decide whether to renew before the release
		startEpoch := releaseEpoch.Sub(common.NewDec(microstaking.RenewalTimeWindowInEpoch - 1))
		if mdByOperator.Renewal().AtStatus(microstaking.Undecided) {
			if !curEpoch.GTE(startEpoch) {
				return errMap3NodeRenewalNotAllowed
			}
		}
	}
	return nil
}
