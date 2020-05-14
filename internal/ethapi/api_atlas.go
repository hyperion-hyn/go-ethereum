package ethapi

import (
    "context"

    "github.com/pkg/errors"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/core"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/rlp"

    staking "github.com/ethereum/go-ethereum/staking/types"
)

// SendRawStakingTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTransactionPoolAPI) SendRawStakingTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {
    if len(encodedTx) >= staking.MaxEncodedPoolTransactionSize {
        err := errors.Wrapf(core.ErrOversizedData, "encoded tx size: %d", len(encodedTx))
        return common.Hash{}, err
    }
    tx := new(types.Transaction)
    if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
        return common.Hash{}, err
    }

    switch(tx.Type()) {
    case types.StakeCreateVal, types.StakeEditVal, types.Delegate, types.Undelegate, types.CollectRewards:
    default:
        err := errors.Wrapf(errors.New("invalid transaction type"),  "got: %d", tx.Type())
        return common.Hash{}, err
    }

    return SubmitTransaction(ctx, s.b, tx)
}

