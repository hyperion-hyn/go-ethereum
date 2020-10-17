package ethclient

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (ec *Client) GetAllValidatorAddresses(ctx context.Context, blockNumber *big.Int) ([]string, error) {
	var result []string
	err := ec.c.CallContext(ctx, &result, "eth_getAllValidatorAddresses", toBlockNumArg(blockNumber))
	return result, err
}

func (ec *Client) GetValidatorInformation(
	ctx context.Context, validatorAddress common.Address, blockNumber *big.Int) (*restaking.PlainValidatorWrapper, error) {
	var validator restaking.PlainValidatorWrapper
	err := ec.c.CallContext(ctx, &validator, "eth_getValidatorInformation", validatorAddress, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return &validator, err
}

func (ec *Client) GetEpochFirstBlockNum(ctx context.Context, epoch uint64) (uint64, error) {

	var result uint64
	err := ec.c.CallContext(ctx, &result, "eth_getEpochFirstBlockNum", epoch)
	return result, err
}

func (ec *Client) GetEpochLastBlockNum(ctx context.Context, epoch uint64) (uint64, error) {

	var result uint64
	err := ec.c.CallContext(ctx, &result, "eth_getEpochLastBlockNum", epoch)
	return result, err
}

func (ec *Client) GetCommitteeAtEpoch(ctx context.Context, epoch uint64) (*restaking.Committee_, error) {
	var committee restaking.Committee_
	err := ec.c.CallContext(ctx, &committee, "eth_getCommitteeAtEpoch", epoch)
	if err != nil {
		return nil, err
	}
	return &committee, err
}

func (ec *Client) GetCommitteeInformationAtEpoch(ctx context.Context, epoch uint64) ([]restaking.PlainValidatorWrapper, error) {
	var validators []restaking.PlainValidatorWrapper

	err := ec.c.CallContext(ctx, &validators, "eth_getCommitteeInformationAtEpoch", epoch)
	if err != nil {
		return nil, err
	}
	return validators, err
}

func (ec *Client) GetValidatorInformationAtEpoch(
	ctx context.Context, validatorAddress common.Address, epoch uint64) (*restaking.PlainValidatorWrapper, error) {

	var validator restaking.PlainValidatorWrapper
	err := ec.c.CallContext(ctx, &validator, "eth_getValidatorInformationAtEpoch", validatorAddress, epoch)
	if err != nil {
		return nil, err
	}
	return &validator, err
}

func (ec *Client) GetValidatorRedelegation(
	ctx context.Context, validatorAddress common.Address, delegatorAddress common.Address, blockNumber *big.Int) (*restaking.Redelegation_, error) {

	var redelegation restaking.Redelegation_
	err := ec.c.CallContext(ctx, &redelegation, "eth_getValidatorRedelegation", validatorAddress, delegatorAddress, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return &redelegation, err
}
