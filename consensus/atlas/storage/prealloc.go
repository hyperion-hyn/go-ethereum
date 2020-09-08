package storage

import (
	"bytes"
	"math/big"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
)

type Signer struct {
	PublicKey *bls.PublicKey
	Coinbase  common.Address
}

func SetupValidatorsInGenesisAt(account *core.GenesisAccount, validators []*Signer) error {
	// Prepare storage wrapper to hold state modifications.
	stateDB, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	if err != nil {
		return err
	}

	signers := make([]*Signer, len(validators))
	for i := 0; i < len(validators); i++ {
		signers[i] = validators[i]
	}

	// Sort the signers
	for i := 0; i < len(signers); i++ {
		for j := i + 1; j < len(signers); j++ {
			if bytes.Compare(signers[i].PublicKey.Serialize()[:], signers[j].PublicKey.Serialize()[:]) > 0 {
				signers[i], signers[j] = signers[j], signers[i]
			}
		}
	}

	tracker := NewStateTracker(stateDB)
	var global Global_t
	storage := New(&global, tracker, common.Address{}, big.NewInt(0))

	for i := 0; i < len(signers); i++ {
		memberStorage := storage.Committee().Members().Get(i)
		memberStorage.PublicKey().SetValue(signers[i].PublicKey.Serialize())
		memberStorage.Coinbase().SetValue(signers[i].Coinbase)
	}

	if account.Storage == nil {
		account.Storage = make(map[common.Hash]common.Hash)
	}

	for k, v := range tracker.Dirties() {
		account.Storage[k] = v
	}

	return nil
}
