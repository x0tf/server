package shared

// Namespace represents a namespace
type Namespace struct {
	ID     string `json:"id"`
	Token  string `json:"token,omitempty"`
	Active bool   `json:"active"`
}

// NamespaceService represents a namespace database service
type NamespaceService interface {
	Namespace(string) (*Namespace, error)
	Namespaces() ([]*Namespace, error)
	CreateOrReplace(*Namespace) error
	Delete(string) error
}
