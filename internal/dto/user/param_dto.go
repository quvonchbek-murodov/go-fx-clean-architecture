package user

type CreateParams struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateParams struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type ListParams struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (p *ListParams) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 || p.Limit > 100 {
		p.Limit = 20
	}
}

func (p ListParams) Offset() int {
	return (p.Page - 1) * p.Limit
}
