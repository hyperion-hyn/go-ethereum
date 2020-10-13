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
	ctx context.Context, map3NodeAddress common.Address, blockNumber *big.Int) (*microstaking.PlainMap3NodeWrapper, error) {

	var bytes hexutil.Bytes

	err := ec.c.CallContext(ctx, &bytes, "eth_getMap3NodeInformation", map3NodeAddress, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	var map3Node microstaking.PlainMap3NodeWrapper
	err = rlp.DecodeBytes(bytes, &map3Node)
	if err != nil {
		return nil, err
	}
	return &map3Node, err

}

func (ec *Client) GetMap3NodeDelegation(
	ctx context.Context, map3NodeAddress common.Address, delegatorAddress common.Address, blockNumber *big.Int) (*microstaking.Microdelegation_, error) {

	var bytes hexutil.Bytes

	err := ec.c.CallContext(ctx, &bytes, "eth_getMap3NodeDelegation", map3NodeAddress, delegatorAddress, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	var redelegation microstaking.Microdelegation_
	err = rlp.DecodeBytes(bytes, &redelegation)
	if err != nil {
		return nil, err
	}
	return &redelegation, err
}

func (ec *Client) GetAllMap3RewardByDelegatorAddress(ctx context.Context, delegatorAddress common.Address, blockNumber *big.Int) (*big.Int, error) {
	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "eth_getAllMap3RewardByDelegatorAddress", delegatorAddress, toBlockNumArg(blockNumber))
	return (*big.Int)(&result), err
}
