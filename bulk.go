package jquants

import (
	"context"
	"encoding/json"
	"net/url"
)

// BulkFile represents a single downloadable whole-dataset CSV file listed by
// the /bulk/list endpoint. The file is a gzip-compressed CSV and is retrieved
// via a signed URL obtained from BulkGet.
type BulkFile struct {
	// Key is the object key identifying the file (e.g.
	// "derivatives/bars/daily/futures/live/derivatives_bars_daily_futures_20260717.csv.gz").
	// It is passed to BulkGet to obtain a signed download URL.
	Key string
	// LastModified is the last-modified timestamp in RFC 3339 format
	// (e.g. "2026-07-17T12:13:51+00:00").
	LastModified string
	// Size is the file size in bytes.
	Size int64
}

func (bf *BulkFile) UnmarshalJSON(b []byte) error {
	var raw struct {
		Key          string  `json:"Key"`
		LastModified string  `json:"LastModified"`
		Size         float64 `json:"Size"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	bf.Key = raw.Key
	bf.LastModified = raw.LastModified
	bf.Size = int64(raw.Size)
	return nil
}

// BulkListRequest specifies filter parameters for the BulkList API.
// Either Endpoint or Date must be provided (but not both). From and To are
// only valid together with Endpoint.
type BulkListRequest struct {
	// Endpoint filters files by the originating API endpoint (e.g. "/equities/bars/daily").
	Endpoint *string
	// Date filters files by date. Accepts YYYY-MM, YYYYMM, YYYY-MM-DD, or YYYYMMDD.
	Date *string
	// From specifies the start date for a date range query (used with Endpoint).
	From *string
	// To specifies the end date for a date range query (used with Endpoint).
	To *string
}

type bulkListParameters struct {
	BulkListRequest
}

func (p bulkListParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Endpoint != nil {
		v.Add("endpoint", *p.Endpoint)
	}
	if p.Date != nil {
		v.Add("date", *p.Date)
	}
	if p.From != nil {
		v.Add("from", *p.From)
	}
	if p.To != nil {
		v.Add("to", *p.To)
	}
	return v, nil
}

type bulkListResponse struct {
	Data []BulkFile `json:"data"`
}

// BulkList retrieves the list of whole-dataset CSV files available for bulk
// download from the /bulk/list endpoint. This endpoint is not paginated.
// Either Endpoint or Date must be provided (but not both); From and To are only
// valid together with Endpoint.
// See https://jpx-jquants.com/en/spec/bulk-list for API details.
func (c *Client) BulkList(ctx context.Context, req BulkListRequest) ([]BulkFile, error) {
	r, err := getJSON[bulkListResponse](ctx, c, "/bulk/list", bulkListParameters{req})
	if err != nil {
		return nil, err
	}
	return r.Data, nil
}

// BulkGetRequest specifies parameters for the BulkGet API.
// Either Key must be provided alone, or Endpoint and Date must be provided together.
type BulkGetRequest struct {
	// Key is a file key returned by BulkList. If set, Endpoint and Date are ignored.
	Key *string
	// Endpoint identifies the originating API endpoint (used together with Date).
	Endpoint *string
	// Date identifies the file date (used together with Endpoint).
	Date *string
}

type bulkGetParameters struct {
	BulkGetRequest
}

func (p bulkGetParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Key != nil {
		v.Add("key", *p.Key)
	}
	if p.Endpoint != nil {
		v.Add("endpoint", *p.Endpoint)
	}
	if p.Date != nil {
		v.Add("date", *p.Date)
	}
	return v, nil
}

type bulkGetResponse struct {
	URL string `json:"url"`
}

// BulkGet retrieves a signed download URL for a whole-dataset CSV file from the
// /bulk/get endpoint. Provide either Key alone, or Endpoint and Date together.
// The returned URL is valid for approximately 5 minutes and yields a
// gzip-compressed CSV file. Tick-level stock trades (available with the add-on
// plan) are delivered exclusively via keys under "equities/trades/".
// See https://jpx-jquants.com/en/spec/bulk-get for API details.
func (c *Client) BulkGet(ctx context.Context, req BulkGetRequest) (string, error) {
	r, err := getJSON[bulkGetResponse](ctx, c, "/bulk/get", bulkGetParameters{req})
	if err != nil {
		return "", err
	}
	return r.URL, nil
}
