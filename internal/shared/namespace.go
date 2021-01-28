package shared

// Namespace represents a namespace
type Namespace struct {
	ID     string
	Token  string
	Active bool
}

// NamespaceService represents a namespace storage service
type NamespaceService interface {
	Namespace(id string) (*Namespace, error)
	Namespaces() ([]*Namespace, error)
	Create(namespace *Namespace) error
	Delete(id string)
}
