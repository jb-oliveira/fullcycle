package entity

import "math"

type Meta struct {
	CurrentPage   int    `json:"current_page"`
	PageSize      int    `json:"page_size"`
	TotalItems    int    `json:"total_items"`
	TotalPages    int    `json:"total_pages"`
	SortField     string `json:"sort_field"`
	SortDirection string `json:"sort_direction"`
}

// Page is the generic struct.
// T is the type parameter (e.g., User, Product).
type Page[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

// NewPage is a helper constructor to calculate total pages automatically
func NewPage[T any](items []T, page, pageSize, totalItems int, sortFld, sortDir string) Page[T] {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	return Page[T]{
		Data: items,
		Meta: Meta{
			CurrentPage:   page,
			PageSize:      pageSize,
			TotalItems:    totalItems,
			TotalPages:    totalPages,
			SortField:     sortFld,
			SortDirection: sortDir,
		},
	}
}
