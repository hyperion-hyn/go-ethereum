package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
	"math/big"
)

type Status byte

const (
	Nil Status = iota
	// Pending means total delegation of this node is still not enough
	Pending
	// Active means allowed in serving map data
	Active
	// Inactive means this node did not keep staking after last stake was due
	Inactive
	// Banned records whether this node is banned from the network
	// because it faked map data it can never be undone
	Banned
	// Dividing means this node is divided from another node but staking not enough
	Dividing
)

const (
	// LockPeriodInEpoch is the number of epochs a undelegated token needs to be before it's released to the delegator's balance
	Map3NodeRenewalPeriodInEpoch = 7
	Million                      = 1000000
	MaxPubKeyAllowed             = 1
)

var (
	Map3NodeLockPeriodInEpoch = numeric.NewDec(180)
	MinSelfDelegation         = numeric.NewDecWithPrec(20, 2) // 20%
	MinDelegation             = numeric.NewDecWithPrec(1, 2)  // 1%
	baseMinTotalNodeStake     = numeric.NewDecFromBigInt(new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(Million)))
	minimumMap3NodeStake      = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(Million))
)

type Map3NodeKey []byte // TODO: fix size?

func (pk Map3NodeKey) Hex() string {
	return common.Bytes2Hex(pk[:])
}

type Map3NodeKeys []Map3NodeKey

type Map3Node struct {
	// ECDSA address of the map3 node
	NodeAddress common.Address `json:"map3node-address"`
	// map3 node's initiator
	InitiatorAddress common.Address `json:"initiator-address"`
	// The node keys of the map3 node (used for communication and identification)
	NodeKeys Map3NodeKeys
	// commission parameters
	Commission Commission
	// description for the validator
	Description Description
	//
	DividedFrom common.Address
}

// Map3NodeWrapper contains map3 node,
// its micro-delegation information
type Map3NodeWrapper struct {
	Map3Node               Map3Node
	Microdelegations       Microdelegations
	RedelegationReference  common.Address
	AccumulatedReward      *big.Int // All the rewarded accumulated so far
	NodeState              NodeState
	TotalDelegation        *big.Int
	TotalPendingDelegation *big.Int
}

type Map3NodeWrappers map[common.Address]Map3NodeWrapper

// Map3NodeSnapshot contains map3 node snapshot and the corresponding epoch
type Map3NodeSnapshot struct {
	Map3Nodes Map3NodeWrappers
	Epoch     *big.Int
}

type Map3NodeSnapshotByEpoch map[uint64]Map3NodeSnapshot

type NodeState struct {
	Status          Status   // map3 node statue
	NodeAge         *big.Int // Node age
	CreationHeight  *big.Int
	ActivationEpoch *big.Int
	ReleaseEpoch    numeric.Dec
}

type AddressSet map[common.Address]struct{}
type PubKeySet map[string]struct{} // node key hex
type Map3NodeAddressSetByDelegator map[common.Address]AddressSet

type Map3NodePool struct {
	Nodes                     Map3NodeWrappers
	NodeSnapshotByEpoch       Map3NodeSnapshotByEpoch
	NodeAddressSetByDelegator Map3NodeAddressSetByDelegator
	NodeKeySet                PubKeySet
	DescriptionIdentitySet    DescriptionIdentitySet
}

// SanityCheck checks basic requirements of a map3 node
func (n *Map3Node) SanityCheck(maxPubKeyAllowed int) error {
	if _, err := n.Description.EnsureLength(); err != nil {
		return err
	}

	if len(n.NodeKeys) == 0 {
		return errNeedAtLeastOneSlotKey
	}

	if c := len(n.NodeKeys); maxPubKeyAllowed != DoNotEnforceMaxBLS &&
		c > maxPubKeyAllowed {
		return errors.Wrapf(
			ErrExcessiveBLSKeys, "have: %d allowed: %d",
			c, maxPubKeyAllowed,
		)
	}

	// TODO: Depend on node state?

	if n.Commission.CommissionRates.Rate.LT(zeroPercent) || n.Commission.CommissionRates.Rate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "rate:%s", n.Commission.CommissionRates.Rate.String(),
		)
	}

	if n.Commission.CommissionRates.MaxRate.LT(zeroPercent) || n.Commission.CommissionRates.MaxRate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "max rate:%s", n.Commission.CommissionRates.MaxRate.String(),
		)
	}

	if n.Commission.CommissionRates.MaxChangeRate.LT(zeroPercent) ||
		n.Commission.CommissionRates.MaxChangeRate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "max change rate:%s", n.Commission.CommissionRates.MaxChangeRate.String(),
		)
	}

	if n.Commission.CommissionRates.Rate.GT(n.Commission.CommissionRates.MaxRate) {
		return errors.Wrapf(
			errCommissionRateTooLarge,
			"rate:%s max rate:%s", n.Commission.CommissionRates.Rate.String(), n.Commission.CommissionRates.MaxRate.String(),
		)
	}

	if n.Commission.CommissionRates.MaxChangeRate.GT(n.Commission.CommissionRates.MaxRate) {
		return errors.Wrapf(
			errCommissionRateTooLarge,
			"rate:%s max change rate:%s", n.Commission.CommissionRates.Rate.String(),
			n.Commission.CommissionRates.MaxChangeRate.String(),
		)
	}

	allKeys := map[string]struct{}{}
	for _, key := range n.NodeKeys {
		if _, ok := allKeys[key.Hex()]; !ok {
			allKeys[key.Hex()] = struct{}{}
		} else {
			return errDuplicateSlotKeys
		}
	}
	return nil
}

func CalcMinTotalNodeStake(blockHeight *big.Int, config *params.ChainConfig) (*big.Int, *big.Int, *big.Int) {
	// TODO: total node state change by time
	return baseMinTotalNodeStake.RoundInt(),
		baseMinTotalNodeStake.Mul(MinSelfDelegation).RoundInt(),
		baseMinTotalNodeStake.Mul(MinDelegation).RoundInt()
}

// CreateValidatorFromNewMsg creates validator from NewValidator message
func CreateMap3NodeFromNewMsg(node *CreateMap3Node, nodeAddr common.Address, blockNum *big.Int) (*Map3Node, error) {
	desc, err := node.Description.EnsureLength()
	if err != nil {
		return nil, err
	}
	commission := Commission{node.CommissionRates, blockNum}
	nodeKeys := append(node.NodeKeys[0:0], node.NodeKeys...)

	v := Map3Node{
		NodeAddress:      nodeAddr,
		InitiatorAddress: node.InitiatorAddress,
		NodeKeys:         nodeKeys,
		Commission:       commission,
		Description:      desc,
	}
	return &v, nil
}

// UpdateValidatorFromEditMsg updates validator from EditValidator message
func UpdateMap3NodeFromEditMsg(map3Node *Map3Node, edit *EditMap3Node) error {
	if map3Node.NodeAddress != edit.Map3NodeAddress {
		return errAddressNotMatch
	}
	newDes, err := UpdateDescription(map3Node.Description, *edit.Description)
	if err != nil {
		return err
	}
	map3Node.Description = newDes

	if !edit.CommissionRate.IsNil() {
		map3Node.Commission.CommissionRates.Rate = edit.CommissionRate
	}

	if edit.NodeKeyToRemove != nil {
		index := -1
		for i := 0; i < len(map3Node.NodeKeys); i++ {
			if edit.NodeKeyToRemove.Hex() == map3Node.NodeKeys[i].Hex() {
				index = i
				break
			}
		}
		// we found key to be removed
		if index >= 0 {
			map3Node.NodeKeys = append(
				map3Node.NodeKeys[:index], map3Node.NodeKeys[index+1:]...,
			)
		} else {
			return errSlotKeyToRemoveNotFound
		}
	}

	if edit.NodeKeyToAdd != nil {
		found := false
		for _, key := range map3Node.NodeKeys {
			if key.Hex() == edit.NodeKeyToAdd.Hex() {
				found = true
				break
			}
		}
		if !found {
			map3Node.NodeKeys = append(map3Node.NodeKeys, *edit.NodeKeyToAdd)
		} else {
			return errSlotKeyToAddExists
		}
	}
	return nil
}
