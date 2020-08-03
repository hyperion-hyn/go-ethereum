package params

import "testing"

func TestAtlasConfig_EpochByBlock(t *testing.T) {
	type fields struct {
		Period         uint64
		BlocksPerEpoch uint64
	}
	type args struct {
		blockNum uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint64
	}{
		// TODO(ATLAS): Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AtlasConfig{
				Period:         tt.fields.Period,
				BlocksPerEpoch: tt.fields.BlocksPerEpoch,
			}
			if got := c.EpochByBlock(tt.args.blockNum); got != tt.want {
				t.Errorf("EpochByBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtlasConfig_EpochFirstBlock(t *testing.T) {
	type fields struct {
		Period         uint64
		BlocksPerEpoch uint64
	}
	type args struct {
		epochNum uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint64
	}{
		// TODO(ATLAS): Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AtlasConfig{
				Period:         tt.fields.Period,
				BlocksPerEpoch: tt.fields.BlocksPerEpoch,
			}
			if got := c.EpochFirstBlock(tt.args.epochNum); got != tt.want {
				t.Errorf("EpochFirstBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtlasConfig_EpochLastBlock(t *testing.T) {
	type fields struct {
		Period         uint64
		BlocksPerEpoch uint64
	}
	type args struct {
		epochNum uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint64
	}{
		// TODO(ATLAS): Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AtlasConfig{
				Period:         tt.fields.Period,
				BlocksPerEpoch: tt.fields.BlocksPerEpoch,
			}
			if got := c.EpochLastBlock(tt.args.epochNum); got != tt.want {
				t.Errorf("EpochLastBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtlasConfig_IsFirstBlock(t *testing.T) {
	type fields struct {
		Period         uint64
		BlocksPerEpoch uint64
	}
	type args struct {
		blockNum uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO(ATLAS): Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AtlasConfig{
				Period:         tt.fields.Period,
				BlocksPerEpoch: tt.fields.BlocksPerEpoch,
			}
			if got := c.IsFirstBlock(tt.args.blockNum); got != tt.want {
				t.Errorf("IsFirstBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtlasConfig_IsLastBlock(t *testing.T) {
	type fields struct {
		Period         uint64
		BlocksPerEpoch uint64
	}
	type args struct {
		blockNum uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO(ATLAS): Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AtlasConfig{
				Period:         tt.fields.Period,
				BlocksPerEpoch: tt.fields.BlocksPerEpoch,
			}
			if got := c.IsLastBlock(tt.args.blockNum); got != tt.want {
				t.Errorf("IsLastBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}