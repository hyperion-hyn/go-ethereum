

# data/ directory
## Compile contract
````shell script
truffle compile all
````

## Generate storage-layout
````shell script
solc contracts/Map3Pool.sol --overwrite --storage-layout --output-dir layout
````

## Generate ABI wrapper
````shell script
cat build/contracts/Map3Pool.json | jq -c '.abi' | abigen --abi - --pkg test --type Map3PoolWrapper --out ../wrapper_test.go
````

## Combine
````shell script
truffle compile all && solc contracts/Map3Pool.sol --overwrite --storage-layout --output-dir layout && cat build/contracts/Map3Pool.json | jq -c '.abi' | abigen --abi - --pkg test --type Map3PoolWrapper --out ../wrapper_test.go
````

# ROOT directory
````shell script
go build ./cmd/storegen/ && ./storegen --storage-layout ./cmd/storegen/test/data/layout/Map3Pool_storage.json --pkg test > ./cmd/storegen/test/map3pool_test.go && go test -v ./cmd/storegen/test/
````