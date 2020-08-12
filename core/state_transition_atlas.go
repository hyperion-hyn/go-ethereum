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
	case types.StakeCreateVal:
		stkMsg := &restaking.CreateValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyCreateValidatorTx(stkMsg, msg.From())
	case types.StakeEditVal:
		stkMsg := &restaking.EditValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyEditValidatorTx(stkMsg, msg.From())
	case types.Redelegate:
		stkMsg := &restaking.Redelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyRedelegateTx(stkMsg, msg.From())
	case types.Unredelegate:
		stkMsg := &restaking.Unredelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		err = st.verifyAndApplyUnredelegateTx(stkMsg, msg.From())
	case types.CollectRedelRewards:
		stkMsg := &restaking.CollectRedelegationRewards{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		_, err := st.verifyAndApplyCollectRedelRewards(stkMsg, msg.From())
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

func (st *StateTransition) verifyAndApplyCreateValidatorTx(msg *restaking.CreateValidator, signer common.Address) error {
	v, err := VerifyCreateValidatorMsg(st.state, st.evm.BlockNumber, msg, signer)
	if err != nil {
		return err
	}
	wrapper := &restaking.ValidatorWrapper_{}
	wrapper.Validator = *v
	wrapper.TotalDelegation = big.NewInt(0).Set(defaultStakingAmount)
	wrapper.TotalDelegationByOperator = big.NewInt(0).Set(defaultStakingAmount)

	validatorPool := st.state.ValidatorPool()
	validatorPool.Validators().Put(v.ValidatorAddress, wrapper)
	wrapperSt, _ := validatorPool.Validators().Get(v.ValidatorAddress)
	wrapperSt.Redelegations().Put(msg.OperatorAddress, &restaking.Redelegation_{
		DelegatorAddress: msg.OperatorAddress,
		Amount:           defaultStakingAmount,
	})

	keySet := validatorPool.SlotKeySet()
	for _, key := range wrapper.Validator.SlotPubKeys.Keys {
		keySet.Get(key.Hex()).SetValue(true)
	}
	if msg.Description.Identity != "" {
		validatorPool.DescriptionIdentitySet().Get(msg.Description.Identity).SetValue(true)
	}

	st.state.SubBalance(msg.OperatorAddress, defaultStakingAmount)

	return nil
}

func (st *StateTransition) verifyAndApplyEditValidatorTx(msg *restaking.EditValidator, signer common.Address) error {
	if err := VerifyEditValidatorMsg(st.state, st.bc, st.evm.EpochNumber, st.evm.BlockNumber, msg, signer); err != nil {
		return err
	}

	validatorPool := st.state.ValidatorPool()
	wrapper, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)

	// update description
	if msg.Description.Identity != "" {
		i := wrapper.Validator().Description().Identity().Value()
		validatorPool.DescriptionIdentitySet().Get(i).SetValue(false)
		validatorPool.DescriptionIdentitySet().Get(msg.Description.Identity).SetValue(true)
	}
	wrapper.Validator().Description().UpdateDescription(msg.Description)

	if !msg.CommissionRate.IsNil() {
		wrapper.Validator().Commission().CommissionRates().Rate().SetValue(*msg.CommissionRate)
		wrapper.Validator().Commission().UpdateHeight().SetValue(st.evm.BlockNumber)
	}

	if msg.SlotKeyToRemove != nil {
		for i := 0; i < wrapper.Validator().SlotPubKeys().Length(); i++ {
			if *msg.SlotKeyToRemove == *wrapper.Validator().SlotPubKeys().Get(i) {
				wrapper.Validator().SlotPubKeys().Remove(i, false)
				validatorPool.SlotKeySet().Get(msg.SlotKeyToRemove.Hex()).SetValue(false)
				break
			}
		}
	}

	if msg.SlotKeyToAdd != nil {
		wrapper.Validator().SlotPubKeys().Push(msg.SlotKeyToAdd)
		validatorPool.SlotKeySet().Get(msg.SlotKeyToAdd.Hex()).SetValue(true)
	}

	if msg.EPOSStatus == restaking.Active || msg.EPOSStatus == restaking.Inactive {
		wrapper.Validator().Status().SetValue(uint8(msg.EPOSStatus))
	}
	return nil
}

func (st *StateTransition) verifyAndApplyRedelegateTx(msg *restaking.Redelegate, signer common.Address) error {
	if err := VerifyRedelegateMsg(st.state, msg, signer); err != nil {
		return err
	}

	wrapper, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	if redelegation, ok := wrapper.Redelegations().Get(msg.DelegatorAddress); ok {
		amt := redelegation.Amount().Value()
		redelegation.Amount().SetValue(big.NewInt(0).Add(amt, defaultStakingAmount))
	} else {
		m := restaking.NewRedelegation(msg.DelegatorAddress, defaultStakingAmount)
		wrapper.Redelegations().Put(msg.DelegatorAddress, m)
	}
	wrapper.AddTotalDelegation(defaultStakingAmount)
	if isOperator(wrapper, msg.DelegatorAddress) {
		wrapper.AddTotalDelegationByOperator(defaultStakingAmount)
	}
	st.state.SubBalance(signer, defaultStakingAmount)
	return nil
}

func (st *StateTransition) verifyAndApplyUnredelegateTx(msg *restaking.Unredelegate, signer common.Address) error {
	if err := VerifyUnredelegateMsg(st.state, st.evm.EpochNumber, msg, signer); err != nil {
		return err
	}

	validator, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	redelegation, _ := validator.Redelegations().Get(msg.DelegatorAddress)
	amount := redelegation.Amount().Value()
	redelegation.Undelegation().Amount().SetValue(amount)
	redelegation.Undelegation().Epoch().SetValue(st.evm.EpochNumber)
	redelegation.Amount().SetValue(common.Big0)
	validator.SubTotalDelegation(amount)
	if isOperator(validator, msg.DelegatorAddress) {
		validator.SubTotalDelegationByOperator(amount)
		// TODO: need 20% ?
	}
	return nil
}

func (st *StateTransition) verifyAndApplyCollectRedelRewards(msg *restaking.CollectRedelegationRewards, signer common.Address) (*big.Int, error) {
	if err := VerifyCollectRedelRewardsMsg(st.state, msg, signer); err != nil {
		return network.NoReward, err
	}
	validator, _ := st.state.ValidatorByAddress(msg.ValidatorAddress)
	redelegation, _ := validator.Redelegations().Get(msg.DelegatorAddress)
	reward := redelegation.Reward().Value()
	handler := RewardToBalance{stateDB: st.state}
	if err := handler.HandleReward(msg.DelegatorAddress, reward, st.evm.EpochNumber); err != nil {
		return nil, err
	}
	return reward, nil
}

type RewardHandler interface {
	HandleReward(delegator common.Address, amount, epoch *big.Int) error
}

type RewardToBalance struct {
	stateDB vm.StateDB
}

func (r *RewardToBalance) HandleReward(delegator common.Address, amount, epoch *big.Int) error {
	r.stateDB.AddBalance(delegator, amount)
	return nil
}

func isOperator(validator *restaking.Storage_ValidatorWrapper_, delegator common.Address) bool {
	if validator.Validator().OperatorAddresses().Set().Get(delegator).Value() {
		return true
	}
	return false
}
