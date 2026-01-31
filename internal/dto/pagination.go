package dto

// PaginationRequest represents pagination request parameters
type PaginationRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1" example:"1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100" example:"10"`
}

// GetOffset calculates the offset for pagination
func (p *PaginationRequest) GetOffset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
	return (p.Page - 1) * p.Limit
}

// Normalize sets default values if not provided
func (p *PaginationRequest) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}
