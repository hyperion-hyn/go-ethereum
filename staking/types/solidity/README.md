## Validator
```
cd solidity
solc ./Validator.sol --overwrite --storage-layout --output-dir layout
storegen --storage-layout ./layout/ValidatorPoolWrapper_storage.json --pkg restaking > ../restaking/entity_gen.go
```

## Map3 Node
```
cd solidity
solc ./Map3.sol --overwrite --storage-layout --output-dir layout
storegen --storage-layout ./layout/Map3NodePoolWrapper_storage.json --pkg microstaking > ../microstaking/entity_gen.go
``` 