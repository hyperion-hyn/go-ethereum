package ethapi

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
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

	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
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

	map3Node := nodeWrapper.Map3Node

	nodeWrapperRpc := microstaking.Map3NodeWrapperRPC{
		Map3Node: microstaking.Map3NodeRPC{
			Map3Address:     map3Node.Map3Address,
			OperatorAddress: map3Node.OperatorAddress,
			NodeKeys: func() []microstaking.NodeKeyRPC {
				nodeKeys := map3Node.NodeKeys.Keys
				nodekeyRPCs := make([]microstaking.NodeKeyRPC, 0)
				for _, nodeKeyTemp := range nodeKeys {
					nodekeyRPCs = append(nodekeyRPCs, nodeKeyTemp.Key)
				}
				return nodekeyRPCs
			}(),
			Commission:      map3Node.Commission,
			Description:     map3Node.Description,
			CreationHeight:  map3Node.CreationHeight,
			Age:             map3Node.Age,
			Status:          map3Node.Status,
			ActivationEpoch: map3Node.ActivationEpoch,
			ReleaseEpoch:    map3Node.ReleaseEpoch,
		},
		Microdelegations: func() []microstaking.Microdelegation_ {
			delegations := make([]microstaking.Microdelegation_, 0)
			delegationKeys := nodeWrapper.Microdelegations.Keys
			for _, keyTemp := range delegationKeys {
				delegationTemp, ok := nodeWrapper.Microdelegations.Get(*keyTemp)
				if ok {
					delegations = append(delegations, delegationTemp)
				}
			}
			return delegations
		}(),
		RedelegationReference:  nodeWrapper.RestakingReference.ValidatorAddress,
		AccumulatedReward:      nodeWrapper.AccumulatedReward,
		TotalDelegation:        nodeWrapper.TotalDelegation,
		TotalPendingDelegation: nodeWrapper.TotalPendingDelegation,
	}

	encodeBytes, err := rlp.EncodeToBytes(nodeWrapperRpc)
	if err != nil {
		return nil, err
	}
	return encodeBytes, nil

}
