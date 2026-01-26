package jquants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// IndexPrice represents daily OHLC (Open, High, Low, Close) data for a market index.
type IndexPrice struct {
	// Date is the trading date in YYYY-MM-DD format.
	Date string
	// Code is the index code (e.g., "0000" for TOPIX, "0001" for TOPIX Core30).
	Code string
	// Open is the opening value of the index.
	Open json.Number
	// High is the highest value of the index for the day.
	High json.Number
	// Low is the lowest value of the index for the day.
	Low json.Number
	// Close is the closing value of the index.
	Close json.Number
}

func (ip *IndexPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date  string      `json:"Date"`
		Code  string      `json:"Code"`
		Open  json.Number `json:"O"`
		High  json.Number `json:"H"`
		Low   json.Number `json:"L"`
		Close json.Number `json:"C"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal index price: %w", err)
	}
	ip.Date = raw.Date
	ip.Code = raw.Code
	ip.Open = raw.Open
	ip.High = raw.High
	ip.Low = raw.Low
	ip.Close = raw.Close
	return nil
}

// IndexPriceRequest specifies filter parameters for the IndexPrice API.
// Either Code or Date must be provided.
type IndexPriceRequest struct {
	// Code filters by index code. Required if Date is not specified.
	Code *string
	// Date filters by a specific date in YYYY-MM-DD format. If specified, Code is ignored.
	Date *string
	// From specifies the start date for a date range query (used with Code).
	From *string
	// To specifies the end date for a date range query (used with Code).
	To *string
}

type indexPriceParameters struct {
	IndexPriceRequest
	PaginationKey *string
}

func (p indexPriceParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Date != nil {
		v.Add("date", *p.Date)
	} else {
		if p.Code == nil {
			return nil, errors.New("code or date is required")
		}
		v.Add("code", *p.Code)
		if p.From != nil {
			v.Add("from", *p.From)
		}
		if p.To != nil {
			v.Add("to", *p.To)
		}
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type indexPriceResponse struct {
	Data          []IndexPrice `json:"data"`
	PaginationKey *string      `json:"pagination_key"`
}

func (r indexPriceResponse) Items() []IndexPrice    { return r.Data }
func (r indexPriceResponse) NextPageKey() *string { return r.PaginationKey }

func (c *Client) sendIndexPriceRequest(ctx context.Context, params indexPriceParameters) (indexPriceResponse, error) {
	var r indexPriceResponse
	resp, err := c.sendRequest(ctx, "/indices/bars/daily", params)
	if err != nil {
		return r, fmt.Errorf("failed to send GET request: %w", err)
	}
	if resp.StatusCode != 200 {
		return r, handleErrorResponse(resp)
	}
	if err = decodeResponse(resp, &r); err != nil {
		return r, fmt.Errorf("failed to decode HTTP response: %w", err)
	}
	return r, nil
}

// IndexPrice retrieves daily index prices from the /indices/bars/daily endpoint.
// It automatically handles pagination to fetch all matching records.
func (c *Client) IndexPrice(ctx context.Context, req IndexPriceRequest) ([]IndexPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (indexPriceResponse, error) {
		params := indexPriceParameters{IndexPriceRequest: req, PaginationKey: paginationKey}
		return c.sendIndexPriceRequest(ctx, params)
	})
}

// TopixPrice represents daily OHLC (Open, High, Low, Close) data for the TOPIX index.
type TopixPrice struct {
	// Date is the trading date in YYYY-MM-DD format.
	Date string
	// Open is the opening value of TOPIX.
	Open json.Number
	// High is the highest value of TOPIX for the day.
	High json.Number
	// Low is the lowest value of TOPIX for the day.
	Low json.Number
	// Close is the closing value of TOPIX.
	Close json.Number
}

func (p *TopixPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date  string      `json:"Date"`
		Open  json.Number `json:"O"`
		High  json.Number `json:"H"`
		Low   json.Number `json:"L"`
		Close json.Number `json:"C"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal topix price: %w", err)
	}
	p.Date = raw.Date
	p.Open = raw.Open
	p.High = raw.High
	p.Low = raw.Low
	p.Close = raw.Close
	return nil
}

// TopixPriceRequest specifies filter parameters for the TopixPrices API.
type TopixPriceRequest struct {
	// From specifies the start date for the query in YYYY-MM-DD format.
	From *string
	// To specifies the end date for the query in YYYY-MM-DD format.
	To *string
}

type topixPriceParameters struct {
	TopixPriceRequest
	PaginationKey *string
}

func (p topixPriceParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.From != nil {
		v.Add("from", *p.From)
	}
	if p.To != nil {
		v.Add("to", *p.To)
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type topixPriceResponse struct {
	Data          []TopixPrice `json:"data"`
	PaginationKey *string      `json:"pagination_key"`
}

func (r topixPriceResponse) Items() []TopixPrice    { return r.Data }
func (r topixPriceResponse) NextPageKey() *string { return r.PaginationKey }

func (c *Client) sendTopixPriceRequest(ctx context.Context, params topixPriceParameters) (topixPriceResponse, error) {
	var r topixPriceResponse
	resp, err := c.sendRequest(ctx, "/indices/bars/daily/topix", params)
	if err != nil {
		return r, fmt.Errorf("failed to send GET request: %w", err)
	}
	if resp.StatusCode != 200 {
		return r, handleErrorResponse(resp)
	}
	if err = decodeResponse(resp, &r); err != nil {
		return r, fmt.Errorf("failed to decode HTTP response: %w", err)
	}
	return r, nil
}

// TopixPrices retrieves daily TOPIX index prices from the /indices/bars/daily/topix endpoint.
// It automatically handles pagination to fetch all matching records.
func (c *Client) TopixPrices(ctx context.Context, req TopixPriceRequest) ([]TopixPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (topixPriceResponse, error) {
		params := topixPriceParameters{TopixPriceRequest: req, PaginationKey: paginationKey}
		return c.sendTopixPriceRequest(ctx, params)
	})
}
