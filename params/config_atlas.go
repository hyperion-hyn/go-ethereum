package params

type AtlasConfig struct {
	Period         uint64 `json:"period"` // Number of seconds between blocks to enforce
	BlocksPerEpoch uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

func (c *AtlasConfig) IsLastBlock(blockNum uint64) bool {
	blocks := c.BlocksPerEpoch
	return blockNum % blocks == blocks - 1
}

func (c *AtlasConfig) EpochLastBlock(epochNum uint64) uint64 {
	blocks := c.BlocksPerEpoch
	if epochNum == 0 {
		return 0
	}
	return blocks * epochNum
}

func (c *AtlasConfig) EpochFirstBlock(epochNum uint64) uint64 {
	return 0
}

func (c *AtlasConfig) EpochByBlock(blockNum uint64) uint64 {
	return 0
}