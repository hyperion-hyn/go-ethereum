package staking

import (
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	isValidatorKeyStr = "Atlas/IsValidator/Key/v1"
	isValidatorStr    = "Atlas/IsValidator/Value/v1"
)

// keys used to retrieve staking related informatio
var (
	IsValidatorKey = crypto.Keccak256Hash([]byte(isValidatorKeyStr))
	IsValidator    = crypto.Keccak256Hash([]byte(isValidatorStr))
)
