package committee

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/staking/availability"
	"github.com/ethereum/go-ethereum/staking/effective"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/pkg/errors"
	"math/big"
)

const (
	CommitteeSize = 88
)

// Committee Provider ..
type Provider interface {
	Compute(epoch *big.Int, reader DataProvider) (*restaking.Committee_, error)
	ReadFromDB(epoch *big.Int, reader DataProvider) (*restaking.Committee_, error)
}

// CandidatesForEPoS ..
type CandidatesForEPoS struct {
	Orders                             map[common.Address]effective.SlotOrder
	OpenSlotCountForExternalValidators int
}

// CompletedEPoSRound ..
type CompletedEPoSRound struct {
	MedianStake       common.Dec               `json:"epos-median-stake"`
	AuctionWinners    []effective.SlotPurchase `json:"epos-slot-winners"`
	AuctionCandidates []*CandidateOrder        `json:"epos-slot-candidates"`
}

// CandidateOrder ..
type CandidateOrder struct {
	*effective.SlotOrder
	StakePerKey *big.Int
	Validator   common.Address
}

// NewEPoSRound runs a fresh computation of EPoS using
// latest data always
func NewEPoSRound(epoch *big.Int, stakedReader DataProvider) (*CompletedEPoSRound, error) {
	eligibleCandidate, err := prepareOrders(epoch, stakedReader)
	if err != nil {
		return nil, err
	}
	median, winners := effective.Apply(
		eligibleCandidate, CommitteeSize,
	)
	auctionCandidates := make([]*CandidateOrder, len(eligibleCandidate))

	i := 0
	for key := range eligibleCandidate {
		// NOTE in principle, a div-by-zero should not
		// happen by this point but the risk of not being explicit about
		// checking is a panic, so the check is worth it
		perKey := big.NewInt(0)
		if l := len(eligibleCandidate[key].SpreadAmong.Keys); l > 0 {
			perKey.Set(
				new(big.Int).Div(
					eligibleCandidate[key].Stake, big.NewInt(int64(l)),
				),
			)
		}
		auctionCandidates[i] = &CandidateOrder{
			SlotOrder:   eligibleCandidate[key],
			StakePerKey: perKey,
			Validator:   key,
		}
		i++
	}

	return &CompletedEPoSRound{
		MedianStake:       median,
		AuctionWinners:    winners,
		AuctionCandidates: auctionCandidates,
	}, nil
}

func prepareOrders(epoch *big.Int, stakedReader DataProvider) (map[common.Address]*effective.SlotOrder, error) {
	candidates := stakedReader.ValidatorList()
	essentials := map[common.Address]*effective.SlotOrder{}
	totalStaked, tempZero := big.NewInt(0), common.ZeroDec()

	for i := range candidates {
		validator, err := stakedReader.ValidatorByAddress(candidates[i])
		if err != nil {
			return nil, err
		}
		snapshot, err := stakedReader.ReadValidatorAtEpoch(epoch, candidates[i])
		if err != nil {
			return nil, err
		}
		if !IsEligibleForEPoSAuction(snapshot, validator, epoch) {
			continue
		}

		validatorStake := validator.TotalDelegation().Value()
		totalStaked.Add(totalStaked, validatorStake)
		keys := validator.Validator().SlotPubKeys().Load()

		essentials[validator.Validator().ValidatorAddress().Value()] = &effective.SlotOrder{
			Stake:       validatorStake,
			SpreadAmong: *keys,
			Percentage:  tempZero,
		}
	}
	totalStakedDec := common.NewDecFromBigInt(totalStaked)

	for _, value := range essentials {
		value.Percentage = common.NewDecFromBigInt(value.Stake).Quo(totalStakedDec)
	}

	return essentials, nil
}

// IsEligibleForEPoSAuction ..
func IsEligibleForEPoSAuction(snapshot, validator *restaking.Storage_ValidatorWrapper_, epoch *big.Int) bool {
	// This original condition to check whether a validator is in last committee is not stable
	// because cross-links may arrive after the epoch ends and it still got counted into the
	// NumBlocksToSign, making this condition to be true when the validator is actually not in committee
	//if snapshot.Counters.NumBlocksToSign.Cmp(validator.Counters.NumBlocksToSign) != 0 {

	// Check whether the validator is in current committee
	if validator.Validator().LastEpochInCommittee().Value().Cmp(epoch) == 0 {
		// validator was in last epoch's committee
		// validator with below-threshold signing activity won't be considered for next epoch
		// and their status will be turned to inactive in FinalizeNewBlock
		computed := availability.ComputeCurrentSigning(snapshot, validator)
		if computed.IsBelowThreshold {
			return false
		}
	}
	// For validators who were not in last epoch's committee
	// or for those who were and signed enough blocks,
	// the decision is based on the status
	switch validator.Validator().Status().Value().Uint64() {
	case uint64(effective.Active):
		return true
	default:
		return false
	}
}

// StakingCandidatesReader ..
type StakingCandidatesReader interface {
	ValidatorByAddress(validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error)
	ValidatorList() []common.Address
}

// ChainReader is a subset of Engine.ChainReader, just enough to do assignment
type ChainReader interface {
	ReadValidatorAtEpoch(epoch *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error)
	// GetCommitteeAtEpoch retrieves sharding state given the epoch number.
	// This api reads the shard state cached or saved on the chaindb.
	// Thus, only should be used to read the shard state of the current chain.
	ReadCommitteeAtEpoch(epoch *big.Int) (*restaking.Storage_Committee_, error)
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig
	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *types.Header
}

// DataProvider ..
type DataProvider interface {
	StakingCandidatesReader
	ChainReader
}

type ChainReaderWithPendingState struct {
	ChainReader
	*state.StateDB
}

type stakingEnabled struct{}

var (
	// WithStakingEnabled ..
	WithStakingEnabled Provider = stakingEnabled{}
	// ErrComputeForEpochInPast ..
	ErrComputeForEpochInPast = errors.New("cannot compute for epoch in past")
)

func eposStakedCommittee(epoch *big.Int, stakerReader DataProvider) (*restaking.Committee_, error) {
	completedEPoSRound, err := NewEPoSRound(epoch, stakerReader)
	if err != nil {
		return nil, err
	}

	committee := &restaking.Committee_{}
	for i := range completedEPoSRound.AuctionWinners {
		purchasedSlot := completedEPoSRound.AuctionWinners[i]
		committee.Slots.Entrys = append(
			committee.Slots.Entrys, &restaking.Slot_{
				EcdsaAddress:   purchasedSlot.Addr,
				BLSPublicKey:   purchasedSlot.Key,
				EffectiveStake: purchasedSlot.EPoSStake,
			},
		)
	}

	// Set the epoch of shard state
	committee.Epoch = big.NewInt(0).Set(epoch) // TODO: epoch + 1?
	return committee, nil
}

// ReadFromDB is a wrapper on GetCommitteeAtEpoch
func (def stakingEnabled) ReadFromDB(epoch *big.Int, reader DataProvider) (newSuperComm *restaking.Committee_, err error) {
	// TODO(storage): read committee by epoch from statedb
	return reader.ReadCommitteeAtEpoch(epoch)
}

// Compute is single entry point for
// computing a new super committee, aka new shard state
func (def stakingEnabled) Compute(epoch *big.Int, stakerReader DataProvider) (newSuperComm *restaking.Committee_, err error) {
	committee, err := eposStakedCommittee(epoch, stakerReader)
	if err != nil {
		return nil, err
	}
	log.Info("computed new super committee", "computed-for-epoch", epoch.Uint64())
	return committee, nil
}
