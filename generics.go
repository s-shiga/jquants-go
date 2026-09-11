package jquants

// page is the common wire envelope used by paginated J-Quants endpoints.
type page[T any] struct {
	Data          []T     `json:"data"`
	PaginationKey *string `json:"pagination_key"`
}
