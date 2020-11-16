package ethapi

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/staking/burning"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/pkg/errors"
	"math/big"
)

var (
	errMicrodelegationNotExist = errors.New("microdelegation does not exist")
	errBurningReceiptNotFound  = errors.New("burning record not found")
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

	keys := state.Map3NodeList()
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

func (s *PublicMicroStakingAPI) GetMutateMap3NodeAtEpoch(ctx context.Context, epoch uint64) ([]string, error) {

	db := s.b.ChainDb()

	mutateMap3Addr := rawdb.ReadMutateMap3Nodes(db, epoch)

	if mutateMap3Addr == nil {
		return nil, nil
	}
	var addrStrs []string

	for _, addrTemp := range mutateMap3Addr {
		addrStrs = append(addrStrs, addrTemp.Hex())
	}
	return addrStrs, nil
}

func (s *PublicMicroStakingAPI) GetUnmicrodelegationReturnRecordsAtBlock(ctx context.Context, block int64) ([]microstaking.MicrostakingReturnRecord, error) {

	db := s.b.ChainDb()
	records := rawdb.ReadUnmicrodelegationReturnRecords(db, big.NewInt(block))
	if records == nil {
		return nil, nil
	}
	return records, nil
}

func (s *PublicMicroStakingAPI) GetTerminateMap3RecordsAtBlock(ctx context.Context, map3Node common.Address, block int64) ([]microstaking.MicrostakingReturnRecord, error) {

	db := s.b.ChainDb()
	records := rawdb.ReadTerminateMap3ReturnRecords(db, big.NewInt(block), map3Node)
	if records == nil {
		return nil, nil
	}
	return records, nil
}

func (s *PublicMicroStakingAPI) GetMap3Requirement(ctx context.Context) (types.Map3Requirement, error) {

	blockNum := s.b.CurrentBlock().Header().Number
	chainContext := s.b.ChainContext()
	requireTotal, requireSelf, requireDel := network.LatestMicrostakingRequirement(blockNum, chainContext.Config())
	zeroPercent := common.NewDec(0)
	hundredPercent := common.NewDec(1)

	map3Requirement := types.Map3Requirement{
		RequireTotal:    requireTotal,
		RequireSelf:     requireSelf,
		RequireDelegate: requireDel,
		MinCommission:   zeroPercent,
		MaxCommission:   hundredPercent,
		Map3LockEpoch:   common.NewDec(microstaking.LockDurationInEpoch),
	}
	return map3Requirement, nil
}

func (s *PublicMicroStakingAPI) GetMicrodelegationIndexByDelegator(ctx context.Context, delegatorAddress common.Address,
	blockNrOrHash rpc.BlockNumberOrHash) ([]microstaking.DelegationIndex_, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	map3NodePool := state.Map3NodePool()
	delegationIndexMap := map3NodePool.DelegationIndexMapByDelegator().Get(delegatorAddress)
	indexes := make([]microstaking.DelegationIndex_, 0)
	for i := 0; i < delegationIndexMap.Keys().Length(); i++ {
		nodeAddr := delegationIndexMap.Keys().Get(i).Value()
		index, ok := delegationIndexMap.Get(nodeAddr)
		if !ok {
			return nil, errors.Errorf("microdelegation index not found, node: %v, delegator: %v", nodeAddr, delegatorAddress)
		}
		indexes = append(indexes, microstaking.DelegationIndex_{
			Map3Address: index.Map3Address().Value(),
			IsOperator:  index.IsOperator().Value(),
		})
	}
	return indexes, nil
}

func (s *PublicMicroStakingAPI) GetTokenBurningReceiptByBlockNum(ctx context.Context, blockNum int64) (*burning.Receipt, error) {
	db := s.b.ChainDb()
	record := rawdb.ReadTokenBurningReceiptByBlockNum(db, big.NewInt(blockNum))
	if record == nil {
		return nil, errBurningReceiptNotFound
	}
	return record, nil
}

func (s *PublicMicroStakingAPI) GetTokenBurningReceiptByReceiptHash(ctx context.Context, hash common.Hash) (*burning.Receipt, error) {
	db := s.b.ChainDb()
	record := rawdb.ReadTokenBurningReceiptByHash(db, hash)
	if record == nil {
		return nil, errBurningReceiptNotFound
	}
	return record, nil
}
