package microstaking

import "github.com/ethereum/go-ethereum/crypto"

const (
	collectRewardsStr = "Microstaking/CollectRewards"
)

// keys used to retrieve staking related information
var (
	CollectRewardsTopic = crypto.Keccak256Hash([]byte(collectRewardsStr))
)