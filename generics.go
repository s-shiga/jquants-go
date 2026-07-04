package jquants

// Response represents a generic API response that contains paginated data
type Response[T any] interface {
	// Items returns the data items from this response
	Items() []T
	// NextPageKey returns the pagination key for the next page, or nil if there are no more pages
	NextPageKey() *string
}
