
# generate storage-layout
````shell script
solc contracts/consortium.sol --overwrite --storage-layout 
````

# generate storage-wrapper
```shell script
storegen --storage-layout FILE --pkg PKG
```
# generate abi-wrapper
```shell script
solc contracts/consortium.sol --abi | grep '^\[' | abigen --abi - --pkg storage --type consortium --out storage/consortium_abi.go
```
# combined
````shell script
solc contracts/consortium.sol --storage-layout | grep '^{' | storegen --storage-layout - --pkg storage --out storage/consortium_storage.go
````
