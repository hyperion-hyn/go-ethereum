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

struct CommissionRates_ {
    Decimal Rate;
    Decimal MaxRate;
    Decimal MaxChangeRate;
}

struct Commission_ {
    CommissionRates_ CommissionRates;
    uint256 UpdateHeight;
}

struct BLSPublicKey_ {
    byte[48] Key;
}

struct BLSPublicKeys_ {
    BLSPublicKey_[] Keys;
}

struct Counters_ {
    // The number of blocks the validator
    // should've signed when in active mode (selected in committee)
    uint256 NumBlocksToSign;
    // The number of blocks the validator actually signed
    uint256 NumBlocksSigned;
}

struct IterableAddressSet_ {
    address[] Keys;
    mapping (address => bool) Map;
}

struct Validator_ {
    // ECDSA address of the validator
    address ValidatorAddress;
    // validator's operator (map3 node address)
    IterableAddressSet_ OperatorAddresses;
    // The BLS public key of the validator for consensus
    BLSPublicKeys_ SlotPubKeys;
    // The number of the last epoch this validator is
    // selected in committee (0 means never selected)
    uint256 LastEpochInCommittee;
    // maximum total delegation allowed
    uint256 MaxTotalDelegation;
    // Is the validator active in participating
    // committee selection process or not
    byte Status;
    // commission parameters
    Commission_ Commission;
    // description for the validator
    Description_ Description;
    // CreationHeight is the height of creation
    uint256 CreationHeight;
}

// Undelegation represents one undelegation entry
struct Undelegation_ {
    uint256 Amount;
    uint256 Epoch;
}

struct Redelegation_ {
    address DelegatorAddress;
    uint256 Amount;
    uint256 Reward;
    Undelegation_ Undelegation;
}

struct RedelegationMapEntry_ {
    Redelegation_ Entry;
    uint256 Index;
}

struct IterableRedelegationMap_ {
    address[] Keys;
    mapping (address => RedelegationMapEntry_) Map;
}

// ValidatorWrapper contains validator, its delegation information
struct ValidatorWrapper_ {
    Validator_ Validator;
    IterableRedelegationMap_ Redelegations;
    Counters_ Counters;
    uint256 BlockReward;    // All the rewarded accumulated so far
    uint256 TotalDelegation;
    uint256 TotalDelegationByOperator;
}

struct ValidatorWrapperMapEntry_ {
    ValidatorWrapper_ Entry;
    uint256 Index;
}

struct IterableValidatorWrapperMap_ {
    address[] Keys;
    mapping (address => ValidatorWrapperMapEntry_) Map;
}

struct Slot_ {
    address EcdsaAddress;
    BLSPublicKey_ BLSPublicKey;
    Decimal EffectiveStake;
}

struct Slots_ {
    Slot_[] Entrys;
}

struct Committee_ {
    uint256 Epoch;
    Slots_ Slots;
}

struct ValidatorPool_ {
    IterableValidatorWrapperMap_ Validators;
    IterableValidatorWrapperMap_ ValidatorSnapshots;
    mapping (string => bool) SlotKeySet;
    mapping (string => bool) DescriptionIdentitySet;
    Committee_ Committee;
}

contract ValidatorPoolWrapper {
    ValidatorPool_ ValidatorPool;
}