package ethapi

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

type PublicRestakingAPI struct {
	b Backend
}

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
	validatorWrapper, err := storageValidatorWarpper.Load()
	if err != nil {
		return nil, err
	}

	validatorWrapperRPC := validatorWrapper.ToPlainValidatorWrapper()
	encodeBytes, err := rlp.EncodeToBytes(&validatorWrapperRPC)
	if err != nil {
		return nil, err
	}
	return encodeBytes, nil
}

func (s *PublicRestakingAPI) GetEpochFirstBlockNum(ctx context.Context, epoch uint64) uint64 {
	blockNumber := s.b.ChainConfig().Atlas.EpochFirstBlock(epoch)
	return blockNumber
}

func (s *PublicRestakingAPI) GetEpochLastBlockNum(
	ctx context.Context, epoch uint64) uint64 {
	blockNumber := s.b.ChainConfig().Atlas.EpochLastBlock(epoch)
	return blockNumber
}

func (s *PublicRestakingAPI) GetCommitteeAtEpoch(ctx context.Context, epoch uint64) (hexutil.Bytes, error) {
	committeeStorage, err := s.b.ChainContext().ReadCommitteeAtEpoch(big.NewInt(int64(epoch)))
	if err != nil {
		return nil, err
	}
	committee, err := committeeStorage.Load()
	encodeBytes, err := rlp.EncodeToBytes(committee)
	if err != nil {
		return nil, err
	}
	return encodeBytes, nil
}

func (s *PublicRestakingAPI) GetCommitteeInformationAtEpoch(ctx context.Context, epoch uint64) (hexutil.Bytes, error) {
	committeeStorage, err := s.b.ChainContext().ReadCommitteeAtEpoch(big.NewInt(int64(epoch)))
	if err != nil {
		return nil, err
	}
	committee, err := committeeStorage.Load()

	slots := committee.Slots.Entrys

	plainValidators := make([]restaking.PlainValidatorWrapper, 0)
	for _, slotTemp := range slots {
		validatorWrapperStorage, err := s.b.ChainContext().ReadValidatorAtEpoch(big.NewInt(int64(epoch)), slotTemp.EcdsaAddress)
		if err != nil {
			return nil, err
		}
		validatorWrapper, err := validatorWrapperStorage.Load()
		if err != nil {
			return nil, err
		}
		validatorWrapperRPC := validatorWrapper.ToPlainValidatorWrapper()
		plainValidators = append(plainValidators, *validatorWrapperRPC)
	}

	encodeBytes, err := rlp.EncodeToBytes(plainValidators)
	if err != nil {
		return nil, err
	}
	return encodeBytes, nil
}

func (s *PublicRestakingAPI) GetValidatorInformationAtEpoch(ctx context.Context, validatorAddress common.Address,
	epoch uint64) (hexutil.Bytes, error) {
	validatorWrapperStorage, err := s.b.ChainContext().ReadValidatorAtEpoch(big.NewInt(int64(epoch)), validatorAddress)
	if err != nil {
		return nil, err
	}

	validatorWrapper, err := validatorWrapperStorage.Load()
	if err != nil {
		return nil, err
	}

	validatorWrapperRPC := validatorWrapper.ToPlainValidatorWrapper()
	encodeBytes, err := rlp.EncodeToBytes(&validatorWrapperRPC)
	if err != nil {
		return nil, err
	}
	return encodeBytes, nil

}

func (s *PublicRestakingAPI) GetValidatorRedelegation(ctx context.Context, validatorAddress common.Address,
	delegatorAddress common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err

	}
	storageValidatorWarpper, err := state.ValidatorByAddress(validatorAddress)
	if err != nil {
		return nil, err
	}

	storageRedelegatetor, ok := storageValidatorWarpper.Redelegations().Get(delegatorAddress)
	if ok {
		redelegation, err := storageRedelegatetor.Load()
		if err != nil {
			return nil, err
		}

		encodeBytes, err := rlp.EncodeToBytes(&redelegation)
		if err != nil {
			return nil, err
		}
		return encodeBytes, nil
	} else {
		return nil, ethereum.NotFound
	}
}
