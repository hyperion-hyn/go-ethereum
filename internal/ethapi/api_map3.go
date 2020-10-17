package ethapi

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
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
	ctx context.Context, map3NodeAddress common.Address, blockNrOrHash rpc.BlockNumberOrHash) (microstaking.PlainMap3NodeWrapper, error) {

	state, header, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return microstaking.PlainMap3NodeWrapper{}, err
	}

	storageNodeWrapper, err := state.Map3NodeByAddress(map3NodeAddress)
	if err != nil {
		return microstaking.PlainMap3NodeWrapper{}, err
	}

	nodeWrapper, err := storageNodeWrapper.Load()
	if err != nil {
		return microstaking.PlainMap3NodeWrapper{}, err
	}

	nodeWrapperRPC := nodeWrapper.ToPlainMap3NodeWrapper()
	nodeWrapperRPC.Map3Node.Age = storageNodeWrapper.Map3Node().CalculateNodeAge(header.Number, s.b.ChainConfig().Atlas)

	return *nodeWrapperRPC, nil
}

func (s *PublicMicroStakingAPI) GetMap3NodeDelegation(
	ctx context.Context, map3NodeAddress common.Address, delegatorAddress common.Address, blockNrOrHash rpc.BlockNumberOrHash) (microstaking.Microdelegation_, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return microstaking.Microdelegation_{}, err
	}

	storageNodeWrapper, err := state.Map3NodeByAddress(map3NodeAddress)
	if err != nil {
		return microstaking.Microdelegation_{}, err
	}
	storageMicrodelegation, ok := storageNodeWrapper.Microdelegations().Get(delegatorAddress)
	if ok {
		microDelegation, err := storageMicrodelegation.Load()
		if err != nil {
			return microstaking.Microdelegation_{}, err
		}
		return *microDelegation, nil

	} else {
		return microstaking.Microdelegation_{}, ethereum.NotFound
	}
}

func (s *PublicMicroStakingAPI) GetAllMap3RewardByDelegatorAddress(
	ctx context.Context, delegatorAddress common.Address, blockNrOrHash rpc.BlockNumberOrHash) (map[common.Address]*hexutil.Big, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	map3NodePool := state.Map3NodePool()
	delegationIndexMap := map3NodePool.DelegationIndexMapByDelegator().Get(delegatorAddress)
	rewardMap := make(map[common.Address]*hexutil.Big)
	for i := 0; i < delegationIndexMap.Keys().Length(); i++ {
		nodeAddr := delegationIndexMap.Keys().Get(i).Value()
		node, err := state.Map3NodeByAddress(nodeAddr)
		if err != nil {
			return nil, err
		}
		if micro, ok := node.Microdelegations().Get(delegatorAddress); ok {
			r := micro.Reward().Value()
			if r.Sign() > 0 {
				rewardMap[nodeAddr] = (*hexutil.Big)(r)
			}
		} else {
			return nil, errMicrodelegationNotExist
		}
	}
	return rewardMap, nil

}

func (s *PublicMicroStakingAPI) GetActiveMap3NodeAtEpoch(ctx context.Context, epoch uint64) ([]ActiveMap3Info, error) {
	blockNum := s.b.ChainConfig().Atlas.EpochFirstBlock(epoch)
	state, _, err := s.b.StateAndHeaderByNumber(ctx, rpc.BlockNumber(blockNum))
	if state == nil || err != nil {
		return nil, err
	}

	db := s.b.ChainDb()
	activeMap3Addr := rawdb.ReadActiveMap3Nodes(db, epoch)
	if activeMap3Addr == nil {
		return nil, nil
	}
	var map3Infos []ActiveMap3Info
	for _, addrTemp := range activeMap3Addr {
		map3WrapperStorage, err := state.Map3NodeByAddress(addrTemp)
		if err != nil {
			return nil, err
		}
		map3Wrapper, err := map3WrapperStorage.Load()
		if err != nil {
			return nil, err
		}
		activeMap3Infor := ActiveMap3Info{
			Address:    addrTemp.Hex(),
			StartEpoch: map3Wrapper.Map3Node.ActivationEpoch.Uint64(),
			EndEpoch:   map3Wrapper.Map3Node.ReleaseEpoch,
			Commission: map3Wrapper.Map3Node.Commission.Rate,
		}
		map3Infos = append(map3Infos, activeMap3Infor)
	}
	return map3Infos, nil
}

func (s *PublicMicroStakingAPI) GetTerminatedMap3NodeAtEpoch(ctx context.Context, epoch uint64) ([]string, error) {

	db := s.b.ChainDb()

	terminatedMap3Addr := rawdb.ReadTerminatedMap3Nodes(db, epoch)

	if terminatedMap3Addr == nil {
		return nil, nil
	}
	var addrStrs []string

	for _, addrTemp := range terminatedMap3Addr {
		addrStrs = append(addrStrs, addrTemp.Hex())
	}
	return addrStrs, nil
}

func (s *PublicMicroStakingAPI) GetMap3Requirement(ctx context.Context) (Map3Requirement, error) {

	blockNum := s.b.CurrentBlock().Header().Number
	chainContext := s.b.ChainContext()
	requireTotal, requireSelf, requireDel := network.LatestMicrostakingRequirement(blockNum, chainContext.Config())
	tenPercent := common.NewDecWithPrec(1, 1)    // 10%
	twentyPercent := common.NewDecWithPrec(2, 1) // 20%

	map3Requirement := Map3Requirement{
		RequireTotal:    requireTotal,
		RequireSelf:     requireSelf,
		RequireDelegate: requireDel,
		MinCommission:   tenPercent,
		MaxCommission:   twentyPercent,
		Map3LockEpoch:   microstaking.Map3NodeLockDurationInEpoch,
	}
	return map3Requirement, nil
}

type ActiveMap3Info struct {
	Address    string     `json:"address"`
	StartEpoch uint64     `json:"start_epoch"`
	EndEpoch   common.Dec `json:"end_epoch"`
	Commission common.Dec `json:"commission"`
}

type Map3Requirement struct {
	RequireTotal    *big.Int   `json:"requireTotal"`
	RequireSelf     *big.Int   `json:"requireSelf"`
	RequireDelegate *big.Int   `json:"requireDelegate"`
	MinCommission   common.Dec `json:"minCommission"`
	MaxCommission   common.Dec `json:"maxCommission"`
	Map3LockEpoch   common.Dec `json:"map3LockEpoch"`
}
