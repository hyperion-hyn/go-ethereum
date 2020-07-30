package reward

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/committee"
	"github.com/ethereum/go-ethereum/staking/types"
	"math/big"
)

// Payout ..
type Payout struct {
	Addr        common.Address
	NewlyEarned *big.Int
	EarningKey  types.BLSPublicKey
}

// CompletedRound ..
type CompletedRound struct {
	Total *big.Int
	Award []Payout
}

// Reader ..
type Reader interface {
	ReadRoundResult() *CompletedRound
	MissingSigners() committee.SlotList
}
