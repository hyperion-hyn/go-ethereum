package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
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

// deprecated
func NewLegacyTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, txType TransactionType) *Transaction {
	if len(data) > 0 {
		data = common.CopyBytes(data)
	}
	d := LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Data:     data,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int),
		V:        new(big.Int),
		R:        new(big.Int),
		S:        new(big.Int),
		Type:     txType,
	}
	if amount != nil {
		d.Value.Set(amount)
	}
	if gasPrice != nil {
		d.GasPrice.Set(gasPrice)
	}
	return &Transaction{
		data: &d,
		time: time.Now(),
	}
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

func (tx *Transaction) Type() TransactionType {
	return tx.data.txType()
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

// LegacyTx is the transaction data of regular ATLAS transactions.
type LegacyTx struct {
	Nonce    uint64          // nonce of sender account
	GasPrice *big.Int        // wei per gas
	Gas      uint64          // gas limit
	To       *common.Address `rlp:"nil"` // nil means contract creation
	Value    *big.Int        // wei amount
	Data     []byte          // contract invocation input data
	V, R, S  *big.Int        // signature values
	Type     TransactionType // ATLAS: include staking tx
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *LegacyTx) copy() TxData {
	cpy := &LegacyTx{
		Nonce: tx.Nonce,
		To:    tx.To, // TODO: copy pointed-to address
		Data:  common.CopyBytes(tx.Data),
		Gas:   tx.Gas,
		// These are initialized below.
		Value:    new(big.Int),
		GasPrice: new(big.Int),
		V:        new(big.Int),
		R:        new(big.Int),
		S:        new(big.Int),
		Type:     tx.Type,
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.GasPrice != nil {
		cpy.GasPrice.Set(tx.GasPrice)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *LegacyTx) chainID() *big.Int       { return deriveChainId(tx.V) }
func (tx *LegacyTx) data() []byte            { return tx.Data }
func (tx *LegacyTx) gas() uint64             { return tx.Gas }
func (tx *LegacyTx) gasPrice() *big.Int      { return tx.GasPrice }
func (tx *LegacyTx) value() *big.Int         { return tx.Value }
func (tx *LegacyTx) nonce() uint64           { return tx.Nonce }
func (tx *LegacyTx) to() *common.Address     { return tx.To }
func (tx *LegacyTx) txType() TransactionType { return tx.Type }

func (tx *LegacyTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *LegacyTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.V, tx.R, tx.S = v, r, s
}
