package types

import (
    "github.com/ethereum/go-ethereum/common"
)
//go:generate gencodec -type Validator -field-override validatorMarshaling -out gen_validator.go

var (
    // ATLAS(zgx): generate a StakingInfoAddress
    StakingInfoAddress = common.StringToAddress("0x0000000000000000000000000000000123456789") // used to save staking state in state db
)

type ValidatorContainer struct {
    Validators []ValidatorWrapper `json:"validators"`
}

func (container ValidatorContainer) IsValidator(addr common.Address) bool {
    return container.Validator(addr) != nil
}

func (container ValidatorContainer) Validator(addr common.Address) *ValidatorWrapper {
    for _, val := range container.Validators {
        if addr == val.Validator.Address {
            return &val
        }
    }
    return nil
}
