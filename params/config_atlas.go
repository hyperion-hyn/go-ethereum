package params

import "math"

type AtlasConfig struct {
	Period                uint64 `json:"period"` // Number of seconds between blocks to enforce
	BlocksPerEpoch        uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
	BlocksPerHalfingCycle uint64 `json:"halfing_cycle"`
}

func (c *AtlasConfig) IsLastBlock(blockNum uint64) bool {
	if blockNum == 0 {
		return false
	}
	blocks := c.BlocksPerEpoch
	return blockNum % blocks == 0
}

func (c *AtlasConfig) IsFirstBlock(blockNum uint64) bool {
	blocks := c.BlocksPerEpoch
	return blockNum % blocks == 1
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
	return blocks * (epochNum - 1) + 1
}

func (c *AtlasConfig) EpochByBlock(blockNum uint64) uint64 {
	blocks := c.BlocksPerEpoch
	return uint64(math.Ceil(float64(blockNum) / float64(blocks)))
}

