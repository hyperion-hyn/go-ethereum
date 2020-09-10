package core

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (pool *TxPool) validateStakingTx(tx *types.Transaction) error {
	currentBlockNumber := pool.chain.CurrentBlock().Number()
	pendingBlockNumber := new(big.Int).Add(currentBlockNumber, big.NewInt(1))
	pendingEpoch := pool.chain.CurrentBlock().Epoch()
	if pool.chainconfig.Atlas.IsLastBlock(currentBlockNumber.Uint64()) {
		pendingEpoch = new(big.Int).Add(pendingEpoch, big.NewInt(1))
	}
	chainContext, ok := pool.chain.(ChainContext)
	if !ok {
		chainContext = nil // might use testing blockchain, set to nil for verifier to handle.
	}
	msg, err := tx.AsMessage(types.MakeSigner(pool.chainconfig, pendingBlockNumber))
	if err != nil {
		return err
	}

	verifier, err := NewStakingVerifier(chainContext)
	if err != nil {
		return err
	}

	switch msg.Type() {
	case types.CreateValidator:
		stkMsg := &restaking.CreateValidator{}
		if err := rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := verifier.VerifyCreateValidatorMsg(pool.currentState, pendingBlockNumber, stkMsg, msg.From())
		return err
	case types.EditValidator:
		stkMsg := &restaking.EditValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := verifier.VerifyEditValidatorMsg(pool.currentState, chainContext, pendingEpoch, pendingBlockNumber, stkMsg, msg.From())
		return err
	case types.Redelegate:
		stkMsg := &restaking.Redelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := verifier.VerifyRedelegateMsg(pool.currentState, stkMsg, msg.From())
		return err
	case types.Unredelegate:
		stkMsg := &restaking.Unredelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := verifier.VerifyUnredelegateMsg(pool.currentState, pendingEpoch, stkMsg, msg.From())
		return err
	case types.CollectRestakingReward:
		stkMsg := &restaking.CollectReward{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := verifier.VerifyCollectRestakingRewardMsg(pool.currentState, stkMsg, msg.From())
		return err
	case types.CreateMap3:
		stkMsg := &microstaking.CreateMap3Node{}
		if err := rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := verifier.VerifyCreateMap3NodeMsg(pool.currentState, chainContext, pendingEpoch, pendingBlockNumber, stkMsg, msg.From())
		return err
	case types.EditMap3:
		stkMsg := &microstaking.EditMap3Node{}
		if err := rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		return verifier.VerifyEditMap3NodeMsg(pool.currentState, pendingEpoch, pendingBlockNumber, stkMsg, msg.From())
	case types.TerminateMap3:
		stkMsg := &microstaking.TerminateMap3Node{}
		if err := rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		return verifier.VerifyTerminateMap3NodeMsg(pool.currentState, pendingEpoch, stkMsg, msg.From())
	case types.Microdelegate:
		stkMsg := &microstaking.Microdelegate{}
		if err := rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		return verifier.VerifyMicrodelegateMsg(pool.currentState, chainContext, pendingBlockNumber, stkMsg, msg.From())
	case types.Unmicrodelegate:
		stkMsg := &microstaking.Unmicrodelegate{}
		if err := rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		return verifier.VerifyUnmicrodelegateMsg(pool.currentState, chainContext, pendingBlockNumber, pendingEpoch, stkMsg, msg.From())
	case types.CollectMicrostakingRewards:
		stkMsg := &microstaking.CollectRewards{}
		if err := rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		return verifier.VerifyCollectMicrostakingRewardsMsg(pool.currentState, stkMsg, msg.From())
	}

	return nil
}
