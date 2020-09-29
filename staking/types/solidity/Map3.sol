// SPDX-License-Identifier: LGPL-3.0 License
pragma solidity ^0.6.4;

struct Decimal {
    uint256 f;
}

struct Description_ {
    string Name;
    string Identity;
    string Website;
    string SecurityContact;
    string Details;
}

struct BLSPublicKey_ {
    byte[48] Key;
}

struct BLSPublicKeys_ {
    BLSPublicKey_[] Keys;
}

struct Commission_ {
    Decimal Rate;
    Decimal RateForNextPeriod;
    uint256 UpdateHeight;
}

struct Map3Node_ {
    // ECDSA address of the map3 node
    address Map3Address;
    // Map3's operator
    address OperatorAddress;
    // The BLS public key of the map3 node
    BLSPublicKeys_ NodeKeys;
    // commission parameter
    Commission_ Commission;
    // description for the map3 node
    Description_ Description;
    // CreationHeight is the height of creation
    uint256 CreationHeight;
    Decimal Age;
    byte Status;
    uint256 PendingEpoch;
    uint256 ActivationEpoch;
    Decimal ReleaseEpoch;
}

// PendingDelegation represents tokens during map3 in pending state
struct PendingDelegation_ {
    uint256 Amount;
    Decimal UnlockedEpoch;
}

// Undelegation represents one undelegation entry
struct Undelegation_ {
    uint256 Amount;
    uint256 Epoch;
}

struct Microdelegation_ {
    address DelegatorAddress;
    uint256 Amount;
    uint256 Reward;
    PendingDelegation_ PendingDelegation;
    Undelegation_ Undelegation;
    Renewal_ Renewal;
}

struct Renewal_ {
    byte Status;
    uint256 UpdateHeight;
}

struct MicrodelegationMapEntry_ {
    Microdelegation_ Entry;
    uint256 Index;
}

struct MicrodelegationMap_ {
    address[] Keys;
    mapping (address => MicrodelegationMapEntry_) Map;
}

struct RestakingReference_ {
    address ValidatorAddress;
}

// Map3NodeWrapper contains map3 node, its micro-delegation information
struct Map3NodeWrapper_ {
    Map3Node_ Map3Node;
    MicrodelegationMap_ Microdelegations;
    RestakingReference_ RestakingReference;
    uint256 AccumulatedReward; // All the rewarded accumulated so far
    uint256 TotalDelegation;
    uint256 TotalPendingDelegation;
}

struct Map3NodeWrapperMapEntry_ {
    Map3NodeWrapper_ Entry;
    uint256 Index;
}

struct Map3NodeWrapperMap_ {
    address[] Keys;
    mapping (address => Map3NodeWrapperMapEntry_) Map;
}

struct DelegationIndex_ {
    address Map3Address;
    bool IsOperator;
}

struct DelegationIndexMapEntry_ {
    DelegationIndex_ Entry;
    uint256 Index;
}

struct DelegationIndexMap_ {
    address[] Keys; // Map3 nodes
    mapping (address => DelegationIndexMapEntry_) Map;
}

struct Map3NodePool_ {
    Map3NodeWrapperMap_ Nodes;
    mapping (address => DelegationIndexMap_) DelegationIndexMapByDelegator;
    mapping (string => bool) NodeKeySet;
    mapping (string => bool) DescriptionIdentitySet;
}

contract Map3NodePoolWrapper {
    Map3NodePool_ Map3NodePool;
}