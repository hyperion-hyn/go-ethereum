package state

import (
	"github.com/ethereum/go-ethereum/common"
	staking "github.com/ethereum/go-ethereum/staking/types"
	"github.com/pkg/errors"
)

var (
	ErrMap3NodeNotExist  = errors.New("staking map3 node does not exist")
	ErrValidatorNotExist = errors.New("staking validator does not exist")
)

func (s *StateDB) Map3NodePool() *staking.Map3NodePoolStorage {
	return s.map3NodePool
}

func (s *StateDB) ValidatorPool() *staking.ValidatorPoolStorage {
	return s.validatorPool
}

func (s *StateDB) Map3NodeByAddress(nodeAddress common.Address) (*staking.Map3NodeWrapperStorage, error) {
	node, ok := s.Map3NodePool().GetNodes().Get(nodeAddress)
	if !ok {
		return nil, ErrMap3NodeNotExist
	}
	return node, nil
}

func (s *StateDB) ValidatorByAddress(validatorAddress common.Address) (*staking.ValidatorWrapperStorage, error) {
	validator, ok := s.ValidatorPool().GetValidators().Get(validatorAddress)
	if !ok {
		return nil, ErrValidatorNotExist
	}
	return validator, nil
}

func (s *StateDB) ValidatorList() []common.Address {
	return s.validatorPool.GetValidators().Keys()
}
