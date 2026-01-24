package jquants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type MarginTradingOutstanding struct {
	Date                     string
	Code                     string
	TotalShortBalance        int64
	TotalLongBalance         int64
	ShortNegotiableBalance   int64
	LongNegotiableBalance    int64
	ShortStandardizedBalance int64
	LongStandardizedBalance  int64
	IssueType                int8
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

type MarginTradingOutstandingRequest struct {
	Code *string
	Date *string
	From *string
	To   *string
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

// MarginTradingOutstanding provides margin trading outstandings.
// https://jpx-jquants.com/en/spec/mkt-margin-int
func (c *Client) MarginTradingOutstanding(ctx context.Context, req MarginTradingOutstandingRequest) ([]MarginTradingOutstanding, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (marginTradingOutstandingResponse, error) {
		params := marginTradingOutstandingParameters{MarginTradingOutstandingRequest: req, PaginationKey: paginationKey}
		return c.sendMarginTradingOutstandingRequest(ctx, params)
	})
}

type ShortSellingValue struct {
	Date                            string
	Sector33Code                    string
	LongSellingValue                int64
	ShortSellingWithRestrictions    int64
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

type ShortSellingValueRequest struct {
	Sector33Code *string
	Date         *string
	From         *string
	To           *string
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

func (c *Client) ShortSellingValue(ctx context.Context, req ShortSellingValueRequest) ([]ShortSellingValue, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (shortSellingValueResponse, error) {
		params := shortSellingValueParameters{ShortSellingValueRequest: req, PaginationKey: paginationKey}
		return c.sendShortSellingValueRequest(ctx, params)
	})
}

// Outstanding Short Selling Positions Reported not implemented

// Margin Trading Outstanding not implemented

// Breakdown Trading not implemented

type TradingCalendar struct {
	Date    string
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

type TradingCalendarRequest struct {
	HolidayDivision *int8
	From            *string
	To              *string
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
