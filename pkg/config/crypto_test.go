package config

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	pass := "test-master-password-123"
	plaintext := []byte("[Interface]\nPrivateKey = abc123\n")

	ciphertext, err := Encrypt(plaintext, pass)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Error("ciphertext should differ from plaintext")
	}

	decrypted, err := Decrypt(ciphertext, pass)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Error("decrypted text does not match original")
	}
}

func TestWrongPassword(t *testing.T) {
	plaintext := []byte("secret data")
	ciphertext, err := Encrypt(plaintext, "correct-pass")
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(ciphertext, "wrong-pass")
	if err == nil {
		t.Error("expected error with wrong password")
	}
}
