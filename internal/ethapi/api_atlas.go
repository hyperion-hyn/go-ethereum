package ethapi

import (
    "context"
    "github.com/pkg/errors"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/log"
    "github.com/ethereum/go-ethereum/rlp"
    "github.com/ethereum/go-ethereum/core"

    staking "github.com/ethereum/go-ethereum/staking/types"
)

// SubmitStakingTransaction is a helper function that submits tx to txPool and logs a message.
func SubmitStakingTransaction(ctx context.Context, b Backend, tx *staking.StakingTransaction) (common.Hash, error) {
    if err := b.SendStakingTx(ctx, tx); err != nil {
        return common.Hash{}, err
    }
    log.Info("Submitted Staking transaction", "fullhash", tx.Hash().Hex())
    return tx.Hash(), nil
}

// SendRawStakingTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTransactionPoolAPI) SendRawStakingTransaction(
    ctx context.Context, encodedTx hexutil.Bytes,
) (common.Hash, error) {
    if len(encodedTx) >= staking.MaxEncodedPoolTransactionSize {
        err := errors.Wrapf(core.ErrOversizedData, "encoded tx size: %d", len(encodedTx))
        return common.Hash{}, err
    }
    tx := new(staking.StakingTransaction)
    if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
        return common.Hash{}, err
    }
    c := s.b.ChainConfig().ChainID
    if id := tx.ChainID(); id.Cmp(c) != 0 {
        return common.Hash{}, errors.Wrapf(
            errInvalidChainID, "blockchain chain id:%s, given %s", c.String(), id.String(),
        )
    }
    return SubmitStakingTransaction(ctx, s.b, tx)
}

