package token

import "github.com/alexedwards/argon2id"

// Hash hashes the given value
func Hash(value string) (string, error) {
	return argon2id.CreateHash(value, argon2id.DefaultParams)
}

// Check compares the given value with the given hash
func Check(hashed, value string) (bool, error) {
	return argon2id.ComparePasswordAndHash(value, hashed)
}
