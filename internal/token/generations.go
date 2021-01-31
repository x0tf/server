package token

import (
	"math/rand"
	"time"
)

// tokenLength represents the length of a token
const tokenLength = 64

// tokenCharacters represents the characters a token may contain
var tokenCharacters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.+-#*")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate generates a new token
func Generate() string {
	runes := make([]rune, tokenLength)
	for i := range runes {
		runes[i] = tokenCharacters[rand.Intn(len(tokenCharacters))]
	}
	return string(runes)
}
