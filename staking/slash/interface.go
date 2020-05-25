package slash

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/params"
	"github.com/ethereum/go-ethereum/shard"
)

// CommitteeReader ..
type CommitteeReader interface {
	Config() *params.ChainConfig
	ReadShardState(epoch *big.Int) (*shard.State, error)
	CurrentBlock() *types.Block
}
