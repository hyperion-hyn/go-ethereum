// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package backend

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/consensus/atlas/storage"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// in this test, we can set n to 1, and it means we can process Atlas and commit a
// block by one node. Otherwise, if n is larger than 1, we have to generate
// other fake events to process Atlas.
func newBlockChain(n int) (*core.BlockChain, *backend, []*bls.SecretKey) {
	genesis, _, signerKeys := getGenesisAndKeys(n)
	memDB := rawdb.NewMemoryDatabase()
	config := atlas.DefaultConfig
	// Use the first key as private key
	b, _ := New(config, memDB).(*backend)
	genesis.MustCommit(memDB)
	blockchain, err := core.NewBlockChain(memDB, nil, genesis.Config, b, vm.Config{}, nil)
	if err != nil {
		panic(err)
	}
	b.Start(blockchain, blockchain.CurrentBlock, blockchain.HasBadBlock)
	snap, err := b.snapshot(blockchain, 0, common.Hash{}, nil)
	if err != nil {
		panic(err)
	}
	if snap == nil {
		panic("failed to get snapshot")
	}
	proposerAddr := snap.ValSet.GetProposer().Signer()

	// find proposer key
	for _, key := range signerKeys {
		addr := crypto.PubkeyToSigner(key.GetPublicKey())
		if addr.String() == proposerAddr.String() {
			b.signHashFn = func(account accounts.Account, hash common.Hash) (signature []byte, publicKey []byte, mask []byte, err error) {
				secrectKey := key
				sign := secrectKey.SignHash(hash.Bytes())

				return sign.Serialize(), secrectKey.GetPublicKey().Serialize(), nil, nil
			}
			b.signer = addr
		}
	}

	return blockchain, b, signerKeys
}

func getGenesisAndKeys(n int) (*core.Genesis, []*ecdsa.PrivateKey, []*bls.SecretKey) {
	// Setup validators
	var nodeKeys = make([]*ecdsa.PrivateKey, n)
	var signerKeys = make([]*bls.SecretKey, n)

	var addrs = make([]common.Address, n)
	for i := 0; i < n; i++ {
		nodeKeys[i], _ = crypto.GenerateKey()
		signerKeys[i], _ = crypto.GenerateBLSKey()
		addrs[i] = crypto.PubkeyToAddress(nodeKeys[i].PublicKey)
	}

	// generate genesis block
	genesis := core.DefaultGenesisBlock()
	genesis.Config = params.TestChainConfig
	// force enable Atlas engine
	genesis.Config.Atlas = &params.AtlasConfig{}
	genesis.Config.Ethash = nil
	genesis.Difficulty = DefaultDifficulty
	genesis.Nonce = emptyNonce.Uint64()
	genesis.Mixhash = types.AtlasDigest

	appendValidators(genesis, signerKeys, addrs)
	return genesis, nodeKeys, signerKeys
}

func appendValidators(genesis *core.Genesis, signers []*bls.SecretKey, addrs []common.Address) error {
	if len(genesis.ExtraData) < types.AtlasExtraVanity {
		genesis.ExtraData = append(genesis.ExtraData, bytes.Repeat([]byte{0x00}, types.AtlasExtraVanity-len(genesis.ExtraData))...)
	}
	genesis.ExtraData = genesis.ExtraData[:types.AtlasExtraVanity]

	account := &core.GenesisAccount{
		Code:       nil,
		Storage:    make(map[common.Hash]common.Hash),
		Balance:    big.NewInt(1),
		Nonce:      0,
		PrivateKey: nil,
	}

	validators := make([]*storage.Signer, len(signers))
	for i := 0; i < len(signers); i++ {
		validators[i] = &storage.Signer{
			PublicKey: signers[i].GetPublicKey(),
			Coinbase:  addrs[i],
		}
	}

	err := storage.SetupValidatorsInGenesisAt(account, validators)
	if err != nil {
		return err
	}

	genesis.Alloc[common.HexToAddress(CONSORTIUM_BOARD)] = *account

	block := genesis.ToBlock(nil)
	hashdata := SealHash(block.Header())

	signatures := make([]*bls.Sign, len(signers))
	publicKeys := make([]*bls.PublicKey, len(signers))
	for i := 0; i < len(signers); i++ {
		signatures[i] = signers[i].SignHash(hashdata.Bytes())
		publicKeys[i] = signers[i].GetPublicKey()
	}

	err = WriteCommittedSealInGenesis(genesis, block.Header().Extra, signatures, publicKeys)
	if err != nil {
		return err
	}

	return nil
}

func makeHeader(parent *types.Block, config *atlas.Config) *types.Header {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent, parent.GasLimit(), parent.GasLimit()),
		GasUsed:    0,
		Extra:      parent.Extra(),
		Time:       parent.Time() + config.BlockPeriod,
		Difficulty: DefaultDifficulty,
	}
	return header
}

func makeBlock(chain *core.BlockChain, engine *backend, parent *types.Block) *types.Block {
	block := makeBlockWithoutSeal(chain, engine, parent)
	stopCh := make(chan struct{})
	resultCh := make(chan *types.Block, 10)
	go engine.Seal(chain, block, resultCh, stopCh)
	blk := <-resultCh
	return blk
}

func makeBlockWithoutSeal(chain *core.BlockChain, engine *backend, parent *types.Block) *types.Block {
	header := makeHeader(parent, engine.config)
	engine.Prepare(chain, header)
	state, _ := chain.StateAt(parent.Root())
	block, _ := engine.FinalizeAndAssemble(chain, header, state, nil, nil, nil)
	return block
}

func TestPrepare(t *testing.T) {
	chain, engine, _ := newBlockChain(1)
	header := makeHeader(chain.Genesis(), engine.config)
	err := engine.Prepare(chain, header)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	header.ParentHash = common.StringToHash("1234567890")
	err = engine.Prepare(chain, header)
	if err != consensus.ErrUnknownAncestor {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrUnknownAncestor)
	}
}

func TestSealStopChannel(t *testing.T) {
	chain, engine, _ := newBlockChain(4)
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	stop := make(chan struct{}, 1)
	eventSub := engine.EventMux().Subscribe(atlas.RequestEvent{})
	eventLoop := func() {
		select {
		case ev := <-eventSub.Chan():
			_, ok := ev.Data.(atlas.RequestEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
			}
			stop <- struct{}{}
		}
		eventSub.Unsubscribe()
	}
	go eventLoop()
	resultCh := make(chan *types.Block, 10)
	go func() {
		err := engine.Seal(chain, block, resultCh, stop)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}()

	finalBlock := <-resultCh
	if finalBlock != nil {
		t.Errorf("block mismatch: have %v, want nil", finalBlock)
	}
}

func getPublicKeys(signers []*bls.SecretKey) []*bls.PublicKey {
	publicKeys := make([]*bls.PublicKey, len(signers))
	for i := 0; i < len(signers); i++ {
		publicKeys[i] = signers[i].GetPublicKey()
	}
	return publicKeys
}

func signWithSecretKeys(signers []*bls.SecretKey, hash common.Hash) (*bls.Sign, *bls.PublicKey, *bls_cosi.Mask, error) {
	var sign bls.Sign
	var publicKey bls.PublicKey
	mask, err := bls_cosi.NewMask(getPublicKeys(signers), nil)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, signer := range signers {
		rv := signer.SignHash(hash.Bytes())
		sign.Add(rv)
		publicKey.Add(signer.GetPublicKey())
		mask.SetKey(signer.GetPublicKey(), true)
	}
	return &sign, &publicKey, mask, nil
}

func TestSealCommittedOtherHash(t *testing.T) {
	chain, engine, signerKeys := newBlockChain(4)
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	otherBlock := makeBlockWithoutSeal(chain, engine, block)

	eventSub := engine.EventMux().Subscribe(atlas.RequestEvent{})
	blockOutputChannel := make(chan *types.Block)
	stopChannel := make(chan struct{})

	go func() {
		select {
		case ev := <-eventSub.Chan():
			event, ok := ev.Data.(atlas.RequestEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
			}

			sign, publicKey, bitmap, err := signWithSecretKeys(signerKeys, SealHash(event.Proposal.(*types.Block).Header()))
			if err != nil {
				t.Errorf("failed to sign with secret keys: %v", err)
			}

			if err := engine.Commit(otherBlock, sign.Serialize(), publicKey.Serialize(), bitmap.Mask()); err != nil {
				t.Error(err.Error())
			}
		}
		eventSub.Unsubscribe()
	}()

	go func() {
		if err := engine.Seal(chain, block, blockOutputChannel, stopChannel); err != nil {
			t.Error(err.Error())
		}
	}()

	select {
	case <-blockOutputChannel:
		t.Error("Wrong block found!")
	default:
		//no block found, stop the sealing
		close(stopChannel)
	}

	select {
	case output := <-blockOutputChannel:
		if output != nil {
			t.Error("Block not nil!")
		}
	}
}

func TestSealCommitted(t *testing.T) {
	chain, engine, _ := newBlockChain(1)
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	expectedBlock, _ := engine.updateBlock(engine.chain.GetHeader(block.ParentHash(), block.NumberU64()-1), block)
	resultCh := make(chan *types.Block, 10)
	go func() {
		err := engine.Seal(chain, block, resultCh, make(chan struct{}))

		if err != nil {
			t.Errorf("error mismatch: have %v, want %v", err, expectedBlock)
		}
	}()

	finalBlock := <-resultCh
	if finalBlock.Hash() != expectedBlock.Hash() {
		t.Errorf("hash mismatch: have %v, want %v", finalBlock.Hash(), expectedBlock.Hash())
	}
}

func TestVerifyHeader(t *testing.T) {
	chain, engine, _ := newBlockChain(1)

	// errEmptyCommittedSeals case
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	block, _ = engine.updateBlock(chain.Genesis().Header(), block)
	err := engine.VerifyHeader(chain, block.Header(), false)
	if err != errEmptyCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, errEmptyCommittedSeals)
	}

	// short extra data
	header := block.Header()
	header.Extra = []byte{}
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidExtraDataFormat {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidExtraDataFormat)
	}
	// incorrect extra format
	header.Extra = []byte("0000000000000000000000000000000012300000000000000000000000000000000000000000000000000000000000000000")
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidExtraDataFormat {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidExtraDataFormat)
	}

	// non zero MixDigest
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.MixDigest = common.StringToHash("123456789")
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidMixDigest {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMixDigest)
	}

	// invalid uncles hash
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.UncleHash = common.StringToHash("123456789")
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidUncleHash {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidUncleHash)
	}

	// invalid difficulty
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Difficulty = big.NewInt(2)
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidDifficulty {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidDifficulty)
	}

	// invalid timestamp
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = chain.Genesis().Time() + engine.config.BlockPeriod - 1
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidTimestamp {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidTimestamp)
	}

	// future block
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = uint64(now().Unix()) + uint64(10)
	err = engine.VerifyHeader(chain, header, false)
	if err != consensus.ErrFutureBlock {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrFutureBlock)
	}

	// invalid nonce
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	copy(header.Nonce[:], hexutil.MustDecode("0x111111111111"))
	header.Number = big.NewInt(int64(engine.config.Epoch))
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidNonce {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidNonce)
	}
}

func TestVerifySeal(t *testing.T) {
	chain, engine, _ := newBlockChain(1)
	genesis := chain.Genesis()
	// cannot verify genesis
	err := engine.VerifySeal(chain, genesis.Header())
	if err != errUnknownBlock {
		t.Errorf("error mismatch: have %v, want %v", err, errUnknownBlock)
	}

	block := makeBlock(chain, engine, genesis)
	// change block content
	header := block.Header()
	header.Number = big.NewInt(4)
	block1 := block.WithSeal(header)
	err = engine.VerifySeal(chain, block1.Header())
	if err != errUnauthorized {
		t.Errorf("error mismatch: have %v, want %v", err, errUnauthorized)
	}

	// unauthorized users but still can get correct signer address
	err = engine.VerifySeal(chain, block.Header())
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
}

func TestVerifyHeaders(t *testing.T) {
	chain, engine, _ := newBlockChain(1)
	genesis := chain.Genesis()

	// success case
	headers := []*types.Header{}
	blocks := []*types.Block{}
	size := 100

	for i := 0; i < size; i++ {
		var b *types.Block
		if i == 0 {
			b = makeBlockWithoutSeal(chain, engine, genesis)
			b, _ = engine.updateBlock(genesis.Header(), b)
		} else {
			b = makeBlockWithoutSeal(chain, engine, blocks[i-1])
			b, _ = engine.updateBlock(blocks[i-1].Header(), b)
		}
		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}
	now = func() time.Time {
		return time.Unix(int64(headers[size-1].Time), 0)
	}
	_, results := engine.VerifyHeaders(chain, headers, nil)
	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	index := 0
OUT1:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
					break OUT1
				}
			}
			index++
			if index == size {
				break OUT1
			}
		case <-timeout.C:
			break OUT1
		}
	}
	// abort cases
	abort, results := engine.VerifyHeaders(chain, headers, nil)
	timeout = time.NewTimer(timeoutDura)
	index = 0
OUT2:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
					break OUT2
				}
			}
			index++
			if index == 5 {
				abort <- struct{}{}
			}
			if index >= size {
				t.Errorf("verifyheaders should be aborted")
				break OUT2
			}
		case <-timeout.C:
			break OUT2
		}
	}
	// error header cases
	headers[2].Number = big.NewInt(100)
	abort, results = engine.VerifyHeaders(chain, headers, nil)
	timeout = time.NewTimer(timeoutDura)
	index = 0
	errors := 0
	expectedErrors := 2
OUT3:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
					errors++
				}
			}
			index++
			if index == size {
				if errors != expectedErrors {
					t.Errorf("error mismatch: have %v, want %v", err, expectedErrors)
				}
				break OUT3
			}
		case <-timeout.C:
			break OUT3
		}
	}
}

func TestPrepareExtra(t *testing.T) {
	validators := make([]atlas.Validator, 0)

	vanity := make([]byte, types.AtlasExtraVanity)
	expectedResult := append(vanity, hexutil.MustDecode("")...)

	h := &types.Header{
		Extra: vanity,
	}

	payload, err := prepareExtra(h, validators)
	if err != nil {
		t.Errorf("error mismatch: have %v, want: nil", err)
	}
	if !reflect.DeepEqual(payload, expectedResult) {
		t.Errorf("payload mismatch: have %v, want %v", payload, expectedResult)
	}

	// append useless information to extra-data
	h.Extra = append(vanity, make([]byte, 15)...)

	payload, err = prepareExtra(h, validators)
	if !reflect.DeepEqual(payload, expectedResult) {
		t.Errorf("payload mismatch: have %v, want %v", payload, expectedResult)
	}
}

func TestWriteSeal(t *testing.T) {
	vanity := bytes.Repeat([]byte{0x00}, types.AtlasExtraVanity)
	istRawData := hexutil.MustDecode("0xf873b8600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009000000000000000000000000000000000")
	expectedSeal := bytes.Repeat([]byte{0x00}, types.AtlasExtraSignature)
	expectedProposer := 0
	expectedIstExtra := &types.AtlasExtra{
		AggSignature: [96]byte{},
		AggBitmap:    [16]byte{},
	}

	copy(expectedIstExtra.AggSignature[:], expectedSeal[:])

	var expectedErr error

	h := &types.Header{
		Extra: append(vanity, istRawData...),
	}

	// normal case
	err := writeSeal(h, expectedSeal, expectedProposer)
	if err != expectedErr {
		t.Errorf("error mismatch: have %v, want %v", err, expectedErr)
	}

	// verify atlas extra-data
	istExtra, err := types.ExtractAtlasExtra(h)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	if !reflect.DeepEqual(istExtra, expectedIstExtra) {
		t.Errorf("extra data mismatch: have %v, want %v", istExtra, expectedIstExtra)
	}

	// invalid seal
	unexpectedSeal := append(expectedSeal, make([]byte, 1)...)
	err = writeSeal(h, unexpectedSeal, expectedProposer)
	if err != errInvalidSignature {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidSignature)
	}
}

func TestWriteCommittedSeals(t *testing.T) {
	vanity := bytes.Repeat([]byte{0x00}, types.AtlasExtraVanity)
	signerKey, _ := crypto.GenerateBLSKey()
	sign := signerKey.SignHash(crypto.Keccak256Hash(signerKey.Serialize()).Bytes())

	expectedCommittedSeal := sign.Serialize()
	expectedPublicKey := signerKey.GetPublicKey().Serialize()
	expectedBitmap := make([]byte, types.AtlasExtraMask)

	data := fmt.Sprintf("0x%x%x%x%x", sign.Serialize(), expectedPublicKey, expectedBitmap, math.PaddedBigBytes(big.NewInt(0), 2))
	istRawData := hexutil.MustDecode(data)

	expectedIstExtra := &types.AtlasExtra{}
	copy(expectedIstExtra.AggSignature[:], expectedCommittedSeal)
	copy(expectedIstExtra.AggBitmap[:], expectedBitmap)

	var expectedErr error

	h := &types.Header{
		Extra: append(vanity, istRawData...),
	}

	// normal case
	err := WriteCommittedSeals(h, expectedCommittedSeal, expectedPublicKey, expectedBitmap)
	if err != expectedErr {
		t.Errorf("error mismatch: have %v, want %v", err, expectedErr)
	}

	// verify atlas extra-data
	istExtra, err := types.ExtractAtlasExtra(h)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	if !reflect.DeepEqual(istExtra, expectedIstExtra) {
		t.Errorf("extra data mismatch: have %v, want %v", istExtra, expectedIstExtra)
	}

	// invalid seal
	unexpectedCommittedSeal := append(expectedCommittedSeal, make([]byte, 1)...)
	err = WriteCommittedSeals(h, unexpectedCommittedSeal, expectedPublicKey, expectedBitmap)
	if err != errInvalidCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidCommittedSeals)
	}
}

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(log.LvlDebug), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))
}
