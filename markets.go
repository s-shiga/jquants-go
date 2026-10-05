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
// Data is daily from 2026-09-25 and weekly before that date. PublicationDate and
// the value fields are nil for earlier records.
type MarginTradingOutstanding struct {
	// PublicationDate is the published date in YYYY-MM-DD format (JSON key "PubDate").
	PublicationDate *string
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
	// TotalShortValue is the total value of short margin positions in yen (JSON key "ShrtVal").
	TotalShortValue *float64
	// TotalLongValue is the total value of long margin positions in yen (JSON key "LongVal").
	TotalLongValue *float64
	// ShortNegotiableValue is the value of negotiable short margin positions in yen (JSON key "ShrtNegVal").
	ShortNegotiableValue *float64
	// LongNegotiableValue is the value of negotiable long margin positions in yen (JSON key "LongNegVal").
	LongNegotiableValue *float64
	// ShortStandardizedValue is the value of standardized short margin positions in yen (JSON key "ShrtStdVal").
	ShortStandardizedValue *float64
	// LongStandardizedValue is the value of standardized long margin positions in yen (JSON key "LongStdVal").
	LongStandardizedValue *float64
	// IssueType is the issue classification (1: margin issue, 2: loan issue, 3: other issue).
	IssueType int8
}

func (mtv *MarginTradingOutstanding) UnmarshalJSON(b []byte) error {
	type StoredRecord MarginTradingOutstanding
	var raw struct {
		*StoredRecord
		PublicationDate                    *string        `json:"PubDate"`
		Date                               string         `json:"Date"`
		Code                               string         `json:"Code"`
		ShortMarginTradeVolume             nullableNumber `json:"ShrtVol"`
		LongMarginTradeVolume              nullableNumber `json:"LongVol"`
		ShortNegotiableMarginTradeVolume   nullableNumber `json:"ShrtNegVol"`
		LongNegotiableMarginTradeVolume    nullableNumber `json:"LongNegVol"`
		ShortStandardizedMarginTradeVolume nullableNumber `json:"ShrtStdVol"`
		LongStandardizedMarginTradeVolume  nullableNumber `json:"LongStdVol"`
		TotalShortValue                    nullableNumber `json:"ShrtVal"`
		TotalLongValue                     nullableNumber `json:"LongVal"`
		ShortNegotiableValue               nullableNumber `json:"ShrtNegVal"`
		LongNegotiableValue                nullableNumber `json:"LongNegVal"`
		ShortStandardizedValue             nullableNumber `json:"ShrtStdVal"`
		LongStandardizedValue              nullableNumber `json:"LongStdVal"`
		IssueType                          string         `json:"IssType"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal margin trading outstanding: %w", err)
	}
	if raw.StoredRecord != nil {
		raw.StoredRecord.Date = raw.Date
		raw.StoredRecord.Code = raw.Code
		*mtv = MarginTradingOutstanding(*raw.StoredRecord)
		return nil
	}
	mtv.PublicationDate = raw.PublicationDate
	mtv.Date = raw.Date
	issueType, err := strconv.ParseInt(raw.IssueType, 10, 8)
	if err != nil {
		return fmt.Errorf("failed to unmarshal margin trading outstanding: %w", err)
	}
	mtv.Code = raw.Code
	u := &unmarshaler{}
	mtv.TotalShortBalance = u.integer(raw.ShortMarginTradeVolume)
	mtv.TotalLongBalance = u.integer(raw.LongMarginTradeVolume)
	mtv.ShortNegotiableBalance = u.integer(raw.ShortNegotiableMarginTradeVolume)
	mtv.LongNegotiableBalance = u.integer(raw.LongNegotiableMarginTradeVolume)
	mtv.ShortStandardizedBalance = u.integer(raw.ShortStandardizedMarginTradeVolume)
	mtv.LongStandardizedBalance = u.integer(raw.LongStandardizedMarginTradeVolume)
	mtv.TotalShortValue = u.float(raw.TotalShortValue)
	mtv.TotalLongValue = u.float(raw.TotalLongValue)
	mtv.ShortNegotiableValue = u.float(raw.ShortNegotiableValue)
	mtv.LongNegotiableValue = u.float(raw.LongNegotiableValue)
	mtv.ShortStandardizedValue = u.float(raw.ShortStandardizedValue)
	mtv.LongStandardizedValue = u.float(raw.LongStandardizedValue)
	mtv.IssueType = int8(issueType)
	return u.err
}

// MarginTradingOutstandingRequest specifies filter parameters for the MarginTradingOutstanding API.
// At least one of Code, Date, or PublishedDate must be provided.
type MarginTradingOutstandingRequest struct {
	// Code filters by security code. Required if neither Date nor PublishedDate is specified.
	Code *string
	// Date filters by a specific record date in YYYY-MM-DD or YYYYMMDD format.
	Date *string
	// From specifies the start date for a date range query (used with Code).
	From *string
	// To specifies the end date for a date range query (used with Code).
	To *string
	// PublishedDate filters by published date in YYYY-MM-DD or YYYYMMDD format.
	// It can be combined with Code, but not with Date, From, or To.
	// Historical records without a publication date are excluded.
	PublishedDate *string
}

type marginTradingOutstandingParameters struct {
	MarginTradingOutstandingRequest
	PaginationKey *string
}

func (p marginTradingOutstandingParameters) values() (url.Values, error) {
	if p.Code == nil && p.Date == nil && p.PublishedDate == nil {
		return nil, errors.New("code, date, or published_date is required")
	}
	if p.PublishedDate == nil {
		return codeDateRangeValues(p.Code, p.Date, p.From, p.To, p.PaginationKey)
	}
	if p.Date != nil || p.From != nil || p.To != nil {
		return nil, errors.New("published_date cannot be combined with date, from, or to")
	}
	v := url.Values{"published_date": {*p.PublishedDate}}
	if p.Code != nil {
		v.Add("code", *p.Code)
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

// MarginTradingOutstanding retrieves margin trading balance data from the /markets/margin-interest endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/mkt-margin-int for API details.
func (c *Client) MarginTradingOutstanding(ctx context.Context, req MarginTradingOutstandingRequest) ([]MarginTradingOutstanding, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[MarginTradingOutstanding], error) {
		params := marginTradingOutstandingParameters{MarginTradingOutstandingRequest: req, PaginationKey: paginationKey}
		return getJSON[page[MarginTradingOutstanding]](ctx, c, "/markets/margin-interest", params)
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
	type StoredRecord ShortSellingValue
	var raw struct {
		*StoredRecord
		Date                                         string         `json:"Date"`
		Sector33Code                                 string         `json:"S33"`
		SellingExcludingShortSellingTurnoverValue    nullableNumber `json:"SellExShortVa"`
		ShortSellingWithRestrictionsTurnoverValue    nullableNumber `json:"ShrtWithResVa"`
		ShortSellingWithoutRestrictionsTurnoverValue nullableNumber `json:"ShrtNoResVa"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal short selling value: %w", err)
	}
	if raw.StoredRecord != nil {
		raw.StoredRecord.Date = raw.Date
		*sst = ShortSellingValue(*raw.StoredRecord)
		return nil
	}
	sst.Date = raw.Date
	sst.Sector33Code = raw.Sector33Code
	u := &unmarshaler{}
	sst.LongSellingValue = u.integer(raw.SellingExcludingShortSellingTurnoverValue)
	sst.ShortSellingWithRestrictions = u.integer(raw.ShortSellingWithRestrictionsTurnoverValue)
	sst.ShortSellingWithoutRestrictions = u.integer(raw.ShortSellingWithoutRestrictionsTurnoverValue)
	return u.err
}

// ShortSellingValueRequest specifies filter parameters for the ShortSellingValue API.
// Either Sector33Code or Date must be provided.
type ShortSellingValueRequest struct {
	// Sector33Code filters by 33-sector classification code.
	Sector33Code *string
	// Date filters by a specific date in YYYY-MM-DD format. It can be combined
	// with Sector33Code, but not with From or To.
	Date *string
	// From specifies the start date for a date range query (used with Sector33Code, not Date).
	From *string
	// To specifies the end date for a date range query (used with Sector33Code, not Date).
	To *string
}

type shortSellingValueParameters struct {
	ShortSellingValueRequest
	PaginationKey *string
}

func (p shortSellingValueParameters) values() (url.Values, error) {
	if p.Sector33Code == nil && p.Date == nil {
		return nil, errors.New("sector33code or date is required")
	}
	if p.Date != nil && (p.From != nil || p.To != nil) {
		return nil, errors.New("date cannot be combined with from or to")
	}
	v := url.Values{}
	if p.Sector33Code != nil {
		v.Add("s33", *p.Sector33Code)
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

// ShortSellingValue retrieves short selling turnover data from the /markets/short-ratio endpoint.
// It automatically handles pagination to fetch all matching records.
func (c *Client) ShortSellingValue(ctx context.Context, req ShortSellingValueRequest) ([]ShortSellingValue, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[ShortSellingValue], error) {
		params := shortSellingValueParameters{ShortSellingValueRequest: req, PaginationKey: paginationKey}
		return getJSON[page[ShortSellingValue]](ctx, c, "/markets/short-ratio", params)
	})
}

// BreakdownTrading represents a breakdown of daily trading value and volume for
// a security by trade type, distinguishing long trades from margin trades and
// splitting margin trades into positions being newly opened versus closed.
type BreakdownTrading struct {
	// Date is the trading date in YYYY-MM-DD format (JSON key "Date").
	Date string
	// Code is the security code (JSON key "Code").
	Code string
	// LongSellValue is the trading value of long selling in yen (JSON key "LongSellVa").
	LongSellValue float64
	// ShortSellWithoutMarginValue is the trading value of short selling that is not
	// margin trading, in yen (JSON key "ShrtNoMrgnVa").
	ShortSellWithoutMarginValue float64
	// MarginSellNewValue is the trading value in yen of sell orders creating new
	// margin sell positions (JSON key "MrgnSellNewVa").
	MarginSellNewValue float64
	// MarginSellCloseValue is the trading value in yen of sell orders closing existing
	// margin buy positions (JSON key "MrgnSellCloseVa").
	MarginSellCloseValue float64
	// LongBuyValue is the trading value of long buying in yen (JSON key "LongBuyVa").
	LongBuyValue float64
	// MarginBuyNewValue is the trading value in yen of buy orders creating new margin
	// buy positions (JSON key "MrgnBuyNewVa").
	MarginBuyNewValue float64
	// MarginBuyCloseValue is the trading value in yen of buy orders closing existing
	// margin sell positions (JSON key "MrgnBuyCloseVa").
	MarginBuyCloseValue float64
	// LongSellVolume is the volume of long selling in shares (JSON key "LongSellVo").
	LongSellVolume int64
	// ShortSellWithoutMarginVolume is the volume of short selling that is not margin
	// trading, in shares (JSON key "ShrtNoMrgnVo").
	ShortSellWithoutMarginVolume int64
	// MarginSellNewVolume is the volume in shares of sell orders creating new margin
	// sell positions (JSON key "MrgnSellNewVo").
	MarginSellNewVolume int64
	// MarginSellCloseVolume is the volume in shares of sell orders closing existing
	// margin buy positions (JSON key "MrgnSellCloseVo").
	MarginSellCloseVolume int64
	// LongBuyVolume is the volume of long buying in shares (JSON key "LongBuyVo").
	LongBuyVolume int64
	// MarginBuyNewVolume is the volume in shares of buy orders creating new margin buy
	// positions (JSON key "MrgnBuyNewVo").
	MarginBuyNewVolume int64
	// MarginBuyCloseVolume is the volume in shares of buy orders closing existing margin
	// sell positions (JSON key "MrgnBuyCloseVo").
	MarginBuyCloseVolume int64
}

func (bt *BreakdownTrading) UnmarshalJSON(b []byte) error {
	type StoredRecord BreakdownTrading
	var raw struct {
		*StoredRecord
		Date            string         `json:"Date"`
		Code            string         `json:"Code"`
		LongSellVa      float64        `json:"LongSellVa"`
		ShrtNoMrgnVa    float64        `json:"ShrtNoMrgnVa"`
		MrgnSellNewVa   float64        `json:"MrgnSellNewVa"`
		MrgnSellCloseVa float64        `json:"MrgnSellCloseVa"`
		LongBuyVa       float64        `json:"LongBuyVa"`
		MrgnBuyNewVa    float64        `json:"MrgnBuyNewVa"`
		MrgnBuyCloseVa  float64        `json:"MrgnBuyCloseVa"`
		LongSellVo      nullableNumber `json:"LongSellVo"`
		ShrtNoMrgnVo    nullableNumber `json:"ShrtNoMrgnVo"`
		MrgnSellNewVo   nullableNumber `json:"MrgnSellNewVo"`
		MrgnSellCloseVo nullableNumber `json:"MrgnSellCloseVo"`
		LongBuyVo       nullableNumber `json:"LongBuyVo"`
		MrgnBuyNewVo    nullableNumber `json:"MrgnBuyNewVo"`
		MrgnBuyCloseVo  nullableNumber `json:"MrgnBuyCloseVo"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal breakdown trading: %w", err)
	}
	if raw.StoredRecord != nil {
		raw.StoredRecord.Date = raw.Date
		raw.StoredRecord.Code = raw.Code
		*bt = BreakdownTrading(*raw.StoredRecord)
		return nil
	}
	bt.Date = raw.Date
	bt.Code = raw.Code
	bt.LongSellValue = raw.LongSellVa
	bt.ShortSellWithoutMarginValue = raw.ShrtNoMrgnVa
	bt.MarginSellNewValue = raw.MrgnSellNewVa
	bt.MarginSellCloseValue = raw.MrgnSellCloseVa
	bt.LongBuyValue = raw.LongBuyVa
	bt.MarginBuyNewValue = raw.MrgnBuyNewVa
	bt.MarginBuyCloseValue = raw.MrgnBuyCloseVa
	u := &unmarshaler{}
	bt.LongSellVolume = u.integer(raw.LongSellVo)
	bt.ShortSellWithoutMarginVolume = u.integer(raw.ShrtNoMrgnVo)
	bt.MarginSellNewVolume = u.integer(raw.MrgnSellNewVo)
	bt.MarginSellCloseVolume = u.integer(raw.MrgnSellCloseVo)
	bt.LongBuyVolume = u.integer(raw.LongBuyVo)
	bt.MarginBuyNewVolume = u.integer(raw.MrgnBuyNewVo)
	bt.MarginBuyCloseVolume = u.integer(raw.MrgnBuyCloseVo)
	return u.err
}

// BreakdownTradingRequest specifies filter parameters for the BreakdownTrading API.
// Either Code or Date must be provided.
type BreakdownTradingRequest struct {
	// Code filters by security code. Required if Date is not specified.
	Code *string
	// Date filters by a specific date in YYYY-MM-DD format. It can be combined
	// with Code, but not with From or To.
	Date *string
	// From specifies the start date for a date range query (used with Code, not Date).
	From *string
	// To specifies the end date for a date range query (used with Code, not Date).
	To *string
}

type breakdownTradingParameters struct {
	BreakdownTradingRequest
	PaginationKey *string
}

func (p breakdownTradingParameters) values() (url.Values, error) {
	return codeDateRangeValues(p.Code, p.Date, p.From, p.To, p.PaginationKey)
}

// BreakdownTrading retrieves the daily breakdown of trading value and volume by
// trade type from the /markets/breakdown endpoint.
// It automatically handles pagination to fetch all matching records.
// This endpoint requires the Premium plan.
// See https://jpx-jquants.com/en/spec/mkt-breakdown for API details.
func (c *Client) BreakdownTrading(ctx context.Context, req BreakdownTradingRequest) ([]BreakdownTrading, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[BreakdownTrading], error) {
		params := breakdownTradingParameters{BreakdownTradingRequest: req, PaginationKey: paginationKey}
		return getJSON[page[BreakdownTrading]](ctx, c, "/markets/breakdown", params)
	})
}

// BreakdownTradingWithChannel retrieves the daily breakdown of trading value and volume
// and streams each record to the provided channel.
// The channel is closed when all records have been sent or an error occurs.
// On error the channel is closed and the error is returned from this method, so callers
// must check the returned error after the channel closes; ranging the channel alone will not surface it.
// This endpoint requires the Premium plan.
// See https://jpx-jquants.com/en/spec/mkt-breakdown for API details.
func (c *Client) BreakdownTradingWithChannel(ctx context.Context, req BreakdownTradingRequest, ch chan<- BreakdownTrading) error {
	return fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, paginationKey *string) (page[BreakdownTrading], error) {
		params := breakdownTradingParameters{BreakdownTradingRequest: req, PaginationKey: paginationKey}
		return getJSON[page[BreakdownTrading]](ctx, c, "/markets/breakdown", params)
	})
}

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
	type StoredRecord OutstandingShortPosition
	var raw struct {
		*StoredRecord
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
	if raw.StoredRecord != nil {
		raw.StoredRecord.Code = raw.Code
		raw.StoredRecord.FundName = raw.FundName
		raw.StoredRecord.Notes = raw.Notes
		*o = OutstandingShortPosition(*raw.StoredRecord)
		return nil
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
// At least one of Code, DisclosureDate, or CalculationDate must be provided.
// DisclosureDate, the DisclosureDateFrom/DisclosureDateTo range, and
// CalculationDate cannot be combined with each other, and the range requires Code.
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
	if p.Code == nil && p.DisclosureDate == nil && p.CalculationDate == nil {
		return nil, errors.New("code, disc_date, or calc_date is required")
	}
	hasRange := p.DisclosureDateFrom != nil || p.DisclosureDateTo != nil
	if hasRange && p.Code == nil {
		return nil, errors.New("disc_date_from and disc_date_to require code")
	}
	selectors := 0
	for _, set := range []bool{p.DisclosureDate != nil, hasRange, p.CalculationDate != nil} {
		if set {
			selectors++
		}
	}
	if selectors > 1 {
		return nil, errors.New("disc_date, disc_date_from/disc_date_to, and calc_date cannot be combined")
	}
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

// OutstandingShortPosition retrieves outstanding short position reports from the /markets/short-sale-report endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/mkt-short-sale for API details.
func (c *Client) OutstandingShortPosition(ctx context.Context, req OutstandingShortPositionRequest) ([]OutstandingShortPosition, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[OutstandingShortPosition], error) {
		params := outstandingShortPositionParameters{OutstandingShortPositionRequest: req, PaginationKey: paginationKey}
		return getJSON[page[OutstandingShortPosition]](ctx, c, "/markets/short-sale-report", params)
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
	type StoredRecord MarginAlert
	var raw struct {
		*StoredRecord
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
		ShrtOutChg    nullableNumber               `json:"ShrtOutChg"`
		LongOutChg    nullableNumber               `json:"LongOutChg"`
		ShrtNegOutChg nullableNumber               `json:"ShrtNegOutChg"`
		ShrtStdOutChg nullableNumber               `json:"ShrtStdOutChg"`
		LongNegOutChg nullableNumber               `json:"LongNegOutChg"`
		LongStdOutChg nullableNumber               `json:"LongStdOutChg"`
		ShrtOutRatio  nullableNumber               `json:"ShrtOutRatio"`
		LongOutRatio  nullableNumber               `json:"LongOutRatio"`
		TSEMrgnRegCls string                       `json:"TSEMrgnRegCls"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal margin alert: %w", err)
	}
	if raw.StoredRecord != nil {
		raw.StoredRecord.Code = raw.Code
		*m = MarginAlert(*raw.StoredRecord)
		return nil
	}
	a := &floatAccumulator{}
	fromNumber := func(v nullableNumber) *float64 {
		if a.err != nil {
			return nil
		}
		result, err := v.float64()
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
	m.ShortOutstandingChange = fromNumber(raw.ShrtOutChg)
	m.LongOutstandingChange = fromNumber(raw.LongOutChg)
	m.ShortNegotiableOutstandingChange = fromNumber(raw.ShrtNegOutChg)
	m.ShortStandardizedOutstandingChange = fromNumber(raw.ShrtStdOutChg)
	m.LongNegotiableOutstandingChange = fromNumber(raw.LongNegOutChg)
	m.LongStandardizedOutstandingChange = fromNumber(raw.LongStdOutChg)
	m.ShortOutstandingRatio = fromNumber(raw.ShrtOutRatio)
	m.LongOutstandingRatio = fromNumber(raw.LongOutRatio)
	m.TSEMarginRegulationClass = raw.TSEMrgnRegCls
	return a.err
}

// MarginAlertRequest specifies filter parameters for the MarginAlert API.
// Either Code or Date must be provided.
type MarginAlertRequest struct {
	// Code filters by security code. Required if Date is not specified.
	Code *string
	// Date filters by publication date. It can be combined with Code, but not
	// with From or To.
	Date *string
	// From specifies the start of a publication date range (used with Code, not Date).
	From *string
	// To specifies the end of a publication date range (used with Code, not Date).
	To *string
}

type marginAlertParameters struct {
	MarginAlertRequest
	PaginationKey *string
}

func (p marginAlertParameters) values() (url.Values, error) {
	return codeDateRangeValues(p.Code, p.Date, p.From, p.To, p.PaginationKey)
}

// MarginAlert retrieves margin trading alert data from the /markets/margin-alert endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/mkt-margin-alert for API details.
func (c *Client) MarginAlert(ctx context.Context, req MarginAlertRequest) ([]MarginAlert, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[MarginAlert], error) {
		params := marginAlertParameters{MarginAlertRequest: req, PaginationKey: paginationKey}
		return getJSON[page[MarginAlert]](ctx, c, "/markets/margin-alert", params)
	})
}

// TradingCalendar represents a trading calendar entry indicating whether a date is a trading day.
type TradingCalendar struct {
	// Date is the calendar date in YYYY-MM-DD format.
	Date string
	// DayType indicates the day type (0: non-business day, 1: business day,
	// 2: TSE half-day trading session, 3: non-business day with holiday trading
	// of OSE derivatives).
	DayType int8
}

func (tc *TradingCalendar) UnmarshalJSON(b []byte) error {
	type StoredRecord TradingCalendar
	var raw struct {
		*StoredRecord
		Date            string `json:"Date"`
		HolidayDivision string `json:"HolDiv"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal trading calendar: %w", err)
	}
	if raw.StoredRecord != nil {
		raw.StoredRecord.Date = raw.Date
		*tc = TradingCalendar(*raw.StoredRecord)
		return nil
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
	// HolidayDivision filters by day type (0: non-business day, 1: business day,
	// 2: TSE half-day trading session, 3: non-business day with holiday trading
	// of OSE derivatives).
	HolidayDivision *int8
	// From specifies the start date for the query in YYYY-MM-DD format.
	From *string
	// To specifies the end date for the query in YYYY-MM-DD format.
	To *string
}

type tradingCalendarParameters struct {
	TradingCalendarRequest
	PaginationKey *string
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
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

// TradingCalendar retrieves the TSE trading calendar from the /markets/calendar endpoint.
func (c *Client) TradingCalendar(ctx context.Context, req TradingCalendarRequest) ([]TradingCalendar, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[TradingCalendar], error) {
		params := tradingCalendarParameters{TradingCalendarRequest: req, PaginationKey: paginationKey}
		return getJSON[page[TradingCalendar]](ctx, c, "/markets/calendar", params)
	})
}
