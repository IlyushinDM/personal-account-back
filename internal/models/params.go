package models

// PaginationParams содержит параметры для пагинации и сортировки.
type PaginationParams struct {
	Page      int
	Limit     int
	SortBy    string // Например: "rating", "experience_years"
	SortOrder string // "asc" или "desc"
}
