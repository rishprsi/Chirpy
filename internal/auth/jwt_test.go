package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidUser(t *testing.T) {
	id := uuid.New()
	token, err := MakeJWT(id, "willthiswork", 60*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	userID, err := ValidateJWT(token, "willthiswork")
	if err != nil {
		t.Fatal(err)
	} else if id != userID {
		t.Fatal("User IDs don't match")
	}
}

func TestInValidSecret(t *testing.T) {
	id := uuid.New()
	token, err := MakeJWT(id, "willthiswork", 60*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	userID, err := ValidateJWT(token, "thishouldnotwork")
	if err == nil {
		t.Fatal(err, userID)
	}
}

func TestTimedOutSecret(t *testing.T) {
	id := uuid.New()
	token, err := MakeJWT(id, "willthiswork", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	userID, err := ValidateJWT(token, "willthiswork")
	if err == nil {
		t.Fatal(err, userID)
	}
}
