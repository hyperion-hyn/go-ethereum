package availability

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/staking/effective"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
	"reflect"
	"testing"
)

func TestBlockSigners(t *testing.T) {
	tests := []struct {
		numSlots               int
		verified               []int
		numPayable, numMissing int
	}{
		{0, []int{}, 0, 0},
		{1, []int{}, 0, 1},
		{1, []int{0}, 1, 0},
		{8, []int{}, 0, 8},
		{8, []int{0}, 1, 7},
		{8, []int{7}, 1, 7},
		{8, []int{1, 3, 5, 7}, 4, 4},
		{8, []int{0, 2, 4, 6}, 4, 4},
		{8, []int{0, 1, 2, 3, 4, 5, 6, 7}, 8, 0},
		{13, []int{0, 1, 4, 5, 6, 9, 12}, 7, 6},
		// TODO: add a real data test case given numSlots of a committee and
		//  number of payable of a certain block
	}
	for i, test := range tests {
		cmt := makeTestCommittee(test.numSlots, 0)
		bm, err := indexesToBitMap(test.verified, test.numSlots)
		if err != nil {
			t.Fatalf("test %d: %v", i, err)
		}
		pSlots, mSlots, err := BlockSigners(bm, cmt)
		if err != nil {
			t.Fatalf("test %d: %v", i, err)
		}
		if len(pSlots.Entrys) != test.numPayable || len(mSlots.Entrys) != test.numMissing {
			t.Errorf("test %d: unexpected result: # pSlots %d/%d, # mSlots %d/%d",
				i, len(pSlots.Entrys), test.numPayable, len(mSlots.Entrys), test.numMissing)
			continue
		}
		if err := checkPayableAndMissing(cmt, test.verified, pSlots, mSlots); err != nil {
			t.Errorf("test %d: %v", i, err)
		}
	}
}

func checkPayableAndMissing(cmt *restaking.Committee_, idxs []int, pSlots, mSlots *restaking.Slots_) error {
	if len(pSlots.Entrys)+len(mSlots.Entrys) != len(cmt.Slots.Entrys) {
		return fmt.Errorf("slots number not expected: %d(payable) + %d(missings) != %d(committee)",
			len(pSlots.Entrys), len(mSlots.Entrys), len(cmt.Slots.Entrys))
	}
	pIndex, mIndex, iIndex := 0, 0, 0
	for i, slot := range cmt.Slots.Entrys {
		if iIndex >= len(idxs) || i != idxs[iIndex] {
			// the slot should be missings and we shall check mSlots[mIndex] == slot
			if mIndex >= len(mSlots.Entrys) || !reflect.DeepEqual(slot, mSlots.Entrys[mIndex]) {
				return fmt.Errorf("addr %v missed from missings slots", slot.EcdsaAddress.String())
			}
			mIndex++
		} else {
			// check pSlots[pIndex] == slot
			if pIndex >= len(pSlots.Entrys) || !reflect.DeepEqual(slot, pSlots.Entrys[pIndex]) {
				return fmt.Errorf("addr %v missed from payable slots", slot.EcdsaAddress.String())
			}
			pIndex++
			iIndex++
		}
	}
	return nil
}

func TestBlockSigners_BitmapOverflow(t *testing.T) {
	tests := []struct {
		numSlots  int
		numBitmap int
		err       error
	}{
		{16, 16, nil},
		{16, 14, nil},
		{16, 8, errors.New("bitmap size too small")},
		{16, 24, errors.New("bitmap size too large")},
	}
	for i, test := range tests {
		cmt := makeTestCommittee(test.numSlots, 0)
		bm, _ := indexesToBitMap([]int{}, test.numBitmap)
		_, _, err := BlockSigners(bm, cmt)
		if (err == nil) != (test.err == nil) {
			t.Errorf("Test %d: BlockSigners got err [%v], expect [%v]", i, err, test.err)
		}
	}
}

func TestBallotResult(t *testing.T) {
	tests := []struct {
		numSlots int
		parVerified, chdVerified int
		parBN, chdBN             int64
		expNumPayable, expNumMissing int
		expErr                       error
	}{
		{1, 1, 1, 10, 11, 1, 0, nil},
		{16, 10, 12,  100, 101, 12, 4, nil},
		{16, 10, 12, 100, 101, 12, 4, errors.New("cannot find shard")},
	}
	for i, test := range tests {
		comm := makeTestCommittee(test.numSlots, 3)
		chdHeader := newTestHeader(test.chdBN, test.numSlots, test.chdVerified)

		slots, payable, missing, err := BallotResult(chdHeader, comm)
		if err != nil {
			if test.expErr == nil {
				t.Errorf("Test %v: unexpected error: %v", i, err)
			}
			continue
		}
		if !reflect.DeepEqual(slots, comm.Slots) {
			t.Errorf("Test %v: Ballot result slots not expected", i)
		}
		if len(payable.Entrys) != test.expNumPayable {
			t.Errorf("Test %v: payable size not expected: %v / %v", i, len(payable.Entrys), test.expNumPayable)
		}
		if len(missing.Entrys) != test.expNumMissing {
			t.Errorf("Test %v: missings size not expected: %v / %v", i, len(missing.Entrys), test.expNumMissing)
		}
	}
}

func TestIncrementValidatorSigningCounts(t *testing.T) {
	tests := []struct {
		numUserSlots int
		verified                  []int
	}{
		{0, []int{0}},
		{1, []int{0}},
		{6, []int{0, 2, 3, 4, 6, 8, 10, 12, 14}},
		{6, []int{1, 3, 5, 7, 9, 11, 13, 15}},
	}
	for _, test := range tests {
		ctx, err := makeIncStateTestCtx(test.numUserSlots, test.verified)
		if err != nil {
			t.Fatal(err)
		}
		if err := IncrementValidatorSigningCounts(nil, ctx.staked, ctx.state, ctx.signers,
			ctx.missings); err != nil {

			t.Fatal(err)
		}
		if err := ctx.checkResult(); err != nil {
			t.Error(err)
		}
	}
}

func TestComputeCurrentSigning(t *testing.T) {
	tests := []struct {
		snapSigned, curSigned, diffSigned int64
		snapToSign, curToSign, diffToSign int64
		pctNum, pctDiv                    int64
		isBelowThreshold                  bool
	}{
		{0, 0, 0, 0, 0, 0, 0, 1, true},
		{0, 1, 1, 0, 1, 1, 1, 1, false},
		{0, 2, 2, 0, 3, 3, 2, 3, true},
		{0, 1, 1, 0, 3, 3, 1, 3, true},
		{100, 225, 125, 200, 350, 150, 5, 6, false},
		{100, 200, 100, 200, 350, 150, 2, 3, true},
		{100, 200, 100, 200, 400, 200, 1, 2, true},
	}
	for i, test := range tests {
		snapWrapper := makeTestStorageWrapper(common.Address{}, test.snapSigned, test.snapToSign)
		curWrapper := makeTestStorageWrapper(common.Address{}, test.curSigned, test.curToSign)

		computed := ComputeCurrentSigning(snapWrapper, curWrapper)

		if computed.Signed.Cmp(new(big.Int).SetInt64(test.diffSigned)) != 0 {
			t.Errorf("test %v: computed signed not expected: %v / %v",
				i, computed.Signed, test.diffSigned)
		}
		if computed.ToSign.Cmp(new(big.Int).SetInt64(test.diffToSign)) != 0 {
			t.Errorf("test %v: computed to sign not expected: %v / %v",
				i, computed.ToSign, test.diffToSign)
		}
		expPct := common.NewDec(test.pctNum).Quo(common.NewDec(test.pctDiv))
		if !computed.Percentage.Equal(expPct) {
			t.Errorf("test %v: computed percentage not expected: %v / %v",
				i, computed.Percentage, expPct)
		}
		if computed.IsBelowThreshold != test.isBelowThreshold {
			t.Errorf("test %v: computed is below threshold not expected: %v / %v",
				i, computed.IsBelowThreshold, test.isBelowThreshold)
		}
	}
}

func TestComputeAndMutateEPOSStatus(t *testing.T) {
	tests := []struct {
		ctx       *computeEPOSTestCtx
		expErr    error
		expStatus effective.Eligibility
	}{
		// active node
		{
			ctx: &computeEPOSTestCtx{
				addr:       common.Address{20, 20},
				snapSigned: 100,
				snapToSign: 100,
				snapEli:    effective.Active,
				curSigned:  200,
				curToSign:  200,
				curEli:     effective.Active,
			},
			expStatus: effective.Active,
		},
		// active -> inactive
		{
			ctx: &computeEPOSTestCtx{
				addr:       common.Address{20, 20},
				snapSigned: 100,
				snapToSign: 100,
				snapEli:    effective.Active,
				curSigned:  200,
				curToSign:  250,
				curEli:     effective.Active,
			},
			expStatus: effective.Inactive,
		},
		// active -> inactive
		{
			ctx: &computeEPOSTestCtx{
				addr:       common.Address{20, 20},
				snapSigned: 100,
				snapToSign: 100,
				snapEli:    effective.Active,
				curSigned:  100,
				curToSign:  200,
				curEli:     effective.Active,
			},
			expStatus: effective.Inactive,
		},
		// status unchanged: inactive -> inactive
		{
			ctx: &computeEPOSTestCtx{
				addr:       common.Address{20, 20},
				snapSigned: 100,
				snapToSign: 100,
				snapEli:    effective.Inactive,
				curSigned:  200,
				curToSign:  200,
				curEli:     effective.Inactive,
			},
			expStatus: effective.Inactive,
		},
		// status unchanged: inactive
		{
			ctx: &computeEPOSTestCtx{
				addr:       common.Address{20, 20},
				snapSigned: 100,
				snapToSign: 100,
				snapEli:    effective.Inactive,
				curSigned:  200,
				curToSign:  200,
				curEli:     effective.Active,
			},
			expStatus: effective.Active,
		},
		// nil validator wrapper in state
		{
			ctx: &computeEPOSTestCtx{
				addr:       common.Address{20, 20},
				snapSigned: 100,
				snapToSign: 100,
				snapEli:    effective.Active,
				curEli:     effective.Nil,
			},
			expErr: errors.New("nil validator wrapper in state"),
		},
		// nil validator wrapper in snapshot
		{
			ctx: &computeEPOSTestCtx{
				addr:      common.Address{20, 20},
				snapEli:   effective.Nil,
				curSigned: 200,
				curToSign: 200,
				curEli:    effective.Active,
			},
			expErr: errors.New("nil validator wrapper in snapshot"),
		},
		// banned node
		{
			ctx: &computeEPOSTestCtx{
				addr:       common.Address{20, 20},
				snapSigned: 100,
				snapToSign: 200,
				snapEli:    effective.Active,
				curSigned:  100,
				curToSign:  200,
				curEli:     effective.Banned,
			},
			expStatus: effective.Banned,
		},
	}
	for i, test := range tests {
		ctx := test.ctx
		ctx.makeStateAndReader()

		err := ComputeAndMutateEPOSStatus(ctx.reader, ctx.stateDB, ctx.addr, ctx.epoch)
		if err != nil {
			if test.expErr == nil {
				t.Errorf("Test %v: unexpected error: %v", i, err)
			}
			continue
		}

		if err := ctx.checkWrapperStatus(test.expStatus); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

// incStateTestCtx is the helper structure for test case TestIncrementValidatorSigningCounts
type incStateTestCtx struct {
	// Initialized fields
	snapState, state  *state.StateDB
	cmt               *restaking.Committee_
	staked            *restaking.StakedSlots
	signers, missings *restaking.Slots_

	// computedSlotMap is parsed map for result checking, which maps from Ecdsa address
	// to the expected behaviour of the address.
	//  typeIncSigned - 0: increase both toSign and signed
	//  typeIncMissing - 1: increase to sign
	computedSlotMap map[common.Address]int
}

const (
	typeIncSigned = iota
	typeIncMissing
)

// makeIncStateTestCtx create and initialize the test context for TestIncrementValidatorSigningCounts
func makeIncStateTestCtx(numSlots int, verified []int) (*incStateTestCtx, error) {
	cmt := makeTestCommittee(numSlots, 0)
	staked := cmt.StakedValidators()
	bitmap, _ := indexesToBitMap(verified, numSlots)
	signers, missing, err := BlockSigners(bitmap, cmt)
	if err != nil {
		return nil, err
	}
	stateDB := newTestStateDBFromCommittee(cmt)
	snapState := stateDB.Copy()

	return &incStateTestCtx{
		snapState: snapState,
		state:     stateDB,
		cmt:       cmt,
		staked:    staked,
		signers:   signers,
		missings:  missing,
	}, nil
}

// checkResult checks the state change result for incStateTestCtx
func (ctx *incStateTestCtx) checkResult() error {
	ctx.computeSlotMaps()

	for addr, typeInc := range ctx.computedSlotMap {
		if err := ctx.checkAddrIncStateByType(addr, typeInc); err != nil {
			return err
		}
	}
	return nil
}

// computeSlotMaps compute for computedSlotMap for incStateTestCtx
func (ctx *incStateTestCtx) computeSlotMaps() {
	ctx.computedSlotMap = make(map[common.Address]int)

	for _, signer := range ctx.signers.Entrys {
		ctx.computedSlotMap[signer.EcdsaAddress] = typeIncSigned
	}
	for _, missing := range ctx.missings.Entrys {
		ctx.computedSlotMap[missing.EcdsaAddress] = typeIncMissing
	}
}

// checkAddrIncStateByType checks whether the state behaviour of a given address follows
// the expected state change rule given typeInc
func (ctx *incStateTestCtx) checkAddrIncStateByType(addr common.Address, typeInc int) error {
	var err error
	switch typeInc {
	case typeIncSigned:
		if err = ctx.checkWrapperChangeByAddr(addr, checkIncWrapperVerified); err != nil {
			err = fmt.Errorf("verified address %s: %v", addr, err)
		}
	case typeIncMissing:
		if err = ctx.checkWrapperChangeByAddr(addr, checkIncWrapperMissing); err != nil {
			err = fmt.Errorf("missing address %s: %v", addr, err)
		}
	default:
		err = errors.New("unknown typeInc")
	}
	return err
}

// checkWrapperChangeByAddr checks whether the wrapper of a given address
// before and after the state change is expected defined by compare function f.
func (ctx *incStateTestCtx) checkWrapperChangeByAddr(addr common.Address,
	f func(w1, w2 *restaking.Storage_ValidatorWrapper_) bool) error {

	snapWrapper, err := ctx.snapState.ValidatorByAddress(addr)
	if err != nil {
		return err
	}
	curWrapper, err := ctx.state.ValidatorByAddress(addr)
	if err != nil {
		return err
	}
	if isExpected := f(snapWrapper, curWrapper); !isExpected {
		return errors.New("validatorWrapper not expected")
	}
	return nil
}

// checkIncWrapperVerified is the compare function to check whether validator wrapper
// is expected for nodes who has verified a block.
func checkIncWrapperVerified(snapWrapper, curWrapper *restaking.Storage_ValidatorWrapper_) bool {
	snapSigned := snapWrapper.Counters().NumBlocksSigned().Value()
	curSigned := curWrapper.Counters().NumBlocksSigned().Value()
	if curSigned.Cmp(new(big.Int).Add(snapSigned, common.Big1)) != 0 {
		return false
	}
	snapToSign := snapWrapper.Counters().NumBlocksToSign().Value()
	curToSign := curWrapper.Counters().NumBlocksToSign().Value()
	return curToSign.Cmp(new(big.Int).Add(snapToSign, common.Big1)) == 0
}

// checkIncWrapperMissing is the compare function to check whether validator wrapper
// is expected for nodes who has missed a block.
func checkIncWrapperMissing(snapWrapper, curWrapper *restaking.Storage_ValidatorWrapper_) bool {
	snapSigned := snapWrapper.Counters().NumBlocksSigned().Value()
	curSigned := curWrapper.Counters().NumBlocksSigned().Value()
	if curSigned.Cmp(snapSigned) != 0 {
		return false
	}
	snapToSign := snapWrapper.Counters().NumBlocksToSign().Value()
	curToSign := curWrapper.Counters().NumBlocksToSign().Value()
	return curToSign.Cmp(new(big.Int).Add(snapToSign, common.Big1)) == 0
}

type computeEPOSTestCtx struct {
	// input arguments
	addr                   common.Address
	epoch                  *big.Int
	snapSigned, snapToSign int64
	snapEli                effective.Eligibility
	curSigned, curToSign   int64
	curEli                 effective.Eligibility

	// computed fields
	stateDB *state.StateDB
	reader  testReader
}

// makeStateAndReader compute for state and reader given the input arguments
func (ctx *computeEPOSTestCtx) makeStateAndReader() {
	ctx.reader = newTestReader()
	if ctx.snapEli != effective.Nil {
		wrapper := makeTestWrapper(ctx.addr, ctx.snapSigned, ctx.snapToSign)
		wrapper.Validator.Status = big.NewInt(int64(ctx.curEli))
		ctx.reader[ctx.epoch.Uint64()].ValidatorPool().Validators().Put(ctx.addr, &wrapper)
	}
	ctx.stateDB = newTestStateDB()
	if ctx.curEli != effective.Nil {
		wrapper := makeTestWrapper(ctx.addr, ctx.curSigned, ctx.curToSign)
		wrapper.Validator.Status = big.NewInt(int64(ctx.curEli))
		ctx.stateDB.ValidatorPool().Validators().Put(ctx.addr, &wrapper)
	}
}

func (ctx *computeEPOSTestCtx) checkWrapperStatus(expStatus effective.Eligibility) error {
	wrapper, err := ctx.stateDB.ValidatorByAddress(ctx.addr)
	if err != nil {
		return err
	}
	status := wrapper.Validator().Status().Value().Uint64()
	if status != uint64(expStatus) {
		return fmt.Errorf("wrapper status unexpected: %v / %v", status, expStatus)
	}
	return nil
}

// testHeader is the fake Header for testing
type testHeader struct {
	number           *big.Int
	lastCommitBitmap []byte
}

func newTestHeader(number int64, numSlots, numVerified int) *testHeader {
	indexes := make([]int, 0, numVerified)
	for i := 0; i != numVerified; i++ {
		indexes = append(indexes, i)
	}
	bitmap, _ := indexesToBitMap(indexes, numSlots)
	return &testHeader{
		number:           new(big.Int).SetInt64(number),
		lastCommitBitmap: bitmap,
	}
}

func (th *testHeader) Number() *big.Int {
	return th.number
}

func (th *testHeader) LastCommitBitmap() []byte {
	return th.lastCommitBitmap
}


// testReader is the fake Reader for testing
type testReader map[uint64]state.StateDB

func (reader testReader) ReadValidatorAtEpoch(epoch *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	stateDB := reader[epoch.Uint64()]
	return stateDB.ValidatorByAddress(validatorAddress)
}

// newTestReader creates an empty test reader
func newTestReader() testReader {
	reader := make(testReader)
	return reader
}

func makeTestCommittee(numSlots int, epoch int64) *restaking.Committee_ {
	slots := make([]*restaking.Slot_, 0, numSlots)
	for i := 0; i != numSlots; i++ {
		slots = append(slots, makeSlot(i, epoch))
	}
	return &restaking.Committee_{
		Epoch: big.NewInt(epoch),
		Slots: restaking.Slots_{Entrys: slots},
	}
}

const testStake = int64(100000000000)

func makeSlot(seed int, epoch int64) *restaking.Slot_ {
	addr := common.BigToAddress(new(big.Int).SetInt64(int64(seed) + epoch*1000000))
	var blsKey restaking.BLSPublicKey_
	copy(blsKey.Key[:], bls.RandPrivateKey().GetPublicKey().Serialize())
	slot := restaking.Slot_{
		EcdsaAddress: addr,
		BLSPublicKey: blsKey,
		EffectiveStake: common.NewDec(testStake),
	}
	return &slot
}

// indexesToBitMap convert the indexes to bitmap. The conversion follows the little-
// endian order.
func indexesToBitMap(idxs []int, n int) ([]byte, error) {
	bSize := (n + 7) >> 3
	res := make([]byte, bSize)
	for _, idx := range idxs {
		byt := idx >> 3
		if byt >= bSize {
			return nil, fmt.Errorf("overflow index when converting to bitmap: %v/%v", byt, bSize)
		}
		msk := byte(1) << uint(idx&7)
		res[byt] ^= msk
	}
	return res, nil
}

func makeTestWrapper(addr common.Address, numSigned, numToSign int64) restaking.ValidatorWrapper_ {
	var val restaking.ValidatorWrapper_
	val.Validator.ValidatorAddress = addr
	val.Counters.NumBlocksToSign = new(big.Int).SetInt64(numToSign)
	val.Counters.NumBlocksSigned = new(big.Int).SetInt64(numSigned)
	return val
}

func makeTestStorageWrapper(addr common.Address, numSigned, numToSign int64) *restaking.Storage_ValidatorWrapper_ {
	wrapper := makeTestWrapper(addr, numSigned, numToSign)
	validators := newTestStateDB().ValidatorPool().Validators()
	validators.Put(addr, &wrapper)
	wrapperSt, _ := validators.Get(addr)
	return wrapperSt
}

// newTestStateDB return an empty test StateDB
func newTestStateDB() *state.StateDB {
	db := rawdb.NewMemoryDatabase()
	sdb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	return sdb
}

// newTestStateDBFromCommittee creates a testStateDB given a shard committee.
// The validator wrappers are only set for user nodes.
func newTestStateDBFromCommittee(cmt *restaking.Committee_) *state.StateDB {
	sdb := newTestStateDB()
	for _, slot := range cmt.Slots.Entrys {
		wrapper := makeTestWrapper(slot.EcdsaAddress, 1, 1)
		wrapper.Validator.SlotPubKeys = restaking.BLSPublicKeys_{Keys: []*restaking.BLSPublicKey_{&slot.BLSPublicKey}}
		sdb.ValidatorPool().Validators().Put(slot.EcdsaAddress, &wrapper)
	}
	return sdb
}
