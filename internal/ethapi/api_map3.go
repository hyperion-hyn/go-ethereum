package ethapi

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

var (
	errMicrodelegationNotExist = errors.New("microdelegation does not exist")
)

type PublicMicroStakingAPI struct {
	b Backend
}

func NewPublicMicroStakingAPI(b Backend) *PublicMicroStakingAPI {
	return &PublicMicroStakingAPI{b}
}

func (s *PublicMicroStakingAPI) GetAllMap3NodeAddresses(
	ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]string, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}

	keys := state.Map3NodePool().Nodes().AllKeys()
	addresses := make([]string, len(keys))
	for i, addr := range keys {
		hexAddr := addr.Hex()
		addresses[i] = hexAddr
	}
	return addresses, nil
}

func (s *PublicMicroStakingAPI) GetMap3NodeInformation(
	ctx context.Context, map3NodeAddress common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {

	state, header, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}

	storageNodeWrapper, err := state.Map3NodeByAddress(map3NodeAddress)
	if err != nil {
		return nil, err
	}

	nodeWrapper, err := storageNodeWrapper.LoadFully()
	if err != nil {
		return nil, err
	}

	nodeWrapperRPC := nodeWrapper.ToPlainMap3NodeWrapper()
	nodeWrapperRPC.Map3Node.Age = storageNodeWrapper.Map3Node().CalculateNodeAge(header.Number, s.b.ChainConfig().Atlas)

	encodeBytes, err := rlp.EncodeToBytes(&nodeWrapperRPC)
	if err != nil {
		return nil, err
	}
	return encodeBytes, nil
}

func (s *PublicMicroStakingAPI) GetMap3NodeDelegation(
	ctx context.Context, map3NodeAddress common.Address, delegatorAddress common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}

	storageNodeWrapper, err := state.Map3NodeByAddress(map3NodeAddress)
	if err != nil {
		return nil, err
	}
	storageMicrodelegation, ok := storageNodeWrapper.Microdelegations().Get(delegatorAddress)
	if ok {
		microDelegation, err := storageMicrodelegation.Load()
		if err != nil {
			return nil, err
		}

		encodeBytes, err := rlp.EncodeToBytes(microDelegation)
		if err != nil {
			return nil, err
		}
		return encodeBytes, nil

	} else {
		return nil, ethereum.NotFound
	}
}

func (s *PublicMicroStakingAPI) GetAllMap3RewardByDelegatorAddress(
	ctx context.Context, delegatorAddress common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	map3NodePool := state.Map3NodePool()
	delegationIndexMap := map3NodePool.DelegationIndexMapByDelegator().Get(delegatorAddress)
	totalRewards := big.NewInt(0)
	for i := 0; i < delegationIndexMap.Keys().Length(); i++ {
		nodeAddr := delegationIndexMap.Keys().Get(i).Value()
		node, err := state.Map3NodeByAddress(nodeAddr)
		if err != nil {
			return nil, err
		}
		if micro, ok := node.Microdelegations().Get(delegatorAddress); ok {
			r := micro.Reward().Value()
			if r.Sign() > 0 {
				totalRewards = totalRewards.Add(totalRewards, r)
				micro.Reward().Clear()
			}
		} else {
			return nil, errMicrodelegationNotExist
		}
	}
	return (*hexutil.Big)(totalRewards), nil

}
