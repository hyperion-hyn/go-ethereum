package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/storage"
	"github.com/ethereum/go-ethereum/numeric"
	"math/big"
)

type MicrodelegationStorage struct {
}

func (m *MicrodelegationStorage) ToMicrodelegation() *Microdelegation {
	return nil
}

func (m *MicrodelegationStorage) GetDelegatorAddress() common.Address {
	return common.Address0
}

func (m *MicrodelegationStorage) SetDelegatorAddress(delegatorAddress common.Address) {
}

func (m *MicrodelegationStorage) SetAmount(amount *big.Int) {
}

func (m *MicrodelegationStorage) GetAmount() *big.Int {
	return nil
}

func (m *MicrodelegationStorage) SetReward(reward *big.Int) {
}

func (m *MicrodelegationStorage) GetReward() *big.Int {
	return nil
}

func (m *MicrodelegationStorage) SetUndelegation(undelegation *Undelegation) {
}

func (m *MicrodelegationStorage) GetUndelegation() *UndelegationStorage {
	return nil
}

func (m *MicrodelegationStorage) SetPendingDelegation(pendingDelegation *PendingDelegation) {
}

func (m *MicrodelegationStorage) GetPendingDelegation() *PendingDelegationStorage {
	return nil
}

func (m *MicrodelegationStorage) SetRenewal(renewal *Renewal) {
}

func (m *MicrodelegationStorage) GetRenewal() *RenewalStorage {
	return nil
}

type RenewalStorage struct {
}

func (r *RenewalStorage) GetUpdateHeight() *big.Int {
	return nil
}

func (r *RenewalStorage) SetUpdateHeight(updateHeight *big.Int) {
}

func (r *RenewalStorage) SetRenew(isRenew bool) {
}

func (r *RenewalStorage) IsRenew() bool {
	return false
}


type PendingDelegationStorage struct {
}

func (p *PendingDelegationStorage) ToPendingDelegation() *PendingDelegation {
	return nil
}

func (p *PendingDelegationStorage) SetAmount(amount *big.Int) {
}

func (p *PendingDelegationStorage) GetAmount() *big.Int {
	return nil
}

func (p *PendingDelegationStorage) SetUnlockedEpoch(epoch *numeric.Dec) {
}

func (p *PendingDelegationStorage) GetUnlockedEpoch() *numeric.Dec {
	return nil
}

type UndelegationStorage struct {
	storage *storage.Storage
}

func (u *UndelegationStorage) ToUndelegation() *Undelegation {
	return nil
}

func (u *UndelegationStorage) SetAmount(amount *big.Int) {
}

func (u *UndelegationStorage) GetAmount() *big.Int {
	return &big.Int{}
}

func (u *UndelegationStorage) SetEpoch(epoch *big.Int) {
}

func (u *UndelegationStorage) GetEpoch() *big.Int {
	return nil
}

type MicrodelegationsStorage struct {
	storage *storage.Storage
}

func (m *MicrodelegationsStorage) ToMicrodelegations() *Microdelegations {
	return nil
}

func (m *MicrodelegationsStorage) Size() int {
	return 0
}

func (m *MicrodelegationsStorage) Contain(key common.Address) bool {
	return false
}

func (m *MicrodelegationsStorage) Put(key common.Address, microdelegation *Microdelegation) bool {
	return false
}

func (m *MicrodelegationsStorage) Remove(key common.Address) {
}

func (m *MicrodelegationsStorage) Get(key common.Address) (*MicrodelegationStorage, bool) {
	return nil, false
}

func (m *MicrodelegationsStorage) Keys() []common.Address {
	return nil
}
