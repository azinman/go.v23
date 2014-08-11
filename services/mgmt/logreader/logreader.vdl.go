// This file was auto-generated by the veyron vdl tool.
// Source: logreader.vdl

// Package logreader contains the interface for reading log files remotely.
package logreader

import (
	// The non-user imports are prefixed with "_gen_" to prevent collisions.
	_gen_io "io"
	_gen_context "veyron2/context"
	_gen_ipc "veyron2/ipc"
	_gen_naming "veyron2/naming"
	_gen_rt "veyron2/rt"
	_gen_vdlutil "veyron2/vdl/vdlutil"
	_gen_wiretype "veyron2/wiretype"
)

// LogLine is a log entry from a log file.
type LogEntry struct {
	// The offset (in bytes) where this entry starts.
	Position int64
	// The content of the log entry.
	Line string
}

const (
	// A special NumEntries value that indicates that all entries should be
	// returned by ReadLog.
	AllEntries = int64(-1)
)

// TODO(bprosnitz) Remove this line once signatures are updated to use typevals.
// It corrects a bug where _gen_wiretype is unused in VDL pacakges where only bootstrap types are used on interfaces.
const _ = _gen_wiretype.TypeIDInvalid

// LogFile can be used to access log files remotely.
// LogFile is the interface the client binds and uses.
// LogFile_ExcludingUniversal is the interface without internal framework-added methods
// to enable embedding without method collisions.  Not to be used directly by clients.
type LogFile_ExcludingUniversal interface {
	// Size returns the number of bytes in the receiving object.
	Size(ctx _gen_context.T, opts ..._gen_ipc.CallOpt) (reply int64, err error)
	// ReadLog receives up to NumEntries log entries starting at the
	// StartPos offset (in bytes) in the receiving object. Each stream chunk
	// contains one log entry.
	//
	// If Follow is true, ReadLog will block and wait for more entries to
	// arrive when it reaches the end of the file.
	//
	// ReadLog returns the position where it stopped reading, i.e. the
	// position where the next entry starts. This value can be used as
	// StartPos for successive calls to ReadLog.
	//
	// The returned error will be io.EOF if and only if ReadLog reached the
	// end of the file and no log entries were returned.
	ReadLog(ctx _gen_context.T, StartPos int64, NumEntries int32, Follow bool, opts ..._gen_ipc.CallOpt) (reply LogFileReadLogCall, err error)
}
type LogFile interface {
	_gen_ipc.UniversalServiceMethods
	LogFile_ExcludingUniversal
}

// LogFileService is the interface the server implements.
type LogFileService interface {

	// Size returns the number of bytes in the receiving object.
	Size(context _gen_ipc.ServerContext) (reply int64, err error)
	// ReadLog receives up to NumEntries log entries starting at the
	// StartPos offset (in bytes) in the receiving object. Each stream chunk
	// contains one log entry.
	//
	// If Follow is true, ReadLog will block and wait for more entries to
	// arrive when it reaches the end of the file.
	//
	// ReadLog returns the position where it stopped reading, i.e. the
	// position where the next entry starts. This value can be used as
	// StartPos for successive calls to ReadLog.
	//
	// The returned error will be io.EOF if and only if ReadLog reached the
	// end of the file and no log entries were returned.
	ReadLog(context _gen_ipc.ServerContext, StartPos int64, NumEntries int32, Follow bool, stream LogFileServiceReadLogStream) (reply int64, err error)
}

// LogFileReadLogCall is the interface for call object of the method
// ReadLog in the service interface LogFile.
type LogFileReadLogCall interface {
	// RecvStream returns the recv portion of the stream
	RecvStream() interface {
		// Advance stages an element so the client can retrieve it
		// with Value.  Advance returns true iff there is an
		// element to retrieve.  The client must call Advance before
		// calling Value. Advance may block if an element is not
		// immediately available.
		Advance() bool

		// Value returns the element that was staged by Advance.
		// Value may panic if Advance returned false or was not
		// called at all.  Value does not block.
		Value() LogEntry

		// Err returns a non-nil error iff the stream encountered
		// any errors.  Err does not block.
		Err() error
	}

	// Finish blocks until the server is done and returns the positional
	// return values for call.
	//
	// If Cancel has been called, Finish will return immediately; the output of
	// Finish could either be an error signalling cancelation, or the correct
	// positional return values from the server depending on the timing of the
	// call.
	//
	// Calling Finish is mandatory for releasing stream resources, unless Cancel
	// has been called or any of the other methods return an error.
	// Finish should be called at most once.
	Finish() (reply int64, err error)

	// Cancel cancels the RPC, notifying the server to stop processing.  It
	// is safe to call Cancel concurrently with any of the other stream methods.
	// Calling Cancel after Finish has returned is a no-op.
	Cancel()
}

type implLogFileReadLogStreamIterator struct {
	clientCall _gen_ipc.Call
	val        LogEntry
	err        error
}

func (c *implLogFileReadLogStreamIterator) Advance() bool {
	c.val = LogEntry{}
	c.err = c.clientCall.Recv(&c.val)
	return c.err == nil
}

func (c *implLogFileReadLogStreamIterator) Value() LogEntry {
	return c.val
}

func (c *implLogFileReadLogStreamIterator) Err() error {
	if c.err == _gen_io.EOF {
		return nil
	}
	return c.err
}

// Implementation of the LogFileReadLogCall interface that is not exported.
type implLogFileReadLogCall struct {
	clientCall _gen_ipc.Call
	readStream implLogFileReadLogStreamIterator
}

func (c *implLogFileReadLogCall) RecvStream() interface {
	Advance() bool
	Value() LogEntry
	Err() error
} {
	return &c.readStream
}

func (c *implLogFileReadLogCall) Finish() (reply int64, err error) {
	if ierr := c.clientCall.Finish(&reply, &err); ierr != nil {
		err = ierr
	}
	return
}

func (c *implLogFileReadLogCall) Cancel() {
	c.clientCall.Cancel()
}

type implLogFileServiceReadLogStreamSender struct {
	serverCall _gen_ipc.ServerCall
}

func (s *implLogFileServiceReadLogStreamSender) Send(item LogEntry) error {
	return s.serverCall.Send(item)
}

// LogFileServiceReadLogStream is the interface for streaming responses of the method
// ReadLog in the service interface LogFile.
type LogFileServiceReadLogStream interface {
	// SendStream returns the send portion of the stream.
	SendStream() interface {
		// Send places the item onto the output stream, blocking if there is no buffer
		// space available.  If the client has canceled, an error is returned.
		Send(item LogEntry) error
	}
}

// Implementation of the LogFileServiceReadLogStream interface that is not exported.
type implLogFileServiceReadLogStream struct {
	writer implLogFileServiceReadLogStreamSender
}

func (s *implLogFileServiceReadLogStream) SendStream() interface {
	// Send places the item onto the output stream, blocking if there is no buffer
	// space available.  If the client has canceled, an error is returned.
	Send(item LogEntry) error
} {
	return &s.writer
}

// BindLogFile returns the client stub implementing the LogFile
// interface.
//
// If no _gen_ipc.Client is specified, the default _gen_ipc.Client in the
// global Runtime is used.
func BindLogFile(name string, opts ..._gen_ipc.BindOpt) (LogFile, error) {
	var client _gen_ipc.Client
	switch len(opts) {
	case 0:
		client = _gen_rt.R().Client()
	case 1:
		switch o := opts[0].(type) {
		case _gen_ipc.Client:
			client = o
		default:
			return nil, _gen_vdlutil.ErrUnrecognizedOption
		}
	default:
		return nil, _gen_vdlutil.ErrTooManyOptionsToBind
	}
	stub := &clientStubLogFile{client: client, name: name}

	return stub, nil
}

// NewServerLogFile creates a new server stub.
//
// It takes a regular server implementing the LogFileService
// interface, and returns a new server stub.
func NewServerLogFile(server LogFileService) interface{} {
	return &ServerStubLogFile{
		service: server,
	}
}

// clientStubLogFile implements LogFile.
type clientStubLogFile struct {
	client _gen_ipc.Client
	name   string
}

func (__gen_c *clientStubLogFile) Size(ctx _gen_context.T, opts ..._gen_ipc.CallOpt) (reply int64, err error) {
	var call _gen_ipc.Call
	if call, err = __gen_c.client.StartCall(ctx, __gen_c.name, "Size", nil, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&reply, &err); ierr != nil {
		err = ierr
	}
	return
}

func (__gen_c *clientStubLogFile) ReadLog(ctx _gen_context.T, StartPos int64, NumEntries int32, Follow bool, opts ..._gen_ipc.CallOpt) (reply LogFileReadLogCall, err error) {
	var call _gen_ipc.Call
	if call, err = __gen_c.client.StartCall(ctx, __gen_c.name, "ReadLog", []interface{}{StartPos, NumEntries, Follow}, opts...); err != nil {
		return
	}
	reply = &implLogFileReadLogCall{clientCall: call, readStream: implLogFileReadLogStreamIterator{clientCall: call}}
	return
}

func (__gen_c *clientStubLogFile) UnresolveStep(ctx _gen_context.T, opts ..._gen_ipc.CallOpt) (reply []string, err error) {
	var call _gen_ipc.Call
	if call, err = __gen_c.client.StartCall(ctx, __gen_c.name, "UnresolveStep", nil, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&reply, &err); ierr != nil {
		err = ierr
	}
	return
}

func (__gen_c *clientStubLogFile) Signature(ctx _gen_context.T, opts ..._gen_ipc.CallOpt) (reply _gen_ipc.ServiceSignature, err error) {
	var call _gen_ipc.Call
	if call, err = __gen_c.client.StartCall(ctx, __gen_c.name, "Signature", nil, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&reply, &err); ierr != nil {
		err = ierr
	}
	return
}

func (__gen_c *clientStubLogFile) GetMethodTags(ctx _gen_context.T, method string, opts ..._gen_ipc.CallOpt) (reply []interface{}, err error) {
	var call _gen_ipc.Call
	if call, err = __gen_c.client.StartCall(ctx, __gen_c.name, "GetMethodTags", []interface{}{method}, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&reply, &err); ierr != nil {
		err = ierr
	}
	return
}

// ServerStubLogFile wraps a server that implements
// LogFileService and provides an object that satisfies
// the requirements of veyron2/ipc.ReflectInvoker.
type ServerStubLogFile struct {
	service LogFileService
}

func (__gen_s *ServerStubLogFile) GetMethodTags(call _gen_ipc.ServerCall, method string) ([]interface{}, error) {
	// TODO(bprosnitz) GetMethodTags() will be replaces with Signature().
	// Note: This exhibits some weird behavior like returning a nil error if the method isn't found.
	// This will change when it is replaced with Signature().
	switch method {
	case "Size":
		return []interface{}{}, nil
	case "ReadLog":
		return []interface{}{}, nil
	default:
		return nil, nil
	}
}

func (__gen_s *ServerStubLogFile) Signature(call _gen_ipc.ServerCall) (_gen_ipc.ServiceSignature, error) {
	result := _gen_ipc.ServiceSignature{Methods: make(map[string]_gen_ipc.MethodSignature)}
	result.Methods["ReadLog"] = _gen_ipc.MethodSignature{
		InArgs: []_gen_ipc.MethodArgument{
			{Name: "StartPos", Type: 37},
			{Name: "NumEntries", Type: 36},
			{Name: "Follow", Type: 2},
		},
		OutArgs: []_gen_ipc.MethodArgument{
			{Name: "", Type: 37},
			{Name: "", Type: 65},
		},

		OutStream: 66,
	}
	result.Methods["Size"] = _gen_ipc.MethodSignature{
		InArgs: []_gen_ipc.MethodArgument{},
		OutArgs: []_gen_ipc.MethodArgument{
			{Name: "", Type: 37},
			{Name: "", Type: 65},
		},
	}

	result.TypeDefs = []_gen_vdlutil.Any{
		_gen_wiretype.NamedPrimitiveType{Type: 0x1, Name: "error", Tags: []string(nil)}, _gen_wiretype.StructType{
			[]_gen_wiretype.FieldType{
				_gen_wiretype.FieldType{Type: 0x25, Name: "Position"},
				_gen_wiretype.FieldType{Type: 0x3, Name: "Line"},
			},
			"veyron2/services/mgmt/logreader.LogEntry", []string(nil)},
	}

	return result, nil
}

func (__gen_s *ServerStubLogFile) UnresolveStep(call _gen_ipc.ServerCall) (reply []string, err error) {
	if unresolver, ok := __gen_s.service.(_gen_ipc.Unresolver); ok {
		return unresolver.UnresolveStep(call)
	}
	if call.Server() == nil {
		return
	}
	var published []string
	if published, err = call.Server().Published(); err != nil || published == nil {
		return
	}
	reply = make([]string, len(published))
	for i, p := range published {
		reply[i] = _gen_naming.Join(p, call.Name())
	}
	return
}

func (__gen_s *ServerStubLogFile) Size(call _gen_ipc.ServerCall) (reply int64, err error) {
	reply, err = __gen_s.service.Size(call)
	return
}

func (__gen_s *ServerStubLogFile) ReadLog(call _gen_ipc.ServerCall, StartPos int64, NumEntries int32, Follow bool) (reply int64, err error) {
	stream := &implLogFileServiceReadLogStream{writer: implLogFileServiceReadLogStreamSender{serverCall: call}}
	reply, err = __gen_s.service.ReadLog(call, StartPos, NumEntries, Follow, stream)
	return
}
