package shared

// Invite represents an invite
type Invite struct {
	Code    string `json:"code"`
	Uses    int    `json:"uses"`
	MaxUses int    `json:"max_uses"`
	Created int64  `json:"created"`
}

// InviteService represents an invite database service
type InviteService interface {
	Invite(code string) (*Invite, error)
	Invites(limit, offset int) ([]*Invite, error)
	CreateOrReplace(invite *Invite) error
	Delete(code string) error
}
