package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashPassword(passwd string) (string, error) {
	return argon2id.CreateHash(passwd, argon2id.DefaultParams)
}

func CheckPasswordHash(passwd, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(passwd, hash)
}

const authorizationBearerPrefix = "bearer"
const authorizationApiKeyPrefix = "apikey"
const authorizationHeaderName = "authorization"

func GetApiKey(h http.Header) (string, error) {
	authHeader := h.Get(authorizationHeaderName)
	parts := strings.Fields(authHeader)
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid authorization header")
	}
	if strings.ToLower(parts[0]) != authorizationApiKeyPrefix {
		return "", fmt.Errorf("Invalid authorization header")
	}
	return parts[1], nil
}

func GetBearerToken(h http.Header) (string, error) {
	authHeader := h.Get(authorizationHeaderName)
	parts := strings.Fields(authHeader)
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid authorization header")
	}
	if strings.ToLower(parts[0]) != authorizationBearerPrefix {
		return "", fmt.Errorf("Invalid authorization header")
	}
	return parts[1], nil
}
