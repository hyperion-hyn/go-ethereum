package availability

import (
	"github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	measure = common.NewDec(2).Quo(common.NewDec(3))
	// ErrDivByZero ..
	ErrDivByZero = errors.New("toSign of availability cannot be 0, mistake in protocol")
)

// BlockSigners ..
func BlockSigners(
	bitmap []byte, parentCommittee *restaking.Committee_,
) (*restaking.Slots_, *restaking.Slots_, error) {
	committerKeys, err := parentCommittee.BLSPublicKeys()
	if err != nil {
		return nil, nil, err
	}
	mask, err := bls.NewMask(committerKeys, nil)
	if err != nil {
		return nil, nil, err
	}
	if err := mask.SetMask(bitmap); err != nil {
		return nil, nil, err
	}

	payable, missing := restaking.Slots_{}, restaking.Slots_{}

	for idx, member := range parentCommittee.Slots.Entrys {
		switch signed, err := mask.IndexEnabled(idx); true {
		case err != nil:
			return nil, nil, err
		case signed:
			payable.Entrys = append(payable.Entrys, member)
		default:
			missing.Entrys = append(missing.Entrys, member)
		}
	}
	return &payable, &missing, nil
}

// BallotResult returns
// (parentCommittee.Slots, payable, missings, err)
func BallotResult(
	header RoundHeader, parentCommittee *restaking.Committee_,
) (*restaking.Slots_, *restaking.Slots_, *restaking.Slots_, error) {
	payable, missing, err := BlockSigners(
		header.LastCommitBitmap(), parentCommittee,
	)
	return &parentCommittee.Slots, payable, missing, err
}

type signerKind struct {
	didSign   bool
	committee *restaking.Slots_
}

func bumpCount(
	bc Reader,
	state ValidatorState,
	signers []signerKind,
	stakedAddrSet map[common.Address]struct{},
) error {
	for _, subset := range signers {
		for i := range subset.committee.Entrys {
			addr := subset.committee.Entrys[i].EcdsaAddress
			// NOTE if the signer address is not part of the staked addrs,
			// then it must be a harmony operated node running,
			// hence keep on going
			if _, isAddrForStaked := stakedAddrSet[addr]; !isAddrForStaked {
				continue
			}

			wrapper, err := state.ValidatorByAddress(addr)
			if err != nil {
				return err
			}

			wrapper.Counters().NumBlocksToSign().SetValue(big.NewInt(0).Add(
				wrapper.Counters().NumBlocksToSign().Value(), common.Big1,
			))

			if subset.didSign {
				wrapper.Counters().NumBlocksSigned().SetValue(big.NewInt(0).Add(
					wrapper.Counters().NumBlocksSigned().Value(), common.Big1,
				))
			}
		}
	}

	return nil
}

// IncrementValidatorSigningCounts ..
func IncrementValidatorSigningCounts(
	bc Reader,
	staked *restaking.StakedSlots,
	state ValidatorState,
	signers, missing *restaking.Slots_,
) error {
	return bumpCount(
		bc, state, []signerKind{{false, missing}, {true, signers}},
		staked.LookupSet,
	)
}

// ComputeCurrentSigning returns (signed, toSign, quotient, error)
func ComputeCurrentSigning(
	snapshot, wrapper *restaking.Storage_ValidatorWrapper_,
) *Computed {
	statsNow, snapSigned, snapToSign :=
		wrapper.Counters(),
		snapshot.Counters().NumBlocksSigned().Value(),
		snapshot.Counters().NumBlocksToSign().Value()

	signed, toSign :=
		new(big.Int).Sub(statsNow.NumBlocksSigned().Value(), snapSigned),
		new(big.Int).Sub(statsNow.NumBlocksToSign().Value(), snapToSign)

	computed := NewComputed(
		signed, toSign, 0, common.ZeroDec(), true,
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

	s1, s2 := common.NewDecFromBigInt(signed), common.NewDecFromBigInt(toSign)
	computed.Percentage = s1.Quo(s2)
	computed.IsBelowThreshold = IsBelowSigningThreshold(computed.Percentage)
	return computed
}

// IsBelowSigningThreshold ..
func IsBelowSigningThreshold(quotient common.Dec) bool {
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
	epoch *big.Int,
) error {
	log.Info("begin compute for availability")

	wrapper, err := state.ValidatorByAddress(addr)
	if err != nil {
		return err
	}
	if wrapper.Validator().Status().Value() == uint8(restaking.Banned) {
		log.Debug("Can't update EPoS status on a banned validator")
		return nil
	}

	snapshot, err := bc.ReadValidatorAtEpoch(epoch, addr)
	if err != nil {
		return err
	}

	computed := ComputeCurrentSigning(snapshot, wrapper)

	log.Info("check if signing percent is meeting required threshold")

	const missedTooManyBlocks = true

	switch computed.IsBelowThreshold {
	case missedTooManyBlocks:
		wrapper.Validator().Status().SetValue(uint8(restaking.Inactive))
		log.Info("validator failed availability threshold, set to inactive", "threshold", measure.String())
	default:
		// Default is no-op so validator who wants
		// to leave the committee can actually leave.
	}
	return nil
}

// Computed represents current epoch
// availability measures, mostly for RPC
type Computed struct {
	Signed            *big.Int   `json:"current-epoch-signed"`
	ToSign            *big.Int   `json:"current-epoch-to-sign"`
	BlocksLeftInEpoch uint64     `json:"-"`
	Percentage        common.Dec `json:"current-epoch-signing-percentage"`
	IsBelowThreshold  bool       `json:"-"`
}

// NewComputed ..
func NewComputed(
	signed, toSign *big.Int,
	blocksLeft uint64,
	percent common.Dec,
	isBelowNow bool) *Computed {
	return &Computed{signed, toSign, blocksLeft, percent, isBelowNow}
}