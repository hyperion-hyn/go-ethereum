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
	lastEpoch := big.NewInt(0).Sub(epoch, common.Big1)
	eligibleCandidate, err := prepareOrders(stakedReader, lastEpoch)
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

func prepareOrders(stakedReader DataProvider, lastEpoch *big.Int) (map[common.Address]*effective.SlotOrder, error) {
	candidates := stakedReader.ValidatorList() // latest validators
	essentials := map[common.Address]*effective.SlotOrder{}
	totalStaked, tempZero := big.NewInt(0), common.ZeroDec()

	for i := range candidates {
		validator, err := stakedReader.ValidatorByAddress(candidates[i])
		if err != nil {
			return nil, err
		}
		// snapshot of validator at the beginning of the given epoch
		snapshot, err := stakedReader.ReadValidatorAtEpoch(lastEpoch, candidates[i])
		if err != nil {
			return nil, err
		}
		if !IsEligibleForEPoSAuction(snapshot, validator, lastEpoch) {
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
func IsEligibleForEPoSAuction(snapshot, validator *restaking.Storage_ValidatorWrapper_, lastEpoch *big.Int) bool {
	// This original condition to check whether a validator is in last committee is not stable
	// because cross-links may arrive after the epoch ends and it still got counted into the
	// NumBlocksToSign, making this condition to be true when the validator is actually not in committee
	//if snapshot.Counters.NumBlocksToSign.Cmp(validator.Counters.NumBlocksToSign) != 0 {

	// Check whether the validator is in current committee
	if validator.Validator().LastEpochInCommittee().Value().Cmp(lastEpoch) == 0 {
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
	switch validator.Validator().Status().Value() {
	case uint8(restaking.Active):
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

	// Set the epoch
	committee.Epoch = big.NewInt(0).Set(epoch)
	return committee, nil
}

// ReadFromDB is a wrapper on GetCommitteeAtEpoch
func (def stakingEnabled) ReadFromDB(epoch *big.Int, reader DataProvider) (*restaking.Committee_, error) {
	committee, err := reader.ReadCommitteeAtEpoch(epoch)
	if err != nil {
		return nil, err
	}
	return committee.Load(), err
}

// Compute is single entry point for computing a new committee for next epoch
func (def stakingEnabled) Compute(epoch *big.Int, stakerReader DataProvider) (newComm *restaking.Committee_, err error) {
	committee, err := eposStakedCommittee(epoch, stakerReader)
	if err != nil {
		return nil, err
	}
	log.Info("computed new committee", "computed-for-epoch", epoch.Uint64())
	return committee, nil
}
