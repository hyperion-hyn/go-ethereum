package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/numeric"
	staking "github.com/ethereum/go-ethereum/staking/types"
	"github.com/pkg/errors"
	"math/big"
)

var (
	errStateDBIsMissing                = errors.New("no stateDB was provided")
	errChainContextMissing             = errors.New("no chain context was provided")
	errEpochMissing                    = errors.New("no epoch was provided")
	errBlockNumMissing                 = errors.New("no block number was provided")
	errNegativeAmount                  = errors.New("amount can not be negative")
	errInvalidSigner                   = errors.New("invalid signer for staking transaction")
	errDupIdentity                     = errors.New("validator identity exists")
	errDupPubKey                       = errors.New("public key exists")
	errInsufficientBalanceForStake     = errors.New("insufficient balance to stake")
	ErrMap3NodeNotExist                = errors.New("staking validator does not exist")
	errMap3NodeSnapshotNotExist        = errors.New("map3 node snapshot not found.")
	errCommissionRateChangeTooHigh     = errors.New("commission rate can not be higher than maximum commission rate")
	errCommissionRateChangeTooFast     = errors.New("change on commission rate can not be more than max change rate within the same epoch")
	errDelegationTooSmall              = errors.New("delegation amount too small")
	errInvalidNodeStateForDelegation   = errors.New("invalid node state for delegation")
	errUnmicrodelegateNotAllowed       = errors.New("not allow to unmicrodelegate in pending status")
	ErrMicrodelegationNotExist         = errors.New("no microdelegation exists")
	errNoRewardsToCollect              = errors.New("no rewards to collect")
	errMap3NodeAlreadyRedelegate       = errors.New("map3 node already redelegated.")
	errInvalidNodeStateForRedelegation = errors.New("invalid node state for redelegation")
	ErrValidatorNotExist               = errors.New("staking validator does not exist")
	errValidatorSnapshotNotExit        = errors.New("validator snapshot not found.")
	ErrRedelegationNotExist            = errors.New("no redelegation exists")
	errMap3NodeRenewalNotAllowed       = errors.New("map3 node renewal not allowed")
	errMap3NodeAlreadyRenewal          = errors.New("map3 node already renewal")
	errMap3NodeNotRenewalByInitiator   = errors.New("map3 node not renewal by initiator")
)

func checkNodeDuplicatedFields(state vm.StateDB, identity string, keys staking.Map3NodeKeys) error {
	map3NodePool := state.Map3NodePool()
	if identity != "" {
		identitySet := map3NodePool.GetDescriptionIdentitySet()
		if identitySet.Contain(identity) {
			return errors.Wrapf(errDupIdentity, "duplicate identity %s", identity)
		}
	}
	if len(keys) != 0 {
		nodeKeySet := map3NodePool.GetNodeKeySet()
		for _, key := range keys {
			if nodeKeySet.Contain(key.Hex()) {
				return errors.Wrapf(errDupPubKey, "duplicate public key %x", key.Hex())
			}
		}
	}
	return nil
}

func checkValDuplicatedFields(state vm.StateDB, identity string, keys staking.BLSPublicKeys) error {
	validatorPool := state.ValidatorPool()
	if identity != "" {
		identitySet := validatorPool.GetDescriptionIdentitySet()
		if identitySet.Contain(identity) {
			return errors.Wrapf(errDupIdentity, "duplicate identity %s", identity)
		}
	}
	if len(keys) != 0 {
		slotKeySet := validatorPool.GetSlotKeySet()
		for _, key := range keys {
			if slotKeySet.Contain(key.Hex()) {
				return errors.Wrapf(errDupPubKey, "duplicate public key %x", key.Hex())
			}
		}
	}
	return nil
}

func VerifyCreateMap3NodeMsg(stateDB vm.StateDB, epoch, blockNum *big.Int, msg *staking.CreateMap3Node,
	signer common.Address, minDel *big.Int) (*staking.Map3Node, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
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
	if msg.InitiatorAddress != signer {
		return nil, errInvalidSigner
	}

	if err := checkNodeDuplicatedFields(
		stateDB,
		msg.Description.Identity,
		msg.NodeKeys); err != nil {
		return nil, err
	}
	if !CanTransfer(stateDB, msg.InitiatorAddress, msg.Amount) {
		return nil, errInsufficientBalanceForStake
	}

	if minDel.Cmp(msg.Amount) > 0 {
		return nil, errDelegationTooSmall
	}

	nodeAddress := crypto.CreateAddress(signer, stateDB.GetNonce(signer))
	node, err := staking.CreateMap3NodeFromNewMsg(msg, nodeAddress, blockNum)
	if err != nil {
		return nil, err
	}
	if err := node.SanityCheck(staking.MaxPubKeyAllowed); err != nil {
		return nil, err
	}

	return node, nil
}

func VerifyEditMap3NodeMsg(
	stateDB vm.StateDB, epoch, blockNum *big.Int, msg *staking.EditMap3Node, signer common.Address,
) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if epoch == nil {
		return errEpochMissing
	}
	if blockNum == nil {
		return errBlockNumMissing
	}

	if err := checkNodeDuplicatedFields(
		stateDB,
		msg.Description.Identity,
		staking.Map3NodeKeys{*msg.NodeKeyToAdd}); err != nil {
		return err
	}
	nodePool := stateDB.Map3NodePool()
	wrapper, ok := nodePool.GetNodes().Get(msg.Map3NodeAddress)
	if !ok {
		return ErrMap3NodeNotExist
	}

	node := wrapper.GetMap3Node().ToMap3Node()
	if node.InitiatorAddress != signer {
		return errInvalidSigner
	}

	if err := staking.UpdateMap3NodeFromEditMsg(node, msg); err != nil {
		return err
	}
	if err := node.SanityCheck(staking.MaxPubKeyAllowed); err != nil {
		return err
	}

	newRate := node.Commission.CommissionRates.Rate
	if newRate.GT(node.Commission.CommissionRates.MaxRate) {
		return errCommissionRateChangeTooHigh
	}

	lastEpoch := big.NewInt(0).Sub(epoch, common.Big1)
	nodeSnapshot, ok := nodePool.GetNodeSnapshotByEpoch().Get(lastEpoch.Uint64())
	if !ok {
		return errMap3NodeSnapshotNotExist
	}
	snapshotMap3Node, ok := nodeSnapshot.GetMap3Nodes().Get(node.NodeAddress)
	if !ok {
		return errMap3NodeSnapshotNotExist
	}

	rateAtBeginningOfEpoch := snapshotMap3Node.GetMap3Node().GetCommission().GetCommissionRates().GetRate()
	if newRate.Sub(*rateAtBeginningOfEpoch).Abs().GT(node.Commission.CommissionRates.MaxChangeRate) {
		return errCommissionRateChangeTooFast
	}
	return nil
}

// VerifyRedelegateMsg verifies the delegate message using the stateDB
// and returns the balance to be deducted by the delegator as well as the
// validatorWrapper with the delegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyMicrodelegateMsg(stateDB vm.StateDB, msg *staking.Microdelegate, minDel *big.Int, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}

	if signer != msg.DelegatorAddress {
		return errInvalidSigner
	}

	if msg.Amount.Sign() == -1 {
		return errNegativeAmount
	}

	map3NodePool := stateDB.Map3NodePool()
	wrapper, ok := map3NodePool.GetNodes().Get(msg.Map3NodeAddress)
	if !ok {
		return ErrMap3NodeNotExist
	}

	// Check if there is enough liquid token to delegate
	if !CanTransfer(stateDB, msg.DelegatorAddress, msg.Amount) {
		return errInsufficientBalanceForStake
	}

	if minDel.Cmp(msg.Amount) > 0 {
		return errDelegationTooSmall
	}

	isInitiator := wrapper.GetMap3Node().GetInitiatorAddress() == signer
	status := wrapper.GetNodeState().GetStatus()
	if !isInitiator && status == staking.Inactive {
		return errInvalidNodeStateForDelegation
	}
	return nil
}

// VerifyUnredelegateMsg verifies the undelegate validator message
// using the stateDB & chainContext and returns the edited validatorWrapper
// with the undelegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyUnmicrodelegateMsg(stateDB vm.StateDB, epoch *big.Int, msg *staking.Unmicrodelegate, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
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

	map3NodePool := stateDB.Map3NodePool()
	wrapper, ok := map3NodePool.GetNodes().Get(msg.Map3NodeAddress)
	if !ok {
		return ErrMap3NodeNotExist
	}

	if wrapper.GetNodeState().GetStatus() != staking.Pending {
		return errUnmicrodelegateNotAllowed
	}

	md, ok := wrapper.GetMicrodelegations().Get(msg.DelegatorAddress)
	if !ok {
		return ErrMicrodelegationNotExist
	}

	p := md.GetPendingDelegation()
	if p == nil {
		return err
	}

	if p.GetUnlockedEpoch().GT(numeric.ZeroDec()) {
		return err
	}
	if p.GetAmount().Cmp(msg.Amount) > 0 {
		return err
	}
	return nil
}

// VerifyCollectRedelRewardsMsg verifies and collects rewards
// from the given delegation slice using the stateDB. It returns all of the
// edited validatorWrappers and the sum total of the rewards.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyCollectMicrodelRewardsDelegation(
	stateDB vm.StateDB, msg *staking.CollectMicrodelegationRewards, signer common.Address,
) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if msg.DelegatorAddress != signer {
		return errInvalidSigner
	}
	map3NodePool := stateDB.Map3NodePool()
	nodeAddressSetByDelegator := map3NodePool.GetNodeAddressSetByDelegator()
	nodeAddressSet, ok := nodeAddressSetByDelegator.Get(signer)
	if !ok {
		return errNoRewardsToCollect
	}

	totalRewards := common.Big0
	for _, nodeAddr := range nodeAddressSet.Keys() {
		if node, ok := map3NodePool.GetNodes().Get(nodeAddr); ok {
			if micro, ok := node.GetMicrodelegations().Get(signer); ok {
				if micro.GetReward().Cmp(common.Big0) > 0 {
					totalRewards.Add(totalRewards, micro.GetReward())
				}
			} else {
				return ErrMicrodelegationNotExist
			}
		} else {
			return ErrMap3NodeNotExist
		}
	}

	if totalRewards.Int64() == 0 {
		return errNoRewardsToCollect
	}
	return nil
}

func VerifyDivideNodeStakeMsg(stateDB vm.StateDB, blockNum *big.Int, msg *staking.DivideMap3NodeStake, singer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if blockNum == nil {
		return errBlockNumMissing
	}

	map3NodePool := stateDB.Map3NodePool()
	wrapper, ok := map3NodePool.GetNodes().Get(msg.Map3NodeAddress)
	if !ok {
		return ErrMap3NodeNotExist
	}
	if wrapper.GetMap3Node().GetInitiatorAddress() != singer {
		return errInvalidSigner
	}

	// TODO: divided node need 20%

	// node exist
	// node state
	// total delegation

	return nil
}

func VerifyRenewNodeStakeMsg(stateDB vm.StateDB, epoch, blockNum *big.Int, msg *staking.RenewMap3NodeStake, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
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

	map3NodePool := stateDB.Map3NodePool()
	wrapper, ok := map3NodePool.GetNodes().Get(msg.Map3NodeAddress)
	if !ok {
		return ErrMap3NodeNotExist
	}
	md, ok := wrapper.GetMicrodelegations().Get(msg.DelegatorAddress)
	if !ok {
		return ErrMicrodelegationNotExist
	}

	if wrapper.GetNodeState().GetStatus() != staking.Active ||
		big.NewInt(0).Sub(wrapper.GetNodeState().GetReleaseEpoch(), epoch).Int64() > staking.Map3NodeRenewalPeriodInEpoch {
		return errMap3NodeRenewalNotAllowed
	}

	if md.GetRenewal() != nil {
		return errMap3NodeAlreadyRenewal
	}

	imd, ok := wrapper.GetMicrodelegations().Get(wrapper.GetMap3Node().GetInitiatorAddress())
	if !ok {
		return ErrMicrodelegationNotExist
	}

	if imd.GetRenewal() != nil && !imd.GetRenewal().IsRenew() {
		return errMap3NodeNotRenewalByInitiator
	}

	return nil
}

// TODO: add unit tests to check staking msg verification

// VerifyCreateValidatorMsg verifies the create validator message using
// the stateDB, epoch, & blocknumber and returns the validatorWrapper created
// in the process.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyCreateValidatorMsg(
	stateDB vm.StateDB, blockNum *big.Int, msg *staking.CreateValidator, signer common.Address,
) (*staking.ValidatorWrapper, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}
	if blockNum == nil {
		return nil, errBlockNumMissing
	}

	map3NodePool := stateDB.Map3NodePool()
	node, ok := map3NodePool.GetNodes().Get(msg.InitiatorAddress)
	if !ok {
		return nil, ErrMap3NodeNotExist
	}

	if node.GetMap3Node().GetInitiatorAddress() != signer {
		return nil, errInvalidSigner
	}

	if node.GetNodeState().GetStatus() != staking.Active {
		return nil, errInvalidNodeStateForRedelegation
	}

	if node.GetRedelegationReference() != common.Address0 {
		return nil, errMap3NodeAlreadyRedelegate
	}

	if err := checkValDuplicatedFields(
		stateDB,
		msg.Description.Identity,
		msg.SlotPubKeys); err != nil {
		return nil, err
	}

	valAddress := crypto.CreateAddress(signer, stateDB.GetNonce(signer))
	v, err := staking.CreateValidatorFromNewMsg(msg, valAddress, blockNum)
	if err != nil {
		return nil, err
	}
	if err := v.SanityCheck(staking.MaxPubKeyAllowed); err != nil {
		return nil, err
	}
	wrapper := &staking.ValidatorWrapper{}
	wrapper.Validator = *v
	wrapper.Redelegations = staking.Redelegations{
		msg.InitiatorAddress: staking.NewRedelegation(msg.InitiatorAddress, node.GetTotalDelegation()),
	}
	wrapper.Counters.NumBlocksSigned = big.NewInt(0)
	wrapper.Counters.NumBlocksToSign = big.NewInt(0)
	wrapper.BlockReward = big.NewInt(0)
	wrapper.TotalDelegation = big.NewInt(0).Set(node.GetTotalDelegation())
	return wrapper, nil
}

// VerifyEditValidatorMsg verifies the edit validator message using
// the stateDB, chainContext and returns the edited validatorWrapper.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyEditValidatorMsg(
	stateDB vm.StateDB, epoch, blockNum *big.Int, msg *staking.EditValidator, signer common.Address,
) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if epoch == nil {
		return errEpochMissing
	}
	if blockNum == nil {
		return errBlockNumMissing
	}

	if err := checkValDuplicatedFields(
		stateDB,
		msg.Description.Identity,
		staking.BLSPublicKeys{*msg.SlotKeyToAdd}); err != nil {
		return err
	}

	nodePool := stateDB.Map3NodePool()
	validatorPool := stateDB.ValidatorPool()
	wrapper, ok := validatorPool.GetValidators().Get(msg.ValidatorAddress)
	if !ok {
		return ErrValidatorNotExist
	}
	validator := wrapper.GetValidator().ToValidator()

	// check signer
	found := false
	for addr := range validator.InitiatorAddresses {
		node, ok := nodePool.GetNodes().Get(addr)
		if !ok {
			return ErrMap3NodeNotExist
		}
		if node.GetMap3Node().GetInitiatorAddress() == signer {
			found = true
			break
		}
	}
	if !found {
		return errInvalidSigner
	}

	if err := staking.UpdateValidatorFromEditMsg(validator, msg); err != nil {
		return err
	}
	if err := validator.SanityCheck(staking.MaxPubKeyAllowed); err != nil {
		return err
	}

	newRate := validator.Commission.CommissionRates.Rate
	if newRate.GT(validator.Commission.CommissionRates.MaxRate) {
		return errCommissionRateChangeTooHigh
	}

	lastEpoch := big.NewInt(0).Sub(epoch, common.Big1)
	validatorSnapshot, ok := validatorPool.GetValidatorSnapshotByEpoch().Get(lastEpoch.Uint64())
	if !ok {
		return errValidatorSnapshotNotExit
	}
	snapshotValidator, ok := validatorSnapshot.GetValidators().Get(validator.ValidatorAddress)
	if !ok {
		return errMap3NodeSnapshotNotExist
	}

	rateAtBeginningOfEpoch := snapshotValidator.GetValidator().GetCommission().GetCommissionRates().GetRate()
	if newRate.Sub(*rateAtBeginningOfEpoch).Abs().GT(validator.Commission.CommissionRates.MaxChangeRate, ) {
		return errCommissionRateChangeTooFast
	}
	return nil
}

// VerifyRedelegateMsg verifies the delegate message using the stateDB
// and returns the balance to be deducted by the delegator as well as the
// validatorWrapper with the delegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyRedelegateMsg(stateDB vm.StateDB, msg *staking.Redelegate, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}

	if msg.Amount.Sign() == -1 {
		return errNegativeAmount
	}

	map3NodePool := stateDB.Map3NodePool()
	node, ok := map3NodePool.GetNodes().Get(msg.DelegatorAddress)
	if !ok {
		return ErrMap3NodeNotExist
	}
	if node.GetMap3Node().GetInitiatorAddress() != signer {
		return errInvalidSigner
	}

	if node.GetNodeState().GetStatus() != staking.Active {
		return errInvalidNodeStateForRedelegation
	}

	if node.GetRedelegationReference() != common.Address0 {
		return errMap3NodeAlreadyRedelegate
	}

	validatorPool := stateDB.ValidatorPool()
	if !validatorPool.GetValidators().Contain(msg.ValidatorAddress) {
		return ErrValidatorNotExist
	}

	return nil
}

// VerifyUnredelegateMsg verifies the undelegate validator message
// using the stateDB & chainContext and returns the edited validatorWrapper
// with the undelegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyUnredelegateMsg(stateDB vm.StateDB, epoch *big.Int, msg *staking.Unredelegate, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if epoch == nil {
		return errEpochMissing
	}

	validatorPool := stateDB.ValidatorPool()
	validator, ok := validatorPool.GetValidators().Get(msg.ValidatorAddress)
	if !ok {
		return ErrValidatorNotExist
	}

	redelegation, ok := validator.GetRedelegations().Get(msg.DelegatorAddress)
	if !ok {
		return ErrRedelegationNotExist
	}

	node, ok := stateDB.Map3NodePool().GetNodes().Get(msg.DelegatorAddress)
	if !ok {
		return ErrMap3NodeNotExist
	}
	if node.GetMap3Node().GetInitiatorAddress() != signer {
		return errInvalidSigner
	}

	if redelegation.GetAmount().Cmp(common.Big0) == 0 {
		return ErrRedelegationNotExist
	}
	return nil
}

// VerifyCollectRedelRewardsMsg verifies and collects rewards
// from the given delegation slice using the stateDB. It returns all of the
// edited validatorWrappers and the sum total of the rewards.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyCollectRedelRewardsMsg(stateDB vm.StateDB, msg *staking.CollectRedelegationRewards, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}

	validatorPool := stateDB.ValidatorPool()
	validator, ok := validatorPool.GetValidators().Get(msg.ValidatorAddress)
	if !ok {
		return ErrValidatorNotExist
	}

	redelegation, ok := validator.GetRedelegations().Get(msg.DelegatorAddress)
	if !ok {
		return ErrRedelegationNotExist
	}

	node, ok := stateDB.Map3NodePool().GetNodes().Get(msg.DelegatorAddress)
	if !ok {
		return ErrMap3NodeNotExist
	}
	if node.GetMap3Node().GetInitiatorAddress() != signer {
		return errInvalidSigner
	}

	if redelegation.GetReward().Cmp(common.Big0) == 0 {
		return errNoRewardsToCollect
	}
	return nil
}
