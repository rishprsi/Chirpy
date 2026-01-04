package auth

import (
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	registeredClaims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}
	regClaims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, regClaims, keyFunc)
	if err != nil {
		log.Println("Failed in parsing claims", tokenString)
		return uuid.UUID{}, err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		log.Println("Failed in getting subject")
		return uuid.UUID{}, err
	}
	uuidID, err := uuid.Parse(id)
	if err != nil {
		log.Println("Failed parsing string to uuid")
		return uuid.UUID{}, err
	}
	return uuidID, nil
}
