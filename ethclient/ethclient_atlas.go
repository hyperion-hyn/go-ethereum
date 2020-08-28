package ethclient

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (ec *Client) GetAllValidatorAddresses(ctx context.Context, blockNumber *big.Int) ([]string, error) {
	var result []string
	err := ec.c.CallContext(ctx, &result, "eth_getAllValidatorAddresses", toBlockNumArg(blockNumber))
	return result, err
}

func (ec *Client) GetValidatorInformation(
	ctx context.Context, validatorAddress common.Address, blockNumber *big.Int) (*restaking.ValidatorWrapperRPC, error) {

	var bytes hexutil.Bytes

	err := ec.c.CallContext(ctx, &bytes, "eth_getValidatorInformation", validatorAddress, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	var validator restaking.ValidatorWrapperRPC
	err = rlp.DecodeBytes(bytes, &validator)
	if err != nil {
		return nil, err
	}
	return &validator, err

}
