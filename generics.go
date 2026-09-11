package jquants

// Response represents a paginated API response. It remains exported for
// compatibility with clients that use it as a generic constraint.
type Response[T any] interface {
	// Items returns the data items from this response.
	Items() []T
	// NextPageKey returns the key for the next page, or nil after the final page.
	NextPageKey() *string
}

// page is the common wire envelope used by paginated J-Quants endpoints.
type page[T any] struct {
	Data          []T     `json:"data"`
	PaginationKey *string `json:"pagination_key"`
}

func (p page[T]) Items() []T           { return p.Data }
func (p page[T]) NextPageKey() *string { return p.PaginationKey }
