package effective

import (
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
	"sort"
)

// medium.com/harmony-one/introducing-harmonys-effective-proof-of-stake-epos-2d39b4b8d58
var (
	two       = common.NewDecFromBigInt(big.NewInt(2))
	c, _      = common.NewDecFromStr("0.15")
	onePlusC  = common.OneDec().Add(c)
	oneMinusC = common.OneDec().Sub(c)
)

func effectiveStake(median, actual common.Dec) common.Dec {
	left := common.MinDec(onePlusC.Mul(median), actual)
	right := oneMinusC.Mul(median)
	return common.MaxDec(left, right)
}

// SlotPurchase ..
type SlotPurchase struct {
	Addr      common.Address
	Key       restaking.BLSPublicKey_
	RawStake  common.Dec
	EPoSStake common.Dec
}

// MarshalJSON ..
func (p SlotPurchase) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Addr      string     `json:"slot-owner"`
		Key       string     `json:"bls-public-key"`
		RawStake  common.Dec `json:"raw-stake"`
		EPoSStake common.Dec `json:"eposed-stake"`
	}{
		p.Addr.Hex(),
		p.Key.Hex(),
		p.RawStake,
		p.EPoSStake,
	})
}

// SlotOrder ..
type SlotOrder struct {
	Stake       *big.Int                 `json:"stake"`
	SpreadAmong restaking.BLSPublicKeys_ `json:"keys-at-auction"`
	Percentage  common.Dec               `json:"percentage-of-total-auction-stake"`
}

// Median ..
func Median(stakes []SlotPurchase) common.Dec {
	if len(stakes) == 0 {
		return common.ZeroDec()
	}

	sort.SliceStable(
		stakes,
		func(i, j int) bool {
			return stakes[i].RawStake.GT(stakes[j].RawStake)
		},
	)
	const isEven = 0
	switch l := len(stakes); l % 2 {
	case isEven:
		left := (l / 2) - 1
		right := l / 2
		return stakes[left].RawStake.Add(stakes[right].RawStake).Quo(two)
	default:
		return stakes[l/2].RawStake
	}
}

// Compute ..
func Compute(
	shortHand map[common.Address]*SlotOrder, pull int,
) (common.Dec, []SlotPurchase) {
	eposedSlots := []SlotPurchase{}
	if len(shortHand) == 0 {
		return common.ZeroDec(), eposedSlots
	}

	type t struct {
		addr common.Address
		slot *SlotOrder
	}

	shorter := []t{}
	for key, value := range shortHand {
		shorter = append(shorter, t{key, value})
	}

	sort.SliceStable(
		shorter,
		func(i, j int) bool {
			return bytes.Compare(
				shorter[i].addr.Bytes(), shorter[j].addr.Bytes(),
			) == -1
		},
	)

	// Expand
	for _, staker := range shorter {
		slotsCount := len(staker.slot.SpreadAmong.Keys)
		if slotsCount == 0 {
			continue
		}
		spread := common.NewDecFromBigInt(staker.slot.Stake).
			QuoInt64(int64(slotsCount))
		for i := 0; i < slotsCount; i++ {
			eposedSlots = append(eposedSlots, SlotPurchase{
				Addr: staker.addr,
				Key:  *staker.slot.SpreadAmong.Keys[i],
				// NOTE these are same because later the .EPoSStake mutated
				RawStake:  spread,
				EPoSStake: spread,
			})
		}
	}

	sort.SliceStable(
		eposedSlots,
		func(i, j int) bool {
			return eposedSlots[i].RawStake.GT(eposedSlots[j].RawStake)
		},
	)

	if l := len(eposedSlots); l < pull {
		pull = l
	}
	picks := eposedSlots[:pull]

	if len(picks) == 0 {
		return common.ZeroDec(), []SlotPurchase{}
	}

	return Median(picks), picks

}

// Apply ..
func Apply(shortHand map[common.Address]*SlotOrder, pull int) (
	common.Dec, []SlotPurchase,
) {
	median, picks := Compute(shortHand, pull)
	for i := range picks {
		picks[i].EPoSStake = effectiveStake(median, picks[i].RawStake)
	}

	return median, picks
}
