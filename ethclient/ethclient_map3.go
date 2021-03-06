package ethclient

import (
	"context"
	"github.com/ethereum/go-ethereum/staking/burning"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
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

func (ec *Client) GetAllMap3RewardByDelegatorAddress(ctx context.Context, delegatorAddress common.Address, blockNumber *big.Int) (map[common.Address]*big.Int, error) {

	resultMap := make(map[common.Address]*hexutil.Big)
	err := ec.c.CallContext(ctx, &resultMap, "eth_getAllMap3RewardByDelegatorAddress", delegatorAddress, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	rewardMap := make(map[common.Address]*big.Int)
	for address, amount := range resultMap {
		rewardMap[address] = (*big.Int)(amount)
	}
	return rewardMap, err
}

func (ec *Client) GetMutateMap3NodeAtEpoch(ctx context.Context, epoch uint64) ([]string, error) {
	var result []string
	err := ec.c.CallContext(ctx, &result, "eth_getMutateMap3NodeAtEpoch", epoch)
	return result, err
}

func (ec *Client) GetMap3Requirement(ctx context.Context) (types.Map3Requirement, error) {
	var requirement types.Map3Requirement
	err := ec.c.CallContext(ctx, &requirement, "eth_getMap3Requirement")
	return requirement, err
}

func (ec *Client) GetUnmicrodelegationReturnRecordsAtBlock(ctx context.Context, block int64) ([]microstaking.MicrostakingReturnRecord, error) {
	var result []microstaking.MicrostakingReturnRecord
	err := ec.c.CallContext(ctx, &result, "eth_getUnmicrodelegationReturnRecordsAtBlock", block)
	return result, err
}

func (ec *Client) GetTerminateMap3RecordsAtBlock(ctx context.Context, map3Node common.Address, block int64) ([]microstaking.MicrostakingReturnRecord, error) {
	var result []microstaking.MicrostakingReturnRecord
	err := ec.c.CallContext(ctx, &result, "eth_getTerminateMap3RecordsAtBlock", map3Node, block)
	return result, err
}

func (ec *Client) GetTokenBurningReceiptByBlockNum(ctx context.Context, block int64) (*burning.Receipt, error) {
	var result burning.Receipt
	err := ec.c.CallContext(ctx, &result, "eth_getTokenBurningReceiptByBlockNum", block)
	return &result, err
}

func (ec *Client) GetTokenBurningReceiptByReceiptHash(ctx context.Context, hash common.Hash) (*burning.Receipt, error) {
	var result burning.Receipt
	err := ec.c.CallContext(ctx, &result, "eth_getTokenBurningReceiptByReceiptHash", hash)
	return &result, err
}
