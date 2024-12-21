package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secretKey, err := generateSecretKey(32)
	if err != nil {
		t.Errorf("couldn't generate secret key")
	}
	fmt.Printf("Here's the secret key: %s/n", secretKey)

	expiresIn := 30 * time.Second
	tokenString, err := MakeJWT(userID, secretKey, expiresIn)
	if err != nil {
		t.Errorf("error making the JWT")
	}
	fmt.Println(tokenString)
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secretKey, err := generateSecretKey(32)
	if err != nil {
		t.Errorf("couldn't generate secret key")
	}

	expiresIn := 30 * time.Second
	tokenString, err := MakeJWT(userID, secretKey, expiresIn)
	if err != nil {
		t.Errorf("error making the JWT")
	}

	id, err := ValidateJWT(tokenString, secretKey)
	if err != nil {
		t.Errorf("error validating the JWT")
	}
	fmt.Printf("Here's the id: %s\n", id.String())

	// Checking with wrong secret
	wrongSecretKey, err := generateSecretKey(32)
	if err != nil {
		t.Errorf("couldn't generate secret key")
	}

	_, err = ValidateJWT(tokenString, wrongSecretKey)
	if err == nil {
		t.Errorf("accepted wrong secret key for validation")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer token123"},
	}
	str, err := GetBearerToken(headers)
	if err != nil {
		t.Errorf("error getting bearer")
	}

	if str != "token123" {
		fmt.Println(str)
		t.Errorf("unexpected return value of GetBearerToken")
	}

	headers = http.Header{
		"Content-Type":     []string{"application/json"},
		"NotAuthorization": []string{"Bearer token123"},
	}

	_, err = GetBearerToken(headers)
	if err == nil {
		t.Errorf("expected error but got no error")
	}
}
