package types

// TransactionType different types of transactions
type TransactionType byte

// Different Transaction Types
const (
	Normal TransactionType = iota
	StakeCreateNode
	StakeEditNode
	Microdelegate
	Unmicrodelegate
	CollectMicrodelRewards
	StakeCreateVal
	StakeEditVal
	Redelegate
	Unredelegate
	CollectMicroreRewards
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
	} else if txType == CollectMicroreRewards {
		return "CollectMicroredelegationRewards"
	}
	return "Normal"
}

