

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
./build-and-test.sh
````
or
````shell script
truffle compile all && go build ../../ && solc contracts/Map3Pool.sol --storage-layout | grep '^{' | ./storegen --storage-layout - --pkg test --out ../map3pool_storage.go && solc contracts/Map3Pool.sol --abi | grep '^\[' | abigen --abi - --pkg test --type Map3PoolWrapper --out ../map3pool_abi.go && go test -v -count=1 ../
````
or
````shell script
go build ../../ && 
solc contracts/Map3Pool.sol --overwrite --pretty-json --combined-json abi,bin,bin-runtime,storage-layout \
  | jq '.contracts[] | {"abi": (.abi | fromjson), "bytecode": .bin, "deployedBytecode": ."bin-runtime", "storage-layout": (."storage-layout" | fromjson) }' \
  | tee \
    >( >build/contracts/Map3Pool.json) \
    >(jq '.abi' | abigen --abi - --pkg storage --type consortium --out ../map3pool_abi.go) \
    >(jq '."storage-layout"' | storegen --storage-layout - --pkg storage --out ../map3pool_storage.go)
&& go test -v -count=1 ../
````

# ROOT directory
````shell script
go build ./cmd/storegen/ && ./storegen --storage-layout ./cmd/storegen/test/data/layout/Map3Pool_storage.json --pkg test > ./cmd/storegen/test/map3pool_storage.go && go test -v ./cmd/storegen/test/
````