package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/pkg/errors"
	"math/big"
)

var (
	ErrSelfDelegationTooSmall          = errors.New("self delegation amount too small")
	errStateDBIsMissing                = errors.New("no stateDB was provided")
	errChainContextMissing             = errors.New("no chain context was provided")
	errEpochMissing                    = errors.New("no epoch was provided")
	errBlockNumMissing                 = errors.New("no block number was provided")
	errNegativeAmount                  = errors.New("amount can not be negative")
	errInvalidSigner                   = errors.New("invalid signer for staking transaction")
	errDupIdentity                     = errors.New("validator identity exists")
	errDuplicateSlotKeys               = errors.New("slot keys can not have duplicates")
	errInsufficientBalanceForStake     = errors.New("insufficient balance to stake")
	errCommissionRateChangeTooFast     = errors.New("change on commission rate can not be more than max change rate within the same epoch")
	errDelegationTooSmall              = errors.New("delegation amount too small")
	errNoRewardsToCollect              = errors.New("no rewards to collect")
	errValidatorNotExist               = errors.New("staking validator does not exist")
	errRedelegationNotExist            = errors.New("redelegation does not exist")
	errInvalidValidatorOperator        = errors.New("invalid validator operator")
	errInvalidTotalDelegation          = errors.New("total delegation can not be bigger than max_total_delegation")
	errInsufficientBalanceToUndelegate = errors.New("insufficient balance to undelegate")
)

var (
	participantVerifier RestakingParticipantVerifier = TokenHolderVerifier{}
)

type RestakingParticipantVerifier interface {
	VerifyForCreatingValidator(stateDB vm.StateDB, msg *restaking.CreateValidator, signer common.Address) (participant, error)
	VerifyForEditingValidator(stateDB vm.StateDB, msg *restaking.EditValidator, signer common.Address) (participant, error)
	VerifyForRedelegating(stateDB vm.StateDB, msg *restaking.Redelegate, signer common.Address) (participant, error)
	VerifyForUnredelegating(stateDB vm.StateDB, msg *restaking.Unredelegate, signer common.Address) (participant, error)
	VerifyForCollectingReward(stateDB vm.StateDB, msg *restaking.CollectReward, signer common.Address) (participant, error)
}

type TokenHolderVerifier struct {
}

func (s TokenHolderVerifier) VerifyForCreatingValidator(stateDB vm.StateDB, msg *restaking.CreateValidator, signer common.Address) (participant, error) {
	if msg.OperatorAddress != signer {
		return nil, errInvalidSigner
	}

	if !CanTransfer(stateDB, signer, defaultStakingAmount) {
		return nil, errInsufficientBalanceForStake
	}

	return &tokenHolder{stateDB: stateDB, holderAddress: signer}, nil
}

func (s TokenHolderVerifier) VerifyForEditingValidator(stateDB vm.StateDB, msg *restaking.EditValidator, signer common.Address) (participant, error) {
	if signer != msg.OperatorAddress {
		return nil, errInvalidSigner
	}
	return &tokenHolder{stateDB: stateDB, holderAddress: signer}, nil
}

func (s TokenHolderVerifier) VerifyForRedelegating(stateDB vm.StateDB, msg *restaking.Redelegate, signer common.Address) (participant, error) {
	if msg.DelegatorAddress != signer {
		return nil, errInvalidSigner
	}

	if !CanTransfer(stateDB, signer, defaultStakingAmount) {
		return nil, errInsufficientBalanceForStake
	}
	return &tokenHolder{stateDB: stateDB, holderAddress: signer}, nil
}

func (s TokenHolderVerifier) VerifyForUnredelegating(stateDB vm.StateDB, msg *restaking.Unredelegate, signer common.Address) (participant, error) {
	if msg.DelegatorAddress != signer {
		return nil, errInvalidSigner
	}
	return &tokenHolder{stateDB: stateDB, holderAddress: signer}, nil
}

func (s TokenHolderVerifier) VerifyForCollectingReward(stateDB vm.StateDB, msg *restaking.CollectReward, signer common.Address) (participant, error) {
	if msg.DelegatorAddress != signer {
		return nil, errInvalidSigner
	}

	validator, err := stateDB.ValidatorByAddress(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	if !validator.Redelegations().Contain(msg.DelegatorAddress) {
		return nil, errRedelegationNotExist
	}
	return &tokenHolder{stateDB: stateDB, holderAddress: signer}, nil
}

func checkValidatorDuplicatedFields(state vm.StateDB, identity string, keys restaking.BLSPublicKeys_) error {
	validatorPool := state.ValidatorPool()
	if identity != "" {
		identitySet := validatorPool.DescriptionIdentitySet()
		if identitySet.Get(identity).Value() {
			return errors.Wrapf(errDupIdentity, "duplicate identity %s", identity)
		}
	}
	if len(keys.Keys) != 0 {
		slotKeySet := validatorPool.SlotKeySet()
		for _, key := range keys.Keys {
			if slotKeySet.Get(key.Hex()).Value() {
				return errors.Wrapf(errDuplicateSlotKeys, "duplicate public key %x", key.Hex())
			}
		}
	}
	return nil
}

// TODO: add unit tests to check staking msg verification

// VerifyCreateValidatorMsg verifies the create validator message using
// the stateDB, epoch, & blocknumber and returns the validatorWrapper created
// in the process.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyCreateValidatorMsg(stateDB vm.StateDB, blockNum *big.Int, msg *restaking.CreateValidator,
	signer common.Address) (*verification, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}
	if blockNum == nil {
		return nil, errBlockNumMissing
	}

	p, err := participantVerifier.VerifyForCreatingValidator(stateDB, msg, signer)
	if err != nil {
		return nil, err
	}

	if err := checkValidatorDuplicatedFields(stateDB, msg.Description.Identity, restaking.NewBLSKeysWithBLSKey(msg.SlotPubKey)); err != nil {
		return nil, err
	}

	valAddress := crypto.CreateAddress(signer, stateDB.GetNonce(signer))
	v, err := restaking.CreateValidatorFromNewMsg(msg, valAddress, blockNum)
	if err != nil {
		return nil, err
	}
	if err := v.SanityCheck(restaking.MaxPubKeyAllowed); err != nil {
		return nil, err
	}

	// TODO(ATLAS): Get staking amount by role of delegator
	amt := defaultStakingAmount
	if err = sanityCheckForDelegation(msg.MaxTotalDelegation, common.Big0, amt); err != nil {
		return nil, err
	}

	wrapper := restaking.ValidatorWrapper_{
		Validator:                 *v,
		Redelegations:             restaking.NewRedelegationMap(),
		TotalDelegation:           big.NewInt(0).Set(amt),
		TotalDelegationByOperator: big.NewInt(0).Set(amt),
		BlockReward:               big.NewInt(0),
	}
	wrapper.Counters.NumBlocksSigned = big.NewInt(0)
	wrapper.Counters.NumBlocksToSign = big.NewInt(0)
	wrapper.Redelegations.Put(msg.OperatorAddress, restaking.NewRedelegation(msg.OperatorAddress, amt))

	return &verification{
		NewValidator:    &wrapper,
		NewRedelegation: amt,
		Participant:     p,
	}, nil
}

// VerifyEditValidatorMsg verifies the edit validator message using
// the stateDB, chainContext and returns the edited validatorWrapper.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyEditValidatorMsg(stateDB vm.StateDB, chainContext ChainContext, epoch, blockNum *big.Int,
	msg *restaking.EditValidator, signer common.Address) (*verification, error) {
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

	blsKeys := restaking.NewEmptyBLSKeys()
	if msg.SlotKeyToAdd != nil {
		blsKeys.Keys = append(blsKeys.Keys, msg.SlotKeyToAdd)
	}
	if err := checkValidatorDuplicatedFields(stateDB, msg.Description.Identity, blsKeys); err != nil {
		return nil, err
	}

	p, err := participantVerifier.VerifyForEditingValidator(stateDB, msg, signer)
	if err != nil {
		return nil, err
	}

	wrapperSt, err := stateDB.ValidatorByAddress(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	if !wrapperSt.IsOperator(msg.OperatorAddress) {
		return nil, errInvalidValidatorOperator
	}
	validator := wrapperSt.Validator().Load()

	// TODO(ATLAS): update block num when updating commission rate?
	if err := restaking.UpdateValidatorFromEditMsg(validator, msg); err != nil {
		return nil, err
	}
	if err := validator.SanityCheck(restaking.MaxPubKeyAllowed); err != nil {
		return nil, err
	}

	// check max change at one epoch
	newRate := validator.Commission.CommissionRates.Rate
	validatorSnapshot, err := chainContext.ReadValidatorAtEpochOrCurrentBlock(epoch, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	rateAtBeginningOfEpoch := validatorSnapshot.Validator().Commission().CommissionRates().Rate().Value()
	if newRate.Sub(rateAtBeginningOfEpoch).Abs().GT(validator.Commission.CommissionRates.MaxChangeRate) {
		return nil, errCommissionRateChangeTooFast
	}

	if msg.MaxTotalDelegation != nil && msg.MaxTotalDelegation.Sign() != 0 {
		if err = sanityCheckForDelegation(msg.MaxTotalDelegation, wrapperSt.TotalDelegation().Value(), common.Big0); err != nil {
			return nil, err
		}
	}
	return &verification{
		Participant: p,
	}, nil
}

// VerifyRedelegateMsg verifies the delegate message using the stateDB
// and returns the balance to be deducted by the delegator as well as the
// validatorWrapper with the delegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyRedelegateMsg(stateDB vm.StateDB, msg *restaking.Redelegate, signer common.Address) (*verification, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}

	p, err := participantVerifier.VerifyForRedelegating(stateDB, msg, signer)
	if err != nil {
		return nil, err
	}

	if _, err := stateDB.ValidatorByAddress(msg.ValidatorAddress); err != nil {
		return nil, err
	}

	// TODO(ATLAS): Get staking amonut by role of delegator
	// TODO(ATLAS): max total delegation && min delegation
	amt := defaultStakingAmount

	return &verification{
		NewRedelegation: amt,
		Participant:     p,
	}, nil
}

// VerifyUnredelegateMsg verifies the undelegate validator message
// using the stateDB & chainContext and returns the edited validatorWrapper
// with the undelegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyUnredelegateMsg(stateDB vm.StateDB, epoch *big.Int, msg *restaking.Unredelegate, signer common.Address) (*verification, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}
	if epoch == nil {
		return nil, errEpochMissing
	}

	p, err := participantVerifier.VerifyForUnredelegating(stateDB, msg, signer)
	if err != nil {
		return nil, err
	}

	validator, err := stateDB.ValidatorByAddress(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	redelegation, ok := validator.Redelegations().Get(msg.DelegatorAddress)
	if !ok {
		return nil, errRedelegationNotExist
	}

	if redelegation.Amount().Value().Cmp(common.Big0) == 0 {
		return nil, errInsufficientBalanceToUndelegate
	}
	return &verification{
		Participant: p,
	}, nil
}

// VerifyCollectRedelRewardMsg verifies and collects rewards
// from the given delegation slice using the stateDB. It returns all of the
// edited validatorWrappers and the sum total of the rewards.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyCollectRedelRewardMsg(stateDB vm.StateDB, msg *restaking.CollectReward, signer common.Address) (*verification, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}

	p, err := participantVerifier.VerifyForCollectingReward(stateDB, msg, signer)
	if err != nil {
		return nil, err
	}

	validator, err := stateDB.ValidatorByAddress(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	redelegation, ok := validator.Redelegations().Get(msg.DelegatorAddress)
	if !ok {
		return nil, errRedelegationNotExist
	}

	if redelegation.Reward().Value().Cmp(common.Big0) == 0 {
		return nil, errNoRewardsToCollect
	}
	return &verification{
		Participant: p,
	}, nil
}

func sanityCheckForDelegation(maxTotalTotalDelegation, currentTotalDelegation, incrementalDelegation *big.Int) error {
	total := big.NewInt(0).Add(currentTotalDelegation, incrementalDelegation)
	if total.Cmp(maxTotalTotalDelegation) > 0 {
		return errors.Wrapf(
			errInvalidTotalDelegation,
			"total %s max-total %s",
			total.String(),
			maxTotalTotalDelegation.String(),
		)
	}
	return nil
}

type verification struct {
	NewValidator    *restaking.ValidatorWrapper_
	NewRedelegation *big.Int
	Participant     participant
}
