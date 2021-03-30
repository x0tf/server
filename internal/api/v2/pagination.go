package v2

// paginatedResponse represents a paginated API response
type paginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination pagination  `json:"pagination"`
}

// pagination holds metadata for a paginated API response
type pagination struct {
	TotalElements     int `json:"total_elements"`
	DisplayedElements int `json:"displayed_elements"`
}
