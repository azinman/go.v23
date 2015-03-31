// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

// RPCVersion represents a version of the RPC protocol.
type RPCVersion uint32

const (
	// UnknownRPCVersion is used for Min/MaxRPCVersion in an Endpoint when
	// we don't know the relevant version numbers.  In this case the RPC
	// implementation will have to guess the correct values.
	UnknownRPCVersion RPCVersion = iota
	// Deprecated versions
	rpcVersion2
	rpcVersion3
	rpcDummyVersion3 // So that the numeric value of RPCVersion4 is 4
	rpcVersion4
	// TODO(ashankar): Unexport all versions except the last before release
	// RPCVersion5 uses the new security model (Principal and Blessings objects),
	// and sends discharges for third-party caveats on the server's blessings
	// during authentication.
	RPCVersion5

	// RPCVersion6 adds control channel encryption to RPCVersion5.
	RPCVersion6

	// RPCVersion7 uses concrete types for security.Discharge during VC
	// authentication.
	RPCVersion7

	// RPCVersion8 uses separate VOM type flow to share VOM types across all flows
	// in a VC.
	RPCVersion8
)

// RPCVersionRange allows you to optionally specify a range of versions to
// use when calling FormatEndpoint
type RPCVersionRange struct {
	Min, Max RPCVersion
}

// EndpointOpt implents the EndpointOpt interface so an RPCVersionRange
// can be used as an option to FormatEndpoint.
func (RPCVersionRange) EndpointOpt() {}