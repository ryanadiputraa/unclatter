package pagination

import "strconv"

const (
	defaultPage = 1
	defaultSize = 20
)

type Meta struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	Size        int   `json:"size"`
	TotalData   int64 `json:"total_data"`
}

type Pagination struct {
	Limit  int
	Offset int
}

func NewPagination(page, size int) *Pagination {
	return &Pagination{
		Limit:  size,
		Offset: (page - 1) * size,
	}
}

func ValidateParam(pageParam, sizeParam string) (pagination *Pagination, errDetail map[string]string, err error) {
	var page int
	var size int
	errDetail = make(map[string]string)

	if len(pageParam) == 0 {
		page = defaultPage
	} else {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			errDetail["page"] = "invalid 'page' param expecting int"
		}
	}
	if len(sizeParam) == 0 {
		size = defaultSize
	} else {
		size, err = strconv.Atoi(sizeParam)
		if err != nil {
			errDetail["size"] = "invalid 'size' param expecting int"
		}
	}

	return NewPagination(page, size), errDetail, err
}
