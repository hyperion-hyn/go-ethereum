// SPDX-License-Identifier: LGPL


pragma solidity ^0.5.1;

contract Storage {

    int version;
    string name;

    constructor() public{

    }

    function Version() public view returns (int) {
        return version;
    }

    function Name() public view returns (string memory) {
        return name;
    }
}
