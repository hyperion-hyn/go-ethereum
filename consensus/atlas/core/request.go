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

import "github.com/ethereum/go-ethereum/consensus/atlas"

func (c *core) handleRequest(request *atlas.Request) error {
	logger := c.logger.New("state", c.state, "seq", c.current.sequence)

	if err := c.checkRequestMsg(request); err != nil {
		if err == errInvalidMessage {
			logger.Warn("invalid request")
			return err
		}
		logger.Warn("unexpected request", "err", err, "number", request.Proposal.Number(), "hash", request.Proposal.SealHash(c.backend))
		return err
	}

	logger.Trace("handleRequest", "number", request.Proposal.Number(), "hash", request.Proposal.SealHash(c.backend))

	// ATLAS(zgx): what if c.state != StateAcceptRequest, should I set pendingRequest? should I put this line in if block and return errFutureMessage?
	c.current.pendingRequest = request
	if c.state == StateAcceptRequest {
		c.sendPreprepare(request)
	}
	return nil
}

// check request state
// return errInvalidMessage if the message is invalid
// return errFutureMessage if the sequence of proposal is larger than current sequence
// return errOldMessage if the sequence of proposal is smaller than current sequence
func (c *core) checkRequestMsg(request *atlas.Request) error {
	if request == nil || request.Proposal == nil {
		return errInvalidMessage
	}

	if c := c.current.sequence.Cmp(request.Proposal.Number()); c > 0 {
		return errOldMessage
	} else if c < 0 {
		return errFutureMessage
	} else {
		return nil
	}
}

func (c *core) storeRequestMsg(request *atlas.Request) {
	logger := c.logger.New("state", c.state)

	logger.Trace("Store future request", "number", request.Proposal.Number(), "hash", request.Proposal.SealHash(c.backend))

	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	c.pendingRequests.Push(request, float32(-request.Proposal.Number().Int64()))
}

func (c *core) processPendingRequests() {
	c.pendingRequestsMu.Lock()
	defer c.pendingRequestsMu.Unlock()

	for !(c.pendingRequests.Empty()) {
		m, prio := c.pendingRequests.Pop()
		r, ok := m.(*atlas.Request)
		if !ok {
			c.logger.Warn("Malformed request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		err := c.checkRequestMsg(r)
		if err != nil {
			if err == errFutureMessage {
				c.logger.Trace("Stop processing request", "number", r.Proposal.Number(), "hash", r.Proposal.SealHash(c.backend))
				c.pendingRequests.Push(m, prio)
				break
			}
			c.logger.Trace("Skip the pending request", "number", r.Proposal.Number(), "hash", r.Proposal.SealHash(c.backend), "err", err)
			continue
		}
		c.logger.Trace("Post pending request", "number", r.Proposal.Number(), "hash", r.Proposal.SealHash(c.backend))

		go c.sendEvent(atlas.RequestEvent{
			Proposal: r.Proposal,
		})
	}
}
