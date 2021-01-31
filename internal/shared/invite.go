package shared

// Invite represents an invite token
type Invite string

// InviteService represents an invite database service
type InviteService interface {
	IsValid(string) (bool, error)
	Invites() ([]Invite, error)
	Create(Invite) error
	Delete(string) error
}
