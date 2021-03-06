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

package core

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/log"
)

func TestCheckMessage(t *testing.T) {
	c := &core{
		logger: log.New(),
		state:  StateAcceptRequest,
		current: newRoundState(&atlas.View{
			Sequence: big.NewInt(1),
			Round:    big.NewInt(0),
		}, newTestValidatorSet(4), common.Hash{}, nil, nil, nil),
	}

	// invalid view format
	err := c.checkMessage(msgPreprepare, nil)
	if err != errInvalidMessage {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidMessage)
	}

	testStates := []State{StateAcceptRequest, StatePreprepared, StatePrepared, StateExpected, StateConfirmed, StateCommitted}
	testCode := []uint64{msgPreprepare, msgPrepare, msgExpect, msgConfirm, msgCommit, msgRoundChange}

	// future sequence
	v := &atlas.View{
		Sequence: big.NewInt(2),
		Round:    big.NewInt(0),
	}
	for i := 0; i < len(testStates); i++ {
		c.state = testStates[i]
		for j := 0; j < len(testCode); j++ {
			err := c.checkMessage(testCode[j], v)
			if err != errFutureMessage {
				t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
			}
		}
	}

	// future round
	v = &atlas.View{
		Sequence: big.NewInt(1),
		Round:    big.NewInt(1),
	}
	for i := 0; i < len(testStates); i++ {
		c.state = testStates[i]
		for j := 0; j < len(testCode); j++ {
			err := c.checkMessage(testCode[j], v)
			if testCode[j] == msgRoundChange {
				if err != nil {
					t.Errorf("error mismatch: have %v, want nil", err)
				}
			} else if err != errFutureMessage {
				t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
			}
		}
	}

	// current view but waiting for round change
	v = &atlas.View{
		Sequence: big.NewInt(1),
		Round:    big.NewInt(0),
	}
	c.waitingForRoundChange = true
	for i := 0; i < len(testStates); i++ {
		c.state = testStates[i]
		for j := 0; j < len(testCode); j++ {
			err := c.checkMessage(testCode[j], v)
			if testCode[j] == msgRoundChange {
				if err != nil {
					t.Errorf("error mismatch: have %v, want nil", err)
				}
			} else if err != errFutureMessage {
				t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
			}
		}
	}
	c.waitingForRoundChange = false

	v = c.currentView()
	// current view, state = StateAcceptRequest
	c.state = StateAcceptRequest
	for i := 0; i < len(testCode); i++ {
		err = c.checkMessage(testCode[i], v)
		if testCode[i] == msgRoundChange {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else if testCode[i] == msgPreprepare {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else {
			if err != errFutureMessage {
				t.Errorf("error mismatch: have %v, want %v", err, errFutureMessage)
			}
		}
	}

	// current view, state = StatePreprepared
	c.state = StatePreprepared
	for i := 0; i < len(testCode); i++ {
		err = c.checkMessage(testCode[i], v)
		if testCode[i] == msgRoundChange {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}

	// current view, state = StatePrepared
	c.state = StatePrepared
	for i := 0; i < len(testCode); i++ {
		err = c.checkMessage(testCode[i], v)
		if testCode[i] == msgRoundChange {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}

	// current view, state = StateExpected
	c.state = StateExpected
	for i := 0; i < len(testCode); i++ {
		err = c.checkMessage(testCode[i], v)
		if testCode[i] == msgRoundChange {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}

	// current view, state = StateConfirmed
	c.state = StateConfirmed
	for i := 0; i < len(testCode); i++ {
		err = c.checkMessage(testCode[i], v)
		if testCode[i] == msgRoundChange {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}

	// current view, state = StateCommitted
	c.state = StateCommitted
	for i := 0; i < len(testCode); i++ {
		err = c.checkMessage(testCode[i], v)
		if testCode[i] == msgRoundChange {
			if err != nil {
				t.Errorf("error mismatch: have %v, want nil", err)
			}
		} else if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}

}

func TestStoreBacklog(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	sys := NewTestSystemWithBackend(N, F)

	// for i, backend := range sys.backends {
	backend := sys.backends[0]
	c := backend.engine.(*core)
	v := &atlas.View{
		Round:    big.NewInt(10),
		Sequence: big.NewInt(10),
	}
	p := c.valSet.GetByIndex(0)
	// push preprepare msg
	preprepare := &atlas.Preprepare{
		View:     v,
		Proposal: makeBlock(1),
	}
	prepreparePayload, _ := Encode(preprepare)
	m := &message{
		Code: msgPreprepare,
		Msg:  prepreparePayload,
	}
	c.storeBacklog(m, p)
	if _, exists := c.backlogs[p.Signer()]; exists {
		t.Errorf("backlogs[signer] should not be existed right now.")
	}

	p = c.valSet.GetByIndex(1) // WARNING: use index(1) different backend(0) to void Backlog from self.
	c.storeBacklog(m, p)
	msg := c.backlogs[p.Signer()].PopItem()
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push prepare msg
	subject := &atlas.Subject{
		View:   v,
		Digest: common.StringToHash("1234567890"),
	}
	signedSubject, err := c.SignSubject(subject)
	if err != nil {
		t.Errorf("failed to sign subject: %v", err)
	}
	subjectPayload, _ := Encode(signedSubject)

	m = &message{
		Code: msgPrepare,
		Msg:  subjectPayload,
	}
	c.storeBacklog(m, p)
	msg = c.backlogs[p.Signer()].PopItem()
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	m = &message{
		Code: msgExpect,
		Msg:  subjectPayload,
	}
	c.storeBacklog(m, p)
	msg = c.backlogs[p.Signer()].PopItem()
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	m = &message{
		Code: msgConfirm,
		Msg:  subjectPayload,
	}
	c.storeBacklog(m, p)
	msg = c.backlogs[p.Signer()].PopItem()
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push commit msg
	m = &message{
		Code: msgCommit,
		Msg:  subjectPayload,
	}
	c.storeBacklog(m, p)
	msg = c.backlogs[p.Signer()].PopItem()
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}

	// push roundChange msg
	m = &message{
		Code: msgRoundChange,
		Msg:  subjectPayload,
	}
	c.storeBacklog(m, p)
	msg = c.backlogs[p.Signer()].PopItem()
	if !reflect.DeepEqual(msg, m) {
		t.Errorf("message mismatch: have %v, want %v", msg, m)
	}
}

func TestProcessFutureBacklog(t *testing.T) {
	N := uint64(4)
	F := uint64(1)
	sys := NewTestSystemWithBackend(N, F)

	backend := sys.backends[0]
	c := backend.engine.(*core)
	c.subscribeEvents()
	defer c.unsubscribeEvents()

	v := &atlas.View{
		Round:    big.NewInt(10),
		Sequence: big.NewInt(10),
	}
	p := c.valSet.GetByIndex(1) // WARNING: use index(1) different backend(0) to void Backlog from self.
	// push a future msg
	subject := &atlas.Subject{
		View:   v,
		Digest: common.StringToHash("1234567890"),
	}
	signedSubject, err := c.SignSubject(subject)
	if err != nil {
		t.Errorf("failed to sign subject")
	}
	subjectPayload, _ := Encode(signedSubject)
	m := &message{
		Code: msgCommit,
		Msg:  subjectPayload,
	}
	c.storeBacklog(m, p)
	c.processBacklog()

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	select {
	case e, ok := <-c.events.Chan():
		if !ok {
			return
		}
		t.Errorf("unexpected events comes: %v", e)
	case <-timeout.C:
		// success
	}
}

func TestProcessBacklog(t *testing.T) {
	v := &atlas.View{
		Round:    big.NewInt(0),
		Sequence: big.NewInt(1),
	}
	preprepare := &atlas.Preprepare{
		View:     v,
		Proposal: makeBlock(1),
	}
	prepreparePayload, _ := Encode(preprepare)

	subject := &atlas.Subject{
		View:   v,
		Digest: common.StringToHash("1234567890"),
	}
	subjectPayload, _ := Encode(subject)

	msgs := []*message{
		&message{
			Code: msgPreprepare,
			Msg:  prepreparePayload,
		},
		&message{
			Code: msgPrepare,
			Msg:  subjectPayload,
		},
		&message{
			Code: msgExpect,
			Msg:  subjectPayload,
		},
		&message{
			Code: msgConfirm,
			Msg:  subjectPayload,
		},
		&message{
			Code: msgCommit,
			Msg:  subjectPayload,
		},
		&message{
			Code: msgRoundChange,
			Msg:  subjectPayload,
		},
	}
	for i := 0; i < len(msgs); i++ {
		testProcessBacklog(t, msgs[i])
	}
}

func testProcessBacklog(t *testing.T, msg *message) {
	N := uint64(4)
	F := uint64(1)
	sys := NewTestSystemWithBackend(N, F)

	backend := sys.backends[0]
	c := backend.engine.(*core)
	c.subscribeEvents()
	defer c.unsubscribeEvents()

	c.state = State(msg.Code)

	p := c.valSet.GetByIndex(1) // WARNING: use index(1) different backend(0) to void Backlog from self.
	c.storeBacklog(msg, p)
	c.processBacklog()

	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	select {
	case ev := <-c.events.Chan():
		e, ok := ev.Data.(backlogEvent)
		if !ok {
			t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
		}
		if e.msg.Code != msg.Code {
			t.Errorf("message code mismatch: have %v, want %v", e.msg.Code, msg.Code)
		}
		// success
	case <-timeout.C:
		t.Error("unexpected timeout occurs")
	}
}
