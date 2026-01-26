package jquants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// MarginTradingOutstanding represents margin trading balance data for a security.
// It shows the outstanding short and long positions broken down by trade type.
type MarginTradingOutstanding struct {
	// Date is the data date in YYYY-MM-DD format.
	Date string
	// Code is the security code (ticker symbol).
	Code string
	// TotalShortBalance is the total short margin trading balance in shares.
	TotalShortBalance int64
	// TotalLongBalance is the total long margin trading balance in shares.
	TotalLongBalance int64
	// ShortNegotiableBalance is the short balance for negotiable margin trades.
	ShortNegotiableBalance int64
	// LongNegotiableBalance is the long balance for negotiable margin trades.
	LongNegotiableBalance int64
	// ShortStandardizedBalance is the short balance for standardized margin trades.
	ShortStandardizedBalance int64
	// LongStandardizedBalance is the long balance for standardized margin trades.
	LongStandardizedBalance int64
	// IssueType indicates the type of issue (1: Prime, 2: Standard, 3: Growth).
	IssueType int8
}

func (mtv *MarginTradingOutstanding) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date                               string  `json:"Date"`
		Code                               string  `json:"Code"`
		ShortMarginTradeVolume             float64 `json:"ShrtVol"`
		LongMarginTradeVolume              float64 `json:"LongVol"`
		ShortNegotiableMarginTradeVolume   float64 `json:"ShrtNegVol"`
		LongNegotiableMarginTradeVolume    float64 `json:"LongNegVol"`
		ShortStandardizedMarginTradeVolume float64 `json:"ShrtStdVol"`
		LongStandardizedMarginTradeVolume  float64 `json:"LongStdVol"`
		IssueType                          string  `json:"IssType"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal margin trading outstanding: %w", err)
	}
	var err error
	mtv.Date = raw.Date
	issueType, err := strconv.ParseInt(raw.IssueType, 10, 8)
	if err != nil {
		return fmt.Errorf("failed to unmarshal margin trading outstanding: %w", err)
	}
	mtv.Code = raw.Code
	mtv.TotalShortBalance = int64(raw.ShortMarginTradeVolume)
	mtv.TotalLongBalance = int64(raw.LongMarginTradeVolume)
	mtv.ShortNegotiableBalance = int64(raw.ShortNegotiableMarginTradeVolume)
	mtv.LongNegotiableBalance = int64(raw.LongNegotiableMarginTradeVolume)
	mtv.ShortStandardizedBalance = int64(raw.ShortStandardizedMarginTradeVolume)
	mtv.LongStandardizedBalance = int64(raw.LongStandardizedMarginTradeVolume)
	mtv.IssueType = int8(issueType)
	return nil
}

// MarginTradingOutstandingRequest specifies filter parameters for the MarginTradingOutstanding API.
// Either Code or Date must be provided.
type MarginTradingOutstandingRequest struct {
	// Code filters by security code. Required if Date is not specified.
	Code *string
	// Date filters by a specific date in YYYY-MM-DD format. If specified, Code is ignored.
	Date *string
	// From specifies the start date for a date range query (used with Code).
	From *string
	// To specifies the end date for a date range query (used with Code).
	To *string
}

type marginTradingOutstandingParameters struct {
	MarginTradingOutstandingRequest
	PaginationKey *string
}

func (p marginTradingOutstandingParameters) values() (url.Values, error) {
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

type marginTradingOutstandingResponse struct {
	Data          []MarginTradingOutstanding `json:"data"`
	PaginationKey *string                    `json:"pagination_key"`
}

func (r marginTradingOutstandingResponse) getData() []MarginTradingOutstanding { return r.Data }
func (r marginTradingOutstandingResponse) getPaginationKey() *string            { return r.PaginationKey }

func (c *Client) sendMarginTradingOutstandingRequest(ctx context.Context, params marginTradingOutstandingParameters) (marginTradingOutstandingResponse, error) {
	var r marginTradingOutstandingResponse
	resp, err := c.sendRequest(ctx, "/markets/margin-interest", params)
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

// MarginTradingOutstanding retrieves margin trading balance data from the /markets/margin-interest endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/mkt-margin-int for API details.
func (c *Client) MarginTradingOutstanding(ctx context.Context, req MarginTradingOutstandingRequest) ([]MarginTradingOutstanding, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (marginTradingOutstandingResponse, error) {
		params := marginTradingOutstandingParameters{MarginTradingOutstandingRequest: req, PaginationKey: paginationKey}
		return c.sendMarginTradingOutstandingRequest(ctx, params)
	})
}

// ShortSellingValue represents short selling turnover data by sector.
// Values are broken down by selling type (long, short with/without restrictions).
type ShortSellingValue struct {
	// Date is the trading date in YYYY-MM-DD format.
	Date string
	// Sector33Code is the 33-sector classification code.
	Sector33Code string
	// LongSellingValue is the turnover value of long selling (non-short) in yen.
	LongSellingValue int64
	// ShortSellingWithRestrictions is the turnover value of short selling with price restrictions in yen.
	ShortSellingWithRestrictions int64
	// ShortSellingWithoutRestrictions is the turnover value of short selling without price restrictions in yen.
	ShortSellingWithoutRestrictions int64
}

func (sst *ShortSellingValue) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date                                         string  `json:"Date"`
		Sector33Code                                 string  `json:"S33"`
		SellingExcludingShortSellingTurnoverValue    float64 `json:"SellExShortVa"`
		ShortSellingWithRestrictionsTurnoverValue    float64 `json:"ShrtWithResVa"`
		ShortSellingWithoutRestrictionsTurnoverValue float64 `json:"ShrtNoResVa"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal short selling value: %w", err)
	}
	sst.Date = raw.Date
	sst.Sector33Code = raw.Sector33Code
	sst.LongSellingValue = int64(raw.SellingExcludingShortSellingTurnoverValue)
	sst.ShortSellingWithRestrictions = int64(raw.ShortSellingWithRestrictionsTurnoverValue)
	sst.ShortSellingWithoutRestrictions = int64(raw.ShortSellingWithoutRestrictionsTurnoverValue)
	return nil
}

// ShortSellingValueRequest specifies filter parameters for the ShortSellingValue API.
// Either Sector33Code or Date must be provided.
type ShortSellingValueRequest struct {
	// Sector33Code filters by 33-sector classification code.
	Sector33Code *string
	// Date filters by a specific date in YYYY-MM-DD format.
	Date *string
	// From specifies the start date for a date range query (used with Sector33Code).
	From *string
	// To specifies the end date for a date range query (used with Sector33Code).
	To *string
}

type shortSellingValueParameters struct {
	ShortSellingValueRequest
	PaginationKey *string
}

func (p shortSellingValueParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Sector33Code != nil {
		v.Add("s33", *p.Sector33Code)
		if p.Date != nil {
			v.Add("date", *p.Date)
		} else {
			if p.From != nil {
				v.Add("from", *p.From)
			}
			if p.To != nil {
				v.Add("to", *p.To)
			}
		}
	} else {
		if p.Date == nil {
			return nil, errors.New("sector33code or date is required")
		}
		v.Add("date", *p.Date)
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type shortSellingValueResponse struct {
	Data          []ShortSellingValue `json:"data"`
	PaginationKey *string             `json:"pagination_key"`
}

func (r shortSellingValueResponse) getData() []ShortSellingValue { return r.Data }
func (r shortSellingValueResponse) getPaginationKey() *string    { return r.PaginationKey }

func (c *Client) sendShortSellingValueRequest(ctx context.Context, params shortSellingValueParameters) (shortSellingValueResponse, error) {
	var r shortSellingValueResponse
	resp, err := c.sendRequest(ctx, "/markets/short-ratio", params)
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

// ShortSellingValue retrieves short selling turnover data from the /markets/short-ratio endpoint.
// It automatically handles pagination to fetch all matching records.
func (c *Client) ShortSellingValue(ctx context.Context, req ShortSellingValueRequest) ([]ShortSellingValue, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (shortSellingValueResponse, error) {
		params := shortSellingValueParameters{ShortSellingValueRequest: req, PaginationKey: paginationKey}
		return c.sendShortSellingValueRequest(ctx, params)
	})
}

// Outstanding Short Selling Positions Reported not implemented

// Margin Trading Outstanding not implemented

// Breakdown Trading not implemented

// TradingCalendar represents a trading calendar entry indicating whether a date is a trading day.
type TradingCalendar struct {
	// Date is the calendar date in YYYY-MM-DD format.
	Date string
	// DayType indicates the day type (0: holiday/non-trading day, 1: trading day, 2: half-day, 3: non-trading day).
	DayType int8
}

func (tc *TradingCalendar) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date            string `json:"Date"`
		HolidayDivision string `json:"HolDiv"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal trading calendar: %w", err)
	}
	tc.Date = raw.Date
	hd, err := strconv.ParseInt(raw.HolidayDivision, 10, 8)
	if err != nil {
		return fmt.Errorf("failed to unmarshal trading calendar: %w", err)
	}
	tc.DayType = int8(hd)
	return nil
}

// TradingCalendarRequest specifies filter parameters for the TradingCalendar API.
type TradingCalendarRequest struct {
	// HolidayDivision filters by day type (0: holiday, 1: trading day, 2: half-day, 3: non-trading day).
	HolidayDivision *int8
	// From specifies the start date for the query in YYYY-MM-DD format.
	From *string
	// To specifies the end date for the query in YYYY-MM-DD format.
	To *string
}

type tradingCalendarParameters struct {
	TradingCalendarRequest
}

func (p tradingCalendarParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.HolidayDivision != nil {
		v.Add("hol_div", strconv.Itoa(int(*p.HolidayDivision)))
	}
	if p.From != nil {
		v.Add("from", *p.From)
	}
	if p.To != nil {
		v.Add("to", *p.To)
	}
	return v, nil
}

type tradingCalendarResponse struct {
	Data []TradingCalendar `json:"data"`
}

// TradingCalendar retrieves the TSE trading calendar from the /markets/calendar endpoint.
func (c *Client) TradingCalendar(ctx context.Context, req TradingCalendarRequest) ([]TradingCalendar, error) {
	var r tradingCalendarResponse
	params := tradingCalendarParameters{TradingCalendarRequest: req}
	resp, err := c.sendRequest(ctx, "/markets/calendar", params)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp)
	}
	if err = decodeResponse(resp, &r); err != nil {
		return nil, fmt.Errorf("failed to decode HTTP response: %w", err)
	}
	return r.Data, nil
}
