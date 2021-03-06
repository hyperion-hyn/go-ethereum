package params

import (
	"math"
	"math/big"
)

// AtlasConfig is the consensus engine configs for Istanbul based sealing.
type AtlasConfig struct {
	RequestTimeout        uint64   `json:"requestTimeout"` // The timeout for each Atlas round in milliseconds.
	Period                uint64   `json:"period"`         // Number of seconds between blocks to enforce
	BlocksPerEpoch        uint64   `json:"epoch"`          // Epoch length to reset votes and checkpoint
	BlocksPerHalfingCycle uint64   `json:"halfingCycle"`
	ScalingCycle          uint64   `json:"scalingCycle"`
	RestakingEnable       bool     `json:"restakingEnable"`
	ProposerPolicy        uint64   `json:"policy"`                   // The policy for proposer selection
	Ceil2Nby3Block        *big.Int `json:"ceil2Nby3Block,omitempty"` // Number of confirmations required to move from one state to next [2F + 1 to Ceil(2N/3)]

	// HIP config
	Map3MigrationBlock              uint64   `json:"map3MigrationBlock"`
	Map3NodeAgeDeadlineBlock        uint64   `json:"map3NodeAgeDeadlineBlock"`
	MicrodelegationIndexRepairBlock uint64   `json:"microdelegationIndexRepairBlock"`
	MicrostakingImproveBlock        uint64   `json:"microstakingImproveBlock"`
	HYNBurningBlock                 uint64   `json:"hynBurningBlock"`
	Athens                          uint64   `json:"athens"`
	ChainID                         *big.Int `json:"chainID"` // ATLAS: atlas chain id is used for replay protection in other chain compatible with Ethereum
}

func (c *AtlasConfig) IsLastBlock(blockNum uint64) bool {
	if blockNum == 0 {
		return true
	}
	blocks := c.BlocksPerEpoch
	return blockNum%blocks == 0
}

func (c *AtlasConfig) IsFirstBlock(blockNum uint64) bool {
	if blockNum == 0 {
		return true
	}
	blocks := c.BlocksPerEpoch
	return blockNum%blocks == 1
}

func (c *AtlasConfig) EpochLastBlock(epochNum uint64) uint64 {
	if epochNum == 0 {
		return 0
	}
	blocks := c.BlocksPerEpoch
	return blocks * epochNum
}

func (c *AtlasConfig) EpochFirstBlock(epochNum uint64) uint64 {
	if epochNum == 0 {
		return 0
	}
	blocks := c.BlocksPerEpoch
	return blocks*(epochNum-1) + 1
}

func (c *AtlasConfig) EpochByBlock(blockNum uint64) uint64 {
	blocksPerEpoch := c.BlocksPerEpoch
	// https://stackoverflow.com/questions/2745074/fast-ceiling-of-an-integer-division-in-c-c
	// epoch = 1 + ((blockNum-1)/blocksPerEpoch)
	return uint64(math.Ceil(float64(blockNum) / float64(blocksPerEpoch)))
}

func (c *AtlasConfig) EpochOfBlock(blockNum uint64) (epoch, firstBlock, lastBlock uint64) {
	epoch = c.EpochByBlock(blockNum)
	firstBlock = c.EpochFirstBlock(epoch)
	lastBlock = c.EpochLastBlock(epoch)
	return
}

func (c *AtlasConfig) String() string {
	return "atlas"
}

func (c *AtlasConfig) IsMicrostakingImprove(num *big.Int) bool {
	return isForked(big.NewInt(int64(c.MicrostakingImproveBlock)), num)
}

func (c *AtlasConfig) IsHYNBurning(num *big.Int) bool {
	return isForked(big.NewInt(int64(c.HYNBurningBlock)), num)
}

func (c *AtlasConfig) IsAthens(num *big.Int) bool {
	return isForked(big.NewInt(int64(c.Athens)), num)
}
