// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package security

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"v.io/v23/context"
	"v.io/v23/verror"
)

var (
	errEmptyChain         = verror.Register(pkgPath+".errEmptyChain", verror.NoRetry, "empty certificate chain in blessings")
	errMisconfiguredRoots = verror.Register(pkgPath+".errMisconfiguredRoots", verror.NoRetry, "recognized root certificates not configured")
	errMultiplePublicKeys = verror.Register(pkgPath+".errMultiplePublicKeys", verror.NoRetry, "invalid blessings: two certificate chains that bind to different public keys")
	errInvalidUnion       = verror.Register(pkgPath+".errInvalidUnion", verror.NoRetry, "cannot create union of blessings bound to different public keys")
)

// Blessings encapsulates all cryptographic operations required to
// prove that a set of (human-readable) blessing names are bound
// to a principal in a specific call.
//
// Blessings objects are meant to be presented to other principals to
// authenticate and authorize actions. The functions 'LocalBlessingNames',
// 'SigningBlessingNames' and 'RemoteBlessingNames' defined in this package
// can be used to uncover the blessing names encapsulated in these objects.
//
// Blessings objects are immutable and multiple goroutines may invoke methods
// on them simultaneously.
//
// See also: https://github.com/vanadium/docs/blob/master/glossary.md#blessing
type Blessings struct {
	chains    [][]Certificate
	publicKey PublicKey
	digests   [][]byte // digests[i] is the digest of chains[i]
	uniqueID  []byte
}

// PublicKey returns the public key of the principal to which
// blessings obtained from this object are bound.
//
// Can return nil if b is the zero value.
func (b Blessings) PublicKey() PublicKey { return b.publicKey }

// ThirdPartyCaveats returns the set of third-party restrictions on the
// scope of the blessings (i.e., the subset of Caveats for which
// ThirdPartyDetails will be non-nil).
func (b Blessings) ThirdPartyCaveats() []Caveat {
	var ret []Caveat
	for _, chain := range b.chains {
		for _, cert := range chain {
			for _, cav := range cert.Caveats {
				if tp := cav.ThirdPartyDetails(); tp != nil {
					ret = append(ret, cav)
				}
			}
		}
	}
	return ret
}

// IsZero returns true if b represents the zero value of blessings (an empty
// set).
func (b Blessings) IsZero() bool {
	// b.publicKey == nil <=> len(b.chains) == 0
	return b.publicKey == nil
}

// Equivalent returns true if 'b' and 'blessings' can be used interchangeably,
// i.e., 'b' will be authorized wherever 'blessings' is and vice-versa.
func (b Blessings) Equivalent(blessings Blessings) bool {
	return bytes.Equal(b.UniqueID(), blessings.UniqueID())
}

// UniqueID returns an identifier of the set of blessings. Two blessings
// objects with the same UniqueID are interchangeable for any authorization
// decisions.
//
// The converse is not guaranteed at this time. Two interchangeable blessings
// objects may in theory have different UniqueIDs.
func (b Blessings) UniqueID() []byte { return b.uniqueID }

func (b *Blessings) init() {
	if len(b.digests) == 0 {
		return
	}
	tohash := make([]byte, 0, len(b.digests[0])*len(b.digests))
	for _, d := range b.digests {
		tohash = append(tohash, d...)
	}
	b.uniqueID = b.publicKey.hash().sum(tohash)
}

// CouldHaveNames returns true iff the blessings 'b' encapsulates the provided
// set of blessing names.
//
// This check ignores all caveats on the blessing name and the recognition status
// of its blessing root.
func (b Blessings) CouldHaveNames(names []string) bool {
	hasNames := make(map[string]bool)
	for _, chain := range b.chains {
		hasNames[claimedName(chain)] = true
	}
	for _, n := range names {
		if !hasNames[n] {
			return false
		}
	}
	return true
}

// Expiry returns the time at which b will no longer be valid, or the zero
// value of time.Time if the blessing does not expire.
func (b Blessings) Expiry() time.Time {
	var min time.Time
	for _, chain := range b.chains {
		for _, cert := range chain {
			for _, cav := range cert.Caveats {
				if t := expiryTime(cav); !t.IsZero() && (min.IsZero() || t.Before(min)) {
					min = t
				}
			}
		}
	}
	return min
}

func (b Blessings) publicKeyDER() []byte {
	chain := b.chains[0]
	return chain[len(chain)-1].PublicKey
}

func (b Blessings) blessingsByNameForPrincipal(p Principal, pattern BlessingPattern) Blessings {
	pKey, err := p.PublicKey().MarshalBinary()
	if err != nil {
		return Blessings{}
	}
	ret := Blessings{publicKey: b.publicKey}
	for _, chain := range b.chains {
		blessing := nameForPrincipal(pKey, p.Roots(), chain)
		if len(blessing) > 0 && pattern.MatchedBy(blessing) {
			ret.chains = append(ret.chains, chain)
		}
	}
	if len(ret.chains) == 0 {
		return Blessings{}
	}
	return ret
}

func (b Blessings) String() string {
	blessings := make([]string, len(b.chains))
	for chainidx, chain := range b.chains {
		onechain := make([]string, len(chain))
		for certidx, cert := range chain {
			onechain[certidx] = cert.Extension
		}
		blessings[chainidx] = fmt.Sprintf("%v", strings.Join(onechain, ChainSeparator))
	}
	return strings.Join(blessings, ",")
}

// claimedName returns the blessing name that the certificate chain claims to
// have (i.e., ignoring any caveats or recognition of the root public key as an
// authority on the namespace).
func claimedName(chain []Certificate) string {
	var buf bytes.Buffer
	buf.WriteString(chain[0].Extension)
	for i := 1; i < len(chain); i++ {
		buf.WriteString(ChainSeparator)
		buf.WriteString(chain[i].Extension)
	}
	return buf.String()
}

func nameForPrincipal(pubkey []byte, roots BlessingRoots, chain []Certificate) string {
	// Verify the chain belongs to this principal
	if !bytes.Equal(chain[len(chain)-1].PublicKey, pubkey) {
		return ""
	}
	// Verify that the root of the chain is recognized as an authority on
	// blessing.
	blessing := claimedName(chain)
	if err := roots.Recognized(chain[0].PublicKey, blessing); err != nil {
		return ""
	}
	return blessing
}

// chainCaveats returns the union of the set of caveats in the  certificates present
// in 'chain'.
func chainCaveats(chain []Certificate) []Caveat {
	var cavs []Caveat
	for _, c := range chain {
		cavs = append(cavs, c.Caveats...)
	}
	return cavs
}

func defaultCaveatValidation(ctx *context.T, call Call, chains [][]Caveat) []error {
	results := make([]error, len(chains))
	for i, chain := range chains {
		for _, cav := range chain {
			if err := cav.Validate(ctx, call); err != nil {
				results[i] = err
				break
			}
		}
	}
	return results
}

// TODO(ashankar): Get rid of this function? It allows users to mess
// with the integrity of 'b'.
func MarshalBlessings(b Blessings) WireBlessings {
	return WireBlessings{b.chains}
}

func wireBlessingsToNative(wire WireBlessings, native *Blessings) error {
	if len(wire.CertificateChains) == 0 {
		return nil
	}
	certchains := wire.CertificateChains
	if len(certchains) == 0 || len(certchains[0]) == 0 {
		return verror.New(errEmptyChain, nil)
	}
	// Public keys should match for all chains.
	marshaledkey := certchains[0][len(certchains[0])-1].PublicKey
	digests := make([][]byte, len(certchains))

	var key PublicKey
	var err error
	if key, digests[0], err = validateCertificateChain(certchains[0]); err != nil {
		return err
	}
	for i := 1; i < len(certchains); i++ {
		chain := certchains[i]
		if len(chain) == 0 {
			return verror.New(errEmptyChain, nil)
		}
		cert := chain[len(chain)-1]
		if !bytes.Equal(marshaledkey, cert.PublicKey) {
			return verror.New(errMultiplePublicKeys, nil)
		}
		if _, digests[i], err = validateCertificateChain(chain); err != nil {
			return err
		}
	}
	native.chains = certchains
	native.publicKey = key
	native.digests = digests
	native.init()
	return nil
}

func wireBlessingsFromNative(wire *WireBlessings, native Blessings) error {
	wire.CertificateChains = native.chains
	return nil
}

// UnionOfBlessings returns a Blessings object that carries the union of the
// provided blessings.
//
// All provided Blessings must have the same PublicKey.
//
// UnionOfBlessings with no arguments returns (nil, nil).
func UnionOfBlessings(blessings ...Blessings) (Blessings, error) {
	for len(blessings) > 0 && blessings[0].IsZero() {
		blessings = blessings[1:]
	}
	switch len(blessings) {
	case 0:
		return Blessings{}, nil
	case 1:
		return blessings[0], nil
	}
	key0 := blessings[0].publicKeyDER()
	var ret Blessings
	for idx, b := range blessings {
		if b.IsZero() {
			continue
		}
		if idx > 0 && !bytes.Equal(key0, b.publicKeyDER()) {
			return Blessings{}, verror.New(errInvalidUnion, nil)
		}
		ret.chains = append(ret.chains, b.chains...)
		ret.digests = append(ret.digests, b.digests...)
	}
	var err error
	if ret.publicKey, err = UnmarshalPublicKey(key0); err != nil {
		return Blessings{}, err
	}
	// For pretty printing, sort the certificate chains so that there is a consistent
	// ordering, irrespective of the ordering of arguments to UnionOfBlessings.
	sort.Stable(blessingsSorter(ret))
	ret.init()
	return ret, nil
}

type blessingsSorter Blessings

func (b blessingsSorter) Len() int { return len(b.chains) }
func (b blessingsSorter) Swap(i, j int) {
	b.chains[i], b.chains[j] = b.chains[j], b.chains[i]
	b.digests[i], b.digests[j] = b.digests[j], b.digests[i]
}
func (b blessingsSorter) Less(i, j int) bool {
	ci := b.chains[i]
	cj := b.chains[j]
	if len(ci) < len(cj) {
		return true
	}
	if len(ci) > len(cj) {
		return false
	}
	// Equal size chains, order by the names in the certificates.
	N := len(ci)
	for idx := 0; idx < N; idx++ {
		ie := ci[idx].Extension
		je := cj[idx].Extension
		if ie < je {
			return true
		}
		if ie > je {
			return false
		}
	}
	return false
}

func (i RejectedBlessing) String() string {
	return fmt.Sprintf("{%q: %v}", i.Blessing, i.Err)
}

// DefaultBlessingPatterns returns the BlessingsPatterns of the Default Blessings
// of the provided Principal.
func DefaultBlessingPatterns(p Principal) (patterns []BlessingPattern) {
	for b, _ := range p.BlessingsInfo(p.BlessingStore().Default()) {
		patterns = append(patterns, BlessingPattern(b))
	}
	return
}

var (
	caveatValidationMu         sync.RWMutex
	caveatValidation           = defaultCaveatValidation
	caveatValidationOverridden = false
)

func overrideCaveatValidation(fn func(ctx *context.T, call Call, sets [][]Caveat) []error) {
	caveatValidationMu.Lock()
	if caveatValidationOverridden {
		panic("security: OverrideCaveatValidation may only be called once")
	}
	caveatValidationOverridden = true
	caveatValidation = fn
	caveatValidationMu.Unlock()
}

func setCaveatValidationForTest(fn func(ctx *context.T, call Call, sets [][]Caveat) []error) {
	// For tests we skip the panic on multiple calls, so that we can easily revert
	// to the default validator.
	caveatValidationMu.Lock()
	caveatValidation = fn
	caveatValidationMu.Unlock()
}

func getCaveatValidation() func(ctx *context.T, call Call, sets [][]Caveat) []error {
	caveatValidationMu.RLock()
	fn := caveatValidation
	caveatValidationMu.RUnlock()
	return fn
}

func isSigningBlessingCaveat(cav Caveat) bool {
	switch cav.Id {
	case ExpiryCaveat.Id:
		return true
	case ConstCaveat.Id:
		return true
	}
	// TODO(ataly): Figure out what revocation caveats can be supported
	// for signing blessings. Our current revocation caveats are third-party
	// caveats and are unsuitable for signing for the following reasons:
	// - They violate the requirement that the caveat should be universally
	//   understandable and validatable. This is because validating the
	//   caveat requires reaching out to a third-party caveat discharger
	//   which may be reachable from some services and unreachable from others
	//   (perhaps due to network partitions).
	// - Given a third-party caveat, there is no way to tell whether it is revocation
	//   caveat.
	return false
}

func validateCaveatsForSigning(ctx *context.T, call Call, chain []Certificate) error {
	for _, cav := range chainCaveats(chain) {
		if !isSigningBlessingCaveat(cav) {
			return NewErrInvalidSigningBlessingCaveat(nil, cav.Id)
		}
		if err := cav.Validate(ctx, call); err != nil {
			return err
		}
	}
	return nil
}

// SigningBlessings returns a blessings object that encapsulates the subset of
// names of the provided blessings object that can be used to sign data at rest.
//
// The names of the returned blessings object can be obtained using the
// 'SigningNames' function.
func SigningBlessings(blessings Blessings) Blessings {
	ret := Blessings{}
	for _, chain := range blessings.chains {
		suitableForSigning := true
		for _, cav := range chainCaveats(chain) {
			if !isSigningBlessingCaveat(cav) {
				suitableForSigning = false
				break
			}
		}
		if suitableForSigning {
			ret.chains = append(ret.chains, chain)
		}
	}
	if len(ret.chains) > 0 {
		ret.publicKey = blessings.publicKey
	}
	return ret
}

// SigningBlessingNames returns the validated set of human-readable blessing
// names encapsulated in the provided signing blessings object, as determined
// by the provided principal.
//
// This function also returns the RejectedBlessings for each blessing name that
// cannot be validated.
// TODO(ataly): While the principal is encapsulated inside context.T, we can't
// extract it due to an import cycle issue. Therefore at the moment we have the
// principal separately passed to this function. We should clean this up.
func SigningBlessingNames(ctx *context.T, p Principal, blessings Blessings) ([]string, []RejectedBlessing) {
	if ctx == nil || p == nil {
		return nil, nil
	}
	if blessings.IsZero() {
		return nil, nil
	}
	var (
		validatedNames []string
		rejected       []RejectedBlessing
		call           = NewCall(&CallParams{LocalPrincipal: p, RemoteBlessings: blessings})
	)
	for _, chain := range blessings.chains {
		name := claimedName(chain)
		if err := p.Roots().Recognized(chain[0].PublicKey, name); err != nil {
			rejected = append(rejected, RejectedBlessing{name, err})
			continue
		}
		if err := validateCaveatsForSigning(ctx, call, chain); err != nil {
			rejected = append(rejected, RejectedBlessing{name, err})
			continue
		}
		validatedNames = append(validatedNames, name)
	}
	return validatedNames, rejected
}

// ReturnBlessingNames returns the validated set of human-readable blessing names
// encapsulated in the blessings object presented by the remote end of a call.
//
// The blessing names are guaranteed to:
//
// (1) Satisfy all the caveats associated with them, in the context of the call.
// (2) Be rooted in call.LocalPrincipal.Roots.
//
// Caveats are considered satisfied for the 'call' if the CaveatValidator implementation
// can be found in the address space of the caller and Validate returns nil.
//
// RemoteBlessingNames also returns the RejectedBlessings for each blessing name that cannot
// be validated.
func RemoteBlessingNames(ctx *context.T, call Call) ([]string, []RejectedBlessing) {
	if ctx == nil || call == nil {
		return nil, nil
	}
	b := call.RemoteBlessings()
	if b.IsZero() {
		return nil, nil
	}
	var (
		validatedNames    []string
		rejected          []RejectedBlessing
		pendingNames      []string
		pendingCaveatSets [][]Caveat
		rootsErr          error
	)
	if p := call.LocalPrincipal(); p == nil || p.Roots() == nil {
		// These conditions should not be possible: Here only for
		// unittests where the Call object is intentionally limited.
		rootsErr = verror.New(errMisconfiguredRoots, ctx)
	}
	for _, chain := range b.chains {
		name := claimedName(chain)
		err := rootsErr
		if err == nil {
			err = call.LocalPrincipal().Roots().Recognized(chain[0].PublicKey, name)
		}
		if err != nil {
			rejected = append(rejected, RejectedBlessing{name, err})
			continue
		}
		cavs := chainCaveats(chain)
		if len(cavs) == 0 {
			validatedNames = append(validatedNames, name) // No caveats to validate, add it to blessingNames.
		} else {
			pendingNames = append(pendingNames, name)
			pendingCaveatSets = append(pendingCaveatSets, cavs)
		}
	}
	if len(pendingCaveatSets) == 0 {
		return validatedNames, rejected
	}

	validationResults := getCaveatValidation()(ctx, call, pendingCaveatSets)
	if g, w := len(validationResults), len(pendingNames); g != w {
		panic(fmt.Sprintf("Got wrong number of validation results. Got %d, expected %d.", g, w))
	}

	for i, resultErr := range validationResults {
		if resultErr == nil {
			validatedNames = append(validatedNames, pendingNames[i])
		} else {
			rejected = append(rejected, RejectedBlessing{pendingNames[i], resultErr})
		}
	}
	return validatedNames, rejected
}

// LocalBlessingNames returns the set of human-readable blessing names encapsulated
// in the blessings object presented by the local end of the call.
//
// The blessing names are guaranteed to be rooted in call.LocalPrincipal.Roots.
//
// LocalBlessingNames does not validate caveats on the blessing names.
func LocalBlessingNames(ctx *context.T, call Call) []string {
	if ctx == nil || call == nil {
		return nil
	}
	var names []string
	for n, _ := range call.LocalPrincipal().BlessingsInfo(call.LocalBlessings()) {
		names = append(names, n)
	}
	return names
}
