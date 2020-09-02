package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/pkg/errors"
	"math"
	"math/big"
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
	gas, err := IntrinsicGasForStaking(st.data, msg.Type() == types.CreateValidator)
	if err != nil {
		return nil, err
	}
	if st.gas < gas {
		return nil, ErrIntrinsicGas
	}
	st.gas -= gas

	// Increment the nonce for the next transaction
	defer st.state.SetNonce(msg.From(), st.state.GetNonce(sender.Address())+1)
	switch msg.Type() {
	case types.CreateValidator:
		stkMsg := &restaking.CreateValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		st.state.IncrementValidatorNonce()
		err = st.verifyAndApplyCreateValidatorTx(stkMsg, msg.From())
	case types.EditValidator:
		stkMsg := &restaking.EditValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyEditValidatorTx(stkMsg, msg.From())
	case types.Redelegate:
		stkMsg := &restaking.Redelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyRedelegateTx(stkMsg, msg.From())
	case types.Unredelegate:
		stkMsg := &restaking.Unredelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		err = st.verifyAndApplyUnredelegateTx(stkMsg, msg.From())
	case types.CollectRedelReward:
		stkMsg := &restaking.CollectReward{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return nil, err
		}
		_, err = st.verifyAndApplyCollectRedelRewardTx(stkMsg, msg.From())
		// TODO: Add log for reward ?
	default:
		return nil, ErrInvalidStakingKind
	}
	st.refundGas()

	// TODO(ATLAS): Txn Fees
	st.state.AddBalance(st.evm.Coinbase, new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice))

	return &ExecutionResult{
		UsedGas: st.gasUsed(),
	}, err
}

func (st *StateTransition) verifyAndApplyCreateValidatorTx(msg *restaking.CreateValidator, signer common.Address) error {
	verified, err := VerifyCreateValidatorMsg(st.state, st.evm.BlockNumber, msg, signer)
	if err != nil {
		return err
	}
	saveNewValidatorToPool(verified.NewValidator, st.state.ValidatorPool())
	return verified.Participant.PostCreateValidator(msg.OperatorAddress, verified.NewRedelegation)
}

func (st *StateTransition) verifyAndApplyEditValidatorTx(msg *restaking.EditValidator, signer common.Address) error {
	if _, err := VerifyEditValidatorMsg(st.state, st.bc, st.evm.EpochNumber, st.evm.BlockNumber, msg, signer); err != nil {
		return err
	}
	validatorPool := st.state.ValidatorPool()
	validator, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	updateValidatorFromPoolByMsg(validator, validatorPool, msg, st.evm.BlockNumber)
	return nil
}

func (st *StateTransition) verifyAndApplyRedelegateTx(msg *restaking.Redelegate, signer common.Address) error {
	verified, err := VerifyRedelegateMsg(st.state, msg, signer)
	if err != nil {
		return err
	}
	wrapper, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	wrapper.AddRedelegation(msg.DelegatorAddress, verified.NewRedelegation)
	return verified.Participant.PostRedelegate(msg.ValidatorAddress, verified.NewRedelegation)
}

func (st *StateTransition) verifyAndApplyUnredelegateTx(msg *restaking.Unredelegate, signer common.Address) error {
	if _, err := VerifyUnredelegateMsg(st.state, st.evm.EpochNumber, msg, signer); err != nil {
		return err
	}

	validator, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	validator.Undelegate(msg.DelegatorAddress, st.evm.EpochNumber)

	// TODO: need 20%? change state to inactive?
	return nil
}

func (st *StateTransition) verifyAndApplyCollectRedelRewardTx(msg *restaking.CollectReward, signer common.Address) (*big.Int, error) {
	if _, err := VerifyCollectRedelRewardMsg(st.state, msg, signer); err != nil {
		return network.NoReward, err
	}
	validator, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	handler := RewardToBalance{StateDB: st.state} // TODO(ATLAS): map3 reward distributor ?
	return payoutRedelegationReward(validator, msg.DelegatorAddress, &handler, st.evm.EpochNumber)
}

func saveNewValidatorToPool(wrapper *restaking.ValidatorWrapper_, validatorPool *restaking.Storage_ValidatorPool_) {
	validatorPool.Validators().Put(wrapper.Validator.ValidatorAddress, wrapper)
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
	epoch *big.Int) (*big.Int, error) {
	redelegation, ok := s.Redelegations().Get(delegator)
	if !ok {
		return nil, errRedelegationNotExist
	}

	r := redelegation.Reward().Value()
	if r.Cmp(common.Big0) == 0 {
		return nil, errNoRewardsToCollect
	}
	redelegation.Reward().SetValue(common.Big0)
	if err := handler.HandleReward(s.Validator().ValidatorAddress().Value(), delegator, r, epoch); err != nil {
		return common.Big0, err
	}
	return r, nil
}

type RestakingRewardHandler interface {
	HandleReward(validator, delegator common.Address, reward, epoch *big.Int) error
}

type RewardToBalance struct {
	StateDB vm.StateDB
}

func (r *RewardToBalance) HandleReward(validator, delegator common.Address, reward, epoch *big.Int) error {
	r.StateDB.AddBalance(delegator, reward)
	return nil
}

type participant interface {
	PostCreateValidator(validator common.Address, amount *big.Int) error
	PostRedelegate(validator common.Address, amount *big.Int) error
}

type tokenHolder struct {
	stateDB       vm.StateDB
	holderAddress common.Address
}

func (t tokenHolder) PostCreateValidator(validator common.Address, amount *big.Int) error {
	t.stateDB.SubBalance(t.holderAddress, amount)
	return nil
}

func (t tokenHolder) PostRedelegate(validator common.Address, amount *big.Int) error {
	t.stateDB.SubBalance(t.holderAddress, amount)
	return nil
}

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGasForStaking(data []byte, isValidatorCreation bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if isValidatorCreation {
		gas = params.TxGasValidatorCreation
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
