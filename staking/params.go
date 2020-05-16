package staking

import (
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	isValidatorKeyStr = "Atlas/IsValidator/Key/v1"
	isValidatorStr    = "Atlas/IsValidator/Value/v1"
	collectRewardsStr = "Atlas/CollectRewards"
)

// keys used to retrieve staking related informatio
var (
	IsValidatorKey      = crypto.Keccak256Hash([]byte(isValidatorKeyStr))
	IsValidator         = crypto.Keccak256Hash([]byte(isValidatorStr))
	CollectRewardsTopic = crypto.Keccak256Hash([]byte(collectRewardsStr))
)
