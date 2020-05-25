package core

import (
    "bytes"
    "math/big"

    "github.com/ethereum/go-ethereum/block"
    "github.com/ethereum/go-ethereum/consensus/votepower"
    "github.com/ethereum/go-ethereum/core/rawdb"
    "github.com/ethereum/go-ethereum/core/state"
    "github.com/ethereum/go-ethereum/core/types"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethdb"
    "github.com/ethereum/go-ethereum/rlp"
    shard "github.com/ethereum/go-ethereum/shardxx"
    "github.com/ethereum/go-ethereum/staking/apr"
    "github.com/ethereum/go-ethereum/staking/effective"
    "github.com/ethereum/go-ethereum/staking/slash"
    staking "github.com/ethereum/go-ethereum/staking/types"
)

// ReadShardState retrieves sharding state given the epoch number.
func (bc *BlockChain) ReadShardState(epoch *big.Int) (*shard.State, error) {
    cacheKey := string(epoch.Bytes())
    if cached, ok := bc.shardStateCache.Get(cacheKey); ok {
        shardState := cached.(*shard.State)
        return shardState, nil
    }
    shardState, err := rawdb.ReadShardState(bc.db, epoch)
    if err != nil {
        return nil, err
    }
    bc.shardStateCache.Add(cacheKey, shardState)
    return shardState, nil
}

// WriteShardStateBytes saves the given sharding state under the given epoch number.
func (bc *BlockChain) WriteShardStateBytes(db rawdb.DatabaseWriter,
    epoch *big.Int, shardState []byte,
) (*shard.State, error) {
    decodeShardState, err := shard.DecodeWrapper(shardState)
    if err != nil {
        return nil, err
    }
    err = rawdb.WriteShardStateBytes(db, epoch, shardState)
    if err != nil {
        return nil, err
    }
    cacheKey := string(epoch.Bytes())
    bc.shardStateCache.Add(cacheKey, decodeShardState)
    return decodeShardState, nil
}

// ReadCommitSig retrieves the commit signature on a block.
func (bc *BlockChain) ReadCommitSig(blockNum uint64) ([]byte, error) {
    if cached, ok := bc.lastCommitsCache.Get("commitSig" + string(blockNum)); ok {
        lastCommits := cached.([]byte)
        return lastCommits, nil
    }
    lastCommits, err := rawdb.ReadBlockCommitSig(bc.db, blockNum)
    if err != nil {
        return nil, err
    }
    return lastCommits, nil
}

// WriteCommitSig saves the commits signatures signed on a block.
func (bc *BlockChain) WriteCommitSig(blockNum uint64, lastCommits []byte) error {
    err := rawdb.WriteBlockCommitSig(bc.db, blockNum, lastCommits)
    if err != nil {
        return err
    }
    bc.lastCommitsCache.Add("commitSig"+string(blockNum), lastCommits)
    return nil
}

// GetVdfByNumber retrieves the rand seed given the block number, return 0 if not exist
func (bc *BlockChain) GetVdfByNumber(number uint64) []byte {
    header := bc.GetHeaderByNumber(number)
    if header == nil {
        return []byte{}
    }

    return header.Vdf()
}

// GetVrfByNumber retrieves the randomness preimage given the block number, return 0 if not exist
func (bc *BlockChain) GetVrfByNumber(number uint64) []byte {
    header := bc.GetHeaderByNumber(number)
    if header == nil {
        return []byte{}
    }
    return header.Vrf()
}

// ChainDb returns the database
func (bc *BlockChain) ChainDb() ethdb.Database { return bc.db }

// GetEpochBlockNumber returns the first block number of the given epoch.
func (bc *BlockChain) GetEpochBlockNumber(epoch *big.Int) (*big.Int, error) {
    // Try cache first
    cacheKey := string(epoch.Bytes())
    if cachedValue, ok := bc.epochCache.Get(cacheKey); ok {
        return (&big.Int{}).SetBytes([]byte(cachedValue.(string))), nil
    }
    blockNum, err := rawdb.ReadEpochBlockNumber(bc.db, epoch)
    if err != nil {
        return nil, errors.Wrapf(
            err, "cannot read epoch block number from database",
        )
    }
    cachedValue := []byte(blockNum.Bytes())
    bc.epochCache.Add(cacheKey, cachedValue)
    return blockNum, nil
}

// StoreEpochBlockNumber stores the given epoch-first block number.
func (bc *BlockChain) StoreEpochBlockNumber(
    epoch *big.Int, blockNum *big.Int,
) error {
    cacheKey := string(epoch.Bytes())
    cachedValue := []byte(blockNum.Bytes())
    bc.epochCache.Add(cacheKey, cachedValue)
    if err := rawdb.WriteEpochBlockNumber(bc.db, epoch, blockNum); err != nil {
        return errors.Wrapf(
            err, "cannot write epoch block number to database",
        )
    }
    return nil
}

// ReadEpochVrfBlockNums retrieves block numbers with valid VRF for the specified epoch
func (bc *BlockChain) ReadEpochVrfBlockNums(epoch *big.Int) ([]uint64, error) {
    vrfNumbers := []uint64{}
    if cached, ok := bc.randomnessCache.Get("vrf-" + string(epoch.Bytes())); ok {
        encodedVrfNumbers := cached.([]byte)
        if err := rlp.DecodeBytes(encodedVrfNumbers, &vrfNumbers); err != nil {
            return nil, err
        }
        return vrfNumbers, nil
    }

    encodedVrfNumbers, err := rawdb.ReadEpochVrfBlockNums(bc.db, epoch)
    if err != nil {
        return nil, err
    }

    if err := rlp.DecodeBytes(encodedVrfNumbers, &vrfNumbers); err != nil {
        return nil, err
    }
    return vrfNumbers, nil
}

// WriteEpochVrfBlockNums saves block numbers with valid VRF for the specified epoch
func (bc *BlockChain) WriteEpochVrfBlockNums(epoch *big.Int, vrfNumbers []uint64) error {
    encodedVrfNumbers, err := rlp.EncodeToBytes(vrfNumbers)
    if err != nil {
        return err
    }

    err = rawdb.WriteEpochVrfBlockNums(bc.db, epoch, encodedVrfNumbers)
    if err != nil {
        return err
    }
    bc.randomnessCache.Add("vrf-"+string(epoch.Bytes()), encodedVrfNumbers)
    return nil
}

// ReadEpochVdfBlockNum retrieves block number with valid VDF for the specified epoch
func (bc *BlockChain) ReadEpochVdfBlockNum(epoch *big.Int) (*big.Int, error) {
    if cached, ok := bc.randomnessCache.Get("vdf-" + string(epoch.Bytes())); ok {
        encodedVdfNumber := cached.([]byte)
        return new(big.Int).SetBytes(encodedVdfNumber), nil
    }

    encodedVdfNumber, err := rawdb.ReadEpochVdfBlockNum(bc.db, epoch)
    if err != nil {
        return nil, err
    }
    return new(big.Int).SetBytes(encodedVdfNumber), nil
}

// WriteEpochVdfBlockNum saves block number with valid VDF for the specified epoch
func (bc *BlockChain) WriteEpochVdfBlockNum(epoch *big.Int, blockNum *big.Int) error {
    err := rawdb.WriteEpochVdfBlockNum(bc.db, epoch, blockNum.Bytes())
    if err != nil {
        return err
    }

    bc.randomnessCache.Add("vdf-"+string(epoch.Bytes()), blockNum.Bytes())
    return nil
}

// WriteCrossLinks saves the hashes of crosslinks by shardID and blockNum combination key
func (bc *BlockChain) WriteCrossLinks(batch rawdb.DatabaseWriter, cls []types.CrossLink) error {
    var err error
    for i := 0; i < len(cls); i++ {
        cl := cls[i]
        err = rawdb.WriteCrossLinkShardBlock(batch, cl.ShardID(), cl.BlockNum(), cl.Serialize())
    }
    return err
}

// DeleteCrossLinks removes the hashes of crosslinks by shardID and blockNum combination key
func (bc *BlockChain) DeleteCrossLinks(cls []types.CrossLink) error {
    var err error
    for i := 0; i < len(cls); i++ {
        cl := cls[i]
        err = rawdb.DeleteCrossLinkShardBlock(bc.db, cl.ShardID(), cl.BlockNum())
    }
    return err
}

// ReadCrossLink retrieves crosslink given shardID and blockNum.
func (bc *BlockChain) ReadCrossLink(shardID uint32, blockNum uint64) (*types.CrossLink, error) {
    bytes, err := rawdb.ReadCrossLinkShardBlock(bc.db, shardID, blockNum)
    if err != nil {
        return nil, err
    }
    crossLink, err := types.DeserializeCrossLink(bytes)

    return crossLink, err
}

// LastContinuousCrossLink saves the last crosslink of a shard
// This function will update the latest crosslink in the sense that
// any previous block's crosslink is received up to this point
// there is no missing hole between genesis to this crosslink of given shardID
func (bc *BlockChain) LastContinuousCrossLink(batch rawdb.DatabaseWriter, shardID uint32) error {
    oldLink, err := bc.ReadShardLastCrossLink(shardID)
    if oldLink == nil || err != nil {
        return err
    }
    newLink := oldLink
    // Starting from last checkpoint, keeping reading immediate next crosslink until there is a gap
    for i := oldLink.BlockNum() + 1; ; i++ {
        tmp, err := bc.ReadCrossLink(shardID, i)
        if err == nil && tmp != nil && tmp.BlockNum() == i {
            newLink = tmp
        } else {
            break
        }
    }

    if newLink.BlockNum() > oldLink.BlockNum() {
        utils.Logger().Debug().Msgf("LastContinuousCrossLink: latest checkpoint blockNum %d", newLink.BlockNum())
        return rawdb.WriteShardLastCrossLink(batch, shardID, newLink.Serialize())
    }
    return nil
}

// ReadShardLastCrossLink retrieves the last crosslink of a shard.
func (bc *BlockChain) ReadShardLastCrossLink(shardID uint32) (*types.CrossLink, error) {
    bytes, err := rawdb.ReadShardLastCrossLink(bc.db, shardID)
    if err != nil {
        return nil, err
    }
    return types.DeserializeCrossLink(bytes)
}

func (bc *BlockChain) writeSlashes(processed slash.Records) error {
    bytes, err := rlp.EncodeToBytes(processed)
    if err != nil {
        const msg = "failed to encode slashing candidates"
        utils.Logger().Error().Msg(msg)
        return err
    }
    if err := rawdb.WritePendingSlashingCandidates(bc.db, bytes); err != nil {
        return err
    }
    return nil
}

// DeleteFromPendingSlashingCandidates ..
func (bc *BlockChain) DeleteFromPendingSlashingCandidates(
    processed slash.Records,
) error {
    bc.pendingSlashingCandidatesMU.Lock()
    defer bc.pendingSlashingCandidatesMU.Unlock()
    current := bc.ReadPendingSlashingCandidates()
    bc.pendingSlashes = processed.SetDifference(current)
    return bc.writeSlashes(bc.pendingSlashes)
}

// ReadPendingSlashingCandidates retrieves pending slashing candidates
func (bc *BlockChain) ReadPendingSlashingCandidates() slash.Records {
    if !bc.Config().IsStaking(bc.CurrentHeader().Epoch()) {
        return slash.Records{}
    }
    return append(bc.pendingSlashes[0:0], bc.pendingSlashes...)
}

// ReadPendingCrossLinks retrieves pending crosslinks
func (bc *BlockChain) ReadPendingCrossLinks() ([]types.CrossLink, error) {
    bytes := []byte{}
    if cached, ok := bc.pendingCrossLinksCache.Get(pendingCLCacheKey); ok {
        bytes = cached.([]byte)
    } else {
        bytes, err := rawdb.ReadPendingCrossLinks(bc.db)
        if err != nil || len(bytes) == 0 {
            return nil, err
        }
    }
    cls := []types.CrossLink{}
    if err := rlp.DecodeBytes(bytes, &cls); err != nil {
        utils.Logger().Error().Err(err).Msg("Invalid pending crosslink RLP decoding")
        return nil, err
    }
    return cls, nil
}

// WritePendingCrossLinks saves the pending crosslinks
func (bc *BlockChain) WritePendingCrossLinks(crossLinks []types.CrossLink) error {
    // deduplicate crosslinks if any
    m := map[uint32]map[uint64](types.CrossLink){}
    for _, cl := range crossLinks {
        if _, ok := m[cl.ShardID()]; !ok {
            m[cl.ShardID()] = map[uint64](types.CrossLink){}
        }
        m[cl.ShardID()][cl.BlockNum()] = cl
    }

    cls := []types.CrossLink{}
    for _, m1 := range m {
        for _, cl := range m1 {
            cls = append(cls, cl)
        }
    }
    utils.Logger().Debug().Msgf("[WritePendingCrossLinks] Before Dedup has %d cls, after Dedup has %d cls", len(crossLinks), len(cls))

    bytes, err := rlp.EncodeToBytes(cls)
    if err != nil {
        utils.Logger().Error().Msg("[WritePendingCrossLinks] Failed to encode pending crosslinks")
        return err
    }
    if err := rawdb.WritePendingCrossLinks(bc.db, bytes); err != nil {
        return err
    }
    by, err := rlp.EncodeToBytes(cls)
    if err == nil {
        bc.pendingCrossLinksCache.Add(pendingCLCacheKey, by)
    }
    return nil

}

// AddPendingSlashingCandidates appends pending slashing candidates
func (bc *BlockChain) AddPendingSlashingCandidates(
    candidates slash.Records,
) error {
    bc.pendingSlashingCandidatesMU.Lock()
    defer bc.pendingSlashingCandidatesMU.Unlock()
    current := bc.ReadPendingSlashingCandidates()

    state, err := bc.State()
    if err != nil {
        return err
    }

    valid := slash.Records{}
    for i := range candidates {
        if err := slash.Verify(bc, state, &candidates[i]); err == nil {
            valid = append(valid, candidates[i])
        }
    }

    pendingSlashes := append(
        bc.pendingSlashes, current.SetDifference(valid)...,
    )

    if l, c := len(pendingSlashes), len(current); l > maxPendingSlashes {
        return errors.Wrapf(
            errExceedMaxPendingSlashes, "current %d with-additional %d", c, l,
        )
    }
    bc.pendingSlashes = pendingSlashes
    return bc.writeSlashes(bc.pendingSlashes)
}

// AddPendingCrossLinks appends pending crosslinks
func (bc *BlockChain) AddPendingCrossLinks(pendingCLs []types.CrossLink) (int, error) {
    bc.pendingCrossLinksMutex.Lock()
    defer bc.pendingCrossLinksMutex.Unlock()

    cls, err := bc.ReadPendingCrossLinks()
    if err != nil || len(cls) == 0 {
        err := bc.WritePendingCrossLinks(pendingCLs)
        return len(pendingCLs), err
    }
    cls = append(cls, pendingCLs...)
    err = bc.WritePendingCrossLinks(cls)
    return len(cls), err
}

// DeleteFromPendingCrossLinks delete pending crosslinks that already committed (i.e. passed in the params)
func (bc *BlockChain) DeleteFromPendingCrossLinks(crossLinks []types.CrossLink) (int, error) {
    bc.pendingCrossLinksMutex.Lock()
    defer bc.pendingCrossLinksMutex.Unlock()

    cls, err := bc.ReadPendingCrossLinks()
    if err != nil || len(cls) == 0 {
        return 0, err
    }

    m := map[uint32]map[uint64](struct{}){}
    for _, cl := range crossLinks {
        if _, ok := m[cl.ShardID()]; !ok {
            m[cl.ShardID()] = map[uint64](struct{}){}
        }
        m[cl.ShardID()][cl.BlockNum()] = struct{}{}
    }

    pendingCLs := []types.CrossLink{}

    for _, cl := range cls {
        if _, ok := m[cl.ShardID()]; ok {
            if _, ok1 := m[cl.ShardID()][cl.BlockNum()]; ok1 {
                continue
            }
        }
        pendingCLs = append(pendingCLs, cl)
    }
    err = bc.WritePendingCrossLinks(pendingCLs)
    return len(pendingCLs), err
}

// IsSameLeaderAsPreviousBlock retrieves a block from the database by number, caching it
func (bc *BlockChain) IsSameLeaderAsPreviousBlock(block *types.Block) bool {
    if IsEpochBlock(block) {
        return false
    }

    previousHeader := bc.GetHeaderByNumber(block.NumberU64() - 1)
    return block.Coinbase() == previousHeader.Coinbase()
}

// ChainDB ...
// TODO(ricl): in eth, this is not exposed. I expose it here because I need it in Harmony object.
// In eth, chainDB is initialized within Ethereum object
func (bc *BlockChain) ChainDB() ethdb.Database {
    return bc.db
}

// GetVMConfig returns the block chain VM config.
func (bc *BlockChain) GetVMConfig() *vm.Config {
    return &bc.vmConfig
}

// ReadCXReceipts retrieves the cross shard transaction receipts of a given shard
func (bc *BlockChain) ReadCXReceipts(shardID uint32, blockNum uint64, blockHash common.Hash) (types.CXReceipts, error) {
    cxs, err := rawdb.ReadCXReceipts(bc.db, shardID, blockNum, blockHash)
    if err != nil || len(cxs) == 0 {
        return nil, err
    }
    return cxs, nil
}

// CXMerkleProof calculates the cross shard transaction merkle proof of a given destination shard
func (bc *BlockChain) CXMerkleProof(toShardID uint32, block *types.Block) (*types.CXMerkleProof, error) {
    proof := &types.CXMerkleProof{BlockNum: block.Number(), BlockHash: block.Hash(), ShardID: block.ShardID(), CXReceiptHash: block.Header().OutgoingReceiptHash(), CXShardHashes: []common.Hash{}, ShardIDs: []uint32{}}

    epoch := block.Header().Epoch()
    shardingConfig := shard.Schedule.InstanceForEpoch(epoch)
    shardNum := int(shardingConfig.NumShards())

    for i := 0; i < shardNum; i++ {
        receipts, err := bc.ReadCXReceipts(uint32(i), block.NumberU64(), block.Hash())
        if err != nil || len(receipts) == 0 {
            continue
        } else {
            hash := types.DeriveSha(receipts)
            proof.CXShardHashes = append(proof.CXShardHashes, hash)
            proof.ShardIDs = append(proof.ShardIDs, uint32(i))
        }
    }
    if len(proof.ShardIDs) == 0 {
        return nil, nil
    }
    return proof, nil
}

// WriteCXReceiptsProofSpent mark the CXReceiptsProof list with given unspent status
// true: unspent, false: spent
func (bc *BlockChain) WriteCXReceiptsProofSpent(db rawdb.DatabaseWriter, cxps []*types.CXReceiptsProof) {
    for _, cxp := range cxps {
        rawdb.WriteCXReceiptsProofSpent(db, cxp)
    }
}

// IsSpent checks whether a CXReceiptsProof is unspent
func (bc *BlockChain) IsSpent(cxp *types.CXReceiptsProof) bool {
    shardID := cxp.MerkleProof.ShardID
    blockNum := cxp.MerkleProof.BlockNum.Uint64()
    by, _ := rawdb.ReadCXReceiptsProofSpent(bc.db, shardID, blockNum)
    return by == rawdb.SpentByte
}

// ReadTxLookupEntry returns where the given transaction resides in the chain,
// as a (block hash, block number, index in transaction list) triple.
// returns 0, 0 if not found
func (bc *BlockChain) ReadTxLookupEntry(txID common.Hash) (common.Hash, uint64, uint64) {
    return rawdb.ReadTxLookupEntry(bc.db, txID)
}

// ReadValidatorInformationAt reads staking
// information of given validatorWrapper at a specific state root
func (bc *BlockChain) ReadValidatorInformationAt(
    addr common.Address, root common.Hash,
) (*staking.ValidatorWrapper, error) {
    state, err := bc.StateAt(root)
    if err != nil || state == nil {
        return nil, errors.Wrapf(err, "at root: %s", root.Hex())
    }
    wrapper, err := state.ValidatorWrapper(addr)
    if err != nil {
        return nil, errors.Wrapf(err, "at root: %s", root.Hex())
    }
    return wrapper, nil
}

// ReadValidatorInformation reads staking information of given validator address
func (bc *BlockChain) ReadValidatorInformation(
    addr common.Address,
) (*staking.ValidatorWrapper, error) {
    return bc.ReadValidatorInformationAt(addr, bc.CurrentBlock().Root())
}

// ReadValidatorSnapshotAtEpoch reads the snapshot
// staking validator information of given validator address
func (bc *BlockChain) ReadValidatorSnapshotAtEpoch(
    epoch *big.Int,
    addr common.Address,
) (*staking.ValidatorSnapshot, error) {
    return rawdb.ReadValidatorSnapshot(bc.db, addr, epoch)
}

// ReadValidatorSnapshot reads the snapshot staking information of given validator address
func (bc *BlockChain) ReadValidatorSnapshot(
    addr common.Address,
) (*staking.ValidatorSnapshot, error) {
    epoch := bc.CurrentBlock().Epoch()
    key := addr.Hex() + epoch.String()
    if cached, ok := bc.validatorSnapshotCache.Get(key); ok {
        return cached.(*staking.ValidatorSnapshot), nil
    }
    return rawdb.ReadValidatorSnapshot(bc.db, addr, epoch)
}

// WriteValidatorSnapshot writes the snapshot of provided validator
func (bc *BlockChain) WriteValidatorSnapshot(
    batch rawdb.DatabaseWriter, snapshot *staking.ValidatorSnapshot,
) error {
    // Batch write the current data as snapshot
    if err := rawdb.WriteValidatorSnapshot(batch, snapshot.Validator, snapshot.Epoch); err != nil {
        return err
    }

    // Update cache
    key := snapshot.Validator.Address.Hex() + snapshot.Epoch.String()
    bc.validatorSnapshotCache.Add(key, snapshot)
    return nil
}

// ReadValidatorStats reads the stats of a validator
func (bc *BlockChain) ReadValidatorStats(
    addr common.Address,
) (*staking.ValidatorStats, error) {
    return rawdb.ReadValidatorStats(bc.db, addr)
}

// UpdateValidatorVotingPower writes the voting power for the committees
func (bc *BlockChain) UpdateValidatorVotingPower(
    batch rawdb.DatabaseWriter,
    block *types.Block,
    newEpochSuperCommittee, currentEpochSuperCommittee *shard.State,
    state *state.DB,
) (map[common.Address]*staking.ValidatorStats, error) {
    if newEpochSuperCommittee == nil {
        return nil, shard.ErrSuperCommitteeNil
    }

    validatorStats := map[common.Address]*staking.ValidatorStats{}

    existing, replacing :=
        currentEpochSuperCommittee.StakedValidators(),
        newEpochSuperCommittee.StakedValidators()

    // TODO could also keep track of the BLS keys which
    // lost a slot because just losing slots doesn't mean that the
    // validator was booted, just that some of their keys lost slots
    for currentValidator := range existing.LookupSet {
        if _, keptSlot := replacing.LookupSet[currentValidator]; !keptSlot {
            // NOTE Think carefully about when time comes to delete offchain things
            // TODO Someone: collect and then delete every 30 epochs
            // rawdb.DeleteValidatorSnapshot(
            // 	bc.db, currentValidator, currentEpochSuperCommittee.Epoch,
            // )
            // rawdb.DeleteValidatorStats(bc.db, currentValidator)
            stats, err := rawdb.ReadValidatorStats(bc.db, currentValidator)
            if err != nil {
                stats = staking.NewEmptyStats()
            }
            // This means it's already in staking epoch
            if currentEpochSuperCommittee.Epoch != nil {
                wrapper, err := state.ValidatorWrapper(currentValidator)
                if err != nil {
                    return nil, err
                }

                if slash.IsBanned(wrapper) {
                    stats.BootedStatus = effective.BannedForDoubleSigning
                } else if wrapper.Status == effective.Inactive {
                    stats.BootedStatus = effective.TurnedInactiveOrInsufficientUptime
                } else {
                    stats.BootedStatus = effective.LostEPoSAuction
                }
            }
            validatorStats[currentValidator] = stats
        }
    }

    rosters := make([]*votepower.Roster, len(newEpochSuperCommittee.Shards))
    for i := range newEpochSuperCommittee.Shards {
        subCommittee := &newEpochSuperCommittee.Shards[i]
        if newEpochSuperCommittee.Epoch == nil {
            return nil, errors.Wrapf(
                errNilEpoch,
                "block epoch %v current-committee-epoch %v",
                block.Epoch(),
                currentEpochSuperCommittee.Epoch,
            )
        }
        roster, err := votepower.Compute(subCommittee, newEpochSuperCommittee.Epoch)
        if err != nil {
            return nil, err
        }
        rosters[i] = roster
    }

    networkWide := votepower.AggregateRosters(rosters)
    for key, value := range networkWide {
        stats, err := rawdb.ReadValidatorStats(bc.db, key)
        if err != nil {
            stats = staking.NewEmptyStats()
        }
        total := numeric.ZeroDec()
        for i := range value {
            total = total.Add(value[i].EffectiveStake)
        }
        stats.TotalEffectiveStake = total
        earningWrapping := make([]staking.VoteWithCurrentEpochEarning, len(value))
        for i := range value {
            earningWrapping[i] = staking.VoteWithCurrentEpochEarning{
                Vote:   value[i],
                Earned: big.NewInt(0),
            }
        }
        stats.MetricsPerShard = earningWrapping

        // fetch raw-stake from snapshot and update per-key metrics
        if snapshot, err := bc.ReadValidatorSnapshotAtEpoch(
            newEpochSuperCommittee.Epoch, key,
        ); err == nil {
            wrapper := snapshot.Validator
            spread := numeric.ZeroDec()
            if len(wrapper.SlotPubKeys) > 0 {
                spread = numeric.NewDecFromBigInt(wrapper.TotalDelegation()).
                    QuoInt64(int64(len(wrapper.SlotPubKeys)))
            }
            for i := range stats.MetricsPerShard {
                stats.MetricsPerShard[i].Vote.RawStake = spread
            }
        }

        // This means it's already in staking epoch
        if currentEpochSuperCommittee.Epoch != nil {
            wrapper, err := state.ValidatorWrapper(key)
            if err != nil {
                return nil, err
            }

            if wrapper.Delegations[0].Amount.Cmp(common.Big0) > 0 {
                if aprComputed, err := apr.ComputeForValidator(
                    bc, block, wrapper,
                ); err != nil {
                    if errors.Cause(err) == apr.ErrInsufficientEpoch {
                        utils.Logger().Info().Err(err).Msg("apr could not be computed")
                    } else {
                        return nil, err
                    }
                } else {
                    now := currentEpochSuperCommittee.Epoch
                    // only insert if APR for current epoch does not exists
                    aprEntry := staking.APREntry{now, *aprComputed}
                    l := len(stats.APRs)
                    // first time inserting apr for validator or
                    // apr for current epoch does not exists
                    // check the last entry's epoch, if not same, insert
                    if l == 0 || stats.APRs[l-1].Epoch.Cmp(now) != 0 {
                        stats.APRs = append(stats.APRs, aprEntry)
                    }
                    // if history is more than staking.APRHistoryLength, pop front
                    if l > staking.APRHistoryLength {
                        stats.APRs = stats.APRs[1:]
                    }
                }
            } else {
                utils.Logger().Info().Msg("zero total delegation, skipping apr computation")
            }
        }
        validatorStats[key] = stats
    }

    return validatorStats, nil
}

// UpdateValidatorSnapshots updates the content snapshot of all validators
// Note: this should only be called within the blockchain insert process.
func (bc *BlockChain) UpdateValidatorSnapshots(
    batch rawdb.DatabaseWriter, epoch *big.Int, state *state.DB, newValidators []common.Address,
) error {
    // Note this is reading the validator list from last block.
    // It's fine since the new validators from this block is already snapshot when created.
    allValidators, err := bc.ReadValidatorList()
    if err != nil {
        return err
    }

    allValidators = append(allValidators, newValidators...)

    // Read all validator's current data and snapshot them
    for i := range allValidators {
        // The snapshot will be captured in the state after the last epoch block is finalized
        validator, err := state.ValidatorWrapper(allValidators[i])
        if err != nil {
            return err
        }

        snapshot := &staking.ValidatorSnapshot{validator, epoch}
        if err := bc.WriteValidatorSnapshot(batch, snapshot); err != nil {
            return err
        }
    }

    return nil
}

// ReadValidatorList reads the addresses of current all validators
func (bc *BlockChain) ReadValidatorList() ([]common.Address, error) {
    if cached, ok := bc.validatorListCache.Get("validatorList"); ok {
        by := cached.([]byte)
        m := []common.Address{}
        if err := rlp.DecodeBytes(by, &m); err != nil {
            return nil, err
        }
        return m, nil
    }
    return rawdb.ReadValidatorList(bc.db)
}

// WriteValidatorList writes the list of validator addresses to database
// Note: this should only be called within the blockchain insert process.
func (bc *BlockChain) WriteValidatorList(
    db rawdb.DatabaseWriter, addrs []common.Address,
) error {
    if err := rawdb.WriteValidatorList(db, addrs); err != nil {
        return err
    }
    bytes, err := rlp.EncodeToBytes(addrs)
    if err == nil {
        bc.validatorListCache.Add("validatorList", bytes)
    }
    return nil
}

// ReadDelegationsByDelegator reads the addresses of validators delegated by a delegator
func (bc *BlockChain) ReadDelegationsByDelegator(
    delegator common.Address,
) (m staking.DelegationIndexes, err error) {
    rawResult := staking.DelegationIndexes{}
    if cached, ok := bc.validatorListByDelegatorCache.Get(string(delegator.Bytes())); ok {
        by := cached.([]byte)
        if err := rlp.DecodeBytes(by, &rawResult); err != nil {
            return nil, err
        }
    } else {
        if rawResult, err = rawdb.ReadDelegationsByDelegator(bc.db, delegator); err != nil {
            return nil, err
        }
    }
    blockNum := bc.CurrentBlock().Number()
    for _, index := range rawResult {
        if index.BlockNum.Cmp(blockNum) <= 0 {
            m = append(m, index)
        } else {
            // Filter out index that's created beyond current height of chain.
            // This only happens when there is a chain rollback.
            utils.Logger().Warn().Msgf("Future delegation index encountered. Skip: %+v", index)
        }
    }
    return m, nil
}

// ReadDelegationsByDelegatorAt reads the addresses of validators delegated by a delegator at a given block
func (bc *BlockChain) ReadDelegationsByDelegatorAt(
    delegator common.Address, blockNum *big.Int,
) (m staking.DelegationIndexes, err error) {
    rawResult := staking.DelegationIndexes{}
    if cached, ok := bc.validatorListByDelegatorCache.Get(string(delegator.Bytes())); ok {
        by := cached.([]byte)
        if err := rlp.DecodeBytes(by, &rawResult); err != nil {
            return nil, err
        }
    } else {
        if rawResult, err = rawdb.ReadDelegationsByDelegator(bc.db, delegator); err != nil {
            return nil, err
        }
    }
    for _, index := range rawResult {
        if index.BlockNum.Cmp(blockNum) <= 0 {
            m = append(m, index)
        } else {
            // Filter out index that's created beyond current height of chain.
            // This only happens when there is a chain rollback.
            utils.Logger().Warn().Msgf("Future delegation index encountered. Skip: %+v", index)
        }
    }
    return m, nil
}

// writeDelegationsByDelegator writes the list of validator addresses to database
func (bc *BlockChain) writeDelegationsByDelegator(
    batch rawdb.DatabaseWriter,
    delegator common.Address,
    indices []staking.DelegationIndex,
) error {
    if err := rawdb.WriteDelegationsByDelegator(
        batch, delegator, indices,
    ); err != nil {
        return err
    }
    bytes, err := rlp.EncodeToBytes(indices)
    if err == nil {
        bc.validatorListByDelegatorCache.Add(string(delegator.Bytes()), bytes)
    }
    return nil
}

// UpdateStakingMetaData updates the metadata of validators and delegations,
// including the full validator list and delegation indexes.
// Note: this should only be called within the blockchain insert process.
func (bc *BlockChain) UpdateStakingMetaData(
    batch rawdb.DatabaseWriter, block *types.Block,
    state *state.DB, epoch, newEpoch *big.Int,
) (newValidators []common.Address, err error) {
    newValidators, newDelegations, err := bc.prepareStakingMetaData(block, state)
    if err != nil {
        utils.Logger().Warn().Msgf("oops, prepareStakingMetaData failed, err: %+v", err)
        return newValidators, err
    }

    if len(newValidators) > 0 {
        list, err := bc.ReadValidatorList()
        if err != nil {
            return newValidators, err
        }

        valMap := map[common.Address]struct{}{}
        for _, addr := range list {
            valMap[addr] = struct{}{}
        }

        newAddrs := []common.Address{}
        for _, addr := range newValidators {
            if _, ok := valMap[addr]; !ok {
                newAddrs = append(newAddrs, addr)
            }

            // Update validator snapshot for the new validator
            validator, err := state.ValidatorWrapper(addr)
            if err != nil {
                return newValidators, err
            }

            if err := bc.WriteValidatorSnapshot(batch, &staking.ValidatorSnapshot{validator, epoch}); err != nil {
                return newValidators, err
            }
            // For validator created at exactly the last block of an epoch, we should create the snapshot
            // for next epoch too.
            if newEpoch.Cmp(epoch) > 0 {
                if err := bc.WriteValidatorSnapshot(batch, &staking.ValidatorSnapshot{validator, newEpoch}); err != nil {
                    return newValidators, err
                }
            }
        }

        // Update validator list
        list = append(list, newAddrs...)
        if err = bc.WriteValidatorList(batch, list); err != nil {
            return newValidators, err
        }
    }

    for addr, delegations := range newDelegations {
        if err := bc.writeDelegationsByDelegator(batch, addr, delegations); err != nil {
            return newValidators, err
        }
    }
    return newValidators, nil
}

// prepareStakingMetaData prepare the updates of validator's
// and the delegator's meta data according to staking transaction.
// The following return values are cached end state to be written to DB.
// The reason for the cached state is to solve the issue that batch DB changes
// won't be reflected immediately so the intermediary state can't be read from DB.
// newValidators - the addresses of the newly created validators
// newDelegations - the map of delegator address and their updated delegation indexes
func (bc *BlockChain) prepareStakingMetaData(
    block *types.Block, state *state.DB,
) (newValidators []common.Address,
    newDelegations map[common.Address]staking.DelegationIndexes,
    err error,
) {
    newDelegations = map[common.Address]staking.DelegationIndexes{}
    blockNum := block.Number()
    for _, txn := range block.StakingTransactions() {
        payload, err := txn.RLPEncodeStakeMsg()
        if err != nil {
            return nil, nil, err
        }
        decodePayload, err := staking.RLPDecodeStakeMsg(payload, txn.StakingType())
        if err != nil {
            return nil, nil, err
        }

        switch txn.StakingType() {
        case staking.DirectiveCreateValidator:
            createValidator := decodePayload.(*staking.CreateValidator)
            newList, appended := utils.AppendIfMissing(
                newValidators, createValidator.ValidatorAddress,
            )
            if !appended {
                return nil, nil, errValidatorExist
            }
            newValidators = newList

            // Add self delegation into the index
            selfIndex := staking.DelegationIndex{
                createValidator.ValidatorAddress,
                uint64(0),
                blockNum,
            }
            delegations, ok := newDelegations[createValidator.ValidatorAddress]
            if !ok {
                // If the cache doesn't have it, load it from DB for the first time.
                delegations, err = bc.ReadDelegationsByDelegator(createValidator.ValidatorAddress)
                if err != nil {
                    return nil, nil, err
                }
            }

            delegations = append(delegations, selfIndex)
            newDelegations[createValidator.ValidatorAddress] = delegations
        case staking.DirectiveEditValidator:
        case staking.DirectiveDelegate:
            delegate := decodePayload.(*staking.Delegate)

            delegations, ok := newDelegations[delegate.DelegatorAddress]
            if !ok {
                // If the cache doesn't have it, load it from DB for the first time.
                delegations, err = bc.ReadDelegationsByDelegator(delegate.DelegatorAddress)
                if err != nil {
                    return nil, nil, err
                }
            }
            if delegations, err = bc.addDelegationIndex(
                delegations, delegate.DelegatorAddress, delegate.ValidatorAddress, state, blockNum,
            ); err != nil {
                return nil, nil, err
            }
            newDelegations[delegate.DelegatorAddress] = delegations
        case staking.DirectiveUndelegate:
        case staking.DirectiveCollectRewards:
        default:
        }
    }

    return newValidators, newDelegations, nil
}

// ReadBlockRewardAccumulator must only be called on beaconchain
func (bc *BlockChain) ReadBlockRewardAccumulator(number uint64) (*big.Int, error) {
    if !bc.chainConfig.IsStaking(shard.Schedule.CalcEpochNumber(number)) {
        return big.NewInt(0), nil
    }
    if cached, ok := bc.blockAccumulatorCache.Get(number); ok {
        return cached.(*big.Int), nil
    }
    return rawdb.ReadBlockRewardAccumulator(bc.db, number)
}

// WriteBlockRewardAccumulator directly writes the BlockRewardAccumulator value
// Note: this should only be called once during staking launch.
func (bc *BlockChain) WriteBlockRewardAccumulator(
    batch rawdb.DatabaseWriter, reward *big.Int, number uint64,
) error {
    if err := rawdb.WriteBlockRewardAccumulator(
        batch, reward, number,
    ); err != nil {
        return err
    }
    bc.blockAccumulatorCache.Add(number, reward)
    return nil
}

// UpdateBlockRewardAccumulator ..
// Note: this should only be called within the blockchain insert process.
func (bc *BlockChain) UpdateBlockRewardAccumulator(
    batch rawdb.DatabaseWriter, diff *big.Int, number uint64,
) error {
    current, err := bc.ReadBlockRewardAccumulator(number - 1)
    if err != nil {
        // one-off fix for pangaea, return after pangaea enter staking.
        current = big.NewInt(0)
        bc.WriteBlockRewardAccumulator(batch, current, number)
    }
    return bc.WriteBlockRewardAccumulator(batch, new(big.Int).Add(current, diff), number)
}

// Note this should read from the state of current block in concern (root == newBlock.root)
func (bc *BlockChain) addDelegationIndex(
    delegations staking.DelegationIndexes,
    delegatorAddress, validatorAddress common.Address, state *state.DB, blockNum *big.Int,
) (staking.DelegationIndexes, error) {
    // If there is an existing delegation, just return
    validatorAddressBytes := validatorAddress.Bytes()
    for _, delegation := range delegations {
        if bytes.Equal(delegation.ValidatorAddress[:], validatorAddressBytes[:]) {
            return delegations, nil
        }
    }

    // Found the delegation from state and add the delegation index
    // Note this should read from the state of current block in concern
    wrapper, err := state.ValidatorWrapper(validatorAddress)
    if err != nil {
        return delegations, err
    }
    for i := range wrapper.Delegations {
        if bytes.Equal(
            wrapper.Delegations[i].DelegatorAddress[:], delegatorAddress[:],
        ) {
            // TODO(audit): change the way of indexing if we allow delegation deletion.
            delegations = append(delegations, staking.DelegationIndex{
                validatorAddress,
                uint64(i),
                blockNum,
            })
        }
    }
    return delegations, nil
}

// ValidatorCandidates returns the up to date validator candidates for next epoch
func (bc *BlockChain) ValidatorCandidates() []common.Address {
    list, err := bc.ReadValidatorList()
    if err != nil {
        return make([]common.Address, 0)
    }
    return list
}

// DelegatorsInformation returns up to date information of delegators of a given validator address
func (bc *BlockChain) DelegatorsInformation(addr common.Address) []*staking.Delegation {
    return make([]*staking.Delegation, 0)
}

// GetECDSAFromCoinbase retrieve corresponding ecdsa address from Coinbase Address
func (bc *BlockChain) GetECDSAFromCoinbase(header *block.Header) (common.Address, error) {
    // backward compatibility: before isStaking epoch, coinbase address is the ecdsa address
    coinbase := header.Coinbase()
    isStaking := bc.Config().IsStaking(header.Epoch())
    if !isStaking {
        return coinbase, nil
    }

    shardState, err := bc.ReadShardState(header.Epoch())
    if err != nil {
        return common.Address{}, errors.Wrapf(
            err, "cannot read shard state",
        )
    }

    committee, err := shardState.FindCommitteeByID(header.ShardID())
    if err != nil {
        return common.Address{}, errors.Wrapf(
            err, "cannot find shard in the shard state",
        )
    }
    for _, member := range committee.Slots {
        // After staking the coinbase address will be the address of bls public key
        if bytes.Equal(member.EcdsaAddress[:], coinbase[:]) {
            return member.EcdsaAddress, nil
        }

        if utils.GetAddressFromBLSPubKeyBytes(member.BLSPublicKey[:]) == coinbase {
            return member.EcdsaAddress, nil
        }
    }
    return common.Address{}, errors.Errorf(
        "cannot find corresponding ECDSA Address for coinbase %s",
        header.Coinbase().Hash().Hex(),
    )
}

// SuperCommitteeForNextEpoch ...
// isVerify=true means validators use it to verify
// isVerify=false means leader is to propose
func (bc *BlockChain) SuperCommitteeForNextEpoch(
    beacon consensus_engine.ChainReader,
    header *block.Header,
    isVerify bool,
) (*shard.State, error) {
    var (
        nextCommittee = new(shard.State)
        err           error
        beaconEpoch   = new(big.Int)
        shardState    = shard.State{}
    )
    switch header.ShardID() {
    case shard.BeaconChainShardID:
        if shard.Schedule.IsLastBlock(header.Number().Uint64()) {
            nextCommittee, err = committee.WithStakingEnabled.Compute(
                new(big.Int).Add(header.Epoch(), common.Big1),
                beacon,
            )
        }
    default:
        // TODO: needs to make sure beacon chain sync works.
        if isVerify {
            //verify
            shardState, err = header.GetShardState()
            if err != nil {
                return &shard.State{}, err
            }
            // before staking epoch
            if shardState.Epoch == nil {
                beaconEpoch = new(big.Int).Add(header.Epoch(), common.Big1)
            } else { // after staking epoch
                beaconEpoch = shardState.Epoch
            }
        } else {
            //propose
            beaconEpoch = beacon.CurrentHeader().Epoch()
        }
        utils.Logger().Debug().Msgf("[SuperCommitteeCalculation] isVerify: %+v, realBeaconEpoch:%+v, beaconEpoch: %+v, headerEpoch:%+v, shardStateEpoch:%+v",
            isVerify, beacon.CurrentHeader().Epoch(), beaconEpoch, header.Epoch(), shardState.Epoch)
        nextEpoch := new(big.Int).Add(header.Epoch(), common.Big1)
        if bc.Config().IsStaking(nextEpoch) {
            // If next epoch is staking epoch, I should wait and listen for beacon chain for epoch changes
            switch beaconEpoch.Cmp(header.Epoch()) {
            case 1:
                // If beacon chain is bigger than shard chain in epoch, it means I should catch up with beacon chain now
                nextCommittee, err = committee.WithStakingEnabled.ReadFromDB(
                    beaconEpoch, beacon,
                )

                utils.Logger().Debug().
                    Uint64("blockNum", header.Number().Uint64()).
                    Uint64("myCurEpoch", header.Epoch().Uint64()).
                    Uint64("beaconEpoch", beaconEpoch.Uint64()).
                    Msg("Propose new epoch as beacon chain's epoch")
            case 0:
                // If it's same epoch, no need to propose new shard state (new epoch change)
            case -1:
                // If beacon chain is behind, shard chain should wait for the beacon chain by not changing epochs.
            }
        } else {
            if bc.Config().IsStaking(beaconEpoch) {
                // If I am not even in the last epoch before staking epoch and beacon chain is already in staking epoch,
                // I should just catch up with beacon chain's epoch
                nextCommittee, err = committee.WithStakingEnabled.ReadFromDB(
                    beaconEpoch, beacon,
                )

                utils.Logger().Debug().
                    Uint64("blockNum", header.Number().Uint64()).
                    Uint64("myCurEpoch", header.Epoch().Uint64()).
                    Uint64("beaconEpoch", beaconEpoch.Uint64()).
                    Msg("Propose entering staking along with beacon chain's epoch")
            } else {
                // If I are not in staking nor has beacon chain proposed a staking-based shard state,
                // do pre-staking committee calculation
                if shard.Schedule.IsLastBlock(header.Number().Uint64()) {
                    nextCommittee, err = committee.WithStakingEnabled.Compute(
                        nextEpoch,
                        bc,
                    )
                }
            }
        }

    }
    return nextCommittee, err
}
