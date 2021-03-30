package utils

var (
	// elementKeyLength represents the length of an element key
	elementKeyLength = 8

	// elementKeyCharacters represents the characters an element key may contain
	elementKeyCharacters      = "abcdefghijklmnopqrstuvwxyz0123456789"
	elementKeyCharactersRunes = []rune(elementKeyCharacters)
)

// GenerateElementKey generates a new element key
func GenerateElementKey() string {
	return GenerateRandomString(elementKeyLength, elementKeyCharactersRunes)
}
