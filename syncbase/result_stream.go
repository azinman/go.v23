// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syncbase

import (
	"sync"

	"v.io/v23/context"
	wire "v.io/v23/services/syncbase"
	"v.io/v23/verror"
	"v.io/v23/vom"
)

type resultStream struct {
	mu sync.Mutex
	// cancel cancels the RPC resultStream.
	cancel context.CancelFunc
	// call is the RPC resultStream object.
	call wire.DatabaseExecClientCall
	// curr is the currently staged result, or nil if nothing is staged.
	curr []*vom.RawBytes
	// err is the first error encountered during streaming. It may also be
	// populated by a call to Cancel.
	err error
	// finished records whether we have called call.Finish().
	finished bool
}

var _ ResultStream = (*resultStream)(nil)

func newResultStream(cancel context.CancelFunc, call wire.DatabaseExecClientCall) *resultStream {
	return &resultStream{
		cancel: cancel,
		call:   call,
	}
}

func (rs *resultStream) Advance() bool {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if rs.err != nil || rs.finished {
		return false
	}
	if !rs.call.RecvStream().Advance() {
		if rs.call.RecvStream().Err() != nil {
			rs.err = rs.call.RecvStream().Err()
		} else {
			rs.err = rs.call.Finish()
			rs.cancel() // TODO(jkline): Copied from stream.go, is this needed?
			rs.finished = true
		}
		return false
	}
	rs.curr = rs.call.RecvStream().Value()
	return true
}

func (rs *resultStream) ResultCount() int {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return len(rs.curr)
}

func (rs *resultStream) Result(i int, value interface{}) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.curr[i].ToValue(value)
}

func (rs *resultStream) Err() error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.err
}

// TODO(jkline): Make Cancel non-blocking (TODO copied from stream.go)
func (rs *resultStream) Cancel() {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.cancel()
	rs.err = verror.New(verror.ErrCanceled, nil)
}