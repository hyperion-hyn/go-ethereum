package ethclient

import (
	"context"
	"math/big"
)

func (ec *Client) GetAllValidatorAddresses(ctx context.Context, blockNumber *big.Int) ([]string, error) {
	var result []string
	err := ec.c.CallContext(ctx, &result, "eth_getAllValidatorAddresses", toBlockNumArg(blockNumber))
	return result, err
}
