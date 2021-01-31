package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateRandomString generates a random string with a specific length
func GenerateRandomString(length int, allowedCharacters []rune) string {
	runes := make([]rune, length)
	for i := range runes {
		runes[i] = allowedCharacters[rand.Intn(len(allowedCharacters))]
	}
	return string(runes)
}
