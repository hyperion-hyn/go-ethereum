package ethclient

import (
	"context"
	"math/big"
)

func (ec *Client) GetAllMap3NodeAddresses(ctx context.Context, blockNumber *big.Int) ([]string, error) {
	var result []string
	err := ec.c.CallContext(ctx, &result, "eth_getAllMap3NodeAddresses", toBlockNumArg(blockNumber))
	return result, err
}
