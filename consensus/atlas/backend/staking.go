package backend

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/atlas/backend/reward"
	"github.com/ethereum/go-ethereum/consensus/atlas/backend/votepower"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/staking/availability"
	"github.com/ethereum/go-ethereum/staking/committee"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"math/big"
)

var (
	releaserFactory = UndelegationReleaserFactory{}
)

// ATLAS
func handleMap3AndAtlasStaking(chain consensus.ChainReader, header *types.Header, stateDB *state.StateDB) (reward.Reader, error) {
	payout, err := accumulateRewardsAndCountSigs(chain, stateDB, header)
	if err != nil {
		return nil, errors.Wrap(err, "cannot pay block reward")
	}
	// TODO(ATLAS): slash

	isNewEpoch := chain.Config().Atlas.IsFirstBlock(header.Number.Uint64())
	isEnd := chain.Config().Atlas.IsLastBlock(header.Number.Uint64())
	if isEnd {
		lastBlockNum := new(big.Int).Sub(header.Number, common.Big1)
		curComm, err := lookupCommitteeAtBlock(lastBlockNum, chain)
		if err != nil {
			return nil, err
		}
		// TODO(ATLAS): disable status change temporarily
		for _, addr := range curComm.StakedValidators().Addrs {
			if err := availability.ComputeAndMutateEPOSStatus(stateDB, addr); err != nil {
				return nil, err
			}
		}

		// renew map3 node and unmicrodelegate and unredelegate
		if err := renewAndActivateMap3Nodes(chain, header, stateDB); err != nil {
			return nil, err
		}

		// update committee
		newComm, err := updateCommitteeForNextEpoch(header, stateDB)
		if err != nil {
			return nil, err
		}
		if err := setLastEpochInCommittee(newComm, stateDB); err != nil {
			return nil, err
		}

		// update validator snapshots
		if err := updateValidatorSnapshots(stateDB); err != nil {
			return nil, err
		}
	}

	if isNewEpoch {
		lastEpoch := new(big.Int).Sub(header.Epoch, common.Big1)
		if err := collectRestakingRewardForRenewedMap3Nodes(stateDB, lastEpoch); err != nil {
			return nil, err
		}

		// Need to be after accumulateRewardsAndCountSigs because unredelegation may release
		releaser, err := releaserFactory.Create(stateDB, chain)
		if err != nil {
			return nil, err
		}
		if err := payoutUnredelegations(header, stateDB, releaser); err != nil {
			return nil, err
		}
		if err := payoutUnmicrodelegations(header, stateDB); err != nil {
			return nil, err
		}
	}
	return payout, nil
}

func collectRestakingRewardForRenewedMap3Nodes(stateDB *state.StateDB, epoch *big.Int) error {
	//foreach node from nodePool
	//if node.isPendAt(lastEpch) && node.IsRestaking {
	//	collect and distrbute restaking reward
	//}
	//rewardHandler := core.RewardToMap3Node{StateDB: stateDB}
	//for _, map3Addr := range stateDB.Map3NodeList() {
	//	node, err := stateDB.Map3NodeByAddress(map3Addr)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if node.IsRenewed(epoch) && node.IsRestaking() {
	//		validatorAddr := node.RestakingReference().ValidatorAddress().Value()
	//		stateDB.ValidatorByAddress()
	//	}
	//}

	return nil
}

func updateValidatorSnapshots(stateDB *state.StateDB) error {
	pool := stateDB.ValidatorPool()
	for _, validatorAddress := range pool.Validators().AllKeys() {
		validatorStr, err := stateDB.ValidatorByAddress(validatorAddress)
		if err != nil {
			return err
		}
		v, err := validatorStr.Load()
		if err != nil {
			return err
		}
		pool.ValidatorSnapshots().Put(validatorAddress, v)
	}
	return nil
}

func renewAndActivateMap3Nodes(chain consensus.ChainReader, header *types.Header, stateDB *state.StateDB) error {
	requireTotal, requireSelf, _ := network.LatestMicrostakingRequirement(header.Number, chain.Config())
	nowEpoch := header.Epoch
	var mutateMap3Addrs []common.Address
	var renewActiveMap3Addrs []common.Address
	for _, nodeAddr := range stateDB.Map3NodeList() {
		node, err := stateDB.Map3NodeByAddress(nodeAddr)
		if err != nil {
			return err
		}

		if node.CanReleaseAt(nowEpoch) {
			nodeAge := node.Map3Node().CalculateNodeAge(header.Number, chain.Config().Atlas)
			node.Map3Node().Age().SetValue(nodeAge)

			mutateMap3Addrs = append(mutateMap3Addrs, nodeAddr)

			isRenewed, notRenewedAmt, err := node.UnmicrodelegateIfNotRenewed(nowEpoch)
			if err != nil {
				return err
			}
			if isRenewed {
				err := node.RenewAndPend(nowEpoch)
				if err != nil {
					return err
				}

				isActive := false
				if node.CanActivate(requireTotal, requireSelf) {
					if err := node.Activate(nowEpoch); err != nil {
						return err
					}
					isActive = true
					renewActiveMap3Addrs = append(renewActiveMap3Addrs, nodeAddr)
				}

				if node.IsRestaking() && (notRenewedAmt.Sign() > 0 || !isActive) {
					validatorAddr := node.RestakingReference().ValidatorAddress().Value()
					validator, err := stateDB.ValidatorByAddress(validatorAddr)
					if err != nil {
						return err
					}

					undelegation := notRenewedAmt
					if !isActive { // undelegate total amount if not active
						undelegation = nil
					}
					validator.Undelegate(nodeAddr, nowEpoch, undelegation)
					validator.InactivateIfSelfDelegationTooLittle()
				}
			} else {
				node.Terminate()
			}
			continue
		}

		if node.CanActivate(requireTotal, requireSelf) {
			if err := node.Activate(nowEpoch); err != nil {
				return err
			}
			mutateMap3Addrs = append(mutateMap3Addrs, nodeAddr)
		}
	}
	log.Info("New mutate map3 nodes", "addresses", mutateMap3Addrs)

	// TODO(ATLAS): write off chain data after inserting block
	// store active and terminated map3 addr to rawdb
	batch := chain.Database().NewBatch()
	rawdb.WriteMutateMap3Nodes(batch, header.Epoch.Uint64(), mutateMap3Addrs)
	rawdb.WriteRenewActiveMap3Nodes(batch, header.Epoch.Uint64(), renewActiveMap3Addrs)
	if err := batch.Write(); err != nil {
		return err
	}
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

func (f UndelegationReleaserFactory) Create(stateDB *state.StateDB, chain consensus.ChainReader) (UndelegationReleaser, error) {
	if chain.Config().Atlas == nil {
		return nil, errors.New("not support to undelegate")
	}
	if chain.Config().Atlas.RestakingEnable {
		return undelegationToMap3Node{
			stateDB:       stateDB,
			rewardHandler: core.RewardToMap3Node{StateDB: stateDB, Chain: chain},
		}, nil
	} else {
		return undelegationToBalance{
			stateDB:       stateDB,
			rewardHandler: core.RewardToBalance{StateDB: stateDB},
		}, nil
	}
}

type UndelegationReleaser interface {
	Release(redelegation *restaking.Storage_Redelegation_, fromValidator common.Address, epoch, blockNum *big.Int) (completed bool, err error)
}

type undelegationToBalance struct {
	stateDB       *state.StateDB
	rewardHandler core.RewardToBalance
}

func (u undelegationToBalance) Release(redelegation *restaking.Storage_Redelegation_, fromValidator common.Address,
	epoch, blockNum *big.Int) (completed bool, err error) {
	// return undelegation
	delegator := redelegation.DelegatorAddress().Value()
	undelegation := redelegation.Undelegation().Amount().Value()
	u.stateDB.AddBalance(delegator, undelegation)
	redelegation.Undelegation().Clear()

	// return reward if redelgation is empty
	if amt := redelegation.Amount().Value(); amt.Sign() == 0 {
		_, err := u.rewardHandler.HandleReward(redelegation, blockNum)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

type undelegationToMap3Node struct {
	stateDB       *state.StateDB
	chain         consensus.ChainReader
	rewardHandler core.RewardToMap3Node
}

func (u undelegationToMap3Node) Release(redelegation *restaking.Storage_Redelegation_, fromValidator common.Address,
	epoch, blockNum *big.Int) (completed bool, err error) {
	// clear undelegation
	redelegation.Undelegation().Clear()

	// return reward if redelgation is empty
	if amt := redelegation.Amount().Value(); amt.Sign() == 0 {
		_, err := u.rewardHandler.HandleReward(redelegation, blockNum)
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

func payoutUnmicrodelegations(header *types.Header, stateDB *state.StateDB) error {
	nowEpoch := header.Epoch
	// Payout undelegated/unlocked tokens
	for _, map3Addr := range stateDB.Map3NodeList() {
		node, err := stateDB.Map3NodeByAddress(map3Addr)
		if err != nil {
			return err
		}

		var toBeRemoved []common.Address
		for _, delegator := range node.Microdelegations().AllKeys() {
			md, ok := node.Microdelegations().Get(delegator)
			if !ok {
				return errMicrodelegationNotExist
			}

			if md.CanReleaseUndelegationAt(nowEpoch) {
				// payout unmicrodelegation
				amt := md.Undelegation().Amount().Value()
				completed := md.Amount().Value().Sign() == 0 &&
					md.PendingDelegation().Amount().Value().Sign() == 0
				if completed {
					amt.Add(amt, md.Reward().Value())
					toBeRemoved = append(toBeRemoved, delegator)
				} else {
					md.Undelegation().Clear()
				}
				stateDB.AddBalance(delegator, amt)
			}
		}

		for _, delegator := range toBeRemoved {
			node.Microdelegations().Remove(delegator)
			stateDB.Map3NodePool().RemoveDelegationIndex(delegator, map3Addr)
		}
	}
	log.Info("paid out unmicrodelegations", "epoch", nowEpoch.Uint64(), "block-number", header.Number.Uint64())
	return nil
}

// Withdraw unlocked tokens to the delegators' accounts
func payoutUnredelegations(header *types.Header, stateDB *state.StateDB, releaser UndelegationReleaser) error {
	nowEpoch, numBlock := header.Epoch, header.Number
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

			if redelegation.CanReleaseUndelegationAt(nowEpoch) {
				completed, err := releaser.Release(redelegation, validatorAddr, nowEpoch, numBlock)
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
	log.Info("paid out unredelegations", "epoch", nowEpoch.Uint64(), "block-number", header.Number.Uint64())
	return nil
}

func updateCommitteeForNextEpoch(header *types.Header, stateDB *state.StateDB) (*restaking.Committee_, error) {
	nextEpoch := big.NewInt(0).Add(header.Epoch, common.Big1)
	nextComm, err := committee.WithStakingEnabled.Compute(nextEpoch, stateDB)
	if err != nil {
		return nil, err
	}
	stateDB.ValidatorPool().UpdateCommittee(nextComm)
	return nextComm, nil
}

// accumulateRewardsAndCountSigs credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward
// This func also do IncrementValidatorSigningCounts for validators
func accumulateRewardsAndCountSigs(chain consensus.ChainReader, state *state.StateDB, header *types.Header) (reward.Reader, error) {
	if header.Number.Uint64() <= 1 {
		// genesis block has no parent to reward.
		return network.EmptyPayout, nil
	}

	// calculate total reward
	lastBlockNum := big.NewInt(0).Sub(header.Number, common.Big1)
	pool := network.NewRewardPool(state)
	totalReward := pool.TakeReward(lastBlockNum, chain.Config())

	// no reward
	if totalReward.Sign() <= 0 {
		return network.EmptyPayout, nil
	}

	payouts := []reward.Payout{}
	comm, payable, missing, err := ballotResult(chain, header) // for last block
	if err != nil {
		return network.EmptyPayout, err
	}

	if err := availability.IncrementValidatorSigningCounts(
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
	lastButOneBlockNum := new(big.Int).Sub(header.Number, common.Big2)
	for i := len(payable.Entrys) - 1; i >= 0; i-- {
		// what to do about share of those that didn't sign
		blsKey := payable.Entrys[i].BLSPublicKey
		voter := votingPower.Voters[blsKey]

		snapshot, err := chain.ReadValidatorSnapshotAtBlock(lastButOneBlockNum, voter.EarningAccount)
		if err != nil {
			return network.EmptyPayout, err
		}

		var due *big.Int
		if i == 0 { // Give out whatever leftover to the first voter/handle
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
	lastButOneBlockNum := new(big.Int).Sub(header.Number, common.Big2)
	comm, err := lookupCommitteeAtBlock(lastButOneBlockNum, bc)
	if err != nil {
		return nil, nil, nil, errors.Errorf("cannot read committee at %v", lastButOneBlockNum)
	}
	reader := availability.CommitBitmapReader{Header: header}
	_, payable, missing, err := availability.BallotResult(reader, comm)
	return comm, payable, missing, err
}

var (
	committeeCache, _     = lru.New(3)
	votingPowerCache, _   = lru.New(3)
	delegateShareCache, _ = lru.New(100)
)

func lookupCommitteeAtBlock(blockNum *big.Int, chain consensus.ChainReader) (*restaking.Committee_, error) {
	commStr, err := chain.ReadCommitteeAtBlock(blockNum)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("committee-%v", commStr.Epoch().Value().String())
	if c, ok := committeeCache.Get(key); ok {
		return c.(*restaking.Committee_), nil
	}

	comm, err := commStr.Load()
	if err != nil {
		return nil, err
	}

	// Put in cache
	committeeCache.Add(key, comm)

	return comm, nil
}

func lookupVotingPower(epoch *big.Int, comm *restaking.Committee_) (*votepower.Roster, error) {
	key := fmt.Sprintf("votingpower-%v", epoch.String())
	if v, ok := votingPowerCache.Get(key); ok {
		return v.(*votepower.Roster), nil
	}

	votingPower, err := votepower.Compute(comm)
	if err != nil {
		return nil, err
	}

	// Put in cache
	votingPowerCache.Add(key, votingPower)

	return votingPower, nil
}

// Lookup or compute the shares of stake for all delegators in a validator
func lookupDelegatorShares(
	epoch *big.Int, snapshot *restaking.Storage_ValidatorWrapper_,
) (map[common.Address]common.Dec, error) {
	valAddr := snapshot.Validator().ValidatorAddress().Value()
	key := fmt.Sprintf("delegatorshares-%s-%s", epoch.String(), valAddr.Hex())
	if d, ok := delegateShareCache.Get(key); ok {
		return d.(map[common.Address]common.Dec), nil
	}

	votingPower := map[common.Address]common.Dec{}
	totalDelegationDec := common.NewDecFromBigInt(snapshot.TotalDelegation().Value())
	if totalDelegationDec.IsZero() {
		log.Info("zero total delegation during AddReward delegation payout",
			"validator-snapshot", valAddr.Hex())
	} else {
		for _, key := range snapshot.Redelegations().AllKeys() {
			delegation, ok := snapshot.Redelegations().Get(key)
			if !ok {
				return nil, errValidatorNotExist
			}
			// NOTE percentage = <this_delegator_amount>/<total_delegation>
			percentage := common.NewDecFromBigInt(delegation.Amount().Value()).Quo(totalDelegationDec)
			votingPower[delegation.DelegatorAddress().Value()] = percentage
		}
	}

	// Put in cache
	delegateShareCache.Add(key, votingPower)

	return votingPower, nil
}

// ATLAS - END
