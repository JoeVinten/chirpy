package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "Lottie"
	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("Error running hash password: %v", err)
	}
	if hash == password {
		t.Errorf("Password should be hashed and not equal to the hash")
	}

	checkedPW, err := CheckPasswordHash(password, hash)

	if err != nil {
		t.Errorf("Error checking the password again the hash: %v", err)
	}

	if checkedPW != true {
		t.Fatalf("Password %s does not match hash %s", password, hash)
	}

}
