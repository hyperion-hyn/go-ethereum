// SPDX-License-Identifier: LGPL-3.0 License
pragma solidity ^0.6.10;
pragma experimental ABIEncoderV2;

contract Map3Pool {
    struct CommissionRates_t {
        uint256 Rate;
        uint256 MaxRate;
        uint256 MaxChangeRate;
    }

    struct Commission_t {
        CommissionRates_t CommissionRates;
        uint256 UpdateHeight;
    }

    struct Description_t {
        string Name;
        string Identity;
        string WebSite;
        string SecurityContract;
        string Details;
        uint256[2] Version; // FOR TEST;
    }

    struct Map3Node_t  {
        address NodeAddress;
        address InitiatorAddress;
        bytes NodeKeys;
        Commission_t Commission;
        Description_t Description;
        address SplittedFrom;
    }

    struct Undelegation_t {
        uint256 Amount;
        uint256 Epoch;
    }

    struct PendingDelegation_t {
        uint256 Amount;
        uint256 Epoch;
    }

    struct Microdelegation_t {
        address DelegatorAddress;
        uint256 Amount;
        uint256 Reward;
        Undelegation_t[] Undelegations;
        PendingDelegation_t[] PendingDelegations;
        PendingDelegation_t[2] PendingDelegationsfixed;
        PendingDelegation_t[2][3] PendingDelegationsfixed2dimension;
        bool AutoRenew;
    }

    struct RedelegationReference_t {
        address ValidatorAddress;
        uint256 ReleasedTotalDelegation;
    }

    struct NodeState_t {
        byte Status;
        uint256 NodeAge;
        uint256 CreationEpoch;
        uint256 ActivationEpoch;
        uint256 ReleaseEpoch;
    }

    struct Map3NodeWrapper_t {
        Map3Node_t Map3Node;
        mapping (address => Microdelegation_t)  Microdelegations;
        //
        RedelegationReference_t RedelegationReference;
        uint256 AccumulatedReward;
        NodeState_t nodeState;
        uint256 TotalDelegation;
        uint256 TotalPendingDelegation;
    }

    struct Map3NodeSnapshot_t {
        mapping (address => Map3NodeWrapper_t) Map3Nodes;
        uint256 Epoch;
    }

    struct Map3NodePool_t {
        mapping (address => Map3NodeWrapper_t) Nodes;
        mapping (uint64 => Map3NodeSnapshot_t) NodeSnapshotByEpoch;
        mapping (address => mapping (address => bool) ) NodeAddressSetByDelegator;
        mapping (string => bool) NodeKeySet;
        mapping (string => bool) DescriptionIdentitySet;
        mapping (string => uint256) NodePriority;
    }

    Map3NodePool_t pool;
    Map3Node_t node;
    int version;
    string name;

    constructor() public {
        version = 666;
        name = "Hyperion";
        node.NodeAddress = 0xA07306b4d845BD243Da172aeE557893172ccd04a;
        node.Commission.CommissionRates.Rate = 0x33 * (10**18);
        node.Description.Version[0]=0xbeef;
        node.Description.Version[1]=0xdead;
        node.Description.Name = "Hyperion - 海伯利安";
        node.Description.Details = "Hyperion, a decentralized map platform, aims to achieve the “One Map” vision - to provide an unified view of global map data and service, and to make it universally accessible just like a public utility for 10B people.\n海伯利安是去中心化的地图生态。";
        node.NodeKeys = bytes("MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDQBkQd2vUJtyNa2MBw4i8S0N9kQAAHwWdr1D5CPWgv/9GsGVCAUmLZhLV6E5JcrsL3fcKpak+oO+X3chffgOANVolvwqPUJif1ciimoMiEOU7+auLhTpRohX44phoCJ7J9C1nklTx1L6YHDrnMpvlAuRf0V6HM5Ro0L56LUMwZmwIDAQAB");
        pool.Nodes[0xA07306b4d845BD243Da172aeE557893172ccd04a].TotalDelegation = 0xdeadbeef;
        pool.Nodes[0xA07306b4d845BD243Da172aeE557893172ccd04a].Microdelegations[0x3CB0B0B6D52885760A5404eb0A593B979c88BcEF].PendingDelegationsfixed2dimension[2][1].Amount = 0xbeef;
        pool.Nodes[0xA07306b4d845BD243Da172aeE557893172ccd04a].Microdelegations[0x3CB0B0B6D52885760A5404eb0A593B979c88BcEF].PendingDelegationsfixed2dimension[0][0].Amount = 0xdead;
        for (uint i = 0; i < 10; i++) {
            pool.Nodes[0xA07306b4d845BD243Da172aeE557893172ccd04a].Microdelegations[0x3CB0B0B6D52885760A5404eb0A593B979c88BcEF].PendingDelegations.push();
        }

        pool.Nodes[0xA07306b4d845BD243Da172aeE557893172ccd04a].Microdelegations[0x3CB0B0B6D52885760A5404eb0A593B979c88BcEF].PendingDelegations[5].Amount = 0x7788;
        pool.NodeKeySet["0xA07306b4d845BD243Da172aeE557893172ccd04a"] = true;
    }

    function Version() public view returns (int) {
        return version;
    }

    function Length() public view returns (uint) {
        return pool.Nodes[0xA07306b4d845BD243Da172aeE557893172ccd04a].Microdelegations[0x3CB0B0B6D52885760A5404eb0A593B979c88BcEF].PendingDelegations.length;
    }
}
