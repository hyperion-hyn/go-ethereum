package microstaking

import (
	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/ethereum/go-ethereum/staking/types/common"
	"math/big"
)

// CreateMap3Node - type for creating a new map3 node
type CreateMap3Node struct {
	OperatorAddress common.Address       `json:"operator-address"`
	Description     Description_         `json:"description"`
	Commission      common.Dec           `json:"commission"`
	NodePubKey      BLSPublicKey_        `json:"node-pub-key"`
	NodeKeySig      common2.BLSSignature `json:"node-key-sig"`
	Amount          *big.Int             `json:"amount"`
}

// EditMap3Node - type for edit existing map3 node
type EditMap3Node struct {
	Map3NodeAddress common.Address        `json:"map3node-address"`
	OperatorAddress common.Address        `json:"operator-address"`
	Description     Description_          `json:"description"`
	NodeKeyToRemove *BLSPublicKey_        `json:"node-key-to_remove" rlp:"nil"`
	NodeKeyToAdd    *BLSPublicKey_        `json:"node-key-to_add" rlp:"nil"`
	NodeKeyToAddSig *common2.BLSSignature `json:"node-key-to-add-sig" rlp:"nil"`
}

type TerminateMap3Node struct {
	Map3NodeAddress common.Address `json:"validator-address"`
	OperatorAddress common.Address `json:"operator-address"`
}

// Microdelegate - type for delegating to a map3 node
type Microdelegate struct {
	DelegatorAddress common.Address `json:"delegator-address"`
	Map3NodeAddress  common.Address `json:"map3node-address"`
	Amount           *big.Int       `json:"amount"`
}

// Unmicrodelegate - type for removing delegation responsibility
type Unmicrodelegate struct {
	DelegatorAddress common.Address `json:"delegator-address"`
	Map3NodeAddress  common.Address `json:"map3node-address"`
	Amount           *big.Int       `json:"amount"`
}

// CollectRewards - type for collecting token rewards
type CollectRewards struct {
	DelegatorAddress common.Address `json:"delegator-address"`
}

type RenewMap3Node struct {
	Map3NodeAddress  common.Address `json:"map3node-address"`
	DelegatorAddress common.Address `json:"delegator_address"`
	IsRenew          bool           `json:"is-renew"`
	// just for map3 operator
	NewCommissionRate common.Dec `json:"new-commission-rate"`
}
