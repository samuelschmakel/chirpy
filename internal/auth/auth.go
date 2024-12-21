package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("couldn't hash the password: %w", err)
	}

	hashedPassword := string(hashedPasswordBytes)

	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	secretKey := []byte(tokenSecret)

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		log.Printf("Token algorithm: %v", token.Header["alg"])
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return uuid.Nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tokenSecret), nil
	})

	// Debug statement:
	log.Printf("Claims: Issuer: %s, IssuedAt: %v, ExpiresAt: %v, Subject: %s", claims.Issuer, claims.IssuedAt, claims.ExpiresAt, claims.Subject)

	// Checking if the token is valid (checks if token is expired)

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return uuid.Nil, fmt.Errorf("error during parsing/validation")
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	id := claims.Subject

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("couldn't parse id to uuid")
	}

	return parsedUUID, nil
}

// This function is used in the testing file
func generateSecretKey(size int) (string, error) {
	// Create a byte slice of the specified size
	bytes := make([]byte, size)

	// Fill the slice with secure random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode the key in base64 for easy storage and usage
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// I should check if the first 7 characters of the string really are "Bearer "!
func GetBearerToken(headers http.Header) (string, error) {
	headerName := "Authorization"

	if values, ok := headers[headerName]; ok {
		fmt.Println(headerName, "found with values: ", values)
		substr := values[0][7:]
		fmt.Println("length of string: ", len(values[0]))
		return substr, nil
	} else {
		return "", fmt.Errorf("header not found")
	}
}
