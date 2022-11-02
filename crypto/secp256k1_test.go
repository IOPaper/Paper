package crypto

import (
	"bytes"
	"testing"
)

func TestNewSecp256k1(t *testing.T) {
	k1, err := NewSecp256k1()
	if err != nil {
		t.Fatal(err)
	}
	sign, err := k1.Sign(bytes.NewBufferString("TEST"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("verify sign: %v", k1.Verify(bytes.NewBufferString("TEST"), sign))
}
