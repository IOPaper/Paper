package crypto

import (
	"bytes"
	"testing"
)

func TestNewFingerprint(t *testing.T) {
	t.Logf("size: %s", NewFingerprint(bytes.NewBufferString("TEST")).Base64())
}
