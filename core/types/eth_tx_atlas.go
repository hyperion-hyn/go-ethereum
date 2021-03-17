package types

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

func NewTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	return newTransaction(nonce, &to, amount, gasLimit, gasPrice, data)
}

func NewContractCreation(nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	return newTransaction(nonce, nil, amount, gasLimit, gasPrice, data)
}

func NewStakingTransaction(nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, txType TransactionType) (*Transaction, error) {
	if txType == Normal {
		return nil, errors.New("invalid transaction type")
	}

	var to common.Address
	switch txType {
	case CreateMap3:
	case EditMap3:
	case TerminateMap3:
	case Microdelegate:
	case Unmicrodelegate:
	case CollectMicrostakingRewards:
	case RenewMap3Node:
		to = Map3Account
	case CreateValidator:
	case EditValidator:
	case Redelegate:
	case Unredelegate:
	case CollectRestakingReward:
		to = ValidatorAccount
	default:
		return nil, errors.New("invalid transaction type")
	}

	payload := append([]byte{byte(txType)}, data...)
	return newTransaction(nonce, &to, amount, gasLimit, gasPrice, payload), nil
}

func newTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	if len(data) > 0 {
		data = common.CopyBytes(data)
	}
	d := ETHTx{
		Nonce:    nonce,
		To:       to,
		Data:     data,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int),
		V:        new(big.Int),
		R:        new(big.Int),
		S:        new(big.Int),
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

// ETHTx is the transaction data of ATLAS transactions compatible with ETH transactions.
type ETHTx struct {
	Nonce    uint64          // nonce of sender account
	GasPrice *big.Int        // wei per gas
	Gas      uint64          // gas limit
	To       *common.Address `rlp:"nil"` // nil means contract creation
	Value    *big.Int        // wei amount
	Data     []byte          // contract invocation input data
	V, R, S  *big.Int        // signature values
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *ETHTx) copy() TxData {
	cpy := &ETHTx{
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
func (tx *ETHTx) chainID() *big.Int   { return deriveChainId(tx.V) }
func (tx *ETHTx) data() []byte        { return tx.Data }
func (tx *ETHTx) gas() uint64         { return tx.Gas }
func (tx *ETHTx) gasPrice() *big.Int  { return tx.GasPrice }
func (tx *ETHTx) value() *big.Int     { return tx.Value }
func (tx *ETHTx) nonce() uint64       { return tx.Nonce }
func (tx *ETHTx) to() *common.Address { return tx.To }
func (tx *ETHTx) txType() TransactionType {
	if tx.to() == nil {
		return Normal
	}
	if *tx.to() == Map3Account || *tx.to() == ValidatorAccount {
		return TransactionType(tx.Data[0])
	}
	return Normal
}

func (tx *ETHTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *ETHTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.V, tx.R, tx.S = v, r, s
}
