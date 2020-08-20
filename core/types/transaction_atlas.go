package types

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
