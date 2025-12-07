package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const issuer = "chirpy"

func MakeJWT(userID uuid.UUID, secret string, validFor time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(validFor)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString([]byte(secret))
	return jwt, err
}

func ValidateJWT(jwtString, secret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(jwtString, &jwt.RegisteredClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Failed token parsing %v", err)
	}
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Failed getting sub claim %v", err)
	}
	id, err := uuid.Parse(sub)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Failed uuid parsing %v", err)
	}
	return id, nil
}

func ValidateJWTFromHeader(h http.Header, secret string) (uuid.UUID, error) {
	if tokenString, err := GetBearerToken(h); err != nil {
		return uuid.UUID{}, err
	} else {
		return ValidateJWT(tokenString, secret)
	}
}

func MakeRefreshToken() string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}
