package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// TransactionType different types of transactions
type TransactionType byte

// Different Transaction Types
const (
	Normal TransactionType = iota

	// validator
	CreateValidator
	EditValidator
	Redelegate
	Unredelegate
	CollectRestakingReward

	// map3 node
	CreateMap3
	EditMap3
	TerminateMap3
	Microdelegate
	Unmicrodelegate
	CollectMicrostakingRewards
	RenewMap3Node
)

func NewStakingTransaction(nonce uint64, gasLimit uint64, gasPrice *big.Int, data []byte, txType TransactionType) *Transaction {
	t := newTransaction(nonce, nil, common.Big0, gasLimit, gasPrice, data)
	t.SetType(txType)
	return t
}

func (txType TransactionType) String() string {
	if txType == Normal {
		return "Normal"
	} else if txType == CreateValidator {
		return "CreateValidator"
	} else if txType == EditValidator {
		return "EditValidator"
	} else if txType == Redelegate {
		return "Redelegate"
	} else if txType == Unredelegate {
		return "Unredelegate"
	} else if txType == CollectRestakingReward {
		return "CollectRestakingReward"
	} else if txType == CreateMap3 {
		return "CreateMap3Node"
	} else if txType == EditMap3 {
		return "EditMap3Node"
	} else if txType == TerminateMap3 {
		return "TerminateMap3Node"
	} else if txType == Microdelegate {
		return "Microdelegate"
	} else if txType == Unmicrodelegate {
		return "Unmicrodelegate"
	} else if txType == CollectMicrostakingRewards {
		return "CollectMicrostakingRewards"
	} else if txType == RenewMap3Node {
		return "RenewMap3Node"
	}
	return "Normal"
}

func (tx *Transaction) SetType(txType TransactionType) {
	tx.data.Type = txType
}

func (tx *Transaction) Type() TransactionType {
	return tx.data.Type
}

func (m Message) Type() TransactionType { return m.txType }

func NewMessageWithTransactionType(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, checkNonce bool, txType TransactionType) Message {
	return Message{
		from:       from,
		to:         to,
		nonce:      nonce,
		amount:     amount,
		gasLimit:   gasLimit,
		gasPrice:   gasPrice,
		data:       data,
		checkNonce: checkNonce,
		txType:     txType,
	}
}
