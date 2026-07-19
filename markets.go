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
	return codeDateRangeValues(p.Code, p.Date, p.From, p.To, p.PaginationKey)
}

type marginTradingOutstandingResponse struct {
	Data          []MarginTradingOutstanding `json:"data"`
	PaginationKey *string                    `json:"pagination_key"`
}

func (r marginTradingOutstandingResponse) Items() []MarginTradingOutstanding { return r.Data }
func (r marginTradingOutstandingResponse) NextPageKey() *string              { return r.PaginationKey }

// MarginTradingOutstanding retrieves margin trading balance data from the /markets/margin-interest endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/mkt-margin-int for API details.
func (c *Client) MarginTradingOutstanding(ctx context.Context, req MarginTradingOutstandingRequest) ([]MarginTradingOutstanding, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (marginTradingOutstandingResponse, error) {
		params := marginTradingOutstandingParameters{MarginTradingOutstandingRequest: req, PaginationKey: paginationKey}
		return getJSON[marginTradingOutstandingResponse](ctx, c, "/markets/margin-interest", params)
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

func (r shortSellingValueResponse) Items() []ShortSellingValue { return r.Data }
func (r shortSellingValueResponse) NextPageKey() *string       { return r.PaginationKey }

// ShortSellingValue retrieves short selling turnover data from the /markets/short-ratio endpoint.
// It automatically handles pagination to fetch all matching records.
func (c *Client) ShortSellingValue(ctx context.Context, req ShortSellingValueRequest) ([]ShortSellingValue, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (shortSellingValueResponse, error) {
		params := shortSellingValueParameters{ShortSellingValueRequest: req, PaginationKey: paginationKey}
		return getJSON[shortSellingValueResponse](ctx, c, "/markets/short-ratio", params)
	})
}

// Breakdown Trading not implemented (Premium plan only)

// OutstandingShortPosition represents an outstanding short position report
// submitted to the exchange when a short seller's position reaches a
// reportable threshold.
type OutstandingShortPosition struct {
	// DisclosureDate is the disclosure date in YYYY-MM-DD format (JSON key "DiscDate").
	DisclosureDate string
	// CalculationDate is the position calculation date in YYYY-MM-DD format (JSON key "CalcDate").
	CalculationDate string
	// Code is the security code (JSON key "Code").
	Code string
	// ShortSellerName is the name of the short seller (JSON key "SSName").
	ShortSellerName string
	// ShortSellerAddress is the address of the short seller (JSON key "SSAddr").
	ShortSellerAddress string
	// DiscretionaryInvestmentContractorName is the name of the discretionary
	// investment contractor, or "-" if none (JSON key "DICName").
	DiscretionaryInvestmentContractorName string
	// DiscretionaryInvestmentContractorAddress is the address of the
	// discretionary investment contractor, or "-" if none (JSON key "DICAddr").
	DiscretionaryInvestmentContractorAddress string
	// FundName is the name of the fund, or "-" if none (JSON key "FundName").
	FundName string
	// ShortPositionToSharesOutstandingRatio is the ratio of the short position
	// to shares outstanding (JSON key "ShrtPosToSO").
	ShortPositionToSharesOutstandingRatio float64
	// ShortPositionShares is the number of shares held short (JSON key "ShrtPosShares").
	ShortPositionShares float64
	// ShortPositionUnits is the number of trading units held short (JSON key "ShrtPosUnits").
	ShortPositionUnits float64
	// PreviousReportDate is the previous report calculation date, or "-" if none (JSON key "PrevRptDate").
	PreviousReportDate string
	// PreviousReportRatio is the short position ratio from the previous report (JSON key "PrevRptRatio").
	PreviousReportRatio float64
	// Notes contains any additional notes, or "-" if none (JSON key "Notes").
	Notes string
}

func (o *OutstandingShortPosition) UnmarshalJSON(b []byte) error {
	var raw struct {
		DiscDate      string  `json:"DiscDate"`
		CalcDate      string  `json:"CalcDate"`
		Code          string  `json:"Code"`
		SSName        string  `json:"SSName"`
		SSAddr        string  `json:"SSAddr"`
		DICName       string  `json:"DICName"`
		DICAddr       string  `json:"DICAddr"`
		FundName      string  `json:"FundName"`
		ShrtPosToSO   float64 `json:"ShrtPosToSO"`
		ShrtPosShares float64 `json:"ShrtPosShares"`
		ShrtPosUnits  float64 `json:"ShrtPosUnits"`
		PrevRptDate   string  `json:"PrevRptDate"`
		PrevRptRatio  float64 `json:"PrevRptRatio"`
		Notes         string  `json:"Notes"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal outstanding short position: %w", err)
	}
	o.DisclosureDate = raw.DiscDate
	o.CalculationDate = raw.CalcDate
	o.Code = raw.Code
	o.ShortSellerName = raw.SSName
	o.ShortSellerAddress = raw.SSAddr
	o.DiscretionaryInvestmentContractorName = raw.DICName
	o.DiscretionaryInvestmentContractorAddress = raw.DICAddr
	o.FundName = raw.FundName
	o.ShortPositionToSharesOutstandingRatio = raw.ShrtPosToSO
	o.ShortPositionShares = raw.ShrtPosShares
	o.ShortPositionUnits = raw.ShrtPosUnits
	o.PreviousReportDate = raw.PrevRptDate
	o.PreviousReportRatio = raw.PrevRptRatio
	o.Notes = raw.Notes
	return nil
}

// OutstandingShortPositionRequest specifies filter parameters for the OutstandingShortPosition API.
// All parameters are optional.
type OutstandingShortPositionRequest struct {
	// Code filters by security code.
	Code *string
	// DisclosureDate filters by disclosure date.
	DisclosureDate *string
	// DisclosureDateFrom specifies the start of a disclosure date range.
	DisclosureDateFrom *string
	// DisclosureDateTo specifies the end of a disclosure date range.
	DisclosureDateTo *string
	// CalculationDate filters by position calculation date.
	CalculationDate *string
}

type outstandingShortPositionParameters struct {
	OutstandingShortPositionRequest
	PaginationKey *string
}

func (p outstandingShortPositionParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Code != nil {
		v.Add("code", *p.Code)
	}
	if p.DisclosureDate != nil {
		v.Add("disc_date", *p.DisclosureDate)
	}
	if p.DisclosureDateFrom != nil {
		v.Add("disc_date_from", *p.DisclosureDateFrom)
	}
	if p.DisclosureDateTo != nil {
		v.Add("disc_date_to", *p.DisclosureDateTo)
	}
	if p.CalculationDate != nil {
		v.Add("calc_date", *p.CalculationDate)
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type outstandingShortPositionResponse struct {
	Data          []OutstandingShortPosition `json:"data"`
	PaginationKey *string                    `json:"pagination_key"`
}

func (r outstandingShortPositionResponse) Items() []OutstandingShortPosition { return r.Data }
func (r outstandingShortPositionResponse) NextPageKey() *string              { return r.PaginationKey }

// OutstandingShortPosition retrieves outstanding short position reports from the /markets/short-sale-report endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/mkt-short-sale for API details.
func (c *Client) OutstandingShortPosition(ctx context.Context, req OutstandingShortPositionRequest) ([]OutstandingShortPosition, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (outstandingShortPositionResponse, error) {
		params := outstandingShortPositionParameters{OutstandingShortPositionRequest: req, PaginationKey: paginationKey}
		return getJSON[outstandingShortPositionResponse](ctx, c, "/markets/short-sale-report", params)
	})
}

// MarginAlertPublicationReason describes why a security was published on the
// margin trading alert list. Each field is a flag reported by the API as the
// string "0" (not applicable) or "1" (applicable).
type MarginAlertPublicationReason struct {
	// Restricted indicates the issue is subject to trading restrictions (JSON key "Restricted").
	Restricted string
	// DailyPublication indicates the issue is subject to daily publication (JSON key "DailyPublication").
	DailyPublication string
	// Monitoring indicates the issue is under monitoring (JSON key "Monitoring").
	Monitoring string
	// RestrictedByJSF indicates the issue is restricted by the Japan Securities Finance Co. (JSON key "RestrictedByJSF").
	RestrictedByJSF string
	// PrecautionByJSF indicates a precaution notice by the Japan Securities Finance Co. (JSON key "PrecautionByJSF").
	PrecautionByJSF string
	// UnclearOrSecOnAlert indicates an unclear status or that the security is on alert (JSON key "UnclearOrSecOnAlert").
	UnclearOrSecOnAlert string
}

// MarginAlert represents a margin trading alert entry, giving outstanding
// margin balances and the reason the issue was published on the alert list.
type MarginAlert struct {
	// PublicationDate is the publication date in YYYY-MM-DD format (JSON key "PubDate").
	PublicationDate string
	// Code is the security code (JSON key "Code").
	Code string
	// ApplicationDate is the application (record) date in YYYY-MM-DD format (JSON key "AppDate").
	ApplicationDate string
	// PublicationReason describes why the issue was published (JSON key "PubReason").
	PublicationReason MarginAlertPublicationReason
	// ShortOutstanding is the outstanding short margin balance (JSON key "ShrtOut").
	ShortOutstanding float64
	// LongOutstanding is the outstanding long margin balance (JSON key "LongOut").
	LongOutstanding float64
	// ShortLongRatio is the ratio of short to long margin balance (JSON key "SLRatio").
	ShortLongRatio float64
	// ShortNegotiableOutstanding is the outstanding negotiable short balance (JSON key "ShrtNegOut").
	ShortNegotiableOutstanding float64
	// ShortStandardizedOutstanding is the outstanding standardized short balance (JSON key "ShrtStdOut").
	ShortStandardizedOutstanding float64
	// LongNegotiableOutstanding is the outstanding negotiable long balance (JSON key "LongNegOut").
	LongNegotiableOutstanding float64
	// LongStandardizedOutstanding is the outstanding standardized long balance (JSON key "LongStdOut").
	LongStandardizedOutstanding float64
	// ShortOutstandingChange is the change in short balance; nil when unavailable ("-") (JSON key "ShrtOutChg").
	ShortOutstandingChange *float64
	// LongOutstandingChange is the change in long balance; nil when unavailable ("-") (JSON key "LongOutChg").
	LongOutstandingChange *float64
	// ShortNegotiableOutstandingChange is the change in negotiable short balance; nil when unavailable ("-") (JSON key "ShrtNegOutChg").
	ShortNegotiableOutstandingChange *float64
	// ShortStandardizedOutstandingChange is the change in standardized short balance; nil when unavailable ("-") (JSON key "ShrtStdOutChg").
	ShortStandardizedOutstandingChange *float64
	// LongNegotiableOutstandingChange is the change in negotiable long balance; nil when unavailable ("-") (JSON key "LongNegOutChg").
	LongNegotiableOutstandingChange *float64
	// LongStandardizedOutstandingChange is the change in standardized long balance; nil when unavailable ("-") (JSON key "LongStdOutChg").
	LongStandardizedOutstandingChange *float64
	// ShortOutstandingRatio is the short balance ratio; nil for ETFs ("*") (JSON key "ShrtOutRatio").
	ShortOutstandingRatio *float64
	// LongOutstandingRatio is the long balance ratio; nil for ETFs ("*") (JSON key "LongOutRatio").
	LongOutstandingRatio *float64
	// TSEMarginRegulationClass is the TSE margin regulation classification code (JSON key "TSEMrgnRegCls").
	TSEMarginRegulationClass string
}

func (m *MarginAlert) UnmarshalJSON(b []byte) error {
	var raw struct {
		PubDate       string                       `json:"PubDate"`
		Code          string                       `json:"Code"`
		AppDate       string                       `json:"AppDate"`
		PubReason     MarginAlertPublicationReason `json:"PubReason"`
		ShrtOut       float64                      `json:"ShrtOut"`
		LongOut       float64                      `json:"LongOut"`
		SLRatio       float64                      `json:"SLRatio"`
		ShrtNegOut    float64                      `json:"ShrtNegOut"`
		ShrtStdOut    float64                      `json:"ShrtStdOut"`
		LongNegOut    float64                      `json:"LongNegOut"`
		LongStdOut    float64                      `json:"LongStdOut"`
		ShrtOutChg    any                          `json:"ShrtOutChg"`
		LongOutChg    any                          `json:"LongOutChg"`
		ShrtNegOutChg any                          `json:"ShrtNegOutChg"`
		ShrtStdOutChg any                          `json:"ShrtStdOutChg"`
		LongNegOutChg any                          `json:"LongNegOutChg"`
		LongStdOutChg any                          `json:"LongStdOutChg"`
		ShrtOutRatio  any                          `json:"ShrtOutRatio"`
		LongOutRatio  any                          `json:"LongOutRatio"`
		TSEMrgnRegCls string                       `json:"TSEMrgnRegCls"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal margin alert: %w", err)
	}
	a := &floatAccumulator{}
	fromAny := func(v any) *float64 {
		if a.err != nil {
			return nil
		}
		result, err := unmarshalFloatFromAny(v)
		a.err = err
		return result
	}
	m.PublicationDate = raw.PubDate
	m.Code = raw.Code
	m.ApplicationDate = raw.AppDate
	m.PublicationReason = raw.PubReason
	m.ShortOutstanding = raw.ShrtOut
	m.LongOutstanding = raw.LongOut
	m.ShortLongRatio = raw.SLRatio
	m.ShortNegotiableOutstanding = raw.ShrtNegOut
	m.ShortStandardizedOutstanding = raw.ShrtStdOut
	m.LongNegotiableOutstanding = raw.LongNegOut
	m.LongStandardizedOutstanding = raw.LongStdOut
	m.ShortOutstandingChange = fromAny(raw.ShrtOutChg)
	m.LongOutstandingChange = fromAny(raw.LongOutChg)
	m.ShortNegotiableOutstandingChange = fromAny(raw.ShrtNegOutChg)
	m.ShortStandardizedOutstandingChange = fromAny(raw.ShrtStdOutChg)
	m.LongNegotiableOutstandingChange = fromAny(raw.LongNegOutChg)
	m.LongStandardizedOutstandingChange = fromAny(raw.LongStdOutChg)
	m.ShortOutstandingRatio = fromAny(raw.ShrtOutRatio)
	m.LongOutstandingRatio = fromAny(raw.LongOutRatio)
	m.TSEMarginRegulationClass = raw.TSEMrgnRegCls
	return a.err
}

// MarginAlertRequest specifies filter parameters for the MarginAlert API.
// All parameters are optional.
type MarginAlertRequest struct {
	// Code filters by security code.
	Code *string
	// Date filters by publication date.
	Date *string
	// From specifies the start of a publication date range.
	From *string
	// To specifies the end of a publication date range.
	To *string
}

type marginAlertParameters struct {
	MarginAlertRequest
	PaginationKey *string
}

func (p marginAlertParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Code != nil {
		v.Add("code", *p.Code)
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
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type marginAlertResponse struct {
	Data          []MarginAlert `json:"data"`
	PaginationKey *string       `json:"pagination_key"`
}

func (r marginAlertResponse) Items() []MarginAlert { return r.Data }
func (r marginAlertResponse) NextPageKey() *string { return r.PaginationKey }

// MarginAlert retrieves margin trading alert data from the /markets/margin-alert endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/mkt-margin-alert for API details.
func (c *Client) MarginAlert(ctx context.Context, req MarginAlertRequest) ([]MarginAlert, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (marginAlertResponse, error) {
		params := marginAlertParameters{MarginAlertRequest: req, PaginationKey: paginationKey}
		return getJSON[marginAlertResponse](ctx, c, "/markets/margin-alert", params)
	})
}

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
	r, err := getJSON[tradingCalendarResponse](ctx, c, "/markets/calendar", tradingCalendarParameters{req})
	if err != nil {
		return nil, err
	}
	return r.Data, nil
}
