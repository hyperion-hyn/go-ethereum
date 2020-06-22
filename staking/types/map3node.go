package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
)

type Map3NodeStatus byte

const (
	Nil Map3NodeStatus = iota
	// Pending means total delegation of this node is still not enough
	Pending
	// Active means allowed in serving map data
	Active
	// Inactive means this node did not keep staking after last stake was due
	Inactive
	// Banned records whether this node is banned from the network
	// because it faked map data it can never be undone
	Banned
)

const (
	// LockPeriodInEpoch is the number of epochs a undelegated token needs to be before it's released to the delegator's balance
	Map3NodeLockPeriodInEpoch = 180
)

type PublicKey []byte

type Map3Node struct {
	// ECDSA address of the map3 node
	Map3NodeAddr common.Address `json:"map3node-address"`
	// map3 node's initiator
	InitiatorAddr common.Address `json:"initiator-address"`
	// The public keys of the map3 node (used for communication and identification)
	PublicKeys []PublicKey
	// map3 node's self declared minimum self delegation
	MinSelfDelegation *big.Int `json:"min-self-delegation"`
	// maximum total delegation allowed
	MaxTotalDelegation *big.Int `json:"max-total-delegation"`
	// commission parameters
	Commission
	// description for the validator
	Description
}

// Map3NodeWrapper contains map3 node,
// its micro-delegation information
type Map3NodeWrapper struct {
	Map3Node
	Microdelegations Microdelegations
	RedelegationReference     RedelegationReference
	// All the rewarded accumulated so far
	AccumulatedReward *big.Int
	NodeState         NodeState
}

type NodeState struct {
	// map3 node statue
	Status Map3NodeStatus
	// Node age
	NodeAge *big.Int
	// CreationHeight is the height of creation
	CreationHeight *big.Int
	ActiveHeight   *big.Int
	ReleaseHeight  *big.Int
}

// SanityCheck checks basic requirements of a validator
func (n *Map3Node) SanityCheck(oneThirdExtrn int) error {
	if _, err := n.EnsureLength(); err != nil {
		return err
	}

	if len(n.SlotPubKeys) == 0 {
		return errNeedAtLeastOneSlotKey
	}

	if c := len(n.SlotPubKeys); oneThirdExtrn != DoNotEnforceMaxBLS &&
		c > oneThirdExtrn {
		return errors.Wrapf(
			ErrExcessiveBLSKeys, "have: %d allowed: %d",
			c, oneThirdExtrn,
		)
	}

	if n.MinSelfDelegation == nil {
		return errNilMinSelfDelegation
	}

	if n.MaxTotalDelegation == nil {
		return errNilMaxTotalDelegation
	}

	// MinSelfDelegation must be >= 10000 ONE
	if n.MinSelfDelegation.Cmp(minimumStake) < 0 {
		return errors.Wrapf(
			errMinSelfDelegationTooSmall,
			"delegation-given %s", n.MinSelfDelegation.String(),
		)
	}

	// MaxTotalDelegation must not be less than MinSelfDelegation
	if n.MaxTotalDelegation.Cmp(n.MinSelfDelegation) < 0 {
		return errors.Wrapf(
			errInvalidMaxTotalDelegation,
			"max-total-delegation %s min-self-delegation %s",
			n.MaxTotalDelegation.String(),
			n.MinSelfDelegation.String(),
		)
	}

	if n.Rate.LT(zeroPercent) || n.Rate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "rate:%s", n.Rate.String(),
		)
	}

	if n.MaxRate.LT(zeroPercent) || n.MaxRate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "max rate:%s", n.MaxRate.String(),
		)
	}

	if n.MaxChangeRate.LT(zeroPercent) || n.MaxChangeRate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "max change rate:%s", n.MaxChangeRate.String(),
		)
	}

	if n.Rate.GT(n.MaxRate) {
		return errors.Wrapf(
			errCommissionRateTooLarge,
			"rate:%s max rate:%s", n.Rate.String(), n.MaxRate.String(),
		)
	}

	if n.MaxChangeRate.GT(n.MaxRate) {
		return errors.Wrapf(
			errCommissionRateTooLarge,
			"rate:%s max change rate:%s", n.Rate.String(), n.MaxChangeRate.String(),
		)
	}

	allKeys := map[shard.BLSPublicKey]struct{}{}
	for i := range n.SlotPubKeys {
		if _, ok := allKeys[n.SlotPubKeys[i]]; !ok {
			allKeys[n.SlotPubKeys[i]] = struct{}{}
		} else {
			return errDuplicateSlotKeys
		}
	}
	return nil
}

// TotalDelegation - return the total amount of token in delegation
func (w *Map3NodeWrapper) TotalDelegation() *big.Int {
	total := big.NewInt(0)
	for _, entry := range w.Microdelegations {
		total.Add(total, entry.Amount)
	}
	return total
}
