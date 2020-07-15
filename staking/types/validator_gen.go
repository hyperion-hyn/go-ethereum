package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/effective"
	"math/big"
)

type BLSPublicKeysStorage struct {
}

func (b *BLSPublicKeysStorage) ToBLSPublicKeys() *BLSPublicKeys {
	return nil
}

func (b *BLSPublicKeysStorage) Get(index int) *BLSPublicKey {
	return nil
}

func (b *BLSPublicKeysStorage) Set(index int, blsPublicKey *BLSPublicKey) {

}

func (b *BLSPublicKeysStorage) Remove(index int, keepOrder bool) {
}

func (b *BLSPublicKeysStorage) Push(blsPublicKey *BLSPublicKey) {
}

func (b *BLSPublicKeysStorage) Pop() *BLSPublicKey {
	return nil
}

func (b *BLSPublicKeysStorage) Len() int {
	return 0
}


type ValidatorStorage struct {
}

func (v *ValidatorStorage) ToValidator() *Validator {
	return nil
}

func (v *ValidatorStorage) GetValidatorAddress() common.Address {
	return common.Address0
}

func (v *ValidatorStorage) SetValidatorAddress(validatorAddress common.Address) {
}

func (v *ValidatorStorage) GetInitiatorAddresses() *AddressSetStorage {
	return nil
}

func (v *ValidatorStorage) SetInitiatorAddress(initiatorAddresses *AddressSet) {
}

func (v *ValidatorStorage) GetSlotPubKeys() *BLSPublicKeysStorage {
	return nil
}

func (v *ValidatorStorage) SetSlotPubKeys(slotPubKeys BLSPublicKeys) {
}

func (v *ValidatorStorage) SetLastEpochInCommittee(lastEpochInCommittee *big.Int) {
}

func (v *ValidatorStorage) GetLastEpochInCommittee() *big.Int {
	return nil
}

func (v *ValidatorStorage) SetStatus(status effective.Eligibility) {
}

func (v *ValidatorStorage) GetStatus() effective.Eligibility {
	return 0
}

func (v *ValidatorStorage) GetCommission() *CommissionStorage {
	return nil
}

func (v *ValidatorStorage) SetCommission(commission *Commission) {
}

func (v *ValidatorStorage) GetDescription() *DescriptionStorage {
	return nil
}

func (v *ValidatorStorage) SetDescription(description *Description) {
}

func (v *ValidatorStorage) SetCreationHeight(creationHeight *big.Int) {
}

func (v *ValidatorStorage) GetCreationHeight() *big.Int {
	return nil
}

type CountersStorage struct {
}

func (c *CountersStorage) ToCounters() *Counters {
	return nil
}

func (c *CountersStorage) SetNumBlocksToSign(numBlocksToSign *big.Int) {
}

func (c *CountersStorage) GetNumBlocksToSign() *big.Int {
	return nil
}

func (c *CountersStorage) SetNumBlocksSigned(numBlocksSigned *big.Int) {
}

func (c *CountersStorage) GetNumBlocksSigned() *big.Int {
	return nil
}

type ValidatorWrapperStorage struct {
}

func (v *ValidatorWrapperStorage) ToValidatorWrapper() *ValidatorWrapper {
	return nil
}

func (v *ValidatorWrapperStorage) SetValidator(validator *Validator) {
}

func (v *ValidatorWrapperStorage) GetValidator() *ValidatorStorage {
	return nil
}

func (v *ValidatorWrapperStorage) SetRedelegations(redelegations *Redelegations) {
}

func (v *ValidatorWrapperStorage) GetRedelegations() *RedelegationsStorage {
	return nil
}

func (v *ValidatorWrapperStorage) SetCounters(counters *Counters) {
}

func (v *ValidatorWrapperStorage) GetCounters() *CountersStorage {
	return nil
}

func (v *ValidatorWrapperStorage) SetBlockReward(blockReward *big.Int) {
}

func (v *ValidatorWrapperStorage) GetBlockReward() *big.Int {
	return nil
}

func (v *ValidatorWrapperStorage) SetTotalDelegation(totalDelegation *big.Int) {
}

func (v *ValidatorWrapperStorage) GetTotalDelegation() *big.Int {
	return nil
}

type ValidatorWrappersStorage struct {
}

func (v *ValidatorWrappersStorage) ToValidatorWrappers() *ValidatorWrappers {
	return nil
}

func (v *ValidatorWrappersStorage) Size() int {
	return 0
}

func (v *ValidatorWrappersStorage) Contain(key common.Address) bool {
	return false
}

func (v *ValidatorWrappersStorage) Put(key common.Address, validatorWrapper *ValidatorWrapper) bool {
	return false
}

func (v *ValidatorWrappersStorage) Remove(key common.Address) {
}

func (v *ValidatorWrappersStorage) Get(key common.Address) (*ValidatorWrapperStorage, bool) {
	return nil, false
}

func (v *ValidatorWrappersStorage) Keys() []common.Address {
	return nil
}

type ValidatorPoolStorage struct {
}

func (v *ValidatorPoolStorage) ToValidatorPool() *ValidatorPool {
	return nil
}

func (v *ValidatorPoolStorage) SetValidators(validators *ValidatorWrappers) {
}

func (v *ValidatorPoolStorage) GetValidators() *ValidatorWrappersStorage {
	return nil
}

func (v *ValidatorPoolStorage) SetSlotKeySet(slotKeySet *PubKeySet) {
}

func (v *ValidatorPoolStorage) GetSlotKeySet() *PubKeySetStorage {
	return nil
}

func (v *ValidatorPoolStorage) SetDescriptionIdentitySet(descriptionIdentitySet *DescriptionIdentitySet) {
}

func (v *ValidatorPoolStorage) GetDescriptionIdentitySet() *DescriptionIdentitySetStorage {
	return nil
}
