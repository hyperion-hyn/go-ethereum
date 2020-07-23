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

# generate bindings
````shell script
cat build/contracts/Storage.json | jq -c '.abi' | abigen --abi - --pkg storage --type Storage --out wrapper.go
````

# generate storage-layout
````shell script
solc contracts/Storage.sol --overwrite --storage-layout
````

