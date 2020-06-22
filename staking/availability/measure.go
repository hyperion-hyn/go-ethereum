package availability

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/committee"
	"github.com/ethereum/go-ethereum/staking/effective"
	staking "github.com/ethereum/go-ethereum/staking/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	measure = numeric.NewDec(2).Quo(numeric.NewDec(3))
	// ErrDivByZero ..
	ErrDivByZero = errors.New("toSign of availability cannot be 0, mistake in protocol")
)

// BlockSigners ..
func BlockSigners(
	bitmap []byte, parentCommittee *committee.Committee,
) (committee.SlotList, committee.SlotList, error) {
	committerKeys, err := parentCommittee.BLSPublicKeys()
	if err != nil {
		return nil, nil, err
	}
	mask, err := bls2.NewMask(committerKeys, nil)
	if err != nil {
		return nil, nil, err
	}
	if err := mask.SetMask(bitmap); err != nil {
		return nil, nil, err
	}

	payable, missing := committee.SlotList{}, committee.SlotList{}

	for idx, member := range parentCommittee.Slots {
		switch signed, err := mask.IndexEnabled(idx); true {
		case err != nil:
			return nil, nil, err
		case signed:
			payable = append(payable, member)
		default:
			missing = append(missing, member)
		}
	}
	return payable, missing, nil
}

// BallotResult returns
// (parentCommittee.Slots, payable, missings, err)
func BallotResult(
	header RoundHeader, parentCommittee *committee.Committee,
) (committee.SlotList, committee.SlotList, committee.SlotList, error) {
	payable, missing, err := BlockSigners(
		// TODO(ATLAS): get LastCommitBitmap from parent header extra
		header.LastCommitBitmap(), parentCommittee,
	)
	return parentCommittee.Slots, payable, missing, err
}

type signerKind struct {
	didSign   bool
	committee committee.SlotList
}

func bumpCount(
	bc Reader,
	state ValidatorState,
	signers []signerKind,
	stakedAddrSet map[common.Address]struct{},
) error {
	for _, subset := range signers {
		for i := range subset.committee {
			addr := subset.committee[i].EcdsaAddress
			// NOTE if the signer address is not part of the staked addrs,
			// then it must be a harmony operated node running,
			// hence keep on going
			if _, isAddrForStaked := stakedAddrSet[addr]; !isAddrForStaked {
				continue
			}

			wrapper, err := state.ValidatorWrapper(addr)
			if err != nil {
				return err
			}

			wrapper.Counters.NumBlocksToSign.Add(
				wrapper.Counters.NumBlocksToSign, common.Big1,
			)

			if subset.didSign {
				wrapper.Counters.NumBlocksSigned.Add(
					wrapper.Counters.NumBlocksSigned, common.Big1,
				)
			}
		}
	}

	return nil
}

// IncrementValidatorSigningCounts ..
func IncrementValidatorSigningCounts(
	bc Reader,
	staked *shard.StakedSlots,
	state ValidatorState,
	signers, missing committee.SlotList,
) error {
	return bumpCount(
		bc, state, []signerKind{{false, missing}, {true, signers}},
		staked.LookupSet,
	)
}

// ComputeCurrentSigning returns (signed, toSign, quotient, error)
func ComputeCurrentSigning(
	snapshot, wrapper *staking.ValidatorWrapper,
) *staking.Computed {
	statsNow, snapSigned, snapToSign :=
		wrapper.Counters,
		snapshot.Counters.NumBlocksSigned,
		snapshot.Counters.NumBlocksToSign

	signed, toSign :=
		new(big.Int).Sub(statsNow.NumBlocksSigned, snapSigned),
		new(big.Int).Sub(statsNow.NumBlocksToSign, snapToSign)

	computed := staking.NewComputed(
		signed, toSign, 0, numeric.ZeroDec(), true,
	)

	if toSign.Cmp(common.Big0) == 0 {
		return computed
	}

	if signed.Sign() == -1 {
		// Shouldn't happen
		log.Error("negative number of signed blocks")
	}

	if toSign.Sign() == -1 {
		// Shouldn't happen
		log.Error("negative number of blocks to sign")
	}

	s1, s2 := numeric.NewDecFromBigInt(signed), numeric.NewDecFromBigInt(toSign)
	computed.Percentage = s1.Quo(s2)
	computed.IsBelowThreshold = IsBelowSigningThreshold(computed.Percentage)
	return computed
}

// IsBelowSigningThreshold ..
func IsBelowSigningThreshold(quotient numeric.Dec) bool {
	return quotient.LTE(measure)
}

// ComputeAndMutateEPOSStatus sets the validator to
// inactive and thereby keeping it out of
// consideration in the pool of validators for
// whenever committee selection happens in future, the
// signing threshold is 66%
func ComputeAndMutateEPOSStatus(
	bc Reader,
	state ValidatorState,
	addr common.Address,
) error {
	log.Info("begin compute for availability")

	wrapper, err := state.ValidatorWrapper(addr)
	if err != nil {
		return err
	}
	if wrapper.Status == effective.Banned {
		utils.Logger().Debug().Msg("Can't update EPoS status on a banned validator")
		return nil
	}

	snapshot, err := bc.ReadValidatorSnapshot(wrapper.Address)
	if err != nil {
		return err
	}

	computed := ComputeCurrentSigning(snapshot.Validator, wrapper)

	log.Info("check if signing percent is meeting required threshold")

	const missedTooManyBlocks = true

	switch computed.IsBelowThreshold {
	case missedTooManyBlocks:
		wrapper.Status = effective.Inactive
		utils.Logger().Info().
			Str("threshold", measure.String()).
			Interface("computed", computed).
			Msg("validator failed availability threshold, set to inactive")
	default:
		// Default is no-op so validator who wants
		// to leave the committee can actually leave.
	}

	return nil
}
