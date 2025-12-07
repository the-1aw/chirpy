package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const issuer = "chirpy"
const authorizationPrefix = "bearer"

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

func GetBearerToken(h http.Header) (string, error) {
	authHeader := h.Get("authorization")
	parts := strings.Fields(authHeader)
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid authorization header")
	}
	if strings.ToLower(parts[0]) != authorizationPrefix {
		return "", fmt.Errorf("Invalid authorization header")
	}
	return parts[1], nil
}

func MakeRefreshToken() string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}
