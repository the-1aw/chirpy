package auth

import "github.com/alexedwards/argon2id"

func HashPassword(passwd string) (string, error) {
	return argon2id.CreateHash(passwd, argon2id.DefaultParams)
}

func CheckPasswordHash(passwd, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(passwd, hash)
}
