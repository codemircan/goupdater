package utils

import (
	"bytes"
	"testing"
)

func TestXOR(t *testing.T) {
	data := []byte("Hello World")
	key := []byte("key")

	encrypted := XOR(data, key)
	decrypted := XOR(encrypted, key)

	if !bytes.Equal(data, decrypted) {
		t.Errorf("Expected %s, got %s", string(data), string(decrypted))
	}

	if bytes.Equal(data, encrypted) {
		t.Errorf("Data should be encrypted")
	}
}
