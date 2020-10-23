package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/pkg/errors"
	"math/big"
)

var (
	errMap3NodeNotExist        = errors.New("map3 node does not exist")
	errMicrodelegationNotExist = errors.New("microdelegation does not exist")

	map3StorageAddress = common.HexToAddress("0x6a7ad21ff076440e39020e289debdcb309e12c23")
)

func (s *StateDB) Map3NodePool() *microstaking.Storage_Map3NodePool_ {
	// singleton
	s.map3Once.Do(func() {
		var g microstaking.Global_t
		globalSt := microstaking.New(&g, s, map3StorageAddress, common.Big0)
		s.map3NodePool = globalSt.Map3NodePool()
	})
	return s.map3NodePool
}

func (s *StateDB) Map3NodeByAddress(map3Address common.Address) (*microstaking.Storage_Map3NodeWrapper_, error) {
	node, ok := s.Map3NodePool().Map3Nodes().Get(map3Address)
	if !ok {
		return nil, errMap3NodeNotExist
	}
	return node, nil
}

func (s *StateDB) Map3NodeSnapshotByAddress(map3Address common.Address) (*microstaking.Storage_Map3NodeWrapper_, error) {
	node, ok := s.Map3NodePool().Map3NodeSnapshots().Get(map3Address)
	if !ok {
		return nil, errMap3NodeNotExist
	}
	return node, nil
}

func (s *StateDB) Map3NodeList() []common.Address {
	return s.Map3NodePool().Map3Nodes().AllKeys()
}

// AddMicrodelegationReward distributes the reward to all the delegators based on stake percentage.
func (s *StateDB) AddMicrodelegationReward(snapshot *microstaking.Storage_Map3NodeWrapper_, reward *big.Int,
	shareLookup map[common.Address]common.Dec) error {
	map3Addr := snapshot.Map3Node().Map3Address().Value()
	if reward.Cmp(common.Big0) == 0 {
		log.Info("0 given as reward", "map3Address", map3Addr)
		return nil
	}

	curNode, err := s.Map3NodeByAddress(map3Addr)
	if err != nil {
		return errors.Wrapf(err, "failed to distribute rewards: map3 node does not exist")
	}
	operatorDelegation, ok := curNode.Microdelegations().Get(curNode.Map3Node().OperatorAddress().Value())
	if !ok {
		return errors.New("missing delegation of operator for reward distribution")
	}
	curNode.AddAccumulatedReward(reward)

	// Payout each delegator's (non-operator) reward
	rewardPool := big.NewInt(0).Set(reward)
	commissionRate := snapshot.Map3Node().Commission().Rate().Value()

	// Payout commission
	if commissionRate.GT(common.ZeroDec()) {
		commissionInt := commissionRate.MulInt(reward).RoundInt()
		operatorDelegation.AddReward(commissionInt)
		rewardPool.Sub(rewardPool, commissionInt)
	}

	// Payout each delegator's reward pro-rata
	totalRewardForDelegators := big.NewInt(0).Set(rewardPool)
	for _, delegatorAddress := range snapshot.Microdelegations().AllKeys() {
		percentage, ok := shareLookup[delegatorAddress]
		if !ok {
			return errors.Wrapf(err, "missing delegation shares for reward distribution")
		}
		if percentage.IsNil() || percentage.IsZero() {
			continue
		}
		rewardInt := percentage.MulInt(totalRewardForDelegators).RoundInt()

		curDelegation, ok := curNode.Microdelegations().Get(delegatorAddress)
		if !ok {
			return errMicrodelegationNotExist
		}
		curDelegation.AddReward(rewardInt)
		rewardPool.Sub(rewardPool, rewardInt)
	}

	// The last remaining bit belongs to the operator
	if rewardPool.Cmp(common.Big0) > 0 {
		operatorDelegation.AddReward(rewardPool)
	}
	return nil
}

/**
 * IncreaseMap3NonceIfZero avoids account state of map3 nodes would be delete if its nonce and balance are zero
 */
func (s *StateDB) IncreaseMap3NonceIfZero() {
	if s.GetNonce(map3StorageAddress) == 0 {
		s.SetNonce(map3StorageAddress, 1)
	}
}
