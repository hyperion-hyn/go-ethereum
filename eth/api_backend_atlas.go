package eth

import (
    "math/big"

    "github.com/harmony-one/harmony/consensus/quorum"
    "github.com/harmony-one/harmony/core/types"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/rpc"
    "github.com/ethereum/go-ethereum/shard"
    "github.com/ethereum/go-ethereum/shard/committee"
)

// SendStakingTx adds a staking transaction
func (b *EthAPIBackend) SendStakingTx(
    ctx context.Context,
    newStakingTx *staking.StakingTransaction) error {
    b.hmy.nodeAPI.AddPendingStakingTransaction(newStakingTx)
    return nil
}

// GetElectedValidatorAddresses returns the address of elected validators for current epoch
func (b *EthAPIBackend) GetElectedValidatorAddresses() []common.Address {
    list, _ := b.hmy.BlockChain().ReadElectedValidatorList()
    return list
}

// GetAllValidatorAddresses returns the up to date validator candidates for next epoch
func (b *EthAPIBackend) GetAllValidatorAddresses() []common.Address {
    return b.hmy.BlockChain().ValidatorCandidates()
}

var (
    zero = numeric.ZeroDec()
)

// GetValidatorInformation returns the information of validator
func (b *EthAPIBackend) GetValidatorInformation(
    addr common.Address,
) (*staking.ValidatorRPCEnchanced, error) {
    wrapper, err := b.hmy.BlockChain().ReadValidatorInformation(addr)
    if err != nil {
        s, _ := internal_common.AddressToBech32(addr)
        return nil, errors.Wrapf(err, "not found address in current state %s", s)
    }

    now := b.hmy.BlockChain().CurrentHeader().Epoch()
    inCommittee := now.Cmp(wrapper.LastEpochInCommittee) == 0
    defaultReply := &staking.ValidatorRPCEnchanced{
        CurrentlyInCommittee: inCommittee,
        Wrapper:              *wrapper,
        Performance:          nil,
        ComputedMetrics:      nil,
        TotalDelegated:       wrapper.TotalDelegation(),
        EPoSStatus: effective.ValidatorStatus(
            inCommittee, wrapper.Status == effective.Active,
        ).String(),
        Lifetime: &staking.AccumulatedOverLifetime{
            wrapper.BlockReward,
            wrapper.Counters,
            zero,
        },
    }

    snapshot, err := b.hmy.BlockChain().ReadValidatorSnapshotAtEpoch(
        now, addr,
    )

    if err != nil {
        return defaultReply, nil
    }

    computed := availability.ComputeCurrentSigning(
        snapshot, wrapper,
    )
    computed.BlocksLeftInEpoch = shard.Schedule.BlocksPerEpoch() - computed.ToSign.Uint64()

    stats, err := b.hmy.BlockChain().ReadValidatorStats(addr)
    if err != nil {
        return defaultReply, nil
    }

    defaultReply.Lifetime.APR = stats.APR

    if defaultReply.CurrentlyInCommittee {
        defaultReply.Performance = &staking.CurrentEpochPerformance{
            CurrentSigningPercentage: *computed,
        }
        defaultReply.ComputedMetrics = stats
    }

    return defaultReply, nil
}

// GetMedianRawStakeSnapshot ..
func (b *EthAPIBackend) GetMedianRawStakeSnapshot() (
    *committee.CompletedEPoSRound, error,
) {
    return committee.NewEPoSRound(b.hmy.BlockChain())
}

// GetTotalStakingSnapshot ..
func (b *EthAPIBackend) GetTotalStakingSnapshot() *big.Int {
    b.TotalStakingCache.Lock()
    defer b.TotalStakingCache.Unlock()
    if b.TotalStakingCache.BlockHeight != -1 &&
        b.TotalStakingCache.BlockHeight > int64(rpc.LatestBlockNumber)-20 {
        return b.TotalStakingCache.TotalStaking
    }
    b.TotalStakingCache.BlockHeight = int64(rpc.LatestBlockNumber)
    candidates := b.hmy.BlockChain().ValidatorCandidates()
    if len(candidates) == 0 {
        b.TotalStakingCache.TotalStaking = big.NewInt(0)
        return b.TotalStakingCache.TotalStaking
    }
    stakes := big.NewInt(0)
    for i := range candidates {
        snapshot, _ := b.hmy.BlockChain().ReadValidatorSnapshot(candidates[i])
        validator, _ := b.hmy.BlockChain().ReadValidatorInformation(candidates[i])
        if !committee.IsEligibleForEPoSAuction(snapshot, validator) {
            continue
        }
        for i := range validator.Delegations {
            stakes.Add(stakes, validator.Delegations[i].Amount)
        }
    }
    b.TotalStakingCache.TotalStaking = stakes
    return b.TotalStakingCache.TotalStaking
}

// GetDelegationsByValidator returns all delegation information of a validator
func (b *EthAPIBackend) GetDelegationsByValidator(validator common.Address) []*staking.Delegation {
    wrapper, err := b.hmy.BlockChain().ReadValidatorInformation(validator)
    if err != nil || wrapper == nil {
        return nil
    }
    delegations := []*staking.Delegation{}
    for i := range wrapper.Delegations {
        delegations = append(delegations, &wrapper.Delegations[i])
    }
    return delegations
}

// GetDelegationsByDelegator returns all delegation information of a delegator
func (b *EthAPIBackend) GetDelegationsByDelegator(
    delegator common.Address,
) ([]common.Address, []*staking.Delegation) {
    addresses := []common.Address{}
    delegations := []*staking.Delegation{}
    delegationIndexes, err := b.hmy.BlockChain().ReadDelegationsByDelegator(delegator)
    if err != nil {
        return nil, nil
    }

    for i := range delegationIndexes {
        wrapper, err := b.hmy.BlockChain().ReadValidatorInformation(
            delegationIndexes[i].ValidatorAddress,
        )
        if err != nil || wrapper == nil {
            return nil, nil
        }

        if uint64(len(wrapper.Delegations)) > delegationIndexes[i].Index {
            delegations = append(delegations, &wrapper.Delegations[delegationIndexes[i].Index])
        } else {
            delegations = append(delegations, nil)
        }
        addresses = append(addresses, delegationIndexes[i].ValidatorAddress)
    }
    return addresses, delegations
}

// GetValidatorSelfDelegation returns the amount of staking after applying all delegated stakes
func (b *EthAPIBackend) GetValidatorSelfDelegation(addr common.Address) *big.Int {
    wrapper, err := b.hmy.BlockChain().ReadValidatorInformation(addr)
    if err != nil || wrapper == nil {
        return nil
    }
    if len(wrapper.Delegations) == 0 {
        return nil
    }
    return wrapper.Delegations[0].Amount
}

// GetShardState ...
func (b *EthAPIBackend) GetShardState() (*shard.State, error) {
    return b.hmy.BlockChain().ReadShardState(b.hmy.BlockChain().CurrentHeader().Epoch())
}

// GetCurrentStakingErrorSink ..
func (b *EthAPIBackend) GetCurrentStakingErrorSink() []staking.RPCTransactionError {
    return b.hmy.nodeAPI.ErroredStakingTransactionSink()
}

// GetCurrentTransactionErrorSink ..
func (b *EthAPIBackend) GetCurrentTransactionErrorSink() []types.RPCTransactionError {
    return b.hmy.nodeAPI.ErroredTransactionSink()
}

// GetPendingCXReceipts ..
func (b *EthAPIBackend) GetPendingCXReceipts() []*types.CXReceiptsProof {
    return b.hmy.nodeAPI.PendingCXReceipts()
}

// GetCurrentUtilityMetrics ..
func (b *EthAPIBackend) GetCurrentUtilityMetrics() (*network.UtilityMetric, error) {
    return network.NewUtilityMetricSnapshot(b.hmy.BlockChain())
}

// GetSuperCommittees ..
func (b *EthAPIBackend) GetSuperCommittees() (*quorum.Transition, error) {
    nowE := b.hmy.BlockChain().CurrentHeader().Epoch()
    thenE := new(big.Int).Sub(nowE, common.Big1)

    var (
        nowCommittee, prevCommittee *shard.State
        err                         error
    )
    nowCommittee, err = b.hmy.BlockChain().ReadShardState(nowE)
    if err != nil {
        return nil, err
    }
    prevCommittee, err = b.hmy.BlockChain().ReadShardState(thenE)
    if err != nil {
        return nil, err
    }

    stakedSlotsNow, stakedSlotsThen :=
        shard.ExternalSlotsAvailableForEpoch(nowE),
        shard.ExternalSlotsAvailableForEpoch(thenE)

    then, now :=
        quorum.NewRegistry(stakedSlotsThen),
        quorum.NewRegistry(stakedSlotsNow)

    for _, comm := range prevCommittee.Shards {
        decider := quorum.NewDecider(quorum.SuperMajorityStake, comm.ShardID)
        if _, err := decider.SetVoters(&comm, prevCommittee.Epoch); err != nil {
            return nil, err
        }
        then.Deciders[fmt.Sprintf("shard-%d", comm.ShardID)] = decider
    }

    for _, comm := range nowCommittee.Shards {
        decider := quorum.NewDecider(quorum.SuperMajorityStake, comm.ShardID)
        if _, err := decider.SetVoters(&comm, nowCommittee.Epoch); err != nil {
            return nil, err
        }
        now.Deciders[fmt.Sprintf("shard-%d", comm.ShardID)] = decider
    }
    return &quorum.Transition{then, now}, nil
}

// GetCurrentBadBlocks ..
func (b *EthAPIBackend) GetCurrentBadBlocks() []core.BadBlock {
    return b.hmy.BlockChain().BadBlocks()
}

// GetLastCrossLinks ..
func (b *EthAPIBackend) GetLastCrossLinks() ([]*types.CrossLink, error) {
    crossLinks := []*types.CrossLink{}
    for i := uint32(1); i < shard.Schedule.InstanceForEpoch(b.CurrentBlock().Epoch()).NumShards(); i++ {
        link, err := b.hmy.BlockChain().ReadShardLastCrossLink(i)
        if err != nil {
            return nil, err
        }
        crossLinks = append(crossLinks, link)
    }

    return crossLinks, nil
}
