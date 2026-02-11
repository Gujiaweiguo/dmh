package utils

import "testing"

func TestHashPasswordAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("P@ssw0rd")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if hash == "" {
		t.Fatalf("hash empty")
	}
	if !CheckPassword("P@ssw0rd", hash) {
		t.Fatalf("expected password match")
	}
	if CheckPassword("wrong", hash) {
		t.Fatalf("expected password mismatch")
	}
}
