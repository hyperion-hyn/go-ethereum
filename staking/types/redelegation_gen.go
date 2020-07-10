package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type RedelegationStorage struct {
}

func (r *RedelegationStorage) ToRedelegation() *Redelegation {
	return nil
}

func (r *RedelegationStorage) GetDelegatorAddress() common.Address {
	return common.Address0
}

func (r *RedelegationStorage) SetDelegatorAddress(delegatorAddress common.Address) {
}

func (r *RedelegationStorage) SetAmount(amount *big.Int) {
}

func (r *RedelegationStorage) GetAmount() *big.Int {
	return nil
}

func (r *RedelegationStorage) SetReward(reward *big.Int) {
}

func (r *RedelegationStorage) GetReward() *big.Int {
	return nil
}

func (r *RedelegationStorage) SetUndelegation(undelegation *Undelegation) {
}

func (r *RedelegationStorage) GetUndelegation() *UndelegationStorage {
	return nil
}

type RedelegationsStorage struct {
}

func (r *RedelegationsStorage) ToRedelegations() *Redelegations {
	return nil
}

func (r *RedelegationsStorage) Size() int {
	return 0
}

func (r *RedelegationsStorage) Contain(key common.Address) bool {
	return false
}

func (r *RedelegationsStorage) Put(key common.Address, redelegations *Redelegation) bool {
	return false
}

func (r *RedelegationsStorage) Remove(key common.Address) {
}

func (r *RedelegationsStorage) Get(key common.Address) (*RedelegationStorage, bool) {
	return nil, false
}

func (r *RedelegationsStorage) Keys() []common.Address {
	return nil
}
