package validation

import (
	"strings"
	"unicode/utf8"
)

var (
	// namespaceIDMinimumLength represents the minimum length of a namespace ID
	namespaceIDMinimumLength = 1

	// namespaceIDMaximumLength represents the maximum length of a namespace ID
	namespaceIDMaximumLength = 32

	// namespaceIDAllowedCharacters contains all allowed characters for a namespace ID
	namespaceIDAllowedCharacters = "abcdefghijklmnopqrstuvwxyz0123456789_"
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
	if length < namespaceIDMinimumLength {
		violations = append(violations, NamespaceIDViolationMinimumLength)
	} else if length > namespaceIDMaximumLength {
		violations = append(violations, NamespaceIDViolationMaximumLength)
	}

	// Validate the strings characters
	for _, char := range []rune(id) {
		if !strings.ContainsRune(namespaceIDAllowedCharacters, char) {
			violations = append(violations, NamespaceIDViolationCharacters)
			break
		}
	}
	return
}
