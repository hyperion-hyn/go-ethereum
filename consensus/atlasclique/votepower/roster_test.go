package votepower

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/harmony-one/bls/ffi/go/bls"
	"math/big"
	"math/rand"
	"testing"
)

var (
	slotList    restaking.Slots_
	totalStake  common.Dec
	stakedNodes = 10

	maxAccountGen = int64(98765654323123134)
	accountGen    = rand.New(rand.NewSource(1337))
	maxKeyGen     = int64(98765654323123134)
	keyGen        = rand.New(rand.NewSource(42))
	maxStakeGen   = int64(200)
	stakeGen      = rand.New(rand.NewSource(541))
)

func init() {
	totalStake = common.ZeroDec()
	for j := 0; j < stakedNodes; j++ {
		newSlot := generateRandomSlot()
		totalStake = totalStake.Add(newSlot.EffectiveStake)
		slotList.Entrys = append(slotList.Entrys, newSlot)
	}
}

func generateRandomSlot() *restaking.Slot_ {
	addr := common.Address{}
	addr.SetBytes(big.NewInt(int64(accountGen.Int63n(maxAccountGen))).Bytes())
	secretKey := bls.SecretKey{}
	secretKey.Deserialize(big.NewInt(int64(keyGen.Int63n(maxKeyGen))).Bytes())
	key := restaking.BLSPublicKey_{}
	key.FromLibBLSPublicKey(secretKey.GetPublicKey())
	stake := common.NewDecFromBigInt(big.NewInt(int64(stakeGen.Int63n(maxStakeGen))))
	return &restaking.Slot_{EcdsaAddress: addr, BLSPublicKey: key, EffectiveStake: stake}
}

func TestCompute(t *testing.T) {
	expectedRoster := NewRoster()
	// Calculated when generated
	expectedRoster.TotalEffectiveStake = totalStake
	staked := slotList.Entrys
	for i := range staked {
		member := PureStakedVote{
			EarningAccount: staked[i].EcdsaAddress,
			Identity:       staked[i].BLSPublicKey,
			OverallPercent: staked[i].EffectiveStake.Quo(expectedRoster.TotalEffectiveStake),
			EffectiveStake: staked[i].EffectiveStake,
		}
		expectedRoster.Voters[staked[i].BLSPublicKey] = &member
	}

	computedRoster, err := Compute(&restaking.Committee_{
		Epoch: big.NewInt(3), Slots: slotList,
	})
	if err != nil {
		t.Error("Computed Roster failed on vote summation to one")
	}

	if !compareRosters(expectedRoster, computedRoster, t) {
		t.Errorf("Compute Roster mismatch with expected Roster")
	}
}

func compareRosters(a, b *Roster, t *testing.T) bool {
	voterMatch := true
	for k, voter := range a.Voters {
		if other, exists := b.Voters[k]; exists {
			if !compareStakedVoter(voter, other) {
				t.Error("voter slot not match")
				voterMatch = false
			}
		} else {
			t.Error("computed roster missing")
			voterMatch = false
		}
	}
	return a.TotalEffectiveStake.Equal(b.TotalEffectiveStake) && voterMatch
}

func compareStakedVoter(a, b *PureStakedVote) bool {
	return a.EarningAccount == b.EarningAccount &&
		a.OverallPercent.Equal(b.OverallPercent) &&
		a.EffectiveStake.Equal(b.EffectiveStake)
}
