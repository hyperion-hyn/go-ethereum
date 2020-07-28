// SPDX-License-Identifier: LGPL-3.0 License
pragma solidity ^0.6.4;

struct Decimal {
    uint256 f;
}

struct Description_t {
    string Name;
    string Identity;
    string WebSite;
    string SecurityContract;
    string Details;
}

struct CommissionRates_t {
    Decimal Rate;
    Decimal MaxRate;
    Decimal MaxChangeRate;
}

struct Commission_t {
    CommissionRates_t CommissionRates;
    uint256 UpdateHeight;
}

struct BLSPublicKey_t {
    byte[48] Key;
}

struct BLSPublicKeys_t {
    BLSPublicKey_t[] Keys;
}

struct Counters_t {
    // The number of blocks the validator
    // should've signed when in active mode (selected in committee)
    uint256 NumBlocksToSign;
    // The number of blocks the validator actually signed
    uint256 NumBlocksSigned;
}

struct AddressSet_t {
    address[] Keys;
    mapping (address => bool) Set;
}

struct Validator_t {
    // ECDSA address of the validatorÒ
    address ValidatorAddress;
    // validator's operator (node address)
    AddressSet_t OperatorAddresses;
    // The BLS public key of the validator for consensus
    BLSPublicKeys_t SlotPubKeys;
    // The number of the last epoch this validator is
    // selected in committee (0 means never selected)
    uint256 LastEpochInCommittee;
    // Is the validator active in participating
    // committee selection process or not
    uint256 Status;
    // commission parameters
    Commission_t Commission;
    // description for the validator
    Description_t Description;
    // CreationHeight is the height of creation
    uint256 CreationHeight;
}

// Undelegation represents one undelegation entry
struct Undelegation_t {
    uint256 Amount;
    uint256 Epoch;
}

struct Redelegation_t {
    address DelegatorAddress;
    uint256 Amount;
    uint256 Reward;
    Undelegation_t Undelegation;
}

struct RedelegationMap_t {
    address[] Keys;
    mapping (address => Redelegation_t) Map;
}

// ValidatorWrapper contains validator, its delegation information
struct ValidatorWrapper_t {
    Validator_t Validator;
    RedelegationMap_t Redelegations;
    Counters_t Counters;
    uint256 BlockReward;    // All the rewarded accumulated so far
    uint256 TotalDelegation;
    uint256 TotalDelegationByOperator;
}

struct ValidatorWrapperMap_t {
    address[] Keys;
    mapping (address => ValidatorWrapper_t) Map;
}

struct Slot_t {
    address EcdsaAddress;
    BLSPublicKey_t BLSPublicKey;
    Decimal EffectiveStake;
}

struct Slots_t {
    Slot_t[] Entrys;
}

struct Committee_t {
    uint256 Epoch;
    Slots_t Slots;
}

struct ValidatorPool_t {
    ValidatorWrapperMap_t Validators;
    mapping (string => bool) PublicKeySet;
    mapping (string => bool) DescriptionIdentitySet;
    Committee_t Committee;
}

contract ValidatorPoolWrapper {
    ValidatorPool_t ValidatorPool;
}