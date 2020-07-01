package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/storage"
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
	return [20]byte{}
}

func (m *Map3NodeStorage) SetNodeAddress(nodeAddress common.Address) {
}

func (m *Map3NodeStorage) GetInitiatorAddress() common.Address {
	return [20]byte{}
}

func (m *Map3NodeStorage) SetInitiatorAddress(initiatorAddress common.Address) {
}

func (m *Map3NodeStorage) GetNodeKeys() Map3NodeKeysStorage {
	return Map3NodeKeysStorage{}
}

func (m *Map3NodeStorage) SetNodeKeys(nodeKeys Map3NodeKeys) {
}

func (m *Map3NodeStorage) GetCommission() CommissionStorage {
	return CommissionStorage{}
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

func (n *NodeStateStorage) SetStatus(status Map3NodeStatus) {
}

func (n *NodeStateStorage) GetStatus() Map3NodeStatus {
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

func (n *NodeStateStorage) SetReleaseEpoch(releaseEpoch *big.Int) {
}

func (n *NodeStateStorage) GetReleaseEpoch() *big.Int {
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

func (m *Map3NodeWrapperStorage) SetRedelegationReference(redelegationReference *RedelegationReference) {
}

func (m *Map3NodeWrapperStorage) GetRedelegationReference() *RedelegationReference {
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

func (m *Map3NodeWrappersStorage) Put(key common.Address, map3NodeWrappers *Map3NodeWrapper) bool {
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

type Map3NodeSnapshotStorage struct {
}

func (m *Map3NodeSnapshotStorage) ToMap3NodeSnapshot() Map3NodeSnapshot {
	return Map3NodeSnapshot{}
}

func (m *Map3NodeSnapshotStorage) SetMap3Nodes(epoch *big.Int) {
}

func (m *Map3NodeSnapshotStorage) GetMap3Nodes() Map3NodeWrappersStorage {
	return Map3NodeWrappersStorage{}
}

func (m *Map3NodeSnapshotStorage) SetEpoch(epoch *big.Int) {
}

func (m *Map3NodeSnapshotStorage) GetEpoch() *big.Int {
	return nil
}

type Map3NodeSnapshotByEpochStorage struct {
}

func (m *Map3NodeSnapshotByEpochStorage) ToMap3NodeSnapshotByEpoch() *Map3NodeSnapshotByEpochStorage {
	return nil
}

func (m *Map3NodeSnapshotByEpochStorage) Size() int {
	return 0
}

func (m *Map3NodeSnapshotByEpochStorage) Contain(key uint64) bool {
	return false
}

func (m *Map3NodeSnapshotByEpochStorage) Put(key uint64, map3NodeWrappers *Map3NodeSnapshot) bool {
	return false
}

func (m *Map3NodeSnapshotByEpochStorage) Remove(key uint64) {
}

func (m *Map3NodeSnapshotByEpochStorage) Get(key uint64) (*Map3NodeSnapshotStorage, bool) {
	return nil, false
}

func (m *Map3NodeSnapshotByEpochStorage) Keys() []uint64 {
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

type Map3NodeAddressSetByDelegatorStorage struct {
}

func (m *Map3NodeAddressSetByDelegatorStorage) ToMap3NodeAddressSetByDelegator() *Map3NodeAddressSetByDelegator {
	return nil
}

func (m *Map3NodeAddressSetByDelegatorStorage) Size() int {
	return 0
}

func (m *Map3NodeAddressSetByDelegatorStorage) Contain(key common.Address) bool {
	return false
}

func (m *Map3NodeAddressSetByDelegatorStorage) Put(key common.Address, map3NodeAddressSet *AddressSet) bool {
	return false
}

func (m *Map3NodeAddressSetByDelegatorStorage) Remove(key common.Address) {
}

func (m *Map3NodeAddressSetByDelegatorStorage) Get(key common.Address) (*AddressSetStorage, bool) {
	return nil, false
}

func (m *Map3NodeAddressSetByDelegatorStorage) Keys() []common.Address {
	return nil
}

type NodeKeySetStorage struct {
}

func (n *NodeKeySetStorage) ToNodeKeySet() NodeKeySet {
	return nil
}

func (n *NodeKeySetStorage) Size() int {
	return 0
}

func (n *NodeKeySetStorage) Contain(key string) bool {
	return false
}

func (n *NodeKeySetStorage) Put(key string) bool {
	return false
}

func (n *NodeKeySetStorage) Remove(key string) {
}

func (n *NodeKeySetStorage) Keys() []string {
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

func (m *Map3NodePoolStorage) SetNodeSnapshotByEpoch(map3NodeSnapshotByEpoch *Map3NodeSnapshotByEpoch) {
}

func (m *Map3NodePoolStorage) GetNodeSnapshotByEpoch() *Map3NodeSnapshotByEpochStorage {
	return nil
}

func (m *Map3NodePoolStorage) SetNodeAddressSetByDelegator(map3NodeAddressSetByDelegator *Map3NodeAddressSetByDelegator) {
}

func (m *Map3NodePoolStorage) GetNodeAddressSetByDelegator() *Map3NodeAddressSetByDelegatorStorage {
	return nil
}

func (m *Map3NodePoolStorage) SetNodeKeySet(nodeKeySet NodeKeySet) {
}

func (m *Map3NodePoolStorage) GetNodeKeySet() *NodeKeySetStorage {
	return nil
}

func (m *Map3NodePoolStorage) SetDescriptionIdentitySet(descriptionIdentitySet *DescriptionIdentitySet) {
}

func (m *Map3NodePoolStorage) GetDescriptionIdentitySet() *DescriptionIdentitySetStorage {
	return nil
}
