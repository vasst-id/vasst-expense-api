package entities

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// SetDefaults sets default values for pagination
func (p *PaginationRequest) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 20
	}
}

// GetOffset calculates the offset for database queries
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// FilterRequest represents base filtering parameters
type FilterRequest struct {
	Search    string `json:"search" form:"search"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// SetDefaults sets default values for filtering
func (f *FilterRequest) SetDefaults() {
	if f.SortOrder == "" {
		f.SortOrder = "desc"
	}
}
