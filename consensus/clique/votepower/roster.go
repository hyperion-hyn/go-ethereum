package votepower

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/committee"
	"github.com/ethereum/go-ethereum/staking/types"
	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/pkg/errors"
	"math/big"
)

var (
	// ErrVotingPowerNotEqualOne ..
	ErrVotingPowerNotEqualOne = errors.New("voting power not equal to one")
)

// Ballot is a vote cast by a validator
type Ballot struct {
	SignerPubKey    types.BLSPublicKey `json:"bls-public-key"`
	BlockHeaderHash common.Hash        `json:"block-header-hash"`
	Signature       []byte             `json:"bls-signature"`
	Height          uint64             `json:"block-height"`
	ViewID          uint64             `json:"view-id"`
}

// MarshalJSON ..
func (b Ballot) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		A string `json:"bls-public-key"`
		B string `json:"block-header-hash"`
		C string `json:"bls-signature"`
		E uint64 `json:"block-height"`
		F uint64 `json:"view-id"`
	}{
		b.SignerPubKey.Hex(),
		b.BlockHeaderHash.Hex(),
		hex.EncodeToString(b.Signature),
		b.Height,
		b.ViewID,
	})
}

// Round is a round of voting in any FBFT phase
type Round struct {
	AggregatedVote *bls.Sign
	BallotBox      map[types.BLSPublicKey]*Ballot
}

func (b Ballot) String() string {
	data, _ := json.Marshal(b)
	return string(data)
}

// NewRound ..
func NewRound() *Round {
	return &Round{
		AggregatedVote: &bls.Sign{},
		BallotBox:      map[types.BLSPublicKey]*Ballot{},
	}
}

// PureStakedVote ..
type PureStakedVote struct {
	EarningAccount common.Address     `json:"earning-account"`
	Identity       types.BLSPublicKey `json:"bls-public-key"`
	OverallPercent numeric.Dec        `json:"overall-percent"`
	EffectiveStake numeric.Dec        `json:"effective-stake"`
	RawStake       numeric.Dec        `json:"raw-stake"`
}

// Roster ..
type Roster struct {
	Voters              map[types.BLSPublicKey]*PureStakedVote
	TotalEffectiveStake numeric.Dec
}

func (r Roster) String() string {
	s, _ := json.Marshal(r)
	return string(s)
}

// Compute creates a new roster based off the shard.SlotList
func Compute(comm *committee.Committee, epoch *big.Int) (*Roster, error) {
	if epoch == nil {
		return nil, errors.New("nil epoch for roster compute")
	}
	roster, staked := NewRoster(), comm.Slots

	for i := range staked {
		if e := staked[i].EffectiveStake; e != nil {
			roster.TotalEffectiveStake = roster.TotalEffectiveStake.Add(*e)
		}
	}

	var lastStakedVoter *PureStakedVote
	totalPercentage := numeric.ZeroDec()
	for i := range staked {
		member := PureStakedVote{
			EarningAccount: staked[i].EcdsaAddress,
			Identity:       staked[i].BLSPublicKey,
			OverallPercent: numeric.ZeroDec(),
			EffectiveStake: numeric.ZeroDec(),
			RawStake:       numeric.ZeroDec(),
		}

		// Real Staker
		if e := staked[i].EffectiveStake; e != nil {
			member.EffectiveStake = member.EffectiveStake.Add(*e)
			member.OverallPercent = e.Quo(roster.TotalEffectiveStake)
			totalPercentage = totalPercentage.Add(member.OverallPercent)
			lastStakedVoter = &member
		}

		if _, ok := roster.Voters[staked[i].BLSPublicKey]; !ok {
			roster.Voters[staked[i].BLSPublicKey] = &member
		} else {
			log.Debug("Duplicate BLS key found", "blsKey", staked[i].BLSPublicKey.Hex())
		}
	}

	// NOTE Enforce voting power sums to one,
	// give diff (expect tiny amt) to last staked voter
	if diff := numeric.OneDec().Sub(totalPercentage); !diff.IsZero() && lastStakedVoter != nil {
		lastStakedVoter.OverallPercent = lastStakedVoter.OverallPercent.Add(diff)
	}
	return roster, nil
}

// NewRoster ..
func NewRoster() *Roster {
	m := map[types.BLSPublicKey]*PureStakedVote{}
	return &Roster{
		Voters: m,
		TotalEffectiveStake:             numeric.ZeroDec(),
	}
}
