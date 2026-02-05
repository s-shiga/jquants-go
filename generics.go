package jquants

import (
	"context"
	"net/url"
)

// Request represents a generic API request that can be sent to the API
type Request[T any] interface {
	// Send returns a set of data
	Send(context.Context, *Client) ([]T, error)
	// Path returns the API endpoint path (e.g., "/equities/bars/daily")
	Path() string
	// Values returns the URL query parameters for this request
	Values() (url.Values, error)
	// SetPaginationKey sets a pagination key
	SetPaginationKey(*string)
}

// Response represents a generic API response that contains paginated data
type Response[T any] interface {
	// Items returns the data items from this response
	Items() []T
	// NextPageKey returns the pagination key for the next page, or nil if there are no more pages
	NextPageKey() *string
}
