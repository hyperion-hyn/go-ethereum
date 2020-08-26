package ethapi

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
)

type PublicRestakingAPI struct {
	b Backend
}

// NewPublicBlockChainAPI creates a new Ethereum blockchain API.
func NewPublicRestakingAPI(b Backend) *PublicRestakingAPI {
	return &PublicRestakingAPI{b}
}

func (s *PublicRestakingAPI) GetAllValidatorAddresses(
	ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]string, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	validatorAddresses := state.ValidatorList()

	// Fetch all validator addresses
	addresses := make([]string, len(validatorAddresses))
	for i, addr := range validatorAddresses {
		hexAddr := addr.Hex()
		addresses[i] = hexAddr
	}
	return addresses, nil
}

func (s *PublicRestakingAPI) GetValidatorInformation(
	ctx context.Context, validatorAddress common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err

	}
	storageValidatorWarpper, err := state.ValidatorByAddress(validatorAddress)
	if err != nil {
		return nil, err
	}
	validatorWrapper := storageValidatorWarpper.Load()

	validator_ := validatorWrapper.Validator

	validatorWrpperRpc := &restaking.ValidatorWrapperRPC{
		Validator: restaking.ValidatorRPC{
			ValidatorAddress: validator_.ValidatorAddress,
			OperatorAddresses: func() []*common.Address {
				return validator_.OperatorAddresses.Keys
			}(),
			SlotPubKeys: func() []restaking.SlotPubKeyRPC {
				pubKeys := make([]restaking.SlotPubKeyRPC, len(validator_.SlotPubKeys.Keys))
				for _, bLSPublicKey := range validator_.SlotPubKeys.Keys {
					pubKeys = append(pubKeys, bLSPublicKey.Key)
				}
				return pubKeys
			}(),
			LastEpochInCommittee: validator_.LastEpochInCommittee,
			MaxTotalDelegation:   validator_.MaxTotalDelegation,
			Status:               validator_.Status,
			Commission:           validator_.Commission,
			Description:          validator_.Description,
			CreationHeight:       validator_.CreationHeight,
		},
		Redelegations: func() []restaking.Redelegation_ {
			redelegations := make([]restaking.Redelegation_, len(validatorWrapper.Redelegations.Keys))
			for _, key := range validatorWrapper.Redelegations.Keys {
				redegation, ok := validatorWrapper.Redelegations.Get(*key)
				if ok {
					redelegations = append(redelegations, redegation)
				}
			}
			return redelegations
		}(),
		Counters:                  validatorWrapper.Counters,
		BlockReward:               validatorWrapper.BlockReward,
		TotalDelegation:           validatorWrapper.TotalDelegation,
		TotalDelegationByOperator: validatorWrapper.TotalDelegationByOperator,
	}
	encodeBytes, err := rlp.EncodeToBytes(validatorWrpperRpc)
	if err != nil {
		return nil, err
	}
	return encodeBytes, nil

}
