```
cd solidity
solc ./Validator.sol --overwrite --storage-layout --output-dir layout
storegen --storage-layout ./layout/ValidatorPoolWrapper_storage.json --pkg validator > ../validator/validator_gen.go
``` 