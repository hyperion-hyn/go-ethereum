package core

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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

	switch msg.Type() {
	case types.StakeCreateVal:
		stkMsg := &restaking.CreateValidator{}
		if err := rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := VerifyCreateValidatorMsg(pool.currentState, pendingBlockNumber, stkMsg, msg.From())
		return err
	case types.StakeEditVal:
		stkMsg := &restaking.EditValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := VerifyEditValidatorMsg(pool.currentState, chainContext, pendingEpoch, pendingBlockNumber, stkMsg, msg.From())
		return err
	case types.Redelegate:
		stkMsg := &restaking.Redelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := VerifyRedelegateMsg(pool.currentState, stkMsg, msg.From())
		return err
	case types.Unredelegate:
		stkMsg := &restaking.Unredelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := VerifyUnredelegateMsg(pool.currentState, pendingEpoch, stkMsg, msg.From())
		return err
	case types.CollectRedelRewards:
		stkMsg := &restaking.CollectReward{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := VerifyCollectRedelRewardsMsg(pool.currentState, stkMsg, msg.From())
		return err
	}
	return nil
}
