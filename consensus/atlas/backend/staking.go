package backend

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/atlas/backend/reward"
	"github.com/ethereum/go-ethereum/consensus/atlas/backend/votepower"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/staking/availability"
	"github.com/ethereum/go-ethereum/staking/committee"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"math/big"
)

// ATLAS
func handleMap3AndAtlasStaking(chain consensus.ChainReader, header *types.Header, stateDB *state.StateDB) (reward.Reader, error) {
	isNewEpoch := chain.Config().Atlas.IsFirstBlock(header.Number.Uint64())
	isEnd := chain.Config().Atlas.IsLastBlock(header.Number.Uint64())
	if isEnd {
		// Needs to be before AccumulateRewardsAndCountSigs because
		// ComputeAndMutateEPOSStatus depends on the signing counts that's
		// consistent with the counts when the new shardState was proposed.
		// Refer to committee.IsEligibleForEPoSAuction()
		curComm, err := lookupCommitteeAtEpoch(header.Epoch, chain)
		if err != nil {
			return nil, err
		}
		// TODO(ATLAS): disable status change temporarily
		for _, _ = range curComm.StakedValidators().Addrs {
			//if err := availability.ComputeAndMutateEPOSStatus(
			//	chain, stateDB, addr, header.Epoch,
			//); err != nil {
			//	return nil, err
			//}
		}

		// renew map3 node and unmicrodelegate and unredelegate
		if err := renewAndActivateMap3Nodes(chain, header, stateDB); err != nil {
			return nil, err
		}

		// update committee
		if _, err := updateCommitteeForNextEpoch(chain, header, stateDB); err != nil {
			return nil, err
		}
	}

	payout, err := accumulateRewardsAndCountSigs(chain, stateDB, header)
	if err != nil {
		return nil, errors.Wrap(err, "cannot pay block reward")
	}
	// TODO(ATLAS): slash

	if isNewEpoch {
		newComm, err := lookupCommitteeAtEpoch(header.Epoch, chain)
		if err != nil {
			return nil, err
		}
		if err := setLastEpochInCommittee(newComm, stateDB); err != nil {
			return nil, err
		}

		// Need to be after accumulateRewardsAndCountSigs because unredelegation may release
		releaser, err := UndelegationReleaserFactory{}.Create(stateDB, chain.Config())
		if err != nil {
			return nil, err
		}
		if err := payoutUnredelegations(header, stateDB, releaser); err != nil {
			return nil, err
		}

		// TODO(ATLAS): payout microdelegation and reward
	}
	return payout, nil
}

func renewAndActivateMap3Nodes(chain consensus.ChainReader, header *types.Header, stateDB *state.StateDB) error {
	requireTotal, requireSelf, _ := network.LatestMap3StakingRequirement(header.Number, chain.Config())
	var addrs []common.Address
	map3NodePool := stateDB.Map3NodePool()
	curEpoch := header.Epoch
	for _, nodeAddr := range map3NodePool.Nodes().AllKeys() {
		node, ok := map3NodePool.Nodes().Get(nodeAddr)
		if !ok {
			log.Error("map3 node should exist", "map3 address", nodeAddr.String())
			continue
		}

		if node.CanRelease(curEpoch) {
			nodeAge := node.Map3Node().CalculateNodeAge(header.Number, chain.Config().Atlas)
			node.Map3Node().Age().SetValue(nodeAge)

			isRenewed, amt, err := node.UnmicrodelegateIfNotRenewed(curEpoch)
			if err != nil {
				return err
			}
			if isRenewed {
				err := node.Pend(curEpoch)
				if err != nil {
					return err
				}
				if node.CanActivate(requireTotal, requireSelf) {
					if err := node.Activate(curEpoch); err != nil {
						return err
					}
				}

				if node.IsAlreadyRestaking() {
					validatorAddr := node.RestakingReference().ValidatorAddress().Value()
					validator, err := stateDB.ValidatorByAddress(validatorAddr)
					if err != nil {
						return err
					}
					validator.Undelegate(nodeAddr, curEpoch, amt)
				}
			} else {
				node.Terminate()
			}
			continue
		}

		if node.CanActivate(requireTotal, requireSelf) {
			if err := node.Activate(curEpoch); err != nil {
				return err
			}
		}
	}
	log.Info("New active map3 nodes", "addresses", addrs)
	return nil
}

func setLastEpochInCommittee(comm *restaking.Committee_, stateDB *state.StateDB) error {
	for _, addr := range comm.StakedValidators().Addrs {
		wrapper, err := stateDB.ValidatorByAddress(addr)
		if err != nil {
			return errors.WithMessage(err, "[Finalize] failed to get validator from state to finalize")
		}
		wrapper.Validator().LastEpochInCommittee().SetValue(comm.Epoch)
	}
	return nil
}

type UndelegationReleaserFactory struct {
}

func (f UndelegationReleaserFactory) Create(stateDB *state.StateDB, config *params.ChainConfig) (UndelegationReleaser, error) {
	if config.Atlas == nil {
		return nil, errors.New("not support to undelegate")
	}
	if config.Atlas.RestakingEnable {
		return undelegationToMap3Node{
			stateDB:       stateDB,
			rewardHandler: core.RewardToMap3Node{StateDB: stateDB},
		}, nil
	} else {
		return undelegationToBalance{
			stateDB:       stateDB,
			rewardHandler: core.RewardToBalance{StateDB: stateDB},
		}, nil
	}
}

type UndelegationReleaser interface {
	Release(redelegation *restaking.Storage_Redelegation_, fromValidator common.Address, epoch *big.Int) (completed bool, err error)
}

type undelegationToBalance struct {
	stateDB *state.StateDB
	rewardHandler core.RewardToBalance
}

func (u undelegationToBalance) Release(redelegation *restaking.Storage_Redelegation_, fromValidator common.Address,
	epoch *big.Int) (completed bool, err error) {
	// return undelegation
	delegator := redelegation.DelegatorAddress().Value()
	undelegation := redelegation.Undelegation().Amount().Value()
	u.stateDB.AddBalance(delegator, undelegation)
	redelegation.Undelegation().Clear()

	// return reward if redelgation is empty
	if amt := redelegation.Amount().Value(); amt.Cmp(common.Big0) == 0 {
		_, err := u.rewardHandler.HandleReward(redelegation, epoch)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

type undelegationToMap3Node struct {
	stateDB *state.StateDB
	rewardHandler core.RewardToMap3Node
}

func (u undelegationToMap3Node) Release(redelegation *restaking.Storage_Redelegation_, fromValidator common.Address,
	epoch *big.Int) (completed bool, err error) {
	// clear undelegation
	redelegation.Undelegation().Clear()

	// return reward if redelgation is empty
	if amt := redelegation.Amount().Value(); amt.Cmp(common.Big0) == 0 {
		_, err := u.rewardHandler.HandleReward(redelegation, epoch)
		if err != nil {
			return false, err
		}

		// clear restaking reference
		map3Addr := redelegation.DelegatorAddress().Value()
		node, err := u.stateDB.Map3NodeByAddress(map3Addr)
		if err != nil {
			return false, err
		}
		node.RestakingReference().Clear()
		return true, nil
	}
	return false, nil
}


// Withdraw unlocked tokens to the delegators' accounts
func payoutUnredelegations(header *types.Header, stateDB *state.StateDB, releaser UndelegationReleaser) error {
	nowEpoch := header.Epoch
	validators := stateDB.ValidatorPool().Validators()
	// Payout undelegated/unlocked tokens
	for _, validatorAddr := range validators.AllKeys() {
		validator, ok := validators.Get(validatorAddr)
		if !ok {
			return errValidatorNotExist
		}

		var toBeRemoved []common.Address
		for _, delegator := range validator.Redelegations().AllKeys() {
			redelegation, ok := validator.Redelegations().Get(delegator)
			if !ok {
				return errRedelegationNotExist
			}

			if redelegation.CanReleaseAt(nowEpoch) {
				completed, err := releaser.Release(redelegation, validatorAddr, nowEpoch)
				if err != nil {
					return err
				}
				if completed {
					toBeRemoved = append(toBeRemoved, delegator)
				}
			}
		}

		for _, delegator := range toBeRemoved {
			validator.Redelegations().Remove(delegator)
		}
	}
	log.Info("paid out delegations", "epoch", nowEpoch.Uint64(), "block-number", header.Number.Uint64())
	return nil
}

func updateCommitteeForNextEpoch(chain consensus.ChainReader, header *types.Header,
	stateDB *state.StateDB) (*restaking.Committee_, error) {
	nextEpoch := big.NewInt(0).Add(header.Epoch, common.Big1)
	nextComm, err := committee.WithStakingEnabled.Compute(nextEpoch, committee.ChainReaderWithPendingState{
		ChainReader: chain,
		StateDB:     stateDB,
	})
	if err != nil {
		return nil, err
	}
	stateDB.ValidatorPool().UpdateCommittee(nextComm)
	return nextComm, nil
}

// accumulateRewardsAndCountSigs credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward
// This func also do IncrementValidatorSigningCounts for validators
func accumulateRewardsAndCountSigs(
	bc consensus.ChainReader, state *state.StateDB, header *types.Header,
) (reward.Reader, error) {
	if header.Number.Uint64() <= 1 {
		// genesis block has no parent to reward.
		return network.EmptyPayout, nil
	}

	// calculate total reward
	lastBlockNum := big.NewInt(0).Sub(header.Number, common.Big1)
	pool := network.NewRewardPool(state)
	totalReward := pool.TakeReward(lastBlockNum, bc.Config())

	// no reward
	if totalReward.Sign() <= 0 {
		return network.EmptyPayout, nil
	}

	payouts := []reward.Payout{}
	comm, payable, missing, err := ballotResult(bc, header) // for last block
	if err != nil {
		return network.EmptyPayout, err
	}

	if err := availability.IncrementValidatorSigningCounts(
		bc,
		comm.StakedValidators(),
		state,
		payable,
		missing,
	); err != nil {
		return network.EmptyPayout, err
	}
	votingPower, err := lookupVotingPower(comm.Epoch, comm)
	if err != nil {
		return network.EmptyPayout, err
	}

	allSignersShare := common.ZeroDec()
	for j := range payable.Entrys {
		voter := votingPower.Voters[payable.Entrys[j].BLSPublicKey]
		voterShare := voter.OverallPercent
		allSignersShare = allSignersShare.Add(voterShare)
	}

	totalRewardDec := common.NewDecFromBigInt(totalReward)
	rewardPool := big.NewInt(0).Set(totalReward)
	for i := len(payable.Entrys) - 1; i >= 0; i-- {
		// what to do about share of those that didn't sign
		blsKey := payable.Entrys[i].BLSPublicKey
		voter := votingPower.Voters[blsKey]
		snapshot, err := bc.ReadValidatorAtEpoch(comm.Epoch, voter.EarningAccount)
		if err != nil {
			return network.EmptyPayout, err
		}

		var due *big.Int
		if i == 0 {	// Give out whatever leftover to the first voter/handle
			due = big.NewInt(0).Set(rewardPool)
		} else {
			due = totalRewardDec.Mul(
				voter.OverallPercent.Quo(allSignersShare),
			).RoundInt()
		}

		shares, err := lookupDelegatorShares(comm.Epoch, snapshot)
		if err != nil {
			return network.EmptyPayout, err
		}
		if err := state.AddRedelegationReward(snapshot, due, shares); err != nil {
			return network.EmptyPayout, err
		}
		payouts = append(payouts, reward.Payout{
			Addr:        voter.EarningAccount,
			NewlyEarned: due,
			EarningKey:  voter.Identity,
		})
		rewardPool.Sub(rewardPool, due)
	}
	return network.NewStakingEraRewardForRound(totalReward, missing, payouts), nil
}

func ballotResult(
	bc consensus.ChainReader, header *types.Header,
) (*restaking.Committee_, *restaking.Slots_, *restaking.Slots_, error) {
	parentHeader := bc.GetHeaderByHash(header.ParentHash)
	if parentHeader == nil {
		return nil, nil, nil, errors.Errorf(
			"cannot find parent block header in DB %s",
			header.ParentHash.Hex(),
		)
	}
	parentCommittee, err := lookupCommitteeAtEpoch(parentHeader.Epoch, bc)
	if err != nil {
		return nil, nil, nil, errors.Errorf(
			"cannot read shard state %v", parentHeader.Epoch,
		)
	}
	reader := availability.CommitBitmapReader{Header: parentHeader}
	_, payable, missing, err := availability.BallotResult(reader, parentCommittee)
	return parentCommittee, payable, missing, err
}

var (
	votingPowerCache   singleflight.Group
	delegateShareCache singleflight.Group
	committeeCache     singleflight.Group
)

func lookupCommitteeAtEpoch(epoch *big.Int, bc consensus.ChainReader) (*restaking.Committee_, error) {
	key := epoch.String()
	results, err, _ := committeeCache.Do(
		key, func() (interface{}, error) {
			// TODO: read from committee provider
			committeeSt, err := bc.ReadCommitteeAtEpoch(epoch)
			if err != nil {
				return nil, err
			}
			comm, err := committeeSt.Load()
			if err != nil {
				return nil, err
			}

			// For new calc, remove old data from 2 epochs ago
			deleteEpoch := big.NewInt(0).Sub(epoch, big.NewInt(2))
			deleteKey := deleteEpoch.String()
			votingPowerCache.Forget(deleteKey)
			return comm, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return results.(*restaking.Committee_), nil
}

func lookupVotingPower(epoch *big.Int, comm *restaking.Committee_) (*votepower.Roster, error) {
	key := epoch.String()
	results, err, _ := votingPowerCache.Do(
		key, func() (interface{}, error) {
			votingPower, err := votepower.Compute(comm)
			if err != nil {
				return nil, err
			}

			// For new calc, remove old data from 3 epochs ago
			deleteEpoch := big.NewInt(0).Sub(comm.Epoch, big.NewInt(3))
			deleteKey := deleteEpoch.String()
			votingPowerCache.Forget(deleteKey)

			return votingPower, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return results.(*votepower.Roster), nil
}

// Lookup or compute the shares of stake for all delegators in a validator
func lookupDelegatorShares(
	epoch *big.Int, snapshot *restaking.Storage_ValidatorWrapper_,
) (map[common.Address]common.Dec, error) {
	valAddr := snapshot.Validator().ValidatorAddress().Value()
	key := fmt.Sprintf("%s-%s", epoch.String(), valAddr.Hex())

	shares, err, _ := delegateShareCache.Do(
		key, func() (interface{}, error) {
			result := map[common.Address]common.Dec{}

			totalDelegationDec := common.NewDecFromBigInt(snapshot.TotalDelegation().Value())
			if totalDelegationDec.IsZero() {
				log.Info("zero total delegation during AddReward delegation payout",
					"validator-snapshot", valAddr.Hex())
				return result, nil
			}

			for _, key := range snapshot.Redelegations().AllKeys() {
				delegation, ok := snapshot.Redelegations().Get(key)
				if !ok {
					return nil, errValidatorNotExist
				}
				// NOTE percentage = <this_delegator_amount>/<total_delegation>
				percentage := common.NewDecFromBigInt(delegation.Amount().Value()).Quo(totalDelegationDec)
				result[delegation.DelegatorAddress().Value()] = percentage
			}

			// For new calc, remove old data from 3 epochs ago
			deleteEpoch := big.NewInt(0).Sub(epoch, big.NewInt(3))
			deleteKey := fmt.Sprintf("%s-%s", deleteEpoch.String(), valAddr.Hex())
			votingPowerCache.Forget(deleteKey)

			return result, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return shares.(map[common.Address]common.Dec), nil
}

// ATLAS - END

