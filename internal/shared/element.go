package shared

// ElementType represents an element type
type ElementType string

const (
	// ElementTypePaste represents the element type for a paste
	ElementTypePaste = ElementType("PASTE")

	// ElementTypeRedirect represents the element type for a redirect
	ElementTypeRedirect = ElementType("REDIRECT")
)

// Element represents an element published on the service
type Element struct {
	Namespace    string                 `json:"namespace"`
	Key          string                 `json:"key"`
	Type         ElementType            `json:"type"`
	InternalData map[string]interface{} `json:"internal_data,omitempty"`
	PublicData   map[string]interface{} `json:"public_data"`
	Views        int                    `json:"views"`
	MaxViews     int                    `json:"max_views"`
	ValidFrom    int64                  `json:"valid_from"`
	ValidUntil   int64                  `json:"valid_until"`
	Created      int64                  `json:"created"`
}

// ElementService represents an element database service
type ElementService interface {
	Element(namespace, key string) (*Element, error)
	Elements(limit, offset int) ([]*Element, error)
	ElementsInNamespace(namespace string, limit, offset int) ([]*Element, error)
	CreateOrReplace(element *Element) error
	Delete(namespace, key string) error
	DeleteInNamespace(namespace string) error
}
