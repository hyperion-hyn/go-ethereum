package storage_test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"math/big"
	"testing"
)

var (
	map3NodeAddr   = makeTestAddr("map3-1")
	operatorAddr   = makeTestAddr("operator-1")
	delegatorAddr  = makeTestAddr("delegator-1")
	delegatorAddr2 = makeTestAddr("delegator-2")
)

func TestStorage_Map3NodeWrapper_AddMicrodelegation(t *testing.T) {
	tests := []struct {
		delegator          common.Address
		amount             *big.Int
		pending            bool
		epoch              *big.Int
		wantIsNewDelegator bool
		wantAmount         *big.Int
		wantPendingAmount  *big.Int
		wantUnlockedEpoch  common.Dec
	}{
		{
			delegator:          operatorAddr,
			amount:             big.NewInt(100),
			pending:            true,
			epoch:              big.NewInt(15),
			wantIsNewDelegator: false,
			wantAmount: big.NewInt(0),
			wantPendingAmount: big.NewInt(200),
			wantUnlockedEpoch: common.NewDecWithPrec(185, 1),
		},
		{
			delegator:          operatorAddr,
			amount:             big.NewInt(60),
			pending:            true,
			epoch:              big.NewInt(7),
			wantIsNewDelegator: false,
			wantAmount: big.NewInt(0),
			wantPendingAmount: big.NewInt(160),
			wantUnlockedEpoch: common.NewDecWithPrec(115, 1),
		},
		{
			delegator:          delegatorAddr,
			amount:             big.NewInt(100),
			pending:            true,
			epoch:              big.NewInt(7),
			wantIsNewDelegator: false,
			wantAmount: big.NewInt(100),
			wantPendingAmount: big.NewInt(100),
			wantUnlockedEpoch: common.NewDec(14),
		},
		{
			delegator:          delegatorAddr,
			amount:             big.NewInt(100),
			pending:            false,
			epoch:              big.NewInt(7),
			wantIsNewDelegator: false,
			wantAmount: big.NewInt(200),
			wantPendingAmount: big.NewInt(0),
			wantUnlockedEpoch: common.NewDec(0),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test-%v", i), func(t *testing.T) {
			node := makeStorageMap3Node(t)
			gotIsNewDelegator := node.AddMicrodelegation(tt.delegator, tt.amount, tt.pending, tt.epoch)
			if gotIsNewDelegator != tt.wantIsNewDelegator {
				t.Errorf("AddMicrodelegation() = %v, want %v", gotIsNewDelegator, tt.wantIsNewDelegator)
			}
			m, _ := node.Microdelegations().Get(tt.delegator)
			if m.Amount().Value().Cmp(tt.wantAmount) != 0 {
				t.Errorf("WantAmount() = %v, want %v", m.Amount().Value(), tt.wantAmount)
			}
			if m.PendingDelegation().Amount().Value().Cmp(tt.wantPendingAmount) != 0 {
				t.Errorf("WantPendingAmount() = %v, want %v", m.PendingDelegation().Amount().Value(), tt.wantPendingAmount)
			}
			if !m.PendingDelegation().UnlockedEpoch().Value().Equal(tt.wantUnlockedEpoch) {
				t.Errorf("WantPendingAmount() = %v, want %v", m.PendingDelegation().UnlockedEpoch().Value(), tt.wantUnlockedEpoch)
			}
		})
	}
}

func makeStorageMap3Node(t *testing.T) *microstaking.Storage_Map3NodeWrapper_ {
	node := makeMap3Node()
	db, err := newTestStateDB()
	if err != nil {
		t.Fatal(err)
	}
	db.Map3NodePool().Map3Nodes().Put(node.Map3Node.Map3Address, node)
	db.IncreaseMap3NonceIfZero()
	db.Commit(true)
	sn, _ := db.Map3NodeByAddress(node.Map3Node.Map3Address)
	return sn
}

func makeMap3Node() *microstaking.Map3NodeWrapper_ {
	n := microstaking.NewMap3NodeWrapperBuilder().
		SetMap3Address(map3NodeAddr).
		SetOperatorAddress(operatorAddr).
		SetCreationHeight(big.NewInt(12306)).
		SetAge(common.NewDecWithPrec(25, 1)).
		SetStatus(microstaking.Active).
		SetActivationEpoch(big.NewInt(1)).
		SetReleaseEpoch(common.NewDec(10)).
		AddMicrodelegation(microstaking.NewMicrodelegation(operatorAddr, big.NewInt(100), common.NewDec(10), true)).
		AddMicrodelegation(microstaking.NewMicrodelegation(delegatorAddr, big.NewInt(100), common.NewDec(10), false)).
		Build()
	return n
}

func newTestStateDB() (*state.StateDB, error) {
	return state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
}

func makeTestAddr(item interface{}) common.Address {
	s := fmt.Sprintf("hyperion-hyn-%v", item)
	return common.BytesToAddress([]byte(s))
}
