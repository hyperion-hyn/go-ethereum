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
	errInvalidTotalDelegation          = errors.New("total delegation can not be bigger than max_total_delegation", )
	errInsufficientBalanceToUndelegate = errors.New("insufficient balance to undelegate")
)

var (
	signerVerifier RestakingSignerQualificationVerifier = signerVerifierForTokenHolder{}
)

type RestakingSignerQualificationVerifier interface {
	VerifyCreateValidatorMsg(stateDB vm.StateDB, msg *restaking.CreateValidator, signer common.Address) error
	VerifyEditValidatorMsg(stateDB vm.StateDB, msg *restaking.EditValidator, signer common.Address) error
	VerifyRedelegateMsg(stateDB vm.StateDB, msg *restaking.Redelegate, signer common.Address) error
	VerifyUnredelegateMsg(stateDB vm.StateDB, msg *restaking.Unredelegate, signer common.Address) error
	VerifyCollectRedelRewardsMsg(stateDB vm.StateDB, msg *restaking.CollectReward, signer common.Address) error
}

type signerVerifierForTokenHolder struct {
}

func (s signerVerifierForTokenHolder) VerifyCreateValidatorMsg(stateDB vm.StateDB, msg *restaking.CreateValidator, signer common.Address) error {
	if msg.OperatorAddress != signer {
		return errInvalidSigner
	}

	if !CanTransfer(stateDB, signer, defaultStakingAmount) {
		return errInsufficientBalanceForStake
	}
	return nil
}

func (s signerVerifierForTokenHolder) VerifyEditValidatorMsg(stateDB vm.StateDB, msg *restaking.EditValidator, signer common.Address) error {
	validator, err := stateDB.ValidatorByAddress(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	if signer != msg.OperatorAddress {
		return errInvalidSigner
	}

	if !validator.Validator().OperatorAddresses().Set().Get(msg.OperatorAddress).Value() {
		return errInvalidValidatorOperator
	}
	return nil
}

func (s signerVerifierForTokenHolder) VerifyRedelegateMsg(stateDB vm.StateDB, msg *restaking.Redelegate, signer common.Address) error {
	if msg.DelegatorAddress != signer {
		return errInvalidSigner
	}

	if !CanTransfer(stateDB, signer, defaultStakingAmount) {
		return errInsufficientBalanceForStake
	}
	return nil
}

func (s signerVerifierForTokenHolder) VerifyUnredelegateMsg(stateDB vm.StateDB, msg *restaking.Unredelegate, signer common.Address) error {
	if msg.DelegatorAddress != signer {
		return errInvalidSigner
	}

	validator, err := stateDB.ValidatorByAddress(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	if !validator.Redelegations().Contain(msg.DelegatorAddress) {
		return errRedelegationNotExist
	}
	return nil
}

func (s signerVerifierForTokenHolder) VerifyCollectRedelRewardsMsg(stateDB vm.StateDB, msg *restaking.CollectReward, signer common.Address) error {
	if msg.DelegatorAddress != signer {
		return errInvalidSigner
	}

	validator, err := stateDB.ValidatorByAddress(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	if !validator.Redelegations().Contain(msg.DelegatorAddress) {
		return errRedelegationNotExist
	}
	return nil
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
	signer common.Address) (*restaking.Validator_, error) {
	if stateDB == nil {
		return nil, errStateDBIsMissing
	}
	if blockNum == nil {
		return nil, errBlockNumMissing
	}

	if err := signerVerifier.VerifyCreateValidatorMsg(stateDB, msg, signer); err != nil {
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

	if err = delegationSanityCheck(msg.MaxTotalDelegation, common.Big0, defaultStakingAmount); err != nil {
		return nil, err
	}

	return v, nil
}

// VerifyEditValidatorMsg verifies the edit validator message using
// the stateDB, chainContext and returns the edited validatorWrapper.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyEditValidatorMsg(stateDB vm.StateDB, chainContext ChainContext, epoch, blockNum *big.Int,
	msg *restaking.EditValidator, signer common.Address) error {
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

	blsKeys := restaking.NewEmptyBLSKeys()
	if msg.SlotKeyToAdd != nil {
		blsKeys.Keys = append(blsKeys.Keys, msg.SlotKeyToAdd)
	}
	if err := checkValidatorDuplicatedFields(stateDB, msg.Description.Identity, blsKeys); err != nil {
		return err
	}

	if err := signerVerifier.VerifyEditValidatorMsg(stateDB, msg, signer); err != nil {
		return err
	}

	wrapperSt, err := stateDB.ValidatorByAddress(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	validator := wrapperSt.Validator().Load()

	if err := restaking.UpdateValidatorFromEditMsg(validator, msg); err != nil {
		return err
	}
	if err := validator.SanityCheck(restaking.MaxPubKeyAllowed); err != nil {
		return err
	}

	// check max change at one epoch
	newRate := validator.Commission.CommissionRates.Rate
	validatorSnapshot, err := chainContext.ReadValidatorAtEpoch(epoch, msg.ValidatorAddress)
	if err != nil {
		return err
	}
	rateAtBeginningOfEpoch := validatorSnapshot.Validator().Commission().CommissionRates().Rate().Value()
	if newRate.Sub(rateAtBeginningOfEpoch).Abs().GT(validator.Commission.CommissionRates.MaxChangeRate) {
		return errCommissionRateChangeTooFast
	}

	if msg.MaxTotalDelegation != nil {
		if err = delegationSanityCheck(msg.MaxTotalDelegation, validator.MaxTotalDelegation, common.Big0); err != nil {
			return err
		}
	}
	return nil
}

// VerifyRedelegateMsg verifies the delegate message using the stateDB
// and returns the balance to be deducted by the delegator as well as the
// validatorWrapper with the delegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyRedelegateMsg(stateDB vm.StateDB, msg *restaking.Redelegate, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}

	if err := signerVerifier.VerifyRedelegateMsg(stateDB, msg, signer); err != nil {
		return err
	}

	if _, err := stateDB.ValidatorByAddress(msg.ValidatorAddress); err != nil {
		return err
	}

	// TODO(ATLAS): max total delegation && min delegation

	return nil
}

// VerifyUnredelegateMsg verifies the undelegate validator message
// using the stateDB & chainContext and returns the edited validatorWrapper
// with the undelegation applied to it.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyUnredelegateMsg(stateDB vm.StateDB, epoch *big.Int, msg *restaking.Unredelegate, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}
	if epoch == nil {
		return errEpochMissing
	}

	if err := signerVerifier.VerifyUnredelegateMsg(stateDB, msg, signer); err != nil {
		return err
	}

	validator, err := stateDB.ValidatorByAddress(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	redelegation, ok := validator.Redelegations().Get(msg.DelegatorAddress)
	if !ok {
		return errRedelegationNotExist
	}

	if redelegation.Amount().Value().Cmp(common.Big0) == 0 {
		return errInsufficientBalanceToUndelegate
	}
	return nil
}

// VerifyCollectRedelRewardsMsg verifies and collects rewards
// from the given delegation slice using the stateDB. It returns all of the
// edited validatorWrappers and the sum total of the rewards.
//
// Note that this function never updates the stateDB, it only reads from stateDB.
func VerifyCollectRedelRewardsMsg(stateDB vm.StateDB, msg *restaking.CollectReward, signer common.Address) error {
	if stateDB == nil {
		return errStateDBIsMissing
	}

	if err := signerVerifier.VerifyCollectRedelRewardsMsg(stateDB, msg, signer); err != nil {
		return err
	}

	validator, err := stateDB.ValidatorByAddress(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	redelegation, ok := validator.Redelegations().Get(msg.DelegatorAddress)
	if !ok {
		return errRedelegationNotExist
	}

	if redelegation.Reward().Value().Cmp(common.Big0) == 0 {
		return errNoRewardsToCollect
	}
	return nil
}

func delegationSanityCheck(maxTotalTotalDelegation, currentTotalDelegation, incrementalDelegation *big.Int) error {
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
