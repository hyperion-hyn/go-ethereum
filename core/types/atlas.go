package types

// TransactionType different types of transactions
type TransactionType byte

// Different Transaction Types
const (
	Normal TransactionType = iota

	// map3 node
	StakeCreateNode
	StakeEditNode
	StopNodeStake
	ResumeNodeStake
	DivideNodeStake
	RenewNodeStake
	Microdelegate
	Unmicrodelegate
	CollectMicrodelRewards

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
		return "CollectMicroredelegationRewards"
	}
	return "Normal"
}

