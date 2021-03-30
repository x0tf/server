package utils

var (
	// tokenLength represents the length of a token
	tokenLength = 64

	// tokenCharacters represents the characters a token may contain
	tokenCharacters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.+-#*"
	tokenCharactersRunes = []rune(tokenCharacters)
)

// GenerateToken generates a new token
func GenerateToken() string {
	return GenerateRandomString(tokenLength, tokenCharactersRunes)
}
