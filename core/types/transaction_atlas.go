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
	StakeCreateVal
	StakeEditVal
	Redelegate
	Unredelegate
	CollectRedelRewards
)

func NewStakingTransaction(nonce uint64, gasLimit uint64, gasPrice *big.Int, data []byte, txType TransactionType) *Transaction {
	t := newTransaction(nonce, nil, common.Big0, gasLimit, gasPrice, data)
	t.SetType(txType)
	return t
}

func (txType TransactionType) String() string {
	// TODO(ATLAS)
	if txType == Normal {
		return "Normal"
	} else if txType == StakeCreateVal {
		return "StakeCreateValidator"
	} else if txType == StakeEditVal {
		return "StakeEditValidator"
	} else if txType == Redelegate {
		return "DelegateValidator"
	} else if txType == Unredelegate {
		return "UndelegateValidator"
	} else if txType == CollectRedelRewards {
		return "CollectRedelegationRewards"
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
