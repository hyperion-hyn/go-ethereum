package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type RedelegationReferenceStorage struct {
}

func (r *RedelegationReferenceStorage) ToRedelegationReference() *RedelegationReference {
	return nil
}

func (r *RedelegationReferenceStorage) GetValidatorAddress() common.Address {
	return [20]byte{}
}

func (r *RedelegationReferenceStorage) SetValidatorAddress(validatorAddress common.Address) {
}

func (r *RedelegationReferenceStorage) SetReleasedTotalDelegation(releasedTotalDelegation *big.Int) {
}

func (r *RedelegationReferenceStorage) GetReleasedTotalDelegation() *big.Int {
	return nil
}
