// SPDX-License-Identifier: LGPL-3.0 License
pragma solidity ^0.6.10;

contract consortium {
    struct Validator_t {
        bytes PubKey; // Signer's Public Key
        address Account; // Signer's Account
    }

    struct Committee_t {
        Validator_t[] members;
    }

    Committee_t committee;

    constructor() public{
    }
}
