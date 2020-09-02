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

import "errors"

var (
	// errInconsistentSubject is returned when received subject is different from
	// current subject.
	errInconsistentSubject = errors.New("inconsistent subjects")
	// errNotFromProposer is returned when received message is supposed to be from
	// proposer.
	errNotFromProposer = errors.New("message does not come from proposer")
	// errNotFromCommitee is returned when received message is supposed to be from
	// committee.
	errNotFromCommittee = errors.New("message does not come from committee")
	// errIgnored is returned when a message was ignored.
	errIgnored = errors.New("message is ignored")
	// errFutureMessage is returned when current view is earlier than the
	// view of the received message.
	errFutureMessage = errors.New("future message")
	// errOldMessage is returned when the received message's view is earlier
	// than current view.
	errOldMessage = errors.New("old message")
	// errInvalidMessage is returned when the message is malformed.
	errInvalidMessage = errors.New("invalid message")
	// errInvalidSignature is returned when the signature is invalid.
	errInvalidSignature = errors.New("invalid signature")
	// errFailedDecodePreprepare is returned when the PRE-PREPARE message is malformed.
	errFailedDecodePreprepare = errors.New("failed to decode PRE-PREPARE")
	// errFailedDecodePrepare is returned when the PREPARE message is malformed.
	errFailedDecodePrepare = errors.New("failed to decode PREPARE")
	// errFailedDecodeExpect is returned when the EXPECT message is malformed.
	errFailedDecodeExpect = errors.New("failed to decode EXPECT")
	// errFailedDecodeConfirm is returned when the CONFIRM message is malformed.
	errFailedDecodeConfirm = errors.New("failed to decode CONFIRM")
	// errFailedDecodeCommit is returned when the COMMIT message is malformed.
	errFailedDecodeCommit = errors.New("failed to decode COMMIT")
	// errFailedDecodeSignPayload is returned when the subject.Payload is malformed.
	errFailedDecodeSignPayload = errors.New("failed to decode subject.SignPayload")
	// errFailedDecodeMessageSet is returned when the message set is malformed.
	errFailedDecodeMessageSet = errors.New("failed to decode message set")
	// errFailedSignData is returned when failed to sign data.
	errFailedSignData = errors.New("failed to sign data")
	// errInvalidSigner is returned when the message is signed by a validator different than message sender
	errInvalidSigner = errors.New("message not signed by the sender")
	// errNotSatisfyQuorum is returned when the message is signed without enough validators
	errNotSatisfyQuorum = errors.New("message not signed by enough validators")
	// errInvalidState is returned when state is invalid
	errInvalidState = errors.New("invalid state")
	// errDuplicateMessage is returned when signature have been aggregated
	errDuplicateMessage = errors.New("duplicate message")
	// errInproperState is returned when handle message in inproper state
	errInproperState = errors.New("inproper state")
)