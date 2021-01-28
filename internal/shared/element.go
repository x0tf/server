package shared

// ElementType represents an element type
type ElementType int

const (
	// ElementTypePaste represents the element type for a paste
	ElementTypePaste = ElementType(0)
)

// Element represents an element published on the service
type Element struct {
	Namespace string
	Key       string
	Type      ElementType
	Data      string
}

// ElementService represents an element database service
type ElementService interface {
	Element(string, string) (*Element, error)
	Elements() ([]*Element, error)
	ElementsInNamespace(string) ([]*Element, error)
	CreateOrReplace(*Element) error
	Delete(string, string) error
}
