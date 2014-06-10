package security

// THIS FILE SHOULD BE DELETED SOON. IT IS A PLACEHOLDER FOR REAL
// IMPLEMENTATIONS OF THE IDENTITY INTERFACES THAT ARE STILL DEVELOPING.

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"veyron2/vom"
)

const chainSeparator = "/"

func fakeName(name string) string {
	const prefix = "fake"
	if len(name) == 0 {
		return prefix
	}
	return prefix + chainSeparator + name
}

// fakeID implements both PublicID and PrivateID.
// Normally, the implementation types would be separated so that encoding
// the PublicID implementation (using VOM perhaps) and sending it across the
// wire does not leak the private key.
//
// However, in this particular case, there is no such concern since:
// (a) The private key is not part of the type, so it is not going to be
//     encoded on the wire anyway, and
// (b) This is a fake identity, no attempt is being made to protect the
//     private key, which is present in this very source file.
type fakeID string

func (id fakeID) Names() []string                             { return []string{string(id)} }
func (id fakeID) Match(s PrincipalPattern) bool               { return matchPrincipalPattern(string(id), s) }
func (fakeID) PublicKey() *ecdsa.PublicKey                    { return &fakeKey.PublicKey }
func (id fakeID) Authorize(context Context) (PublicID, error) { return id, nil }
func (fakeID) ThirdPartyCaveats() []ServiceCaveat             { return nil }
func (id fakeID) PublicID() PublicID                          { return id }
func (fakeID) PrivateKey() *ecdsa.PrivateKey                  { return &fakeKey }
func (id fakeID) Bless(blessee PublicID, blessingName string, duration time.Duration, caveats []ServiceCaveat) (PublicID, error) {
	return fakeID(string(id) + chainSeparator + blessingName), nil
}
func (id fakeID) Derive(pub PublicID) (PrivateID, error) {
	fakePub, ok := pub.(fakeID)
	if !ok {
		return nil, fmt.Errorf("PrivateID of type %T cannot be obtained from PublicID of type %T", id, pub)
	}
	return fakePub, nil
}

func (id fakeID) MintDischarge(caveat ThirdPartyCaveat, ctx Context, duration time.Duration, caveats []ServiceCaveat) (ThirdPartyDischarge, error) {
	return nil, fmt.Errorf("discharge cannot be constructed for ThirdPartyCaveat of type %T from PrivateID of type %T", caveat, id)
}

// FakePublicID returns an implementation of the veyron PublicID interface that
// uses a fixed public key and the provided name.
func FakePublicID(name string) PublicID { return fakeID(fakeName(name)) }

// FakePrivateID returns an implementation of the veyron PrivateID interface
// that uses a fixed private key which is not kept secret and the provided
// name.
func FakePrivateID(name string) PrivateID { return fakeID(fakeName(name)) }

// matchesPattern checks if the provided name conforms to the provided pattern.
// This function assumes pattern to be of the one of the following forms:
// - Pattern is a chained name of the form p_0/.../p_k; in this case the check
//   succeeds iff the provided name is of the form n_0/.../n_m such that m <= k
//   and for all i from 0 to m, p_i = n_i.
// - Pattern is a chained name of the form p_0/.../p_k/*; in this case the check
//   succeeds iff the provided name is of the form n_0/.../n_m such that for all i
//   from 0 to min(m, k), p_i = n_i.
func matchesPattern(name, pattern string) bool {
	patternParts := strings.Split(pattern, chainSeparator)
	patternLen := len(patternParts)
	nameParts := strings.Split(name, chainSeparator)
	nameLen := len(nameParts)

	if patternParts[patternLen-1] != AllPrincipals && nameLen > patternLen {
		return false
	}

	min := nameLen
	if patternParts[patternLen-1] == AllPrincipals && nameLen > patternLen-1 {
		min = patternLen - 1
	}

	for i := 0; i < min; i++ {
		if patternParts[i] != nameParts[i] {
			return false
		}
	}
	return true
}

func matchPrincipalPattern(name string, pattern PrincipalPattern) bool {
	return pattern == AllPrincipals || matchesPattern(name, string(pattern))
}

func generateAndPrintFakeKey() {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal("ERROR:", err)
		return
	}
	var buf bytes.Buffer
	b64 := base64.NewEncoder(base64.URLEncoding, &buf)
	if err := vom.NewEncoder(b64).Encode(priv); err != nil {
		log.Fatal("ERROR:", err)
		return
	}
	b64.Close()
	fmt.Println(buf.String())
}

var fakeKey ecdsa.PrivateKey

func init() {
	// The string was obtained from a call to generateAndPrintFakeKey()
	b64 := []byte("_4EEGgFCAP-DMRgBAgFDAQlQdWJsaWNLZXkAAUQBAUQAARdjcnlwdG8vZWNkc2EuUHJpdmF0ZUtleQD_hTIYAQMBRQEFQ3VydmUAAUQBAVgAAUQBAVkAARZjcnlwdG8vZWNkc2EuUHVibGljS2V5AP-JGxABAQEVY3J5cHRvL2VsbGlwdGljLkN1cnZlAP-HBBoBRgD_ixIQAQQBDG1hdGgvYmlnLkludAD_jS8YAQEBSAELQ3VydmVQYXJhbXMAARljcnlwdG8vZWxsaXB0aWMucDI1NkN1cnZlAP-PBBoBSQD_kU0YAQYBRAEBUAABRAEBTgABRAEBQgABRAECR3gAAUQBAkd5AAEhAQdCaXRTaXplAAEbY3J5cHRvL2VsbGlwdGljLkN1cnZlUGFyYW1zAP-C_gEvAQEB_44BAwEFIQL_____AAAAAQAAAAAAAAAAAAAAAP_______________wEHIQL_____AAAAAP__________vOb6racXnoTzucrC_GMlUQEJIQJaxjXYqjqT57PrvVV2mIa8ZR0GsMxTsPY7zjw-J9JgSwELIQJrF9Hy4SxCR_i85uVjpEDydwN9gS3rM6D0oTlF2JjClgENIQJP40Li_hp_m47n60p8D54WK84zV2sxXs7LtkBoN79R9QH-AgAAAAEPIQLWGT97fh69CMt4-GSGZ1i59_2X66MfiLiZy0xVPw4wUgERIQLy3MnucrogQbqYVDIHZbHBiYvlFttg6aq03cTTsbgqrgABEyEC8Wy60qJXxntBAuK5sf5ejTGuSwZ4ivLgE99rQab2z4MA")
	// Register the elliptic curve type for the encoded public key with VOM.
	vom.Register(elliptic.P256())
	if err := vom.NewDecoder(base64.NewDecoder(base64.URLEncoding, bytes.NewReader(b64))).Decode(&fakeKey); err != nil {
		panic(err)
	}

	var fakeID fakeID
	vom.Register(fakeID)
}
