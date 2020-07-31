package core

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	staking "github.com/ethereum/go-ethereum/staking/types"
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
		stkMsg := &staking.CreateValidator{}
		if err := rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		_, err := VerifyCreateValidatorMsg(pool.currentState, pendingBlockNumber, stkMsg, msg.From())
		return err
	case types.StakeEditVal:
		stkMsg := &staking.EditValidator{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		return VerifyEditValidatorMsg(pool.currentState, chainContext, pendingBlockNumber, pendingBlockNumber, stkMsg, msg.From())
	case types.Redelegate:
		stkMsg := &staking.Redelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		return VerifyRedelegateMsg(pool.currentState, stkMsg, msg.From())
	case types.Unredelegate:
		stkMsg := &staking.Unredelegate{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		return VerifyUnredelegateMsg(pool.currentState, pendingEpoch, stkMsg, msg.From())
	case types.CollectRedelRewards:
		stkMsg := &staking.CollectRedelegationRewards{}
		if err = rlp.DecodeBytes(msg.Data(), stkMsg); err != nil {
			return err
		}
		return VerifyCollectRedelRewardsMsg(pool.currentState, stkMsg, msg.From())
	}
	return nil
}