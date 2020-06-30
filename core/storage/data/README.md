# compile
````shell script
truffle compile
````
# abi
````shell script
cat build/contracts/Storage.json | jq -c '.abi'
````

# bytecode
````shell script
cat build/contracts/Storage.json | jq '.bytecode'
````

# deployedBytecode

````shell script
cat build/contracts/Storage.json | jq '.deployedBytecode'
````
