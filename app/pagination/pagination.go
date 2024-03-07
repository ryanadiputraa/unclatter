package pagination

const (
	DefaultPage = 1
	DefaultSize = 20
)

type Meta struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	Size        int `json:"size"`
	TotalData   int `json:"total_data"`
}

type Pagination struct {
	Limit  int
	Offset int
}

type PaginationParam struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func NewPagination(arg PaginationParam) *Pagination {
	return &Pagination{
		Limit:  arg.Size,
		Offset: (arg.Page - 1) * arg.Size,
	}
}
