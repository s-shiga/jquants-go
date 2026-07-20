package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// TimelyDisclosure represents a single TDnet timely disclosure index record.
// It describes a disclosure document published by a listed company, including
// its identifiers, timing, title, and the set of attached document types.
type TimelyDisclosure struct {
	// DisclosureNumber is the 14-digit disclosure number (JSON key "DiscNo").
	DisclosureNumber string
	// Code is the security code (JSON key "Code").
	Code string
	// CompanyName is the company name (JSON key "Name").
	CompanyName string
	// DisclosureDate is the disclosure date in YYYY-MM-DD format (JSON key "DiscDate").
	DisclosureDate string
	// DisclosureTime is the disclosure time in HH:MM format (JSON key "DiscTime").
	DisclosureTime string
	// Title is the disclosure title (JSON key "Title").
	Title string
	// DisclosureStatus is the status of the disclosure (JSON key "DiscStatus").
	// It is nil for new disclosures, "revision" for corrected disclosures, and
	// "delete" for deleted disclosures.
	DisclosureStatus *string
	// RevisionNumber is the revision number, returned by the API as a JSON string (JSON key "RevNo").
	RevisionNumber string
	// DisclosureItems holds the public item codes classifying the disclosure (JSON key "DiscItems").
	DisclosureItems []string
	// Documents holds the available document type codes: "g" full PDF, "s" summary PDF,
	// "x" XBRL (JSON key "Docs").
	Documents []string
}

func (td *TimelyDisclosure) UnmarshalJSON(b []byte) error {
	var raw struct {
		DiscNo     string   `json:"DiscNo"`
		Code       string   `json:"Code"`
		Name       string   `json:"Name"`
		DiscDate   string   `json:"DiscDate"`
		DiscTime   string   `json:"DiscTime"`
		Title      string   `json:"Title"`
		DiscStatus *string  `json:"DiscStatus"`
		RevNo      string   `json:"RevNo"`
		DiscItems  []string `json:"DiscItems"`
		Docs       []string `json:"Docs"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal timely disclosure: %w", err)
	}
	td.DisclosureNumber = raw.DiscNo
	td.Code = raw.Code
	td.CompanyName = raw.Name
	td.DisclosureDate = raw.DiscDate
	td.DisclosureTime = raw.DiscTime
	td.Title = raw.Title
	td.DisclosureStatus = raw.DiscStatus
	td.RevisionNumber = raw.RevNo
	td.DisclosureItems = raw.DiscItems
	td.Documents = raw.Docs
	return nil
}

// TimelyDisclosureRequest specifies filter parameters for the TimelyDisclosure API.
// Either Date or Code must be provided. From and To narrow the date range when
// querying by Code and must be specified together. DiscItems is an optional
// comma-separated list of public item codes filtered with AND semantics.
type TimelyDisclosureRequest struct {
	// Date filters by disclosure date in YYYY-MM-DD format.
	Date *string
	// Code filters by security code.
	Code *string
	// From is the start of the disclosure date range (used with Code).
	From *string
	// To is the end of the disclosure date range (used with Code).
	To *string
	// DiscItems is a comma-separated list of public item codes (AND filter).
	DiscItems *string
}

type timelyDisclosureParameters struct {
	TimelyDisclosureRequest
	PaginationKey *string
}

func (p timelyDisclosureParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Date != nil {
		v.Add("date", *p.Date)
	}
	if p.Code != nil {
		v.Add("code", *p.Code)
	}
	if p.From != nil {
		v.Add("from", *p.From)
	}
	if p.To != nil {
		v.Add("to", *p.To)
	}
	if p.DiscItems != nil {
		v.Add("discItems", *p.DiscItems)
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type timelyDisclosureResponse struct {
	Data          []TimelyDisclosure `json:"data"`
	PaginationKey *string            `json:"pagination_key"`
}

func (r timelyDisclosureResponse) Items() []TimelyDisclosure { return r.Data }
func (r timelyDisclosureResponse) NextPageKey() *string      { return r.PaginationKey }

// TimelyDisclosure retrieves the TDnet timely disclosure index list from the /td/list endpoint.
// It automatically handles pagination to fetch all matching records.
// This endpoint requires the TimelyDisclosure add-on plan.
// See https://jpx-jquants.com/en/spec/td-list for API details.
func (c *Client) TimelyDisclosure(ctx context.Context, req TimelyDisclosureRequest) ([]TimelyDisclosure, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (timelyDisclosureResponse, error) {
		params := timelyDisclosureParameters{TimelyDisclosureRequest: req, PaginationKey: paginationKey}
		return getJSON[timelyDisclosureResponse](ctx, c, "/td/list", params)
	})
}

// TimelyDisclosureWithChannel retrieves the TDnet timely disclosure index list and
// streams each record to the provided channel.
// The channel is closed when all records have been sent or an error occurs.
// On error the channel is closed and the error is returned from this method, so callers
// must check the returned error after the channel closes; ranging the channel alone will not surface it.
// This endpoint requires the TimelyDisclosure add-on plan.
// See https://jpx-jquants.com/en/spec/td-list for API details.
func (c *Client) TimelyDisclosureWithChannel(ctx context.Context, req TimelyDisclosureRequest, ch chan<- TimelyDisclosure) error {
	return fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, paginationKey *string) (timelyDisclosureResponse, error) {
		params := timelyDisclosureParameters{TimelyDisclosureRequest: req, PaginationKey: paginationKey}
		return getJSON[timelyDisclosureResponse](ctx, c, "/td/list", params)
	})
}

// TimelyDisclosureFileURLs holds the signed download URLs for a disclosure's
// document files. Each field is nil when the corresponding document type does
// not exist or was filtered out via the Docs request parameter.
type TimelyDisclosureFileURLs struct {
	// PDF is the signed URL for the full disclosure PDF.
	PDF *string `json:"pdf"`
	// SummaryPDF is the signed URL for the summary PDF.
	SummaryPDF *string `json:"summaryPdf"`
	// XBRL is the signed URL for the XBRL document.
	XBRL *string `json:"xbrl"`
}

// TimelyDisclosureFiles represents the signed download URLs for a single
// disclosure's document files. The URLs expire 15 minutes after issuance.
type TimelyDisclosureFiles struct {
	// DisclosureNumber is the 14-digit disclosure number.
	DisclosureNumber string `json:"discNo"`
	// Files holds the signed download URLs for the disclosure's documents.
	Files TimelyDisclosureFileURLs `json:"files"`
}

// TimelyDisclosureFilesRequest specifies parameters for the TimelyDisclosureFiles API.
type TimelyDisclosureFilesRequest struct {
	// DisclosureNumber is the 14-digit disclosure number (required).
	DisclosureNumber string
	// Docs is an optional comma-separated subset of document types (g/s/x).
	// When nil, all available documents are returned.
	Docs *string
}

type timelyDisclosureFilesParameters struct {
	TimelyDisclosureFilesRequest
}

func (p timelyDisclosureFilesParameters) values() (url.Values, error) {
	v := url.Values{}
	v.Add("discNo", p.DisclosureNumber)
	if p.Docs != nil {
		v.Add("docs", *p.Docs)
	}
	return v, nil
}

// TimelyDisclosureFiles retrieves the signed download URLs for a disclosure's
// document files from the /td/files endpoint.
// The returned URLs expire 15 minutes after issuance.
// This endpoint requires the TimelyDisclosure add-on plan.
// See https://jpx-jquants.com/en/spec/td-files for API details.
func (c *Client) TimelyDisclosureFiles(ctx context.Context, req TimelyDisclosureFilesRequest) (TimelyDisclosureFiles, error) {
	return getJSON[TimelyDisclosureFiles](ctx, c, "/td/files", timelyDisclosureFilesParameters{req})
}

// TimelyDisclosureBulk represents a signed download URL for the bulk TDnet
// timely disclosure CSV covering five years of disclosures. The URL points to
// a gzip-compressed CSV file and expires shortly after issuance.
type TimelyDisclosureBulk struct {
	// LastUpdated is the timestamp the bulk file was last updated (RFC 3339).
	LastUpdated string `json:"lastUpdated"`
	// URL is the signed download URL for the gzip-compressed CSV file.
	URL string `json:"url"`
}

type timelyDisclosureBulkParameters struct{}

func (p timelyDisclosureBulkParameters) values() (url.Values, error) {
	return url.Values{}, nil
}

// TimelyDisclosureBulk retrieves a signed download URL for the bulk TDnet timely
// disclosure CSV from the /td/bulk endpoint. The returned URL points to a
// gzip-compressed CSV file covering five years of disclosures and expires shortly
// after issuance.
// This endpoint requires the TimelyDisclosure add-on plan.
// See https://jpx-jquants.com/en/spec/td-bulk for API details.
func (c *Client) TimelyDisclosureBulk(ctx context.Context) (TimelyDisclosureBulk, error) {
	return getJSON[TimelyDisclosureBulk](ctx, c, "/td/bulk", timelyDisclosureBulkParameters{})
}
