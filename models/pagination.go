package models

type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewPaginatedResponse[T any](data []T, page, limit, total int) PaginatedResponse[T] {
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}
	return PaginatedResponse[T]{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
