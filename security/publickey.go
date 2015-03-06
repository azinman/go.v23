package security

import (
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/x509"
	"encoding"
	"fmt"
	"v.io/v23/verror"
)

var (
	errUnrecognizedKey = verror.Register(pkgPath+".errUnrecognizedKey", verror.NoRetry, "{1:}{2:}unrecognized PublicKey type({3}){:_}")
)

// PublicKey represents a public key using an unspecified algorithm.
//
// MarshalBinary returns the DER-encoded PKIX representation of the public key,
// while UnmarshalPublicKey creates a PublicKey object from the marshaled bytes.
//
// String returns a human-readable representation of the public key.
type PublicKey interface {
	encoding.BinaryMarshaler
	fmt.Stringer

	// hash returns a cryptographic hash function whose security strength is
	// appropriate for creating message digests to sign with this public key.
	// For example, an ECDSA public key with a 512-bit curve would require a
	// 512-bit hash function, whilst a key with a 256-bit curve would be
	// happy with a 256-bit hash function.
	hash() Hash
	implementationsOnlyInThisPackage()
}

type ecdsaPublicKey struct {
	key *ecdsa.PublicKey
}

func (pk *ecdsaPublicKey) MarshalBinary() ([]byte, error)    { return x509.MarshalPKIXPublicKey(pk.key) }
func (pk *ecdsaPublicKey) String() string                    { return publicKeyString(pk) }
func (pk *ecdsaPublicKey) implementationsOnlyInThisPackage() {}
func (pk *ecdsaPublicKey) hash() Hash {
	if nbits := pk.key.Curve.Params().BitSize; nbits <= 160 {
		return SHA1Hash
	} else if nbits <= 256 {
		return SHA256Hash
	} else if nbits <= 384 {
		return SHA384Hash
	} else {
		return SHA512Hash
	}
}

func publicKeyString(pk PublicKey) string {
	bytes, err := pk.MarshalBinary()
	if err != nil {
		return fmt.Sprintf("<invalid public key: %v>", err)
	}
	const hextable = "0123456789abcdef"
	hash := md5.Sum(bytes)
	var repr [md5.Size * 3]byte
	for i, v := range hash {
		repr[i*3] = hextable[v>>4]
		repr[i*3+1] = hextable[v&0x0f]
		repr[i*3+2] = ':'
	}
	return string(repr[:len(repr)-1])
}

// UnmarshalPublicKey returns a PublicKey object from the DER-encoded PKIX represntation of it
// (typically obtianed via PublicKey.MarshalBinary).
func UnmarshalPublicKey(bytes []byte) (PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return nil, err
	}
	switch v := key.(type) {
	case *ecdsa.PublicKey:
		return &ecdsaPublicKey{v}, nil
	default:
		return nil, verror.New(errUnrecognizedKey, nil, fmt.Sprintf("%T", key))
	}
}

// NewECDSAPublicKey creates a PublicKey object that uses the ECDSA algorithm and the provided ECDSA public key.
func NewECDSAPublicKey(key *ecdsa.PublicKey) PublicKey {
	return &ecdsaPublicKey{key}
}
