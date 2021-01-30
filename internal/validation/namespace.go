package validation

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	// namespaceIDMinimumLength represents the minimum length of a namespace ID
	namespaceIDMinimumLength = 1

	// namespaceIDMaximumLength represents the maximum length of a namespace ID
	namespaceIDMaximumLength = 32

	// namespaceIDAllowedCharacters contains all allowed characters for a namespace ID
	namespaceIDAllowedCharacters = "abcdefghijklmnopqrstuvwxyz0123456789_"
)

var (
	// ErrNamespaceIDTooShort is used when a namespace ID is too short
	ErrNamespaceIDTooShort = fmt.Errorf("the given namespace ID is too short (minimum is %d)", namespaceIDMinimumLength)

	// ErrNamespaceIDTooShort is used when a namespace ID is too long
	ErrNamespaceIDTooLong = fmt.Errorf("the given namespace ID is too long (maximum is %d)", namespaceIDMaximumLength)

	// ErrNamespaceIDTooShort is used when a namespace ID contains at least one illegal character
	ErrNamespaceIDContainsIllegalCharacter = fmt.Errorf("the given namespace ID contains an illegal character (allowed are '%s')", namespaceIDAllowedCharacters)
)

// ValidateNamespaceID validates a given namespace ID
func ValidateNamespaceID(id string) (errors []error) {
	// Validate the length of the ID
	length := utf8.RuneCountInString(id)
	if length < namespaceIDMinimumLength {
		errors = append(errors, ErrNamespaceIDTooShort)
	} else if length > namespaceIDMaximumLength {
		errors = append(errors, ErrNamespaceIDTooLong)
	}

	// Validate the strings characters
	for _, char := range []rune(id) {
		if !strings.ContainsRune(namespaceIDAllowedCharacters, char) {
			errors = append(errors, ErrNamespaceIDContainsIllegalCharacter)
			break
		}
	}
	return
}
