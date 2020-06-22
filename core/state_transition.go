// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	staking "github.com/ethereum/go-ethereum/staking/types"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

var (
	errInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all the necessary work to work out a valid new state root.

1) Nonce handling
2) Pre pay gas
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==
  4a) Attempt to run transaction data
  4b) If valid, use result as code for the new state object
== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	gp         *GasPool
	msg        Message
	gas        uint64
	gasPrice   *big.Int
	initialGas uint64
	value      *big.Int
	data       []byte
	state      vm.StateDB
	evm        *vm.EVM
}

// Message represents a message sent to a contract.
type Message interface {
	From() common.Address
	//FromFrontier() (common.Address, error)
	To() *common.Address

	GasPrice() *big.Int
	Gas() uint64
	Value() *big.Int

	Nonce() uint64
	CheckNonce() bool
	Data() []byte
	Type() types.TransactionType	// ATLAS
}

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, contractCreation, isHomestead bool, isEIP2028 bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if contractCreation && isHomestead {
		gas = params.TxGasContractCreation
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
		nonZeroGas := params.TxDataNonZeroGasFrontier
		if isEIP2028 {
			nonZeroGas = params.TxDataNonZeroGasEIP2028
		}
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

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *vm.EVM, msg Message, gp *GasPool) *StateTransition {
	return &StateTransition{
		gp:       gp,
		evm:      evm,
		msg:      msg,
		gasPrice: msg.GasPrice(),
		value:    msg.Value(),
		data:     msg.Data(),
		state:    evm.StateDB,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(evm *vm.EVM, msg Message, gp *GasPool) ([]byte, uint64, bool, error) {
	return NewStateTransition(evm, msg, gp).TransitionDb()
}

// ATLAS: ApplyStakingMessage computes the new state for staking message
func ApplyStakingMessage(evm *vm.EVM, msg Message, gp *GasPool) (uint64, error) {
	return NewStateTransition(evm, msg, gp).StakingTransitionDb()
}

// to returns the recipient of the message.
func (st *StateTransition) to() common.Address {
	if st.msg == nil || st.msg.To() == nil /* contract creation */ {
		return common.Address{}
	}
	return *st.msg.To()
}

func (st *StateTransition) useGas(amount uint64) error {
	if st.gas < amount {
		return vm.ErrOutOfGas
	}
	st.gas -= amount

	return nil
}

func (st *StateTransition) buyGas() error {
	mgval := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.Gas()), st.gasPrice)
	if st.state.GetBalance(st.msg.From()).Cmp(mgval) < 0 {
		return errInsufficientBalanceForGas
	}
	if err := st.gp.SubGas(st.msg.Gas()); err != nil {
		return err
	}
	st.gas += st.msg.Gas()

	st.initialGas = st.msg.Gas()
	st.state.SubBalance(st.msg.From(), mgval)
	return nil
}

func (st *StateTransition) preCheck() error {
	// Make sure this transaction's nonce is correct.
	if st.msg.CheckNonce() {
		nonce := st.state.GetNonce(st.msg.From())
		if nonce < st.msg.Nonce() {
			return ErrNonceTooHigh
		} else if nonce > st.msg.Nonce() {
			return ErrNonceTooLow
		}
	}
	return st.buyGas()
}

// TransitionDb will transition the state by applying the current message and
// returning the result including the used gas. It returns an error if failed.
// An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, usedGas uint64, failed bool, err error) {
	if err = st.preCheck(); err != nil {
		return
	}
	msg := st.msg
	sender := vm.AccountRef(msg.From())
	homestead := st.evm.ChainConfig().IsHomestead(st.evm.BlockNumber)
	istanbul := st.evm.ChainConfig().IsIstanbul(st.evm.BlockNumber)
	contractCreation := msg.To() == nil

	// Pay intrinsic gas
	gas, err := IntrinsicGas(st.data, contractCreation, homestead, istanbul)
	if err != nil {
		return nil, 0, false, err
	}
	if err = st.useGas(gas); err != nil {
		return nil, 0, false, err
	}

	var (
		evm = st.evm
		// vm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error.
		vmerr error
	)
	if contractCreation {
		ret, _, st.gas, vmerr = evm.Create(sender, st.data, st.gas, st.value)
	} else {
		// Increment the nonce for the next transaction
		st.state.SetNonce(msg.From(), st.state.GetNonce(sender.Address())+1)
		ret, st.gas, vmerr = evm.Call(sender, st.to(), st.data, st.gas, st.value)
	}
	if vmerr != nil {
		log.Debug("VM returned with error", "err", vmerr)
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		if vmerr == vm.ErrInsufficientBalance {
			return nil, 0, false, vmerr
		}
	}
	st.refundGas()
	st.state.AddBalance(st.evm.Coinbase, new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice))

	return ret, st.gasUsed(), vmerr != nil, err
}

func (st *StateTransition) refundGas() {
	// Apply refund counter, capped to half of the used gas.
	refund := st.gasUsed() / 2
	if refund > st.state.GetRefund() {
		refund = st.state.GetRefund()
	}
	st.gas += refund

	// Return ETH for remaining gas, exchanged at the original rate.
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
	st.state.AddBalance(st.msg.From(), remaining)

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(st.gas)
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gas
}

// ATLAS
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
	case types.StakeEditNode:
	case types.Microdelegate:
	case types.Unmicrodelegate:
	case types.CollectMicrodelRewards:
	case types.StakeCreateVal:
		stkMsg := &staking.CreateValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		//utils.Logger().Info().
		//	Msgf("[DEBUG STAKING] staking type: %s, gas: %d, txn: %+v", msg.Type(), gas, stkMsg)
		//if msg.From() != stkMsg.ValidatorAddress {
		//	return 0, errInvalidSigner
		//}
		err = st.verifyAndApplyCreateValidatorTx(stkMsg, st.evm.BlockNumber)
	case types.StakeEditVal:
		stkMsg := &staking.EditValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		//utils.Logger().Info().
		//	Msgf("[DEBUG STAKING] staking type: %s, gas: %d, txn: %+v", msg.Type(), gas, stkMsg)
		//if msg.From() != stkMsg.ValidatorAddress {
		//	return 0, errInvalidSigner
		//}
		err = st.verifyAndApplyEditValidatorTx(stkMsg, st.evm.BlockNumber)
	case types.Redelegate:
		stkMsg := &staking.Redelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		//utils.Logger().Info().Msgf("[DEBUG STAKING] staking type: %s, gas: %d, txn: %+v", msg.Type(), gas, stkMsg)
		//if msg.From() != stkMsg.MicrodelegatorAddress {
		//	return 0, errInvalidSigner
		//}
		err = st.verifyAndApplyDelegateTx(stkMsg)
	case types.Unredelegate:
		stkMsg := &staking.Unredelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		//utils.Logger().Info().Msgf("[DEBUG STAKING] staking type: %s, gas: %d, txn: %+v", msg.Type(), gas, stkMsg)
		//if msg.From() != stkMsg.MicrodelegatorAddress {
		//	return 0, errInvalidSigner
		//}
		err = st.verifyAndApplyUndelegateTx(stkMsg)
	case types.CollectMicroreRewards:
		stkMsg := &staking.CollectRedelegationRewards{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return 0, err
		}
		//utils.Logger().Info().Msgf("[DEBUG STAKING] staking type: %s, gas: %d, txn: %+v", msg.Type(), gas, stkMsg)
		//if msg.From() != stkMsg.MicrodelegatorAddress {
		//	return 0, errInvalidSigner
		//}
		collectedRewards, tempErr := st.verifyAndApplyCollectRewards(stkMsg)
		err = tempErr
		if err == nil {
			st.state.AddLog(&types.Log{
				Address:     stkMsg.DelegatorAddress,
				Topics:      []common.Hash{staking2.CollectRewardsTopic},
				Data:        collectedRewards.Bytes(),
				BlockNumber: st.evm.BlockNumber.Uint64(),
			})
		}
	default:
		return 0, staking.ErrInvalidStakingKind
	}
	st.refundGas()

	// TODO(ATLAS) Txn Fees
	//txFee := new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice)
	//st.state.AddBalance(st.evm.Coinbase, txFee)

	return st.gasUsed(), err
}

func (st *StateTransition) verifyAndApplyCreateMap3NodeTx(createMap3Node *staking.CreateMap3Node) error {
	return nil
}

func (st *StateTransition) verifyAndApplyEditMap3NodeTx(editMap3Node *staking.EditMap3Node) error {
	return nil
}

func (st *StateTransition) verifyAndApplyMicrodelegateTx(microdelegate *staking.Microdelegate) error {
	return nil
}

func (st *StateTransition) verifyAndApplyUnmicrodelegateTx(unmicrodelegate *staking.Unmicrodelegate) error {
	return nil
}

func (st *StateTransition) verifyAndApplyCollectMicrodelRewardsTx(
	collectMicrodelegationRewards *staking.CollectMicrodelegationRewards,
) error {
	return nil
}




func (st *StateTransition) verifyAndApplyCreateValidatorTx(
	createValidator *staking.CreateValidator, blockNum *big.Int,
) error {
	wrapper, err := VerifyAndCreateValidatorFromMsg(
		st.state, st.evm.EpochNumber, blockNum, createValidator,
	)
	if err != nil {
		return err
	}
	// TODO(Storage)
	if err := st.state.UpdateValidator(wrapper); err != nil {
		return err
	}
	// TODO(ATLAS): update node state
	return nil
}

func (st *StateTransition) verifyAndApplyEditValidatorTx(
	editValidator *staking.EditValidator, blockNum *big.Int,
) error {
	wrapper, err := VerifyAndEditValidatorFromMsg(
		st.state, st.bc, st.evm.EpochNumber, blockNum, editValidator,
	)
	if err != nil {
		return err
	}
	return st.state.UpdateValidatorWrapper(wrapper.Address, wrapper)
}

func (st *StateTransition) verifyAndApplyDelegateTx(delegate *staking.Redelegate) error {
	wrapper, balanceToBeDeducted, err := VerifyAndDelegateFromMsg(st.state, delegate)
	if err != nil {
		return err
	}

	st.state.SubBalance(delegate.DelegatorAddress, balanceToBeDeducted)

	return st.state.UpdateValidatorWrapper(wrapper.Address, wrapper)
}

func (st *StateTransition) verifyAndApplyUndelegateTx(
	undelegate *staking.Unredelegate,
) error {
	wrapper, err := VerifyAndUndelegateFromMsg(st.state, st.evm.EpochNumber, undelegate)
	if err != nil {
		return err
	}
	return st.state.UpdateValidatorWrapper(wrapper.Address, wrapper)
}

func (st *StateTransition) verifyAndApplyCollectRewards(collectRewards *staking.CollectRedelegationRewards) (*big.Int, error) {
	if st.bc == nil {
		return network.NoReward, errors.New("[CollectRedelegationRewards] No chain context provided")
	}
	delegations, err := st.bc.ReadDelegationsByDelegator(collectRewards.DelegatorAddress)
	if err != nil {
		return network.NoReward, err
	}
	updatedValidatorWrappers, totalRewards, err := VerifyAndCollectRewardsFromDelegation(
		st.state, delegations,
	)
	if err != nil {
		return network.NoReward, err
	}
	for _, wrapper := range updatedValidatorWrappers {
		if err := st.state.UpdateValidatorWrapper(wrapper.Address, wrapper); err != nil {
			return network.NoReward, err
		}
	}
	st.state.AddBalance(collectRewards.DelegatorAddress, totalRewards)
	return totalRewards, nil
}

// ATLAS - END