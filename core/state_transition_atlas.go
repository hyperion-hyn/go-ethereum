package core

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"math"
	"math/big"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
)

var (
	ErrInvalidStakingKind = errors.New("bad staking kind")

	defaultStakingAmount = big.NewInt(0).Mul(big.NewInt(params.Ether), big.NewInt(1000000)) // 1million * 10^18
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
func ApplyStakingMessage(evm *vm.EVM, msg Message, gp *GasPool, bc ChainContext) (*ExecutionResult, error) {
	return NewStateTransitionWithChainContext(evm, msg, gp, bc).StakingTransitionDb()
}

// StakingTransitionDb will transition the state by applying the staking message and
// returning the result including the used gas. It returns an error if failed.
// It is used for staking transaction only
func (st *StateTransition) StakingTransitionDb() (*ExecutionResult, error) {
	if err := st.preCheck(); err != nil {
		return nil, err
	}
	msg := st.msg
	sender := vm.AccountRef(msg.From())

	// Pay intrinsic gas
	gas, err := IntrinsicGasForStaking(st.data, msg.Type() == types.CreateValidator || msg.Type() == types.CreateMap3)
	if err != nil {
		return nil, err
	}
	if st.gas < gas {
		return nil, ErrIntrinsicGas
	}
	st.gas -= gas

	// Increment the nonce for the next transaction
	defer st.state.SetNonce(msg.From(), st.state.GetNonce(sender.Address())+1)

	verifier, err := NewStakingVerifier(st.bc)
	if err != nil {
		return nil, err
	}

	// staking errors as evm execution errors, not consensus errors
	switch msg.Type() {
	case types.CreateValidator:
		stkMsg := &restaking.CreateValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		st.state.IncreaseValidatorNonceIfZero()
		err = st.verifyAndApplyCreateValidatorTx(verifier, stkMsg, msg.From())
	case types.EditValidator:
		stkMsg := &restaking.EditValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyEditValidatorTx(verifier, stkMsg, msg.From())
	case types.Redelegate:
		stkMsg := &restaking.Redelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyRedelegateTx(verifier, stkMsg, msg.From())
	case types.Unredelegate:
		stkMsg := &restaking.Unredelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyUnredelegateTx(verifier, stkMsg, msg.From())
	case types.CollectRestakingReward:
		stkMsg := &restaking.CollectReward{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		_, err = st.verifyAndApplyCollectRedelRewardTx(verifier, stkMsg, msg.From())
	case types.CreateMap3:
		stkMsg := &microstaking.CreateMap3Node{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		st.state.IncreaseMap3NonceIfZero()
		err = st.verifyAndApplyCreateMap3NodeTx(verifier, stkMsg, msg.From())
	case types.EditMap3:
		stkMsg := &microstaking.EditMap3Node{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyEditMap3NodeTx(verifier, stkMsg, msg.From())
	case types.TerminateMap3:
		stkMsg := &microstaking.TerminateMap3Node{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyTerminateMap3NodeTx(verifier, stkMsg, msg.From())
	case types.Microdelegate:
		stkMsg := &microstaking.Microdelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyMicrodelegateTx(verifier, stkMsg, msg.From())
	case types.Unmicrodelegate:
		stkMsg := &microstaking.Unmicrodelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyUnmicrodelegateTx(verifier, stkMsg, msg.From())
	case types.CollectMicrostakingRewards:
		stkMsg := &microstaking.CollectRewards{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		_, err = st.verifyAndApplyCollectMicrodelRewardsTx(verifier, stkMsg, msg.From())
	case types.RenewMap3Node:
		stkMsg := &microstaking.RenewMap3Node{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyRenewMap3NodeTx(verifier, stkMsg, msg.From())
	default:
		return nil, ErrInvalidStakingKind
	}
	st.refundGas()

	if _, ok := st.bc.Engine().(consensus.Atlas); ok {
		// fee should not be given to evm.Coinbase directly, it should be distributed between validators.
		fee := new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice)
		pool := network.NewRewardPool(st.state)
		pool.AddTxFeeAsReward(st.evm.BlockNumber, fee)
	} else {
		st.state.AddBalance(st.evm.Coinbase, new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice))
	}

	result := ExecutionResult{
		UsedGas: st.gasUsed(),
	}
	if err != nil {
		// revert reason
		revert, err := abi.PackRevert(err.Error())
		if err != nil {
			return nil, err
		}
		result.Err = vm.ErrExecutionReverted
		result.ReturnData = revert
	}
	return &result, nil
}

func (st *StateTransition) verifyAndApplyCreateValidatorTx(verifier StakingVerifier, msg *restaking.CreateValidator, signer common.Address) error {
	verified, err := verifier.VerifyCreateValidatorMsg(st.state, st.evm.BlockNumber, msg, signer)
	if err != nil {
		return err
	}
	saveNewValidatorToPool(verified.NewValidator, st.state.ValidatorPool())
	return verified.Participant.postCreateValidator(verified.NewValidator.Validator.ValidatorAddress, verified.NewRedelegation)
}

func (st *StateTransition) verifyAndApplyEditValidatorTx(verifier StakingVerifier, msg *restaking.EditValidator, signer common.Address) error {
	if _, err := verifier.VerifyEditValidatorMsg(st.state, st.evm.BlockNumber, msg, signer); err != nil {
		return err
	}
	validatorPool := st.state.ValidatorPool()
	validator, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	updateValidatorFromPoolByMsg(validator, validatorPool, msg, st.evm.BlockNumber)
	return nil
}

func (st *StateTransition) verifyAndApplyRedelegateTx(verifier StakingVerifier, msg *restaking.Redelegate, signer common.Address) error {
	verified, err := verifier.VerifyRedelegateMsg(st.state, msg, signer)
	if err != nil {
		return err
	}
	wrapper, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	wrapper.AddRedelegation(msg.DelegatorAddress, verified.NewRedelegation)
	return verified.Participant.postRedelegate(msg.ValidatorAddress, verified.NewRedelegation)
}

func (st *StateTransition) verifyAndApplyUnredelegateTx(verifier StakingVerifier, msg *restaking.Unredelegate, signer common.Address) error {
	if _, err := verifier.VerifyUnredelegateMsg(st.state, st.evm.EpochNumber, msg, signer); err != nil {
		return err
	}

	validator, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	validator.Undelegate(msg.DelegatorAddress, st.evm.EpochNumber, nil)
	validator.InactivateIfSelfDelegationTooLittle()
	return nil
}

func (st *StateTransition) verifyAndApplyCollectRedelRewardTx(verifier StakingVerifier, msg *restaking.CollectReward, signer common.Address) (*big.Int, error) {
	verified, err := verifier.VerifyCollectRestakingRewardMsg(st.state, msg, signer)
	if err != nil {
		return network.NoReward, err
	}
	validator, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	reward, err := payoutRedelegationReward(validator, msg.DelegatorAddress, verified.Participant.rewardHandler(), st.evm.BlockNumber)
	if err != nil {
		return network.NoReward, err
	}

	// Add log if everything is good
	st.state.AddLog(&types.Log{
		Address:     msg.DelegatorAddress,
		Topics:      []common.Hash{restaking.CollectRewardTopic},
		Data:        reward.Bytes(),
		BlockNumber: st.evm.BlockNumber.Uint64(),
	})
	return reward, nil
}

func saveNewValidatorToPool(wrapper *restaking.ValidatorWrapper_, validatorPool *restaking.Storage_ValidatorPool_) {
	validatorPool.Validators().Put(wrapper.Validator.ValidatorAddress, wrapper)
	validatorPool.ValidatorSnapshots().Put(wrapper.Validator.ValidatorAddress, wrapper)
	keySet := validatorPool.SlotKeySet()
	for _, key := range wrapper.Validator.SlotPubKeys.Keys {
		keySet.Get(key.Hex()).SetValue(true)
	}

	if identity := wrapper.Validator.Description.Identity; identity != "" {
		validatorPool.DescriptionIdentitySet().Get(identity).SetValue(true)
	}
}

func updateValidatorFromPoolByMsg(validator *restaking.Storage_ValidatorWrapper_, pool *restaking.Storage_ValidatorPool_,
	msg *restaking.EditValidator, blockNum *big.Int) {
	// update description
	if msg.Description.Identity != "" {
		i := validator.Validator().Description().Identity().Value()
		pool.DescriptionIdentitySet().Get(i).SetValue(false)
		pool.DescriptionIdentitySet().Get(msg.Description.Identity).SetValue(true)
	}
	validator.Validator().Description().IncrementalUpdateFrom(msg.Description)

	if msg.CommissionRate != nil {
		curRate := validator.Validator().Commission().CommissionRates().Rate().Value()
		if !curRate.Equal(*msg.CommissionRate) {
			validator.Validator().Commission().CommissionRates().Rate().SetValue(*msg.CommissionRate)
			validator.Validator().Commission().UpdateHeight().SetValue(blockNum)
		}
	}

	if msg.SlotKeyToRemove != nil {
		for i := 0; i < validator.Validator().SlotPubKeys().Length(); i++ {
			if validator.Validator().SlotPubKeys().Get(i).Equal(msg.SlotKeyToRemove) {
				validator.Validator().SlotPubKeys().Remove(i)
				pool.SlotKeySet().Get(msg.SlotKeyToRemove.Hex()).SetValue(false)
				break
			}
		}
	}

	if msg.SlotKeyToAdd != nil {
		validator.Validator().SlotPubKeys().Push(msg.SlotKeyToAdd)
		pool.SlotKeySet().Get(msg.SlotKeyToAdd.Hex()).SetValue(true)
	}

	if msg.EPOSStatus == restaking.Active || msg.EPOSStatus == restaking.Inactive {
		validator.Validator().Status().SetValue(uint8(msg.EPOSStatus))
	}

	if msg.MaxTotalDelegation != nil && msg.MaxTotalDelegation.Sign() > 0 {
		validator.Validator().MaxTotalDelegation().SetValue(msg.MaxTotalDelegation)
	}
}

func payoutRedelegationReward(s *restaking.Storage_ValidatorWrapper_, delegator common.Address, handler RestakingRewardHandler,
	blockNum *big.Int) (*big.Int, error) {
	redelegation, ok := s.Redelegations().Get(delegator)
	if !ok {
		return nil, errMicrodelegationNotExist
	}

	r, err := handler.HandleReward(redelegation, blockNum)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type RestakingRewardHandler interface {
	HandleReward(redelegation *restaking.Storage_Redelegation_, blockNum *big.Int) (*big.Int, error)
}

type RewardToBalance struct {
	StateDB vm.StateDB
}

func (handler RewardToBalance) HandleReward(redelegation *restaking.Storage_Redelegation_, blockNum *big.Int) (*big.Int, error) {
	r := redelegation.Reward().Value()
	handler.StateDB.AddBalance(redelegation.DelegatorAddress().Value(), r)
	redelegation.Reward().Clear()
	return r, nil
}

type participant interface {
	restakingAmount() *big.Int
	postCreateValidator(validator common.Address, amount *big.Int) error
	postRedelegate(validator common.Address, amount *big.Int) error
	rewardHandler() RestakingRewardHandler
}

type tokenHolder struct {
	stateDB       vm.StateDB
	holderAddress common.Address
	amount        *big.Int
}

func (t tokenHolder) restakingAmount() *big.Int {
	return t.amount
}

func (t tokenHolder) postCreateValidator(validator common.Address, amount *big.Int) error {
	t.stateDB.SubBalance(t.holderAddress, amount)
	return nil
}

func (t tokenHolder) postRedelegate(validator common.Address, amount *big.Int) error {
	t.stateDB.SubBalance(t.holderAddress, amount)
	return nil
}

func (t tokenHolder) rewardHandler() RestakingRewardHandler {
	return &RewardToBalance{StateDB: t.stateDB}
}

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGasForStaking(data []byte, isMap3OrValidatorCreation bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if isMap3OrValidatorCreation {
		gas = params.TxGasMap3OrValidatorCreation
	} else {
		gas = params.TxGas
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		nonZeroGas := params.TxDataNonZeroGasEIP2028
		if (math.MaxUint64-gas)/nonZeroGas < nz {
			return 0, vm.ErrOutOfGas
		}
		gas += nz * nonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/params.TxDataZeroGas < z {
			return 0, vm.ErrOutOfGas
		}
		gas += z * params.TxDataZeroGas
	}
	return gas, nil
}
