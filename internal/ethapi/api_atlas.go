package ethapi

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
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
	ctx context.Context, validatorAddress common.Address, blockNrOrHash rpc.BlockNumberOrHash) (restaking.PlainValidatorWrapper, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return restaking.PlainValidatorWrapper{}, err

	}
	storageValidatorWarpper, err := state.ValidatorByAddress(validatorAddress)
	if err != nil {
		return restaking.PlainValidatorWrapper{}, err
	}
	validatorWrapper, err := storageValidatorWarpper.Load()
	if err != nil {
		return restaking.PlainValidatorWrapper{}, err
	}

	validatorWrapperRPC := validatorWrapper.ToPlainValidatorWrapper()
	//encodeBytes, err := rlp.EncodeToBytes(&validatorWrapperRPC)
	//if err != nil {
	//	return nil, err
	//}
	return *validatorWrapperRPC, nil
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

func (s *PublicRestakingAPI) GetCommitteeAtEpoch(ctx context.Context, epoch uint64) (restaking.Committee_, error) {
	blockNum := s.b.ChainConfig().Atlas.EpochFirstBlock(epoch)
	committeeStorage, err := s.b.ChainContext().ReadCommitteeAtBlock(big.NewInt(int64(blockNum)))
	if err != nil {
		return restaking.Committee_{}, err
	}
	committee, err := committeeStorage.Load()
	return *committee, nil
}

func (s *PublicRestakingAPI) GetCommitteeInformationAtEpoch(ctx context.Context, epoch uint64) ([]restaking.PlainValidatorWrapper, error) {
	blockNum := s.b.ChainConfig().Atlas.EpochFirstBlock(epoch)
	committeeStorage, err := s.b.ChainContext().ReadCommitteeAtBlock(big.NewInt(int64(blockNum)))
	if err != nil {
		return nil, err
	}
	committee, err := committeeStorage.Load()

	slots := committee.Slots.Entrys

	plainValidators := make([]restaking.PlainValidatorWrapper, 0)
	for _, slotTemp := range slots {
		validatorWrapperStorage, err := s.b.ChainContext().ReadValidatorSnapshotAtBlock(big.NewInt(int64(blockNum)), slotTemp.EcdsaAddress)
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
	return plainValidators, nil
}

func (s *PublicRestakingAPI) GetValidatorInformationAtEpoch(ctx context.Context, validatorAddress common.Address,
	epoch uint64) (restaking.PlainValidatorWrapper, error) {
	blockNum := s.b.ChainConfig().Atlas.EpochFirstBlock(epoch)
	validatorWrapperStorage, err := s.b.ChainContext().ReadValidatorSnapshotAtBlock(big.NewInt(int64(blockNum)), validatorAddress)
	if err != nil {
		return restaking.PlainValidatorWrapper{}, err
	}
	validatorWrapper, err := validatorWrapperStorage.Load()
	if err != nil {
		return restaking.PlainValidatorWrapper{}, err
	}
	validatorWrapperRPC := validatorWrapper.ToPlainValidatorWrapper()
	return *validatorWrapperRPC, nil
}

func (s *PublicRestakingAPI) GetValidatorRedelegation(ctx context.Context, validatorAddress common.Address,
	delegatorAddress common.Address, blockNrOrHash rpc.BlockNumberOrHash) (restaking.Redelegation_, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return restaking.Redelegation_{}, err

	}
	storageValidatorWarpper, err := state.ValidatorByAddress(validatorAddress)
	if err != nil {
		return restaking.Redelegation_{}, err
	}

	storageRedelegatetor, ok := storageValidatorWarpper.Redelegations().Get(delegatorAddress)
	if ok {
		redelegation, err := storageRedelegatetor.Load()
		if err != nil {
			return restaking.Redelegation_{}, err
		}
		return *redelegation, nil
	} else {
		return restaking.Redelegation_{}, ethereum.NotFound
	}
}
