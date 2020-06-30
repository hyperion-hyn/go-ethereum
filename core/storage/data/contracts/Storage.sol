// SPDX-License-Identifier: LGPL


pragma solidity ^0.5.1;

contract Storage {

    int version;
    string name;

    constructor() public{

    }

    function Hello() public pure returns (string memory res) {
        return "hello world";
    }

    function Version() public view returns (int v) {
        return version;
    }

    function Name() public view returns (string memory) {
        return name;
    }
}
