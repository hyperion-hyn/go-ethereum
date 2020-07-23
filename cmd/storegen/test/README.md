

# data/ directory
## Compile contract
```` 
truffle compile all
````

## Generate storage-layout
````
solc contracts/Map3Pool.sol --overwrite --storage-layout --output-dir layout
````

## Generate ABI wrapper
````
cat build/contracts/Map3Pool.json | jq -c '.abi' | abigen --abi - --pkg test --type Map3PoolWrapper --out ../wrapper_test.go
````

# ROOT directory
````
go build ./cmd/storegen/ && ./storegen --storage-layout ./cmd/storegen/test/data/layout/Map3Pool_storage.json --pkg test > ./cmd/storegen/test/map3pool_test.go && go test -v ./cmd/storegen/test/
````