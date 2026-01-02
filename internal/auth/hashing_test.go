package auth

import "testing"

func TestValidHash(t *testing.T) {
	hash, err := HashPassword("pasword")
	if err != nil {
		t.Fatal(err)
	}
	success, err := CheckPasswordAndHash("password", hash)
	if err != nil {
		t.Fatal(err)
	}
	if !success {
		t.Errorf("Password should be valid but invalid")
	}
}

func TestInvalidHash(t *testing.T) {
	hash, err := HashPassword("pasword")
	if err != nil {
		t.Fatal(err)
	}
	success, err := CheckPasswordAndHash("p@ssword", hash)
	if err != nil {
		t.Fatal(err)
	}
	if success {
		t.Errorf("Password should be valid but invalid")
	}
}
