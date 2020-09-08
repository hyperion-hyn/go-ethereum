# generate json
````shell script
mkdir build/contracts
solc contracts/consortium.sol --pretty-json --combined-json abi,bin,bin-runtime,storage-layout | jq '.contracts[] | {"abi": (.abi | fromjson), "bytecode": .bin, "deployedBytecode": ."bin-runtime", "storage-layout": (."storage-layout" | fromjson) }' > build/contracts/consortium.json
```` 

# generate storage-wrapper
```shell script
solc contracts/consortium.sol --storage-layout | grep '^{' | storegen --storage-layout - --pkg storage --out storage/consortium_storage.go
```

# generate abi-wrapper
```shell script
solc contracts/consortium.sol --abi | grep '^\[' | abigen --abi - --pkg storage --type consortium --out storage/consortium_abi.go
````

# combined
````shell script
solc contracts/consortium.sol --overwrite --pretty-json --combined-json abi,bin,bin-runtime,storage-layout \
  | jq '.contracts[] | {"abi": (.abi | fromjson), "bytecode": .bin, "deployedBytecode": ."bin-runtime", "storage-layout": (."storage-layout" | fromjson) }' \
  | tee \
    >( >build/contracts/consortium.json) \
    >(jq '.abi' | abigen --abi - --pkg storage --type consortium --out storage/consortium_abi.go) \
    >(jq '."storage-layout"' | storegen --storage-layout - --pkg storage --out storage/consortium_storage.go)
````
