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
solc contracts/Storage.sol --storage-layout
````

# combine
````shell script
cd cmd/storegen/test/data
````

````shell script
truffle compile all && go build ../../ && solc contracts/Map3Pool.sol --storage-layout | grep '^{' | ./storegen --storage-layout - --pkg test --out ../map3pool_storage_test.go && solc contracts/Map3Pool.sol --abi | grep '^\[' | abigen --abi - --pkg test --type Map3PoolWrapper --out ../map3pool_abi_test.go && go test -v -count=1 ../
````

