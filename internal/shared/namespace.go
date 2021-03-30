package shared

// Namespace represents a namespace
type Namespace struct {
	ID      string `json:"id"`
	Token   string `json:"token,omitempty"`
	Active  bool   `json:"active"`
	Created int64  `json:"created"`
}

// NamespaceService represents a namespace database service
type NamespaceService interface {
	Count() (int, error)
	Namespace(id string) (*Namespace, error)
	Namespaces(limit, offset int) ([]*Namespace, error)
	CreateOrReplace(namespace *Namespace) error
	Delete(id string) error
}
