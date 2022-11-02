package crypto

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
)

type Fingerprint []byte

func NewFingerprint(src io.Reader) Fingerprint {
	hash256 := sha256.New()
	io.Copy(hash256, src)
	return hash256.Sum(nil)
}

func (f Fingerprint) Hex() string {
	return strings.ToUpper(hex.EncodeToString(f[:len(f)/2]))
}

func (f Fingerprint) Base64() string {
	return base64.StdEncoding.EncodeToString(f[:len(f)/2])
}

func (f Fingerprint) Bytes() []byte {
	return f[:len(f)/2]
}

func (f Fingerprint) Equal(fp []byte) bool {
	return subtle.ConstantTimeCompare(f, fp) == 1
}
