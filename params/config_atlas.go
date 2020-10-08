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
	BlocksPerHalfingCycle uint64   `json:"halfing_cycle"`
	ScalingCycle          uint64   `json:"scaling_cycle"`
	RestakingEnable       bool     `json:"restaking_enable"`
	ProposerPolicy        uint64   `json:"policy"`                   // The policy for proposer selection
	Ceil2Nby3Block        *big.Int `json:"ceil2Nby3Block,omitempty"` // Number of confirmations required to move from one state to next [2F + 1 to Ceil(2N/3)]
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
	blocks := c.BlocksPerEpoch
	return uint64(math.Ceil(float64(blockNum) / float64(blocks)))
}

func (c *AtlasConfig) String() string {
	return "atlas"
}
