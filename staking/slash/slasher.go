package slash

import (
	"github.com/ethereum/go-ethereum/staking/types"
)

const (
	// UnavailabilityInConsecutiveBlockSigning is how many blocks in a row
	// before "slashing by unavailability" occurs
	UnavailabilityInConsecutiveBlockSigning = 1380
)

// Slasher ..
type Slasher interface {
	ShouldSlash(types.BlsPublicKey) bool
}

// ThresholdDecider ..
type ThresholdDecider interface {
	SlashThresholdMet(types.BlsPublicKey) bool
}
