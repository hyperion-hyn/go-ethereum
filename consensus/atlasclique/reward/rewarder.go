package reward

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

// Payout ..
type Payout struct {
	Addr        common.Address
	NewlyEarned *big.Int
	EarningKey  restaking.BLSPublicKey_
}

// CompletedRound ..
type CompletedRound struct {
	Total *big.Int
	Award []Payout
}

// Reader ..
type Reader interface {
	ReadRoundResult() *CompletedRound
	MissingSigners() restaking.Slots_
}
