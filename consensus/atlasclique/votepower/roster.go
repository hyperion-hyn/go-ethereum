package votepower

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/pkg/errors"
)

var (
	// ErrVotingPowerNotEqualOne ..
	ErrVotingPowerNotEqualOne = errors.New("voting power not equal to one")
)

// Ballot is a vote cast by a validator
type Ballot struct {
	SignerPubKey    restaking.BLSPublicKey_ `json:"bls-public-key"`
	BlockHeaderHash common.Hash             `json:"block-header-hash"`
	Signature       []byte                  `json:"bls-signature"`
	Height          uint64                  `json:"block-height"`
	ViewID          uint64                  `json:"view-id"`
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
	BallotBox      map[restaking.BLSPublicKey_]*Ballot
}

func (b Ballot) String() string {
	data, _ := json.Marshal(b)
	return string(data)
}

// NewRound ..
func NewRound() *Round {
	return &Round{
		AggregatedVote: &bls.Sign{},
		BallotBox:      map[restaking.BLSPublicKey_]*Ballot{},
	}
}

// PureStakedVote ..
type PureStakedVote struct {
	EarningAccount common.Address          `json:"earning-account"`
	Identity       restaking.BLSPublicKey_ `json:"bls-public-key"`
	OverallPercent common.Dec              `json:"overall-percent"`
	EffectiveStake common.Dec              `json:"effective-stake"`
	RawStake       common.Dec              `json:"raw-stake"`
}

// Roster ..
type Roster struct {
	Voters              map[restaking.BLSPublicKey_]*PureStakedVote
	TotalEffectiveStake common.Dec
}

func (r Roster) String() string {
	s, _ := json.Marshal(r)
	return string(s)
}

// Compute creates a new roster based off the shard.SlotList
func Compute(comm *restaking.Committee_) (*Roster, error) {
	roster, staked := NewRoster(), comm.Slots

	for i := range staked.Entrys {
		if e := staked.Entrys[i].EffectiveStake; e.IsPositive() {
			roster.TotalEffectiveStake = roster.TotalEffectiveStake.Add(e)
		}
	}

	var lastStakedVoter *PureStakedVote
	totalPercentage := common.ZeroDec()
	for i := range staked.Entrys {
		member := PureStakedVote{
			EarningAccount: staked.Entrys[i].EcdsaAddress,
			Identity:       staked.Entrys[i].BLSPublicKey,
			OverallPercent: common.ZeroDec(),
			EffectiveStake: common.ZeroDec(),
			RawStake:       common.ZeroDec(),
		}

		// Real Staker
		if e := staked.Entrys[i].EffectiveStake; e.IsPositive() {
			member.EffectiveStake = member.EffectiveStake.Add(e)
			member.OverallPercent = e.Quo(roster.TotalEffectiveStake)
			totalPercentage = totalPercentage.Add(member.OverallPercent)
			lastStakedVoter = &member
		}

		if _, ok := roster.Voters[staked.Entrys[i].BLSPublicKey]; !ok {
			roster.Voters[staked.Entrys[i].BLSPublicKey] = &member
		} else {
			log.Debug("Duplicate BLS key found", "blsKey", staked.Entrys[i].BLSPublicKey.Hex())
		}
	}

	// NOTE Enforce voting power sums to one,
	// give diff (expect tiny amt) to last staked voter
	if diff := common.OneDec().Sub(totalPercentage); !diff.IsZero() && lastStakedVoter != nil {
		lastStakedVoter.OverallPercent = lastStakedVoter.OverallPercent.Add(diff)
	}
	return roster, nil
}

// NewRoster ..
func NewRoster() *Roster {
	m := map[restaking.BLSPublicKey_]*PureStakedVote{}
	return &Roster{
		Voters:              m,
		TotalEffectiveStake: common.ZeroDec(),
	}
}
