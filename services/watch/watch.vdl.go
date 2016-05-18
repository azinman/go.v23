// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: watch

// Package watch defines interfaces for watching a sequence of changes.
//
// API Overview
//
// Watcher service allows a client to watch for updates that match a
// query.  For each watched query, the client will receive a reliable
// stream of watch events without re-ordering.
//
// The watching is done by starting a streaming RPC. The argument to
// the RPC contains the query. The result stream consists of a
// never-ending sequence of Change messages (until the call fails or
// is cancelled).
//
// Root Entity
//
// The Object name that receives the Watch RPC is called the root
// entity.  The root entity is the parent of all entities that the
// client cares about.  Therefore, the query is confined to children
// of the root entity, and the names in the Change messages are all
// relative to the root entity.
//
// Watch Request
//
// When a client makes a watch request, it can indicate whether it
// wants to receive the initial states of the entities that match the
// query, just new changes to the entities, or resume watching from a
// particular point in a previous watch stream.  On receiving a watch
// request, the server sends one or more messages to the client. The
// first message informs the client that the server has registered the
// client's request; the instant of time when the client receives the
// event is referred to as the client's "watch point" for that query.
//
// Atomic Delivery
//
// The response stream consists of a sequence of Change messages. Each
// Change message contains an optional continued bit
// (default=false). A sub-sequence of Change messages with
// continued=true followed by a Change message with continued=false
// forms an "atomic group". Systems that support multi-entity atomic
// updates may guarantee that all changes resulting from a single
// atomic update are delivered in the same "atomic group". It is up to
// the documentation of a particular system that implements the Watch
// API to document whether or not it supports such grouping. We expect
// that most callers will ignore the notion of atomic delivery and the
// continued bit, i.e., they will just process each Change message as
// it is received.
//
// Initial State
//
// The first atomic group delivered by a watch call is special. It is
// delivered as soon as possible and contains the initial state of the
// entities being watched.  The client should consider itself caught up
// after processing this first atomic group.  The messages in this first
// atomic group depend on the value of ResumeMarker.
//
//   (1) ResumeMarker is "" or not specified: For every entity P that
//       matches the query and exists, there will be at least one message
//       delivered with entity == P and the last such message will contain
//       the current state of P.  For every entity Q (including the entity
//       itself) that matches the query but does not exist, either no
//       message will be delivered, or the last message for Q will have
//       state == DOES_NOT_EXIST. At least one message for entity="" will
//       be delivered.
//
//   (2) ResumeMarker == "now": there will be exactly one message with
//       entity = "" and state INITIAL_STATE_SKIPPED.  The client cannot
//       assume whether or not the entity exists after receiving this
//       message.
//
//   (3) ResumeMarker has a value R from a preceding watch call on this
//       entity: The same messages as described in (1) will be delivered
//       to the client except that any information implied by messages
//       received on the preceding call up to and including R may not be
//       delivered. The expectation is that the client will start with
//       state it had built up from the preceding watch call, apply the
//       changes received from this call and build an up-to-date view of
//       the entities without having to fetch a potentially large amount
//       of information that has not changed.  Note that some information
//       that had already been delivered by the preceding call might be
//       delivered again.
//
// Ordering and Reliability
//
// The Change messages that apply to a particular element of the
// entity will be delivered eventually in order without loss for the
// duration of the RPC. Note however that if multiple Changes apply to
// the same element, the implementation is free to suppress them and
// deliver just the last one.  The underlying system must provide the
// guarantee that any relevant update received for an entity E after a
// client's watch point for E MUST be delivered to that client.
//
// These tight guarantees allow for the following simplifications in
// the client:
//
//   (1) The client does not need to have a separate polling loop to
//       make up for missed updates.
//
//   (2) The client does not need to manage timestamps/versions
//       manually; the last update delivered corresponds to the
//       eventual state of the entity.
package watch

import (
	"fmt"
	"io"
	"v.io/v23"
	"v.io/v23/context"
	"v.io/v23/i18n"
	"v.io/v23/rpc"
	"v.io/v23/security/access"
	"v.io/v23/vdl"
	"v.io/v23/verror"
	"v.io/v23/vom"
)

var _ = __VDLInit() // Must be first; see __VDLInit comments for details.

//////////////////////////////////////////////////
// Type definitions

// ResumeMarker specifies how much of the existing underlying state
// is delivered to the client when the watch request is received by
// the system. The client can set this marker in one of the
// following ways to get different semantics:
//
// (A) Parameter is left empty.
//     Semantics: Fetch initial state.
//     The client wants the entities' initial states to be delivered.
//     See the description in "Initial State".
//
// (B) Parameter is set to the string "now" (UTF-8 encoding).
//     Semantics: Fetch new changes only.
//     The client just wants to get the changes received by the
//     system after the watch point. The system may deliver changes
//     from before the watch point as well.
//
// (C) Parameter is set to a value received in an earlier
//     Change.ResumeMarker field while watching the same entity with
//     the same query.
//     Semantics: Resume from a specific point.
//     The client wants to receive the changes from a specific point
//     - this value must correspond to a value received in the
//     Change.ResumeMarker field. The system may deliver changes
//     from before the ResumeMarker as well.  If the system cannot
//     resume the stream from this point (e.g., if it is too far
//     behind in the stream), it can return the
//     ErrUnknownResumeMarker error.
//
// An implementation MUST support the empty string "" marker
// (initial state fetching) and the "now" marker. It need not
// support resuming from a specific point.
type ResumeMarker []byte

func (ResumeMarker) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/services/watch.ResumeMarker"`
}) {
}

func (x ResumeMarker) VDLIsZero() bool {
	return len(x) == 0
}

func (x ResumeMarker) VDLWrite(enc vdl.Encoder) error {
	if err := enc.StartValue(vdl.TypeOf((*ResumeMarker)(nil))); err != nil {
		return err
	}
	if err := enc.EncodeBytes([]byte(x)); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *ResumeMarker) VDLRead(dec vdl.Decoder) error {
	if err := dec.StartValue(); err != nil {
		return err
	}
	var bytes []byte
	if err := dec.DecodeBytes(-1, &bytes); err != nil {
		return err
	}
	*x = bytes
	return dec.FinishValue()
}

// GlobRequest specifies which entities should be watched and, optionally,
// how to resume from a previous Watch call.
type GlobRequest struct {
	// Pattern specifies the subset of the children of the root entity
	// for which the client wants updates.
	Pattern string
	// ResumeMarker specifies how to resume from a previous Watch call.
	// See the ResumeMarker type for detailed comments.
	ResumeMarker ResumeMarker
}

func (GlobRequest) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/services/watch.GlobRequest"`
}) {
}

func (x GlobRequest) VDLIsZero() bool {
	if x.Pattern != "" {
		return false
	}
	if len(x.ResumeMarker) != 0 {
		return false
	}
	return true
}

func (x GlobRequest) VDLWrite(enc vdl.Encoder) error {
	if err := enc.StartValue(vdl.TypeOf((*GlobRequest)(nil)).Elem()); err != nil {
		return err
	}
	if x.Pattern != "" {
		if err := enc.NextField("Pattern"); err != nil {
			return err
		}
		if err := enc.StartValue(vdl.StringType); err != nil {
			return err
		}
		if err := enc.EncodeString(x.Pattern); err != nil {
			return err
		}
		if err := enc.FinishValue(); err != nil {
			return err
		}
	}
	if len(x.ResumeMarker) != 0 {
		if err := enc.NextField("ResumeMarker"); err != nil {
			return err
		}
		if err := x.ResumeMarker.VDLWrite(enc); err != nil {
			return err
		}
	}
	if err := enc.NextField(""); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *GlobRequest) VDLRead(dec vdl.Decoder) error {
	*x = GlobRequest{}
	if err := dec.StartValue(); err != nil {
		return err
	}
	if (dec.StackDepth() == 1 || dec.IsAny()) && !vdl.Compatible(vdl.TypeOf(*x), dec.Type()) {
		return fmt.Errorf("incompatible struct %T, from %v", *x, dec.Type())
	}
	for {
		f, err := dec.NextField()
		if err != nil {
			return err
		}
		switch f {
		case "":
			return dec.FinishValue()
		case "Pattern":
			if err := dec.StartValue(); err != nil {
				return err
			}
			var err error
			if x.Pattern, err = dec.DecodeString(); err != nil {
				return err
			}
			if err := dec.FinishValue(); err != nil {
				return err
			}
		case "ResumeMarker":
			if err := x.ResumeMarker.VDLRead(dec); err != nil {
				return err
			}
		default:
			if err := dec.SkipValue(); err != nil {
				return err
			}
		}
	}
}

// Change is the new value for a watched entity.
type Change struct {
	// Name is the Object name of the entity that changed.  This name is relative
	// to the root entity (i.e. the name of the Watcher service).
	Name string
	// State must be one of Exists, DoesNotExist, or InitialStateSkipped.
	State int32
	// Value contains the service-specific data for the entity.
	Value *vom.RawBytes
	// If present, provides a compact representation of all the messages
	// that have been received by the caller for the given Watch call.
	// For example, it could be a sequence number or a multi-part
	// timestamp/version vector. This marker can be provided in the
	// Request message to allow the caller to resume the stream watching
	// at a specific point without fetching the initial state.
	ResumeMarker ResumeMarker
	// If true, this Change is followed by more Changes that are in the
	// same group as this Change.
	Continued bool
}

func (Change) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/services/watch.Change"`
}) {
}

func (x Change) VDLIsZero() bool {
	if x.Name != "" {
		return false
	}
	if x.State != 0 {
		return false
	}
	if x.Value != nil && !x.Value.VDLIsZero() {
		return false
	}
	if len(x.ResumeMarker) != 0 {
		return false
	}
	if x.Continued {
		return false
	}
	return true
}

func (x Change) VDLWrite(enc vdl.Encoder) error {
	if err := enc.StartValue(vdl.TypeOf((*Change)(nil)).Elem()); err != nil {
		return err
	}
	if x.Name != "" {
		if err := enc.NextField("Name"); err != nil {
			return err
		}
		if err := enc.StartValue(vdl.StringType); err != nil {
			return err
		}
		if err := enc.EncodeString(x.Name); err != nil {
			return err
		}
		if err := enc.FinishValue(); err != nil {
			return err
		}
	}
	if x.State != 0 {
		if err := enc.NextField("State"); err != nil {
			return err
		}
		if err := enc.StartValue(vdl.Int32Type); err != nil {
			return err
		}
		if err := enc.EncodeInt(int64(x.State)); err != nil {
			return err
		}
		if err := enc.FinishValue(); err != nil {
			return err
		}
	}
	if x.Value != nil && !x.Value.VDLIsZero() {
		if err := enc.NextField("Value"); err != nil {
			return err
		}
		if err := x.Value.VDLWrite(enc); err != nil {
			return err
		}
	}
	if len(x.ResumeMarker) != 0 {
		if err := enc.NextField("ResumeMarker"); err != nil {
			return err
		}
		if err := x.ResumeMarker.VDLWrite(enc); err != nil {
			return err
		}
	}
	if x.Continued {
		if err := enc.NextField("Continued"); err != nil {
			return err
		}
		if err := enc.StartValue(vdl.BoolType); err != nil {
			return err
		}
		if err := enc.EncodeBool(x.Continued); err != nil {
			return err
		}
		if err := enc.FinishValue(); err != nil {
			return err
		}
	}
	if err := enc.NextField(""); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *Change) VDLRead(dec vdl.Decoder) error {
	*x = Change{
		Value: vom.RawBytesOf(vdl.ZeroValue(vdl.AnyType)),
	}
	if err := dec.StartValue(); err != nil {
		return err
	}
	if (dec.StackDepth() == 1 || dec.IsAny()) && !vdl.Compatible(vdl.TypeOf(*x), dec.Type()) {
		return fmt.Errorf("incompatible struct %T, from %v", *x, dec.Type())
	}
	for {
		f, err := dec.NextField()
		if err != nil {
			return err
		}
		switch f {
		case "":
			return dec.FinishValue()
		case "Name":
			if err := dec.StartValue(); err != nil {
				return err
			}
			var err error
			if x.Name, err = dec.DecodeString(); err != nil {
				return err
			}
			if err := dec.FinishValue(); err != nil {
				return err
			}
		case "State":
			if err := dec.StartValue(); err != nil {
				return err
			}
			tmp, err := dec.DecodeInt(32)
			if err != nil {
				return err
			}
			x.State = int32(tmp)
			if err := dec.FinishValue(); err != nil {
				return err
			}
		case "Value":
			x.Value = new(vom.RawBytes)
			if err := x.Value.VDLRead(dec); err != nil {
				return err
			}
		case "ResumeMarker":
			if err := x.ResumeMarker.VDLRead(dec); err != nil {
				return err
			}
		case "Continued":
			if err := dec.StartValue(); err != nil {
				return err
			}
			var err error
			if x.Continued, err = dec.DecodeBool(); err != nil {
				return err
			}
			if err := dec.FinishValue(); err != nil {
				return err
			}
		default:
			if err := dec.SkipValue(); err != nil {
				return err
			}
		}
	}
}

//////////////////////////////////////////////////
// Const definitions

// The entity exists and its full value is included in Value.
const Exists = int32(0)

// The entity does not exist.
const DoesNotExist = int32(1)

// The root entity and its children may or may not exist. Used only
// for initial state delivery when the client is not interested in
// fetching the initial state. See the "Initial State" section
// above.
const InitialStateSkipped = int32(2)

//////////////////////////////////////////////////
// Error definitions

var (
	ErrUnknownResumeMarker = verror.Register("v.io/v23/services/watch.UnknownResumeMarker", verror.NoRetry, "{1:}{2:} unknown resume marker {_}")
)

// NewErrUnknownResumeMarker returns an error with the ErrUnknownResumeMarker ID.
func NewErrUnknownResumeMarker(ctx *context.T) error {
	return verror.New(ErrUnknownResumeMarker, ctx)
}

//////////////////////////////////////////////////
// Interface definitions

// GlobWatcherClientMethods is the client interface
// containing GlobWatcher methods.
//
// GlobWatcher allows a client to receive updates for changes to objects
// that match a pattern.  See the package comments for details.
type GlobWatcherClientMethods interface {
	// WatchGlob returns a stream of changes that match a pattern.
	WatchGlob(_ *context.T, req GlobRequest, _ ...rpc.CallOpt) (GlobWatcherWatchGlobClientCall, error)
}

// GlobWatcherClientStub adds universal methods to GlobWatcherClientMethods.
type GlobWatcherClientStub interface {
	GlobWatcherClientMethods
	rpc.UniversalServiceMethods
}

// GlobWatcherClient returns a client stub for GlobWatcher.
func GlobWatcherClient(name string) GlobWatcherClientStub {
	return implGlobWatcherClientStub{name}
}

type implGlobWatcherClientStub struct {
	name string
}

func (c implGlobWatcherClientStub) WatchGlob(ctx *context.T, i0 GlobRequest, opts ...rpc.CallOpt) (ocall GlobWatcherWatchGlobClientCall, err error) {
	var call rpc.ClientCall
	if call, err = v23.GetClient(ctx).StartCall(ctx, c.name, "WatchGlob", []interface{}{i0}, opts...); err != nil {
		return
	}
	ocall = &implGlobWatcherWatchGlobClientCall{ClientCall: call}
	return
}

// GlobWatcherWatchGlobClientStream is the client stream for GlobWatcher.WatchGlob.
type GlobWatcherWatchGlobClientStream interface {
	// RecvStream returns the receiver side of the GlobWatcher.WatchGlob client stream.
	RecvStream() interface {
		// Advance stages an item so that it may be retrieved via Value.  Returns
		// true iff there is an item to retrieve.  Advance must be called before
		// Value is called.  May block if an item is not available.
		Advance() bool
		// Value returns the item that was staged by Advance.  May panic if Advance
		// returned false or was not called.  Never blocks.
		Value() Change
		// Err returns any error encountered by Advance.  Never blocks.
		Err() error
	}
}

// GlobWatcherWatchGlobClientCall represents the call returned from GlobWatcher.WatchGlob.
type GlobWatcherWatchGlobClientCall interface {
	GlobWatcherWatchGlobClientStream
	// Finish blocks until the server is done, and returns the positional return
	// values for call.
	//
	// Finish returns immediately if the call has been canceled; depending on the
	// timing the output could either be an error signaling cancelation, or the
	// valid positional return values from the server.
	//
	// Calling Finish is mandatory for releasing stream resources, unless the call
	// has been canceled or any of the other methods return an error.  Finish should
	// be called at most once.
	Finish() error
}

type implGlobWatcherWatchGlobClientCall struct {
	rpc.ClientCall
	valRecv Change
	errRecv error
}

func (c *implGlobWatcherWatchGlobClientCall) RecvStream() interface {
	Advance() bool
	Value() Change
	Err() error
} {
	return implGlobWatcherWatchGlobClientCallRecv{c}
}

type implGlobWatcherWatchGlobClientCallRecv struct {
	c *implGlobWatcherWatchGlobClientCall
}

func (c implGlobWatcherWatchGlobClientCallRecv) Advance() bool {
	c.c.valRecv = Change{}
	c.c.errRecv = c.c.Recv(&c.c.valRecv)
	return c.c.errRecv == nil
}
func (c implGlobWatcherWatchGlobClientCallRecv) Value() Change {
	return c.c.valRecv
}
func (c implGlobWatcherWatchGlobClientCallRecv) Err() error {
	if c.c.errRecv == io.EOF {
		return nil
	}
	return c.c.errRecv
}
func (c *implGlobWatcherWatchGlobClientCall) Finish() (err error) {
	err = c.ClientCall.Finish()
	return
}

// GlobWatcherServerMethods is the interface a server writer
// implements for GlobWatcher.
//
// GlobWatcher allows a client to receive updates for changes to objects
// that match a pattern.  See the package comments for details.
type GlobWatcherServerMethods interface {
	// WatchGlob returns a stream of changes that match a pattern.
	WatchGlob(_ *context.T, _ GlobWatcherWatchGlobServerCall, req GlobRequest) error
}

// GlobWatcherServerStubMethods is the server interface containing
// GlobWatcher methods, as expected by rpc.Server.
// The only difference between this interface and GlobWatcherServerMethods
// is the streaming methods.
type GlobWatcherServerStubMethods interface {
	// WatchGlob returns a stream of changes that match a pattern.
	WatchGlob(_ *context.T, _ *GlobWatcherWatchGlobServerCallStub, req GlobRequest) error
}

// GlobWatcherServerStub adds universal methods to GlobWatcherServerStubMethods.
type GlobWatcherServerStub interface {
	GlobWatcherServerStubMethods
	// Describe the GlobWatcher interfaces.
	Describe__() []rpc.InterfaceDesc
}

// GlobWatcherServer returns a server stub for GlobWatcher.
// It converts an implementation of GlobWatcherServerMethods into
// an object that may be used by rpc.Server.
func GlobWatcherServer(impl GlobWatcherServerMethods) GlobWatcherServerStub {
	stub := implGlobWatcherServerStub{
		impl: impl,
	}
	// Initialize GlobState; always check the stub itself first, to handle the
	// case where the user has the Glob method defined in their VDL source.
	if gs := rpc.NewGlobState(stub); gs != nil {
		stub.gs = gs
	} else if gs := rpc.NewGlobState(impl); gs != nil {
		stub.gs = gs
	}
	return stub
}

type implGlobWatcherServerStub struct {
	impl GlobWatcherServerMethods
	gs   *rpc.GlobState
}

func (s implGlobWatcherServerStub) WatchGlob(ctx *context.T, call *GlobWatcherWatchGlobServerCallStub, i0 GlobRequest) error {
	return s.impl.WatchGlob(ctx, call, i0)
}

func (s implGlobWatcherServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implGlobWatcherServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{GlobWatcherDesc}
}

// GlobWatcherDesc describes the GlobWatcher interface.
var GlobWatcherDesc rpc.InterfaceDesc = descGlobWatcher

// descGlobWatcher hides the desc to keep godoc clean.
var descGlobWatcher = rpc.InterfaceDesc{
	Name:    "GlobWatcher",
	PkgPath: "v.io/v23/services/watch",
	Doc:     "// GlobWatcher allows a client to receive updates for changes to objects\n// that match a pattern.  See the package comments for details.",
	Methods: []rpc.MethodDesc{
		{
			Name: "WatchGlob",
			Doc:  "// WatchGlob returns a stream of changes that match a pattern.",
			InArgs: []rpc.ArgDesc{
				{"req", ``}, // GlobRequest
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Resolve"))},
		},
	},
}

// GlobWatcherWatchGlobServerStream is the server stream for GlobWatcher.WatchGlob.
type GlobWatcherWatchGlobServerStream interface {
	// SendStream returns the send side of the GlobWatcher.WatchGlob server stream.
	SendStream() interface {
		// Send places the item onto the output stream.  Returns errors encountered
		// while sending.  Blocks if there is no buffer space; will unblock when
		// buffer space is available.
		Send(item Change) error
	}
}

// GlobWatcherWatchGlobServerCall represents the context passed to GlobWatcher.WatchGlob.
type GlobWatcherWatchGlobServerCall interface {
	rpc.ServerCall
	GlobWatcherWatchGlobServerStream
}

// GlobWatcherWatchGlobServerCallStub is a wrapper that converts rpc.StreamServerCall into
// a typesafe stub that implements GlobWatcherWatchGlobServerCall.
type GlobWatcherWatchGlobServerCallStub struct {
	rpc.StreamServerCall
}

// Init initializes GlobWatcherWatchGlobServerCallStub from rpc.StreamServerCall.
func (s *GlobWatcherWatchGlobServerCallStub) Init(call rpc.StreamServerCall) {
	s.StreamServerCall = call
}

// SendStream returns the send side of the GlobWatcher.WatchGlob server stream.
func (s *GlobWatcherWatchGlobServerCallStub) SendStream() interface {
	Send(item Change) error
} {
	return implGlobWatcherWatchGlobServerCallSend{s}
}

type implGlobWatcherWatchGlobServerCallSend struct {
	s *GlobWatcherWatchGlobServerCallStub
}

func (s implGlobWatcherWatchGlobServerCallSend) Send(item Change) error {
	return s.s.Send(item)
}

var __VDLInitCalled bool

// __VDLInit performs vdl initialization.  It is safe to call multiple times.
// If you have an init ordering issue, just insert the following line verbatim
// into your source files in this package, right after the "package foo" clause:
//
//    var _ = __VDLInit()
//
// The purpose of this function is to ensure that vdl initialization occurs in
// the right order, and very early in the init sequence.  In particular, vdl
// registration and package variable initialization needs to occur before
// functions like vdl.TypeOf will work properly.
//
// This function returns a dummy value, so that it can be used to initialize the
// first var in the file, to take advantage of Go's defined init order.
func __VDLInit() struct{} {
	if __VDLInitCalled {
		return struct{}{}
	}
	__VDLInitCalled = true

	// Register types.
	vdl.Register((*ResumeMarker)(nil))
	vdl.Register((*GlobRequest)(nil))
	vdl.Register((*Change)(nil))

	// Set error format strings.
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrUnknownResumeMarker.ID), "{1:}{2:} unknown resume marker {_}")

	return struct{}{}
}
