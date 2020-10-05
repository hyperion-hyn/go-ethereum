package restaking

import "github.com/ethereum/go-ethereum/crypto"

const (
	collectRewardStr = "Restaking/CollectReward"
)

// keys used to retrieve staking related information
var (
	CollectRewardTopic = crypto.Keccak256Hash([]byte(collectRewardStr))
)