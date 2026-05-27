package domain

type PaginationRequest struct {
    Page  int `form:"page"`
    Limit int `form:"limit"`
}

type PaginationResponse struct {
    Total       int64       `json:"total"`
    Page        int         `json:"page"`
    Limit       int         `json:"limit"`
    TotalPages  int         `json:"total_pages"`
    Data        interface{} `json:"data"`
}

func (p *PaginationRequest) Normalize() {
    if p.Page <= 0 {
        p.Page = 1
    }
    if p.Limit <= 0 || p.Limit > 100 {
        p.Limit = 10
    }
}

func (p *PaginationRequest) Offset() int {
    return (p.Page - 1) * p.Limit
}