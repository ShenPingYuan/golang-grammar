package pagination

type Params struct {
	Page int
	Size int
}

type Result struct {
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
	Items      interface{} `json:"items"`
}

func NewParams(page, size int) Params {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return Params{Page: page, Size: size}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.Size
}

func NewResult(params Params, total int, items interface{}) Result {
	totalPages := total / params.Size
	if total%params.Size > 0 {
		totalPages++
	}
	return Result{
		Page:       params.Page,
		Size:       params.Size,
		Total:      total,
		TotalPages: totalPages,
		Items:      items,
	}
}