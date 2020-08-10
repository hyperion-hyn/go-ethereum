
# generate storage-layout
````shell script
solc contracts/consortium.sol --overwrite --storage-layout 
````

# generate storage-wrapper
```shell script
storegen --storage-layout FILE --pkg PKG
```

# combined
````shell script
solc contracts/consortium.sol --overwrite --storage-layout | grep '^{' | storegen --storage-layout - --pkg storage --out storage/consortium.go
````
