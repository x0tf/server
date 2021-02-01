package utils

// tokenLength represents the length of a token
var tokenLength = 64

// tokenCharacters represents the characters a token may contain
var tokenCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.+-#*"

var tokenCharactersRunes = []rune(tokenCharacters)

// GenerateToken generates a new token
func GenerateToken() string {
	return GenerateRandomString(tokenLength, tokenCharactersRunes)
}
