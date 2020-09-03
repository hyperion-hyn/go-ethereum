package microstaking

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	common2 "github.com/ethereum/go-ethereum/staking/types/common"
	"github.com/pkg/errors"
	"math/big"
)

const (
	DoNotEnforceMaxBLS = -1
	MaxPubKeyAllowed   = 1
)

var (
	errNeedAtLeastOneSlotKey   = errors.New("need at least one slot key")
	ErrExcessiveBLSKeys        = errors.New("more slot keys provided than allowed")
	errDuplicateNodeKeys       = errors.New("map3 node keys can not have duplicates")
	errAddressNotMatch         = errors.New("validator key not match")
	errNodeKeyToRemoveNotFound = errors.New("map3 node key to remove not found")
	errNodeKeyToAddExists      = errors.New("map3 node key to add already exists")
	errMicrodelegationNotExist = errors.New("microdelegation does not exist")
)

var (
	Map3NodeLockDurationInEpoch = common.NewDec(180)
)

type Map3Status byte

const (
	// Nil is a default state that represents a no-op
	Nil Map3Status = iota
	// Pending means total delegation of this map3 node is still not enough
	Pending
	Active
	Terminated
)

func (e Map3Status) String() string {
	switch e {
	case Pending:
		return "pending"
	case Active:
		return "active"
	case Terminated:
		return "terminated"
	default:
		return "unknown"
	}
}

// SanityCheck checks basic requirements of a validator
func (n *Map3Node_) SanityCheck(maxPubKeyAllowed int) error {
	if err := n.Description.EnsureLength(); err != nil {
		return err
	}

	if len(n.NodeKeys.Keys) == 0 {
		return errNeedAtLeastOneSlotKey
	}

	if c := len(n.NodeKeys.Keys); maxPubKeyAllowed != DoNotEnforceMaxBLS &&
		c > maxPubKeyAllowed {
		return errors.Wrapf(
			ErrExcessiveBLSKeys, "have: %d allowed: %d",
			c, maxPubKeyAllowed,
		)
	}

	if err := n.Commission.SanityCheck(); err != nil {
		return err
	}

	allKeys := map[string]struct{}{}
	for i := range n.NodeKeys.Keys {
		key := n.NodeKeys.Keys[i].Hex()
		if _, ok := allKeys[key]; !ok {
			allKeys[key] = struct{}{}
		} else {
			return errDuplicateNodeKeys
		}
	}
	return nil
}

// Storage_Map3NodeWrapper_
func (s *Storage_Map3NodeWrapper_) AddTotalDelegation(amount *big.Int) {
	totalDelegation := s.TotalDelegation().Value()
	totalDelegation = totalDelegation.Add(totalDelegation, amount)
	s.TotalDelegation().SetValue(totalDelegation)
}

func (s *Storage_Map3NodeWrapper_) SubTotalDelegation(amount *big.Int) {
	totalDelegation := s.TotalDelegation().Value()
	totalDelegation = totalDelegation.Sub(totalDelegation, amount)
	s.TotalDelegation().SetValue(totalDelegation)
}

func (s *Storage_Map3NodeWrapper_) AddTotalPendingDelegation(amount *big.Int) {
	totalPendingDelegation := s.TotalPendingDelegation().Value()
	totalPendingDelegation = totalPendingDelegation.Add(totalPendingDelegation, amount)
	s.TotalPendingDelegation().SetValue(totalPendingDelegation)
}

func (s *Storage_Map3NodeWrapper_) SubTotalPendingDelegation(amount *big.Int) {
	totalPendingDelegation := s.TotalPendingDelegation().Value()
	totalPendingDelegation = totalPendingDelegation.Sub(totalPendingDelegation, amount)
	s.TotalPendingDelegation().SetValue(totalPendingDelegation)
}

func (s *Storage_Map3NodeWrapper_) AddMicrodelegation(delegator common.Address, amount *big.Int,
	pending bool, epoch *big.Int) (isNewDelegator bool) {
	isExist := s.Microdelegations().Contain(delegator)
	if !isExist {
		s.Microdelegations().Put(delegator, &Microdelegation_{})
	}
	md, _ := s.Microdelegations().Get(delegator)
	if pending {
		md.PendingDelegation().AddAmount(amount, epoch)
		s.AddTotalDelegation(amount)
	} else {
		md.AddAmount(amount)
		s.AddTotalPendingDelegation(amount)
	}
	return !isExist
}

func (s *Storage_Map3NodeWrapper_) IsOperator(delegator common.Address) bool {
	return s.Map3Node().OperatorAddress().Value() == delegator
}

func (s *Storage_Map3NodeWrapper_) IsAlreadyRestaking() bool {
	addr0 := common.Address{}
	return s.RestakingReference().ValidatorAddress().Value() != addr0
}

func (s *Storage_Map3NodeWrapper_) Unmicrodelegate(delegator common.Address, amount *big.Int) (toReturn *big.Int, completed bool) {
	if md, ok := s.Microdelegations().Get(delegator); ok {
		if pd := md.PendingDelegation().Amount().Value(); pd.Cmp(amount) < 0 {
			amount = big.NewInt(0).Set(pd)
		}
		md.PendingDelegation().SubAmount(amount)
		s.SubTotalPendingDelegation(amount)
		toReturn = big.NewInt(0).Set(amount)

		if md.Amount().Value().Uint64() == 0 || md.Undelegation().Amount().Value().Uint64() == 0 {
			toReturn = toReturn.Add(toReturn, md.Reward().Value())
			s.Microdelegations().Remove(delegator)
			return toReturn, true
		}
		return toReturn, false
	}
	return common.Big0, true
}

func (s *Storage_Map3NodeWrapper_) CanActivateMap3Node(requireTotal, requireSelf *big.Int) bool {
	if s.Map3Node().Status().Value() != uint8(Pending) {
		return false
	}

	total := big.NewInt(0).Add(s.TotalPendingDelegation().Value(), s.TotalDelegation().Value())
	if total.Cmp(requireTotal) >= 0 {
		operator := s.Map3Node().OperatorAddress().Value()
		m, ok := s.Microdelegations().Get(operator)
		if !ok {
			log.Error("operator's delegation should exist", "map3", s.Map3Node().Map3Address().Value().String())
			return false
		}

		self := big.NewInt(0).Add(m.Amount().Value(), m.PendingDelegation().Amount().Value())
		if self.Cmp(requireSelf) >= 0 {
			return true
		}
	}
	return false
}

func (s *Storage_Map3NodeWrapper_) ActivateMap3Node(epoch *big.Int) error {
	// change pending delegation
	for _, delegator := range s.Microdelegations().AllKeys() {
		delegation, ok := s.Microdelegations().Get(delegator)
		if !ok {
			return errors.Wrapf(errMicrodelegationNotExist, "delegation should exist, map3: %v, delegator: %v",
				s.Map3Node().Map3Address().Value().String(), delegator.String())
		}
		pd := delegation.PendingDelegation().Amount().Value()
		delegation.AddAmount(pd)
		delegation.PendingDelegation().Clear()
	}
	s.AddTotalDelegation(s.TotalPendingDelegation().Value())
	s.TotalPendingDelegation().SetValue(common.Big0)

	// update state
	status := s.Map3Node().Status().Value()
	if status == uint8(Pending) {
		time := common.OneDec().Mul(Map3NodeLockDurationInEpoch)
		s.Map3Node().ReleaseEpoch().SetValue(time)
	}
	s.Map3Node().Status().SetValue(uint8(Active))
	s.Map3Node().ActivationEpoch().SetValue(epoch)
	return nil
}

func (s *Storage_Map3NodeWrapper_) LoadFully() (*Map3NodeWrapper_, error) {
	s.Map3Node().load()
	if _, err := s.Microdelegations().LoadFully(); err != nil {
		return nil, err
	}
	s.RestakingReference().load()
	s.AccumulatedReward().Value()
	s.TotalDelegation().Value()
	s.TotalPendingDelegation().Value()

	// copy
	src := s.obj
	des := Map3NodeWrapper_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

// Storage_ValidatorWrapperMap_
func (s *Storage_Map3NodeWrapperMap_) AllKeys() []common.Address {
	addressSlice := make([]common.Address, 0)
	addressLength := s.Keys().Length()
	for i := 0; i < addressLength; i++ {
		addressSlice = append(addressSlice, s.Keys().Get(i).Value())
	}
	return addressSlice
}

func (s *Storage_Map3NodeWrapperMap_) Put(key common.Address, map3Node *Map3NodeWrapper_) {
	if s.Contain(key) {
		s.Map().Get(key).Entry().Save(map3Node)
	} else {
		keysLength := s.Keys().Length()
		//set keys
		s.Keys().Get(keysLength).SetValue(key)
		//set map
		sValidatorWrapper := s.Map().Get(key)
		//set map entity
		sValidatorWrapperEntity := sValidatorWrapper.Entry()
		sValidatorWrapperEntity.Save(map3Node)
		//set map index
		sValidatorWrapper.Index().SetValue(big.NewInt(0).Add(big.NewInt(int64(keysLength)), common.Big1)) //because index start with 1
	}
}

func (s *Storage_Map3NodeWrapperMap_) Contain(key common.Address) bool {
	return s.Map().Get(key).Index().Value().Cmp(common.Big0) > 0
}

func (s *Storage_Map3NodeWrapperMap_) Get(key common.Address) (*Storage_Map3NodeWrapper_, bool) {
	if s.Contain(key) {
		return s.Map().Get(key).Entry(), true
	}
	return nil, false
}

func (s *Storage_Map3NodePool_) UpdateDelegationIndex(delegator common.Address, index *DelegationIndex_) {
	indexMap := s.DelegationIndexMapByDelegator().Get(delegator)
	indexMap.Put(delegator, index)
}

func (s *Storage_Map3NodePool_) RemoveDelegationIndex(delegator, map3Addr common.Address) {
	indexMap := s.DelegationIndexMapByDelegator().Get(delegator)
	indexMap.Remove(map3Addr)
}

// CreateValidatorFromNewMsg creates validator from NewValidator message
func CreateMap3NodeFromNewMsg(msg *CreateMap3Node, map3Addr common.Address, blockNum *big.Int) *Map3Node_ {
	v := Map3Node_{
		Map3Address:     map3Addr,
		OperatorAddress: msg.OperatorAddress,
		NodeKeys:        NewBLSKeysWithBLSKey(msg.NodePubKey),
		Commission: Commission_{
			Rate:              msg.Commission,
			RateForNextPeriod: msg.Commission,
		},
		Description:    msg.Description,
		CreationHeight: blockNum,
		Status:         uint8(Pending),
	}
	return &v
}

// UpdateValidatorFromEditMsg updates validator from EditValidator message
func UpdateMap3NodeFromEditMsg(map3Node *Map3Node_, edit *EditMap3Node) error {
	if map3Node.Map3Address != edit.Map3NodeAddress {
		return errAddressNotMatch
	}
	if err := map3Node.Description.UpdateFrom(edit.Description); err != nil {
		return err
	}

	if edit.NodeKeyToRemove != nil {
		index := -1
		for i, key := range map3Node.NodeKeys.Keys {
			if *key == *edit.NodeKeyToRemove {
				index = i
				break
			}
		}
		// we found key to be removed
		if index >= 0 {
			map3Node.NodeKeys.Keys = append(
				map3Node.NodeKeys.Keys[:index], map3Node.NodeKeys.Keys[index+1:]...,
			)
		} else {
			return errNodeKeyToRemoveNotFound
		}
	}

	if edit.NodeKeyToAdd != nil {
		found := false
		for _, key := range map3Node.NodeKeys.Keys {
			if *key == *edit.NodeKeyToAdd {
				found = true
				break
			}
		}
		if !found {
			if err := common2.VerifyBLSKey(edit.NodeKeyToAdd, edit.NodeKeyToAddSig); err != nil {
				return err
			}
			map3Node.NodeKeys.Keys = append(map3Node.NodeKeys.Keys, edit.NodeKeyToAdd)
		} else {
			return errNodeKeyToAddExists
		}
	}
	return nil
}
