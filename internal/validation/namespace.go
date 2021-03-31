package validation

import (
	"strings"
	"unicode/utf8"
)

var (
	// NamespaceIDMinimumLength represents the minimum length of a namespace ID
	NamespaceIDMinimumLength = 1

	// NamespaceIDMaximumLength represents the maximum length of a namespace ID
	NamespaceIDMaximumLength = 32

	// NamespaceIDAllowedCharacters contains all allowed characters for a namespace ID
	NamespaceIDAllowedCharacters = "abcdefghijklmnopqrstuvwxyz0123456789_"
)

// NamespaceIDViolation represents a violation of given namespace ID rules
type NamespaceIDViolation string

const (
	// NamespaceIDViolationMinimumLength is used when a namespace ID is too short
	NamespaceIDViolationMinimumLength = NamespaceIDViolation("MINIMUM_LENGTH")

	// NamespaceIDViolationMaximumLength is used when a namespace ID is too long
	NamespaceIDViolationMaximumLength = NamespaceIDViolation("MAXIMUM_LENGTH")

	// NamespaceIDViolationCharacters is used when a namespace ID contains illegal characters
	NamespaceIDViolationCharacters = NamespaceIDViolation("CHARACTERS")
)

// ValidateNamespaceID validates a given namespace ID
func ValidateNamespaceID(id string) (violations []NamespaceIDViolation) {
	// Validate the length of the ID
	length := utf8.RuneCountInString(id)
	if length < NamespaceIDMinimumLength {
		violations = append(violations, NamespaceIDViolationMinimumLength)
	} else if length > NamespaceIDMaximumLength {
		violations = append(violations, NamespaceIDViolationMaximumLength)
	}

	// Validate the strings characters
	for _, char := range []rune(id) {
		if !strings.ContainsRune(NamespaceIDAllowedCharacters, char) {
			violations = append(violations, NamespaceIDViolationCharacters)
			break
		}
	}
	return
}
