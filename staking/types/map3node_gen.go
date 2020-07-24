package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/storage"
	"github.com/ethereum/go-ethereum/numeric"
	"math/big"
)

type Map3NodeKeysStorage struct {
	storage *storage.Storage
}

func (m *Map3NodeKeysStorage) ToMap3NodeKeys() *Map3NodeKeys {
	return nil
}

func (m *Map3NodeKeysStorage) Get(index int) *Map3NodeKey {
	return nil
}

func (m *Map3NodeKeysStorage) Set(index int, map3NodeKey *Map3NodeKey) {

}

func (m *Map3NodeKeysStorage) Remove(index int, keepOrder bool) {
}

func (m *Map3NodeKeysStorage) Push(map3NodeKey *Map3NodeKey) {
}

func (m *Map3NodeKeysStorage) Pop() *Map3NodeKey {
	return nil
}

func (m *Map3NodeKeysStorage) Len() int {
	return 0
}

type Map3NodeStorage struct {
	storage *storage.Storage
}

func (m *Map3NodeStorage) ToMap3Node() *Map3Node {
	return nil
}

func (m *Map3NodeStorage) GetNodeAddress() common.Address {
	return common.Address0
}

func (m *Map3NodeStorage) SetNodeAddress(nodeAddress common.Address) {
}

func (m *Map3NodeStorage) GetInitiatorAddress() common.Address {
	return common.Address0
}

func (m *Map3NodeStorage) SetInitiatorAddress(initiatorAddress common.Address) {
}

func (m *Map3NodeStorage) GetNodeKeys() *Map3NodeKeysStorage {
	return nil
}

func (m *Map3NodeStorage) SetNodeKeys(nodeKeys Map3NodeKeys) {
}

func (m *Map3NodeStorage) GetCommission() *CommissionStorage {
	return nil
}

func (m *Map3NodeStorage) SetCommission(commission Commission) {
}

func (m *Map3NodeStorage) GetDescription() *DescriptionStorage {
	return nil
}

func (m *Map3NodeStorage) SetDescription(description *Description) {
}

type NodeStateStorage struct {
}

func (n *NodeStateStorage) ToNodeState() *NodeState {
	return nil
}

func (n *NodeStateStorage) SetStatus(status Status) {
}

func (n *NodeStateStorage) GetStatus() Status {
	return 0
}

func (n *NodeStateStorage) SetNodeAge(nodeAge *big.Int) {
}

func (n *NodeStateStorage) GetNodeAge() *big.Int {
	return nil
}

func (n *NodeStateStorage) SetCreationEpoch(creationEpoch *big.Int) {
}

func (n *NodeStateStorage) GetCreationEpoch() *big.Int {
	return nil
}

func (n *NodeStateStorage) SetActivationEpoch(activationEpoch *big.Int) {
}

func (n *NodeStateStorage) GetActivationEpoch() *big.Int {
	return nil
}

func (n *NodeStateStorage) SetReleaseEpoch(releaseEpoch *numeric.Dec) {
}

func (n *NodeStateStorage) GetReleaseEpoch() *numeric.Dec {
	return nil
}

type Map3NodeWrapperStorage struct {
}

func (m *Map3NodeWrapperStorage) ToMap3NodeWrapper() *Map3NodeWrapper {
	return nil
}

func (m *Map3NodeWrapperStorage) SetMap3Node(map3Node *Map3Node) {
}

func (m *Map3NodeWrapperStorage) GetMap3Node() *Map3NodeStorage {
	return nil
}

func (m *Map3NodeWrapperStorage) SetMicrodelegations(microdelegations *Microdelegations) {
}

func (m *Map3NodeWrapperStorage) GetMicrodelegations() *MicrodelegationsStorage {
	return nil
}

func (m *Map3NodeWrapperStorage) SetRedelegationReference() {
}

func (m *Map3NodeWrapperStorage) GetRedelegationReference() *RedelegationReferenceStorage {
	return nil
}

func (m *Map3NodeWrapperStorage) SetAccumulatedReward(accumulatedReward *big.Int) {
}

func (m *Map3NodeWrapperStorage) GetAccumulatedReward() *big.Int {
	return nil
}

func (m *Map3NodeWrapperStorage) SetNodeState(nodeState *NodeState) {
}

func (m *Map3NodeWrapperStorage) GetNodeState() *NodeStateStorage {
	return nil
}

func (m *Map3NodeWrapperStorage) SetTotalDelegation(totalDelegation *big.Int) {
}

func (m *Map3NodeWrapperStorage) GetTotalDelegation() *big.Int {
	return nil
}

func (m *Map3NodeWrapperStorage) SetTotalPendingDelegation(totalPendingDelegation *big.Int) {
}

func (m *Map3NodeWrapperStorage) GetTotalPendingDelegation() *big.Int {
	return nil
}

type Map3NodeWrappersStorage struct {
}

func (m *Map3NodeWrappersStorage) ToMap3NodeWrappers() *Map3NodeWrappers {
	return nil
}

func (m *Map3NodeWrappersStorage) Size() int {
	return 0
}

func (m *Map3NodeWrappersStorage) Contain(key common.Address) bool {
	return false
}

func (m *Map3NodeWrappersStorage) Put(key common.Address, map3NodeWrapper *Map3NodeWrapper) bool {
	return false
}

func (m *Map3NodeWrappersStorage) Remove(key common.Address) {
}

func (m *Map3NodeWrappersStorage) Get(key common.Address) (*Map3NodeWrapperStorage, bool) {
	return nil, false
}

func (m *Map3NodeWrappersStorage) Keys() []common.Address {
	return nil
}

type AddressSetStorage struct {
}

func (m *AddressSetStorage) ToAddressSet() *AddressSet {
	return nil
}

func (m *AddressSetStorage) Size() int {
	return 0
}

func (m *AddressSetStorage) Contain(key common.Address) bool {
	return false
}

func (m *AddressSetStorage) Put(key common.Address) bool {
	return false
}

func (m *AddressSetStorage) Remove(key common.Address) {
}

func (m *AddressSetStorage) Keys() []common.Address {
	return nil
}

type AddressToAddressSetMapStorage struct {
}

func (m *AddressToAddressSetMapStorage) Size() int {
	return 0
}

func (m *AddressToAddressSetMapStorage) Contain(key common.Address) bool {
	return false
}

func (m *AddressToAddressSetMapStorage) Put(key common.Address, addressSet *AddressSet) bool {
	return false
}

func (m *AddressToAddressSetMapStorage) Remove(key common.Address) {
}

func (m *AddressToAddressSetMapStorage) Get(key common.Address) (*AddressSetStorage, bool) {
	return nil, false
}

func (m *AddressToAddressSetMapStorage) Keys() []common.Address {
	return nil
}

type PubKeySetStorage struct {
}

func (p *PubKeySetStorage) ToPubKeySet() *PubKeySet {
	return nil
}

func (p *PubKeySetStorage) Size() int {
	return 0
}

func (p *PubKeySetStorage) Contain(key string) bool {
	return false
}

func (p *PubKeySetStorage) Put(key string) bool {
	return false
}

func (p *PubKeySetStorage) Remove(key string) {
}

func (p *PubKeySetStorage) Keys() []string {
	return nil
}

type Map3NodePoolStorage struct {
}

func (m *Map3NodePoolStorage) ToMap3NodePool() *Map3NodePool {
	return nil
}

func (m *Map3NodePoolStorage) SetNodes(map3NodeWrappers *Map3NodeWrappers) {
}

func (m *Map3NodePoolStorage) GetNodes() *Map3NodeWrappersStorage {
	return nil
}

func (m *Map3NodePoolStorage) SetNodeAddressSetByDelegator(map3NodeAddressSetByDelegator *AddressToAddressSetMap) {
}

func (m *Map3NodePoolStorage) GetNodeAddressSetByDelegator() *AddressToAddressSetMapStorage {
	return nil
}

func (m *Map3NodePoolStorage) SetNodeAddressSetByInitiator(map3NodeAddressSetByInitiator *AddressToAddressSetMap) {
}

func (m *Map3NodePoolStorage) GetNodeAddressSetByInitiator() *AddressToAddressSetMapStorage {
	return nil
}

func (m *Map3NodePoolStorage) SetNodeKeySet(nodeKeySet *PubKeySet) {
}

func (m *Map3NodePoolStorage) GetNodeKeySet() *PubKeySetStorage {
	return nil
}

func (m *Map3NodePoolStorage) SetDescriptionIdentitySet(descriptionIdentitySet *DescriptionIdentitySet) {
}

func (m *Map3NodePoolStorage) GetDescriptionIdentitySet() *DescriptionIdentitySetStorage {
	return nil
}
