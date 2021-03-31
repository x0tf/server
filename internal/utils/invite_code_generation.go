package utils

var (
	// inviteCodeLength represents the length of an invite code
	inviteCodeLength = 32

	// inviteCodeCharacters represents the characters an invite code may contain
	inviteCodeCharacters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.+-#*"
	inviteCodeCharactersRunes = []rune(inviteCodeCharacters)
)

// GenerateInviteCode generates a new invite code
func GenerateInviteCode() string {
	return GenerateRandomString(inviteCodeLength, inviteCodeCharactersRunes)
}
