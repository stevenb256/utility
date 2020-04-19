package utl

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	_l "github.com/stevenb256/log"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

// NonceSize size of nonce
const NonceSize = 24

// KeySize lengtho of crypto key
const KeySize = 32

// Key is a crypto key
type Key *[KeySize]byte

// ErrInvalidCryptoKey invalid crypto key
var ErrInvalidCryptoKey = _l.NewError(100, "crypto", "invalid crypto key length")

// ErrCantOpenSealedBytes can't open sealed bytes; crypto problem
var ErrCantOpenSealedBytes = _l.NewError(101, "crypto", "unable to open/unseal bytes")

// ErrCantDecryptBytes can't decrypt bytes encrypted in other function
var ErrCantDecryptBytes = _l.NewError(102, "crypto", "unable to decrypt bytes")

// NewKey - returns new crypto key from a buffer
func NewKey(buf []byte) Key {
	if len(buf) != KeySize {
		return nil
	}
	key := new([KeySize]byte)
	copy(key[:], buf[:KeySize])
	return key
}

// KeyToBase64 - converts key to base64 encoding
func KeyToBase64(key Key) string {
	return base64.StdEncoding.EncodeToString(key[:])
}

// GenerateCryptoKeys returns public, private keys or error
func GenerateCryptoKeys() (*[KeySize]byte, *[KeySize]byte, error) {
	return box.GenerateKey(rand.Reader)
}

// CryptoKeyFromBase64 get crypto key from base64 string
func CryptoKeyFromBase64(key64 string) (Key, error) {
	buf, err := base64.StdEncoding.DecodeString(key64)
	if _l.Check(err) {
		return nil, err
	}
	if len(buf) != KeySize {
		return nil, _l.Fail(ErrInvalidCryptoKey, key64)
	}
	var key [KeySize]byte
	copy((key)[:], buf)
	return &key, nil
}

// CryptoKeyToBase64 converts a crypto key to base64
func CryptoKeyToBase64(key *[KeySize]byte) string {
	return base64.StdEncoding.EncodeToString(key[:])
}

// SealBytes encrypts/signs buffer with a public key of recipient and private key
// of the sender
func SealBytes(buf []byte, public, private Key) ([]byte, error) {
	var nonce [NonceSize]byte
	io.ReadFull(rand.Reader, nonce[:])
	return box.Seal(nonce[:], buf, &nonce, public, private), nil
}

// OpenSealedBytes - decrypts bytes with public key of the sender and
// private key of the recipient
func OpenSealedBytes(buf []byte, public, private Key) ([]byte, error) {

	// locals
	var nonce [NonceSize]byte

	// check args
	if nil == public {
		return nil, _l.Fail(_l.ErrInvalidArg, "nil public key")
	}
	if nil == private {
		return nil, _l.Fail(_l.ErrInvalidArg, "nil private key")
	}
	if len(buf) < len(nonce) {
		return nil, _l.Fail(_l.ErrInvalidArg, "sealed buffer smaller than nonce")
	}

	// copy over nonce
	copy(nonce[:], buf[:NonceSize])

	// open it
	clear, b := box.Open(nil, buf[NonceSize:], &nonce, public, private)
	if false == b {
		return nil, _l.Fail(ErrCantOpenSealedBytes)
	}

	// done
	return clear, nil
}

// EncryptBytes used to just encrypt bytes with a random key
func EncryptBytes(in []byte, key Key) ([]byte, error) {
	var nonce [NonceSize]byte
	io.ReadFull(rand.Reader, nonce[:])
	out := make([]byte, NonceSize)
	copy(out, nonce[:])
	return secretbox.Seal(out, in, &nonce, key), nil
}

// DecryptBytes used to decrypt bytes with a key used in EncryptBytes
func DecryptBytes(in []byte, key Key) ([]byte, error) {
	var nonce [NonceSize]byte
	_l.Assert(len(in) >= NonceSize)
	copy(nonce[:], in[:NonceSize])
	out, worked := secretbox.Open(nil, in[NonceSize:], &nonce, key)
	if false == worked {
		return nil, _l.Fail(ErrCantDecryptBytes)
	}
	return out, nil
}
