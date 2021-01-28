package shared

// Namespace represents a namespace
type Namespace struct {
	ID     string
	Token  string
	Active bool
}

// NamespaceService represents a namespace database service
type NamespaceService interface {
	Namespace(string) (*Namespace, error)
	Namespaces() ([]*Namespace, error)
	CreateOrReplace(*Namespace) error
	Delete(string) error
}
