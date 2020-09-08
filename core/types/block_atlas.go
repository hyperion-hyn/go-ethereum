package types

import (
    "math/big"
)

const (
    // 5760 * 15 = 86400 seconds = 1 day
    blocksPerEpoch = 5760
)

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (h *Header) Epoch() *big.Int {
    return new(big.Int).Div(h.Number, new(big.Int).SetUint64(blocksPerEpoch))
}
