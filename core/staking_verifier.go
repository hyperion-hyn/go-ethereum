package core

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/effective"
	staking "github.com/ethereum/go-ethereum/staking/types"
	"github.com/pkg/errors"
	"math/big"
)

var (
	ErrInvalidSelfDelegation         = errors.New("self delegation can not be less than min_self_delegation")
	errStateDBIsMissing              = errors.New("no stateDB was provided")
	errChainContextMissing           = errors.New("no chain context was provided")
	errEpochMissing                  = errors.New("no epoch was provided")
	errBlockNumMissing               = errors.New("no block number was provided")
	errNegativeAmount                = errors.New("amount can not be negative")
	errInvalidSigner                 = errors.New("invalid signer for staking transaction")
	errDupIdentity                   = errors.New("validator identity exists")
	errDupPubKey                     = errors.New("public key exists")
	errInsufficientBalanceForStake   = errors.New("insufficient balance to stake")
	errMap3NodeNotExist              = errors.New("staking validator does not exist")
	errMap3NodeSnapshotNotExist      = errors.New("map3 node snapshot not found.")
	errCommissionRateChangeTooHigh   = errors.New("commission rate can not be higher than maximum commission rate")
	errCommissionRateChangeTooFast   = errors.New("change on commission rate can not be more than max change rate within the same epoch")
	errDelegationTooSmall            = errors.New("minimum delegation amount too small")
	errInvalidNodeStateForDelegation = errors.New("invalid node state for delegation")
	errUnmicrodelegateNotAllowed     = errors.New("not allow to unmicrodelegate in pending status")
	errMicrodelegationNotExist       = errors.New("no microdelegation exists")
	errNoRewardsToCollect            = errors.New("no rewards to collect")
)

func checkNodeDuplicatedFields(state vm.StateDB, identity string, keys []staking.Map3NodeKey) error {
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

func VerifyAndCreateMap3NodeFromMsg(
	stateDB vm.StateDB, epoch, blockNum *big.Int, msg *staking.CreateMap3Node, signer common.Address, minSelf numeric.Dec,
) (*staking.Map3NodeWrapper, error) {
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

	if minSelf.GT(numeric.NewDecFromBigInt(msg.Amount)) {
		return nil, ErrInvalidSelfDelegation
	}

	nodeAddress := crypto.CreateAddress(signer, stateDB.GetNonce(signer))
	node, err := staking.CreateMap3NodeFromNewMsg(msg, nodeAddress, blockNum)
	if err != nil {
		return nil, err
	}
	if err := node.SanityCheck(staking.MaxPubKeyAllowed); err != nil {
		return nil, err
	}
	wrapper := &staking.Map3NodeWrapper{}
	wrapper.Map3Node = *node
	wrapper.Microdelegations = staking.Microdelegations{
		node.InitiatorAddress: staking.NewMicrodelegation(node.InitiatorAddress, msg.Amount, epoch, msg.AutoRenew, true),
	}
	wrapper.AccumulatedReward = big.NewInt(0)
	wrapper.NodeState = staking.NodeState{
		Status:          staking.Pending,
		NodeAge:         big.NewInt(0),
		CreationEpoch:   epoch,
		ActivationEpoch: big.NewInt(0),
		ReleaseEpoch:    big.NewInt(0),
	}
	wrapper.TotalDelegation = big.NewInt(0)
	wrapper.TotalPendingDelegation = big.NewInt(0).Set(msg.Amount)
	return wrapper, nil
}

func VerifyAndEditMap3NodeFromMsg(
	stateDB vm.StateDB, epoch, blockNum *big.Int, msg *staking.EditMap3Node, signer common.Address,
) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if blockNum == nil {
		return errBlockNumMissing
	}

	if err := checkNodeDuplicatedFields(
		stateDB,
		msg.Description.Identity,
		[]staking.Map3NodeKey{*msg.NodeKeyToAdd}); err != nil {
		return err
	}
	map3NodePool := stateDB.Map3NodePool()
	wrapperSt, ok := map3NodePool.GetNodes().Get(msg.Map3NodeAddress)
	if !ok {
		return errMap3NodeNotExist
	}
	if err := staking.UpdateMap3NodeFromEditMsg(wrapperSt.GetMap3Node(), map3NodePool.GetNodeKeySet(),
		map3NodePool.GetDescriptionIdentitySet(), msg); err != nil {
		return err
	}

	wrapper := wrapperSt.ToMap3NodeWrapper()
	if wrapper.Map3Node.InitiatorAddress != signer {
		return errInvalidSigner
	}
	if err := wrapper.Map3Node.SanityCheck(staking.MaxPubKeyAllowed); err != nil {
		return err
	}

	newRate := wrapper.Map3Node.Commission.CommissionRates.Rate
	if newRate.GT(wrapper.Map3Node.Commission.CommissionRates.MaxRate) {
		return errCommissionRateChangeTooHigh
	}

	lastEpoch := big.NewInt(0).Sub(epoch, common.Big1)
	nodeSnapshot, ok := map3NodePool.GetNodeSnapshotByEpoch().Get(lastEpoch.Uint64())
	if !ok {
		return errMap3NodeSnapshotNotExist
	}
	snapshotMap3Node, ok := nodeSnapshot.GetMap3Nodes().Get(wrapper.Map3Node.NodeAddress)
	if !ok {
		return errMap3NodeSnapshotNotExist
	}

	rateAtBeginningOfEpoch := snapshotMap3Node.GetMap3Node().GetCommission().GetCommissionRates().GetRate()
	if rateAtBeginningOfEpoch.IsNil() ||
		(!newRate.IsNil() && !rateAtBeginningOfEpoch.Equal(newRate)) {
		wrapperSt.GetMap3Node().GetCommission().SetUpdateHeight(blockNum)
	}

	if newRate.Sub(rateAtBeginningOfEpoch).Abs().GT(wrapper.Map3Node.Commission.CommissionRates.MaxChangeRate) {
		return errCommissionRateChangeTooFast
	}
	return nil
}

// VerifyAndDelegateFromMsg verifies the delegate message using the stateDB
// and returns the balance to be deducted by the delegator as well as the
// validatorWrapper with the delegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyAndMicrodelegateFromMsg(
	stateDB vm.StateDB, epoch *big.Int, msg *staking.Microdelegate, minDel numeric.Dec, signer common.Address,
) (*big.Int, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}

	if signer != msg.DelegatorAddress {
		return nil, errInvalidSigner
	}

	if msg.Amount.Sign() == -1 {
		return nil, errNegativeAmount
	}

	if minDel.GT(numeric.NewDecFromBigInt(msg.Amount)) {
		return nil, errDelegationTooSmall
	}

	map3NodePool := stateDB.Map3NodePool()
	wrapperSt, ok := map3NodePool.GetNodes().Get(msg.Map3NodeAddress)
	if !ok {
		return nil, errMap3NodeNotExist
	}

	// Check if there is enough liquid token to delegate
	if !CanTransfer(stateDB, msg.DelegatorAddress, msg.Amount) {
		return nil, errInsufficientBalanceForStake
	}

	status := wrapperSt.GetNodeState().GetStatus()
	if !(status == staking.Pending || status == staking.Active) {
		return nil, errInvalidNodeStateForDelegation
	}

	if status == staking.Active {
		// TODO(ATLAS): collect reward from validator as initiator
	}

	if microdelegation, ok := wrapperSt.GetMicrodelegations().Get(msg.DelegatorAddress); ok {
		if status == staking.Pending {
			pendingDelegations := microdelegation.GetPendingDelegations()
			if pendingDelegations.Len() != 0 &&
				pendingDelegations.Get(pendingDelegations.Len() - 1).GetEpoch().Cmp(epoch) == 0 {
				lastPending := pendingDelegations.Get(pendingDelegations.Len() - 1)
				amount := lastPending.GetAmount()
				lastPending.SetAmount(big.NewInt(0).Add(amount, msg.Amount))
			} else {
				pendingDelegations.Push(&staking.PendingDelegation{
					Amount: msg.Amount,
					Epoch:  epoch,
				})
			}
		} else {	// Active
			microdelegation.SetAmount(big.NewInt(0).Add(microdelegation.GetAmount(), msg.Amount))
		}
		microdelegation.SetAutoRenew(msg.AutoRenew)
	} else {
		wrapperSt.GetMicrodelegations().Put(
			msg.DelegatorAddress,
			staking.NewMicrodelegation(msg.DelegatorAddress, msg.Amount, epoch, msg.AutoRenew, status == staking.Pending),
		)
	}

	if status == staking.Active {
		wrapperSt.SetTotalPendingDelegation(big.NewInt(0).Add(wrapperSt.GetTotalPendingDelegation(), msg.Amount))
	} else {
		wrapperSt.SetTotalDelegation(big.NewInt(0).Add(wrapperSt.GetTotalDelegation(), msg.Amount))
	}
	return msg.Amount, nil
}

// VerifyAndUndelegateFromMsg verifies the undelegate validator message
// using the stateDB & chainContext and returns the edited validatorWrapper
// with the undelegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyAndUnmicrodelegateFromMsg(
	stateDB vm.StateDB, epoch *big.Int, msg *staking.Unmicrodelegate, signer common.Address,
) (*big.Int, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}
	if epoch == nil {
		return nil, errEpochMissing
	}
	if msg.Amount.Sign() == -1 {
		return nil, errNegativeAmount
	}
	if msg.DelegatorAddress != signer {
		return nil, errInvalidSigner
	}

	map3NodePool := stateDB.Map3NodePool()
	wrapperSt, ok := map3NodePool.GetNodes().Get(msg.Map3NodeAddress)
	if !ok {
		return nil, errMap3NodeNotExist
	}

	if wrapperSt.GetNodeState().GetStatus() != staking.Pending {
		return nil, errUnmicrodelegateNotAllowed
	}

	if microdelegation, ok := wrapperSt.GetMicrodelegations().Get(msg.DelegatorAddress); ok {
		if err := Unmicrodelgate(microdelegation, msg.Amount, epoch); err != nil {
			return nil, err
		}
		// TODO(ATLAS): delete delegation index
	} else {
		return nil, errMicrodelegationNotExist
	}
	return msg.Amount, nil
}

// VerifyAndCollectRewardsFromDelegation verifies and collects rewards
// from the given delegation slice using the stateDB. It returns all of the
// edited validatorWrappers and the sum total of the rewards.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyAndCollectMicrodelRewardsFromDelegation(
	stateDB vm.StateDB, msg *staking.CollectMicrodelegationRewards, signer common.Address,
) (*big.Int, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}
	if msg.DelegatorAddress != signer {
		return nil, errInvalidSigner
	}
	map3NodePool := stateDB.Map3NodePool()
	nodeAddressSetByDelegator := map3NodePool.GetNodeAddressSetByDelegator()
	nodeAddressSet, ok := nodeAddressSetByDelegator.Get(signer)
	if !ok {
		return common.Big0, errNoRewardsToCollect
	}

	totalRewards := common.Big0
	for _, nodeAddr := range nodeAddressSet.Keys() {
		if node, ok := map3NodePool.GetNodes().Get(nodeAddr); ok {
			if micro, ok := node.GetMicrodelegations().Get(signer); ok {
				if micro.GetReward().Cmp(common.Big0) > 0 {
					totalRewards.Add(totalRewards, micro.GetReward())
					micro.SetReward(common.Big0)
				}
			} else {
				return common.Big0, errMicrodelegationNotExist
			}
		} else {
			return common.Big0, errMap3NodeNotExist
		}
	}

	if totalRewards.Int64() == 0 {
		return nil, errNoRewardsToCollect
	}
	return totalRewards, nil
}



// TODO: add unit tests to check staking msg verification

// VerifyAndCreateValidatorFromMsg verifies the create validator message using
// the stateDB, epoch, & blocknumber and returns the validatorWrapper created
// in the process.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyAndCreateValidatorFromMsg(
	stateDB vm.StateDB, epoch *big.Int, blockNum *big.Int, msg *staking.CreateValidator,
) (*staking.ValidatorWrapper, error) {
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
	if stateDB.IsValidator(msg.ValidatorAddress) {
		return nil, errors.Wrapf(
			errValidatorExist, common2.MustAddressToBech32(msg.ValidatorAddress),
		)
	}
	if err := checkDuplicateFields(
		stateDB,
		msg.ValidatorAddress,
		msg.Identity,
		msg.SlotPubKeys); err != nil {
		return nil, err
	}
	if !CanTransfer(stateDB, msg.ValidatorAddress, msg.Amount) {
		return nil, errInsufficientBalanceForStake
	}
	v, err := staking.CreateValidatorFromNewMsg(msg, blockNum, epoch)
	if err != nil {
		return nil, err
	}
	wrapper := &staking.ValidatorWrapper{}
	wrapper.Validator = *v
	wrapper.Redelegations = []staking.Redelegation{
		staking.NewDelegation(v.Address, msg.Amount),
	}
	wrapper.Counters.NumBlocksSigned = big.NewInt(0)
	wrapper.Counters.NumBlocksToSign = big.NewInt(0)
	wrapper.BlockReward = big.NewInt(0)
	maxBLSKeyAllowed := shard.ExternalSlotsAvailableForEpoch(epoch) / 3
	if err := wrapper.SanityCheck(maxBLSKeyAllowed); err != nil {
		return nil, err
	}
	return wrapper, nil
}

// VerifyAndEditValidatorFromMsg verifies the edit validator message using
// the stateDB, chainContext and returns the edited validatorWrapper.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyAndEditValidatorFromMsg(
	stateDB vm.StateDB, chainContext ChainContext,
	epoch, blockNum *big.Int, msg *staking.EditValidator,
) (*staking.ValidatorWrapper, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}
	if chainContext == nil {
		return nil, errChainContextMissing
	}
	if blockNum == nil {
		return nil, errBlockNumMissing
	}
	if !stateDB.IsValidator(msg.ValidatorAddress) {
		return nil, errValidatorNotExist
	}
	newBlsKeys := []shard.BLSPublicKey{}
	if msg.SlotKeyToAdd != nil {
		newBlsKeys = append(newBlsKeys, *msg.SlotKeyToAdd)
	}
	if err := checkDuplicateFields(
		chainContext, stateDB,
		msg.ValidatorAddress,
		msg.Identity,
		newBlsKeys); err != nil {
		return nil, err
	}
	wrapper, err := stateDB.ValidatorWrapperCopy(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	if err := staking.UpdateValidatorFromEditMsg(&wrapper.Validator, msg, epoch); err != nil {
		return nil, err
	}
	newRate := wrapper.Validator.Rate
	if newRate.GT(wrapper.Validator.MaxRate) {
		return nil, errCommissionRateChangeTooHigh
	}

	snapshotValidator, err := chainContext.ReadValidatorSnapshot(wrapper.Address)
	if err != nil {
		return nil, errors.WithMessage(err, "validator snapshot not found.")
	}
	rateAtBeginningOfEpoch := snapshotValidator.Validator.Rate

	if rateAtBeginningOfEpoch.IsNil() ||
		(!newRate.IsNil() && !rateAtBeginningOfEpoch.Equal(newRate)) {
		wrapper.Validator.UpdateHeight = blockNum
	}

	if newRate.Sub(rateAtBeginningOfEpoch).Abs().GT(
		wrapper.Validator.MaxChangeRate,
	) {
		return nil, errCommissionRateChangeTooFast
	}
	maxBLSKeyAllowed := shard.ExternalSlotsAvailableForEpoch(epoch) / 3
	if err := wrapper.SanityCheck(maxBLSKeyAllowed); err != nil {
		return nil, err
	}
	return wrapper, nil
}


// VerifyAndDelegateFromMsg verifies the delegate message using the stateDB
// and returns the balance to be deducted by the delegator as well as the
// validatorWrapper with the delegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyAndDelegateFromMsg(
	stateDB vm.StateDB, msg *staking.Redelegate,
) (*staking.ValidatorWrapper, *big.Int, error) {
	if stateDB == nil {
		return nil, nil, errStateDBIsMissing
	}
	if !stateDB.IsValidator(msg.ValidatorAddress) {
		return nil, nil, errValidatorNotExist
	}
	if msg.Amount.Sign() == -1 {
		return nil, nil, errNegativeAmount
	}
	if msg.Amount.Cmp(minimumDelegation) < 0 {
		return nil, nil, errDelegationTooSmall
	}
	wrapper, err := stateDB.ValidatorWrapperCopy(msg.ValidatorAddress)
	if err != nil {
		return nil, nil, err
	}

	// Check if there is enough liquid token to delegate
	if !CanTransfer(stateDB, msg.DelegatorAddress, msg.Amount) {
		return nil, nil, errors.Wrapf(
			errInsufficientBalanceForStake, "had %v, tried to stake %v",
			stateDB.GetBalance(msg.DelegatorAddress), msg.Amount)
	}

	// Check for existing delegation
	for i := range wrapper.Delegations {
		delegation := &wrapper.Delegations[i]
		if bytes.Equal(delegation.DelegatorAddress.Bytes(), msg.DelegatorAddress.Bytes()) {
			delegation.Amount.Add(delegation.Amount, msg.Amount)
			if err := wrapper.SanityCheck(
				staking.DoNotEnforceMaxBLS,
			); err != nil {
				return nil, nil, err
			}
			return wrapper, msg.Amount, nil
		}
	}

	// Add new delegation
	wrapper. = append(
		wrapper.Delegations, staking.NewDelegation(
			msg.DelegatorAddress, msg.Amount,
		),
	)
	if err := wrapper.SanityCheck(staking.DoNotEnforceMaxBLS); err != nil {
		return nil, nil, err
	}
	return wrapper, msg.Amount, nil
}

// VerifyAndUndelegateFromMsg verifies the undelegate validator message
// using the stateDB & chainContext and returns the edited validatorWrapper
// with the undelegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyAndUndelegateFromMsg(
	stateDB vm.StateDB, epoch *big.Int, msg *staking.Unredelegate,
) (*staking.ValidatorWrapper, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}
	if epoch == nil {
		return nil, errEpochMissing
	}

	if msg.Amount.Sign() == -1 {
		return nil, errNegativeAmount
	}

	if !stateDB.IsValidator(msg.ValidatorAddress) {
		return nil, errValidatorNotExist
	}

	wrapper, err := stateDB.ValidatorWrapperCopy(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	for i := range wrapper.Delegations {
		delegation := &wrapper.Delegations[i]
		if bytes.Equal(delegation.DelegatorAddress.Bytes(), msg.DelegatorAddress.Bytes()) {
			if err := delegation.Undelegate(epoch, msg.Amount); err != nil {
				return nil, err
			}
			if err := wrapper.SanityCheck(
				staking.DoNotEnforceMaxBLS,
			); err != nil {
				// allow self delegation to go below min self delegation
				// but set the status to inactive
				if errors.Cause(err) == staking.ErrInvalidSelfDelegation {
					wrapper.Status = effective.Inactive
				} else {
					return nil, err
				}
			}
			return wrapper, nil
		}
	}
	return nil, errNoDelegationToUndelegate
}

// VerifyAndCollectRewardsFromDelegation verifies and collects rewards
// from the given delegation slice using the stateDB. It returns all of the
// edited validatorWrappers and the sum total of the rewards.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyAndCollectRewardsFromDelegation(
	stateDB vm.StateDB, delegations []staking.MicrodelegationIndex,
) ([]*staking.ValidatorWrapper, *big.Int, error) {
	if stateDB == nil {
		return nil, nil, errStateDBIsMissing
	}
	updatedValidatorWrappers := []*staking.ValidatorWrapper{}
	totalRewards := big.NewInt(0)
	for i := range delegations {
		delegation := &delegations[i]
		wrapper, err := stateDB.ValidatorWrapperCopy(delegation.ValidatorAddress)
		if err != nil {
			return nil, nil, err
		}
		if uint64(len(wrapper.Delegations)) > delegation.Index {
			delegation := &wrapper.Delegations[delegation.Index]
			if delegation.Reward.Cmp(common.Big0) > 0 {
				totalRewards.Add(totalRewards, delegation.Reward)
				delegation.Reward.SetUint64(0)
			}
		} else {
			utils.Logger().Warn().
				Str("validator", delegation.ValidatorAddress.String()).
				Uint64("delegation index", delegation.Index).
				Int("delegations length", len(wrapper.Delegations)).
				Msg("Redelegation index out of bound")
			return nil, nil, errors.New("Redelegation index out of bound")
		}
		updatedValidatorWrappers = append(updatedValidatorWrappers, wrapper)
	}
	if totalRewards.Int64() == 0 {
		return nil, nil, errNoRewardsToCollect
	}
	return updatedValidatorWrappers, totalRewards, nil
}

func Unmicrodelgate(microdelegation *staking.MicrodelegationStorage, amount, epoch *big.Int) error {
	amt := big.NewInt(0).Set(amount)
	var removedIndexes []int
	for i := 0; i < microdelegation.GetPendingDelegations().Len(); i++ {
		pd := microdelegation.GetPendingDelegations().Get(i)
		if big.NewInt(0).Sub(epoch, pd.GetEpoch()).Int64() > int64(staking.PendingDelegationLockPeriodInEpoch) {
			amt.Sub(amt, pd.GetAmount())
			removedIndexes = append(removedIndexes, i)
		} else {
			break
		}
	}
	return nil
}