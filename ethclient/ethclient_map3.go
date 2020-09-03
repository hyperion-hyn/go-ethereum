package ethclient

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"math/big"
)

func (ec *Client) GetAllMap3NodeAddresses(ctx context.Context, blockNumber *big.Int) ([]string, error) {
	var result []string
	err := ec.c.CallContext(ctx, &result, "eth_getAllMap3NodeAddresses", toBlockNumArg(blockNumber))
	return result, err
}

func (ec *Client) GetMap3NodeInformation(
	ctx context.Context, map3NodeAddress common.Address, blockNumber *big.Int) (*microstaking.Map3NodeWrapperRPC, error) {

	var bytes hexutil.Bytes

	err := ec.c.CallContext(ctx, &bytes, "eth_getMap3NodeInformation", map3NodeAddress, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	var map3Node microstaking.Map3NodeWrapperRPC
	err = rlp.DecodeBytes(bytes, &map3Node)
	if err != nil {
		return nil, err
	}
	return &map3Node, err

}
