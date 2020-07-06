package committee

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/types"
	"math/big"
)

type SlotStorage struct {
}

func (s *SlotStorage) ToSlot() *Slot {
	return nil
}

func (s *SlotStorage) GetEcdsaAddress() common.Address {
	return [20]byte{}
}

func (s *SlotStorage) SetEcdsaAddress(ecdsaAddress common.Address) {
}

func (s *SlotStorage) GetBLSPublicKey() types.BLSPublicKey {
	return [48]byte{}
}

func (s *SlotStorage) SetBLSPublicKey(blsPublicKey types.BLSPublicKey) {
}

func (s *SlotStorage) GetEffectiveStake() *numeric.Dec {
	return nil
}

func (s *SlotStorage) SetEffectiveStake(effectiveStake *numeric.Dec) {
}

type SlotListStorage struct {
}

func (s *SlotListStorage) ToSlotList() *SlotList {
	return nil
}

func (s *SlotListStorage) Get(index int) *SlotStorage {
	return nil
}

func (s *SlotListStorage) Set(index int, slot *Slot) {

}

func (s *SlotListStorage) Remove(index int, keepOrder bool) {
}

func (s *SlotListStorage) Push(slot *Slot) {
}

func (s *SlotListStorage) Pop() *Slot {
	return nil
}

func (s *SlotListStorage) Len() int {
	return 0
}

type CommitteeStorage struct {
}

func (c *CommitteeStorage) ToCommittee() *Committee {
	return nil
}

func (c *CommitteeStorage) GetSlots() *SlotList {
	return nil
}

func (c *CommitteeStorage) SetSlots(slots *SlotList) {
}

func (c *CommitteeStorage) SetEpoch(epoch *big.Int) {
}

func (c *CommitteeStorage) GetEpoch() *big.Int {
	return nil
}

type CommitteeByEpochStorage struct {
}

func (c *CommitteeByEpochStorage) ToCommitteeByEpoch() *CommitteeByEpoch {
	return nil
}

func (c *CommitteeByEpochStorage) Size() int {
	return 0
}

func (c *CommitteeByEpochStorage) Contain(key uint64) bool {
	return false
}

func (c *CommitteeByEpochStorage) Put(key uint64, committee *Committee) bool {
	return false
}

func (c *CommitteeByEpochStorage) Remove(key uint64) {
}

func (c *CommitteeByEpochStorage) Get(key uint64) (*CommitteeStorage, bool) {
	return nil, false
}

func (c *CommitteeByEpochStorage) Keys() []uint64 {
	return nil
}
