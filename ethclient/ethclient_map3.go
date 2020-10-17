package ethclient

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

	var map3Node microstaking.PlainMap3NodeWrapper
	err := ec.c.CallContext(ctx, &map3Node, "eth_getMap3NodeInformation", map3NodeAddress, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return &map3Node, err

}

func (ec *Client) GetMap3NodeDelegation(
	ctx context.Context, map3NodeAddress common.Address, delegatorAddress common.Address, blockNumber *big.Int) (*microstaking.Microdelegation_, error) {

	var redelegation microstaking.Microdelegation_
	err := ec.c.CallContext(ctx, &redelegation, "eth_getMap3NodeDelegation", map3NodeAddress, delegatorAddress, toBlockNumArg(blockNumber))
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

func (ec *Client) GetActiveMap3NodeAtEpoch(ctx context.Context, epoch uint64) ([]string, error) {
	var result []string
	err := ec.c.CallContext(ctx, &result, "eth_getActiveMap3NodeAtEpoch", epoch)
	return result, err
}

func (ec *Client) GetTerminatedMap3NodeAtEpoch(ctx context.Context, epoch uint64) ([]string, error) {
	var result []string
	err := ec.c.CallContext(ctx, &result, "eth_getTerminatedMap3NodeAtEpoch", epoch)
	return result, err
}
