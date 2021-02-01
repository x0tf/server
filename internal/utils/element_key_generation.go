package utils

// elementKeyLength represents the length of an element key
var elementKeyLength = 8

// elementKeyCharacters represents the characters an element key may contain
var elementKeyCharacters = "abcdefghijklmnopqrstuvwxyz0123456789"

var elementKeyCharactersRunes = []rune(elementKeyCharacters)

// GenerateElementKey generates a new element key
func GenerateElementKey() string {
	return GenerateRandomString(elementKeyLength, elementKeyCharactersRunes)
}
