package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// IssueInformation represents master data for a listed security.
// It contains company information, sector classifications, and market details.
type IssueInformation struct {
	// Date is the date of the information in YYYY-MM-DD format.
	Date string
	// Code is the security code (ticker symbol).
	Code string
	// CompanyName is the company name in Japanese.
	CompanyName string
	// CompanyNameEnglish is the company name in English.
	CompanyNameEnglish string
	// Sector17Code is the 17-sector classification code.
	Sector17Code int8
	// Sector17Name is the name of the 17-sector classification.
	Sector17Name string
	// Sector33Code is the 33-sector classification code.
	Sector33Code string
	// Sector33Name is the name of the 33-sector classification.
	Sector33Name string
	// ScaleCategory is the market capitalization scale category.
	ScaleCategory string
	// MarketCode is the market section code.
	MarketCode string
	// MarketName is the name of the market section.
	MarketName string
	// MarginCode is the margin trading classification code (nil if not applicable).
	MarginCode *int8
	// MarginName is the name of the margin trading classification.
	MarginName *string
	// ProductCategory is the product category code (JSON key "ProdCat").
	ProductCategory string
}

func (ii *IssueInformation) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date               string  `json:"Date"`
		Code               string  `json:"Code"`
		CompanyName        string  `json:"CoName"`
		CompanyNameEnglish string  `json:"CoNameEn"`
		Sector17Code       string  `json:"S17"`
		Sector17CodeName   string  `json:"S17Nm"`
		Sector33Code       string  `json:"S33"`
		Sector33CodeName   string  `json:"S33Nm"`
		ScaleCategory      string  `json:"ScaleCat"`
		MarketCode         string  `json:"Mkt"`
		MarketCodeName     string  `json:"MktNm"`
		MarginCode         *string `json:"Mrgn"`
		MarginCodeName     *string `json:"MrgnNm"`
		ProductCategory    string  `json:"ProdCat"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	ii.Date = raw.Date
	ii.Code = raw.Code
	ii.CompanyName = raw.CompanyName
	ii.CompanyNameEnglish = raw.CompanyNameEnglish
	sector17Code, err := strconv.ParseInt(raw.Sector17Code, 10, 8)
	if err != nil {
		return err
	}
	ii.Sector17Code = int8(sector17Code)
	ii.Sector17Name = raw.Sector17CodeName
	ii.Sector33Code = raw.Sector33Code
	ii.Sector33Name = raw.Sector33CodeName
	ii.ScaleCategory = raw.ScaleCategory
	ii.MarketCode = raw.MarketCode
	ii.MarketName = raw.MarketCodeName
	ii.MarginCode = nil
	if raw.MarginCode != nil {
		marginCode, err := strconv.ParseInt(*raw.MarginCode, 10, 8)
		if err != nil {
			return err
		}
		v := int8(marginCode)
		ii.MarginCode = &v
	}
	ii.MarginName = raw.MarginCodeName
	ii.ProductCategory = raw.ProductCategory
	return nil
}

// IssueInformationRequest specifies filter parameters for the IssueInformation API.
type IssueInformationRequest struct {
	// Code filters by security code. If nil, returns all securities.
	Code *string
	// Date filters by date in YYYY-MM-DD format. If nil, returns the latest data.
	Date *string
}

type issueInformationParameters struct {
	IssueInformationRequest
	PaginationKey *string
}

func (p issueInformationParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Code != nil {
		v.Add("code", *p.Code)
	}
	if p.Date != nil {
		v.Add("date", *p.Date)
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

// IssueInformation retrieves master data for listed securities from the /equities/master endpoint.
// It returns company information, sector classifications, and market details.
func (c *Client) IssueInformation(ctx context.Context, req IssueInformationRequest) ([]IssueInformation, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[IssueInformation], error) {
		params := issueInformationParameters{IssueInformationRequest: req, PaginationKey: paginationKey}
		return getJSON[page[IssueInformation]](ctx, c, "/equities/master", params)
	})
}

// StockPrice represents daily OHLCV (Open, High, Low, Close, Volume) data for a security.
// It includes both unadjusted and split-adjusted price data.
type StockPrice struct {
	// Date is the trading date in YYYY-MM-DD format.
	Date string
	// Code is the security code (ticker symbol).
	Code string
	// Open is the opening price (nil if no trading occurred).
	Open *json.Number
	// High is the highest price of the day (nil if no trading occurred).
	High *json.Number
	// Low is the lowest price of the day (nil if no trading occurred).
	Low *json.Number
	// Close is the closing price (nil if no trading occurred).
	Close *json.Number
	// UpperLimit indicates whether the stock hit the daily price limit up.
	UpperLimit bool
	// LowerLimit indicates whether the stock hit the daily price limit down.
	LowerLimit bool
	// Volume is the trading volume in shares (nil if no trading occurred).
	Volume *int64
	// TurnoverValue is the total trading value in yen (nil if no trading occurred).
	TurnoverValue *int64
	// AdjustmentFactor is the cumulative adjustment factor for stock splits.
	AdjustmentFactor json.Number
	// AdjustedOpen is the split-adjusted opening price.
	AdjustedOpen *json.Number
	// AdjustedHigh is the split-adjusted highest price.
	AdjustedHigh *json.Number
	// AdjustedLow is the split-adjusted lowest price.
	AdjustedLow *json.Number
	// AdjustedClose is the split-adjusted closing price.
	AdjustedClose *json.Number
	// AdjustedVolume is the split-adjusted trading volume, rounded by the API to
	// one decimal place, so it can be fractional.
	AdjustedVolume *json.Number
	// MorningOpen is the morning-session opening price.
	MorningOpen *json.Number
	// MorningHigh is the morning-session highest price.
	MorningHigh *json.Number
	// MorningLow is the morning-session lowest price.
	MorningLow *json.Number
	// MorningClose is the morning-session closing price.
	MorningClose *json.Number
	// MorningUpperLimit indicates whether the morning session hit the daily price limit up.
	MorningUpperLimit *bool
	// MorningLowerLimit indicates whether the morning session hit the daily price limit down.
	MorningLowerLimit *bool
	// MorningVolume is the morning-session trading volume in shares.
	MorningVolume *int64
	// MorningTurnoverValue is the morning-session trading value in yen.
	MorningTurnoverValue *int64
	// MorningAdjustedOpen is the split-adjusted morning-session opening price.
	MorningAdjustedOpen *json.Number
	// MorningAdjustedHigh is the split-adjusted morning-session highest price.
	MorningAdjustedHigh *json.Number
	// MorningAdjustedLow is the split-adjusted morning-session lowest price.
	MorningAdjustedLow *json.Number
	// MorningAdjustedClose is the split-adjusted morning-session closing price.
	MorningAdjustedClose *json.Number
	// MorningAdjustedVolume is the split-adjusted morning-session trading volume,
	// rounded by the API to one decimal place, so it can be fractional.
	MorningAdjustedVolume *json.Number
	// AfternoonOpen is the afternoon-session opening price.
	AfternoonOpen *json.Number
	// AfternoonHigh is the afternoon-session highest price.
	AfternoonHigh *json.Number
	// AfternoonLow is the afternoon-session lowest price.
	AfternoonLow *json.Number
	// AfternoonClose is the afternoon-session closing price.
	AfternoonClose *json.Number
	// AfternoonUpperLimit indicates whether the afternoon session hit the daily price limit up.
	AfternoonUpperLimit *bool
	// AfternoonLowerLimit indicates whether the afternoon session hit the daily price limit down.
	AfternoonLowerLimit *bool
	// AfternoonVolume is the afternoon-session trading volume in shares.
	AfternoonVolume *int64
	// AfternoonTurnoverValue is the afternoon-session trading value in yen.
	AfternoonTurnoverValue *int64
	// AfternoonAdjustedOpen is the split-adjusted afternoon-session opening price.
	AfternoonAdjustedOpen *json.Number
	// AfternoonAdjustedHigh is the split-adjusted afternoon-session highest price.
	AfternoonAdjustedHigh *json.Number
	// AfternoonAdjustedLow is the split-adjusted afternoon-session lowest price.
	AfternoonAdjustedLow *json.Number
	// AfternoonAdjustedClose is the split-adjusted afternoon-session closing price.
	AfternoonAdjustedClose *json.Number
	// AfternoonAdjustedVolume is the split-adjusted afternoon-session trading volume,
	// rounded by the API to one decimal place, so it can be fractional.
	AfternoonAdjustedVolume *json.Number
}

func (sp *StockPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date                    string         `json:"Date"`
		Code                    string         `json:"Code"`
		Open                    nullableNumber `json:"O"`
		High                    nullableNumber `json:"H"`
		Low                     nullableNumber `json:"L"`
		Close                   nullableNumber `json:"C"`
		UpperLimit              string         `json:"UL"`
		LowerLimit              string         `json:"LL"`
		Volume                  nullableNumber `json:"Vo"`
		TurnoverValue           nullableNumber `json:"Va"`
		AdjustmentFactor        json.Number    `json:"AdjFactor"`
		AdjustedOpen            nullableNumber `json:"AdjO"`
		AdjustedHigh            nullableNumber `json:"AdjH"`
		AdjustedLow             nullableNumber `json:"AdjL"`
		AdjustedClose           nullableNumber `json:"AdjC"`
		AdjustedVolume          nullableNumber `json:"AdjVo"`
		MorningOpen             nullableNumber `json:"MO"`
		MorningHigh             nullableNumber `json:"MH"`
		MorningLow              nullableNumber `json:"ML"`
		MorningClose            nullableNumber `json:"MC"`
		MorningUpperLimit       string         `json:"MUL"`
		MorningLowerLimit       string         `json:"MLL"`
		MorningVolume           nullableNumber `json:"MVo"`
		MorningTurnoverValue    nullableNumber `json:"MVa"`
		MorningAdjustedOpen     nullableNumber `json:"MAdjO"`
		MorningAdjustedHigh     nullableNumber `json:"MAdjH"`
		MorningAdjustedLow      nullableNumber `json:"MAdjL"`
		MorningAdjustedClose    nullableNumber `json:"MAdjC"`
		MorningAdjustedVolume   nullableNumber `json:"MAdjVo"`
		AfternoonOpen           nullableNumber `json:"AO"`
		AfternoonHigh           nullableNumber `json:"AH"`
		AfternoonLow            nullableNumber `json:"AL"`
		AfternoonClose          nullableNumber `json:"AC"`
		AfternoonUpperLimit     string         `json:"AUL"`
		AfternoonLowerLimit     string         `json:"ALL"`
		AfternoonVolume         nullableNumber `json:"AVo"`
		AfternoonTurnoverValue  nullableNumber `json:"AVa"`
		AfternoonAdjustedOpen   nullableNumber `json:"AAdjO"`
		AfternoonAdjustedHigh   nullableNumber `json:"AAdjH"`
		AfternoonAdjustedLow    nullableNumber `json:"AAdjL"`
		AfternoonAdjustedClose  nullableNumber `json:"AAdjC"`
		AfternoonAdjustedVolume nullableNumber `json:"AAdjVo"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	u := &unmarshaler{}
	upperLimit, err := unmarshalLimit(raw.UpperLimit)
	if err != nil {
		return err
	}
	lowerLimit, err := unmarshalLimit(raw.LowerLimit)
	if err != nil {
		return err
	}
	morningUpperLimit, err := unmarshalOptionalLimit(raw.MorningUpperLimit)
	if err != nil {
		return err
	}
	morningLowerLimit, err := unmarshalOptionalLimit(raw.MorningLowerLimit)
	if err != nil {
		return err
	}
	afternoonUpperLimit, err := unmarshalOptionalLimit(raw.AfternoonUpperLimit)
	if err != nil {
		return err
	}
	afternoonLowerLimit, err := unmarshalOptionalLimit(raw.AfternoonLowerLimit)
	if err != nil {
		return err
	}
	sp.Date = raw.Date
	sp.Code = raw.Code
	sp.Open = u.jsonNumber(raw.Open)
	sp.High = u.jsonNumber(raw.High)
	sp.Low = u.jsonNumber(raw.Low)
	sp.Close = u.jsonNumber(raw.Close)
	sp.UpperLimit = upperLimit
	sp.LowerLimit = lowerLimit
	sp.Volume = u.volume(raw.Volume)
	sp.TurnoverValue = u.volume(raw.TurnoverValue)
	sp.AdjustmentFactor = raw.AdjustmentFactor
	sp.AdjustedOpen = u.jsonNumber(raw.AdjustedOpen)
	sp.AdjustedHigh = u.jsonNumber(raw.AdjustedHigh)
	sp.AdjustedLow = u.jsonNumber(raw.AdjustedLow)
	sp.AdjustedClose = u.jsonNumber(raw.AdjustedClose)
	sp.AdjustedVolume = u.jsonNumber(raw.AdjustedVolume)
	sp.MorningOpen = u.jsonNumber(raw.MorningOpen)
	sp.MorningHigh = u.jsonNumber(raw.MorningHigh)
	sp.MorningLow = u.jsonNumber(raw.MorningLow)
	sp.MorningClose = u.jsonNumber(raw.MorningClose)
	sp.MorningUpperLimit = morningUpperLimit
	sp.MorningLowerLimit = morningLowerLimit
	sp.MorningVolume = u.volume(raw.MorningVolume)
	sp.MorningTurnoverValue = u.volume(raw.MorningTurnoverValue)
	sp.MorningAdjustedOpen = u.jsonNumber(raw.MorningAdjustedOpen)
	sp.MorningAdjustedHigh = u.jsonNumber(raw.MorningAdjustedHigh)
	sp.MorningAdjustedLow = u.jsonNumber(raw.MorningAdjustedLow)
	sp.MorningAdjustedClose = u.jsonNumber(raw.MorningAdjustedClose)
	sp.MorningAdjustedVolume = u.jsonNumber(raw.MorningAdjustedVolume)
	sp.AfternoonOpen = u.jsonNumber(raw.AfternoonOpen)
	sp.AfternoonHigh = u.jsonNumber(raw.AfternoonHigh)
	sp.AfternoonLow = u.jsonNumber(raw.AfternoonLow)
	sp.AfternoonClose = u.jsonNumber(raw.AfternoonClose)
	sp.AfternoonUpperLimit = afternoonUpperLimit
	sp.AfternoonLowerLimit = afternoonLowerLimit
	sp.AfternoonVolume = u.volume(raw.AfternoonVolume)
	sp.AfternoonTurnoverValue = u.volume(raw.AfternoonTurnoverValue)
	sp.AfternoonAdjustedOpen = u.jsonNumber(raw.AfternoonAdjustedOpen)
	sp.AfternoonAdjustedHigh = u.jsonNumber(raw.AfternoonAdjustedHigh)
	sp.AfternoonAdjustedLow = u.jsonNumber(raw.AfternoonAdjustedLow)
	sp.AfternoonAdjustedClose = u.jsonNumber(raw.AfternoonAdjustedClose)
	sp.AfternoonAdjustedVolume = u.jsonNumber(raw.AfternoonAdjustedVolume)
	return u.err
}

func unmarshalLimit(s string) (bool, error) {
	switch s {
	case "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, fmt.Errorf("unknown value: %s", s)
	}
}

func unmarshalOptionalLimit(s string) (*bool, error) {
	if s == "" {
		return nil, nil
	}
	value, err := unmarshalLimit(s)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// StockPriceRequest specifies filter parameters for the StockPrice API.
// Either Code or Date must be provided.
type StockPriceRequest struct {
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

type stockPriceParameters struct {
	StockPriceRequest
	PaginationKey *string
}

func (p stockPriceParameters) values() (url.Values, error) {
	return codeDateRangeValues(p.Code, p.Date, p.From, p.To, p.PaginationKey)
}

// StockPrice retrieves daily stock prices from the /equities/bars/daily endpoint.
// It automatically handles pagination to fetch all matching records.
func (c *Client) StockPrice(ctx context.Context, req StockPriceRequest) ([]StockPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[StockPrice], error) {
		params := stockPriceParameters{StockPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[page[StockPrice]](ctx, c, "/equities/bars/daily", params)
	})
}

// StockPriceWithChannel retrieves daily stock prices and streams each record to the provided channel.
// The channel is closed when all records have been sent or an error occurs.
// On error the channel is closed and the error is returned from this method, so callers
// must check the returned error after the channel closes; ranging the channel alone will not surface it.
func (c *Client) StockPriceWithChannel(ctx context.Context, req StockPriceRequest, ch chan<- StockPrice) error {
	return fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, paginationKey *string) (page[StockPrice], error) {
		params := stockPriceParameters{StockPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[page[StockPrice]](ctx, c, "/equities/bars/daily", params)
	})
}

// MinuteStockPrice represents one-minute OHLCV data for a security.
type MinuteStockPrice struct {
	// Date is the trading date in YYYY-MM-DD format (JSON key "Date").
	Date string
	// Time is the start time of the one-minute bar in HH:MM format (JSON key "Time").
	Time string
	// Code is the security code (JSON key "Code").
	Code string
	// Open is the opening price of the minute (nil if no trading occurred) (JSON key "O").
	Open *json.Number
	// High is the highest price of the minute (nil if no trading occurred) (JSON key "H").
	High *json.Number
	// Low is the lowest price of the minute (nil if no trading occurred) (JSON key "L").
	Low *json.Number
	// Close is the closing price of the minute (nil if no trading occurred) (JSON key "C").
	Close *json.Number
	// Volume is the trading volume in shares for the minute (nil if no trading occurred) (JSON key "Vo").
	Volume *int64
	// TurnoverValue is the total trading value in yen for the minute (nil if no trading occurred) (JSON key "Va").
	TurnoverValue *int64
}

func (m *MinuteStockPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date string         `json:"Date"`
		Time string         `json:"Time"`
		Code string         `json:"Code"`
		Open nullableNumber `json:"O"`
		High nullableNumber `json:"H"`
		Low  nullableNumber `json:"L"`
		C    nullableNumber `json:"C"`
		Vo   nullableNumber `json:"Vo"`
		Va   nullableNumber `json:"Va"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal minute stock price: %w", err)
	}
	u := &unmarshaler{}
	m.Date = raw.Date
	m.Time = raw.Time
	m.Code = raw.Code
	m.Open = u.jsonNumber(raw.Open)
	m.High = u.jsonNumber(raw.High)
	m.Low = u.jsonNumber(raw.Low)
	m.Close = u.jsonNumber(raw.C)
	m.Volume = u.volume(raw.Vo)
	m.TurnoverValue = u.volume(raw.Va)
	return u.err
}

// MinuteStockPriceRequest specifies filter parameters for the MinuteStockPrice API.
// Either Code or Date must be provided.
type MinuteStockPriceRequest struct {
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

type minuteStockPriceParameters struct {
	MinuteStockPriceRequest
	PaginationKey *string
}

func (p minuteStockPriceParameters) values() (url.Values, error) {
	return codeDateRangeValues(p.Code, p.Date, p.From, p.To, p.PaginationKey)
}

// MinuteStockPrice retrieves one-minute stock prices from the /equities/bars/minute endpoint.
// It automatically handles pagination to fetch all matching records.
// This endpoint requires the minute-bars add-on plan.
// See https://jpx-jquants.com/en/spec/eq-bars-minute for API details.
func (c *Client) MinuteStockPrice(ctx context.Context, req MinuteStockPriceRequest) ([]MinuteStockPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[MinuteStockPrice], error) {
		params := minuteStockPriceParameters{MinuteStockPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[page[MinuteStockPrice]](ctx, c, "/equities/bars/minute", params)
	})
}

// MinuteStockPriceWithChannel retrieves one-minute stock prices and streams each record to the provided channel.
// The channel is closed when all records have been sent or an error occurs.
// On error the channel is closed and the error is returned from this method, so callers
// must check the returned error after the channel closes; ranging the channel alone will not surface it.
// This endpoint requires the minute-bars add-on plan.
// See https://jpx-jquants.com/en/spec/eq-bars-minute for API details.
func (c *Client) MinuteStockPriceWithChannel(ctx context.Context, req MinuteStockPriceRequest, ch chan<- MinuteStockPrice) error {
	return fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, paginationKey *string) (page[MinuteStockPrice], error) {
		params := minuteStockPriceParameters{MinuteStockPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[page[MinuteStockPrice]](ctx, c, "/equities/bars/minute", params)
	})
}

// MorningSessionStockPrice represents the current day's morning-session (前場)
// OHLCV data for a security.
type MorningSessionStockPrice struct {
	// Date is the trading date in YYYY-MM-DD format (JSON key "Date").
	Date string
	// Code is the security code (JSON key "Code").
	Code string
	// Open is the morning-session opening price (nil if no trading occurred) (JSON key "MO").
	Open *json.Number
	// High is the morning-session highest price (nil if no trading occurred) (JSON key "MH").
	High *json.Number
	// Low is the morning-session lowest price (nil if no trading occurred) (JSON key "ML").
	Low *json.Number
	// Close is the morning-session closing price (nil if no trading occurred) (JSON key "MC").
	Close *json.Number
	// Volume is the morning-session trading volume in shares (nil if no trading occurred) (JSON key "MVo").
	Volume *int64
	// TurnoverValue is the morning-session total trading value in yen (nil if no trading occurred) (JSON key "MVa").
	TurnoverValue *int64
}

func (m *MorningSessionStockPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date string         `json:"Date"`
		Code string         `json:"Code"`
		Open nullableNumber `json:"MO"`
		High nullableNumber `json:"MH"`
		Low  nullableNumber `json:"ML"`
		C    nullableNumber `json:"MC"`
		Vo   nullableNumber `json:"MVo"`
		Va   nullableNumber `json:"MVa"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal morning session stock price: %w", err)
	}
	u := &unmarshaler{}
	m.Date = raw.Date
	m.Code = raw.Code
	m.Open = u.jsonNumber(raw.Open)
	m.High = u.jsonNumber(raw.High)
	m.Low = u.jsonNumber(raw.Low)
	m.Close = u.jsonNumber(raw.C)
	m.Volume = u.volume(raw.Vo)
	m.TurnoverValue = u.volume(raw.Va)
	return u.err
}

// MorningSessionStockPriceRequest specifies filter parameters for the MorningSessionStockPrice API.
type MorningSessionStockPriceRequest struct {
	// Code filters by security code. If nil, returns all securities.
	Code *string
}

type morningSessionStockPriceParameters struct {
	MorningSessionStockPriceRequest
	PaginationKey *string
}

func (p morningSessionStockPriceParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Code != nil {
		v.Add("code", *p.Code)
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

// MorningSessionStockPrice retrieves the current day's morning-session OHLCV data
// from the /equities/bars/daily/am endpoint.
// It automatically handles pagination to fetch all matching records.
// This endpoint requires the Premium plan.
// Outside the morning-session publication window the API responds with HTTP 210,
// which is surfaced as a NoContent error.
// See https://jpx-jquants.com/en/spec/eq-bars-daily-am for API details.
func (c *Client) MorningSessionStockPrice(ctx context.Context, req MorningSessionStockPriceRequest) ([]MorningSessionStockPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[MorningSessionStockPrice], error) {
		params := morningSessionStockPriceParameters{MorningSessionStockPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[page[MorningSessionStockPrice]](ctx, c, "/equities/bars/daily/am", params)
	})
}

// EarningsCalendar represents a scheduled or announced earnings release for a listed company.
type EarningsCalendar struct {
	// Date is the scheduled announcement date in YYYY-MM-DD format, or "" if undecided (JSON key "Date").
	Date string
	// Code is the security code (JSON key "Code").
	Code string
	// CompanyName is the company name in Japanese (JSON key "CoName").
	CompanyName string
	// FiscalYear is the fiscal year-end description (JSON key "FY").
	FiscalYear string
	// SectorName is the sector name in Japanese (JSON key "SectorNm").
	SectorName string
	// FiscalQuarter is the fiscal quarter description (JSON key "FQ").
	FiscalQuarter string
	// Section is the market section name in Japanese (JSON key "Section").
	Section string
}

func (e *EarningsCalendar) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date     string `json:"Date"`
		Code     string `json:"Code"`
		CoName   string `json:"CoName"`
		FY       string `json:"FY"`
		SectorNm string `json:"SectorNm"`
		FQ       string `json:"FQ"`
		Section  string `json:"Section"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal earnings calendar: %w", err)
	}
	e.Date = raw.Date
	e.Code = raw.Code
	e.CompanyName = raw.CoName
	e.FiscalYear = raw.FY
	e.SectorName = raw.SectorNm
	e.FiscalQuarter = raw.FQ
	e.Section = raw.Section
	return nil
}

// EarningsCalendarRequest specifies filter parameters for the EarningsCalendar API.
// The endpoint takes no filter parameters.
type EarningsCalendarRequest struct{}

type earningsCalendarParameters struct {
	EarningsCalendarRequest
	PaginationKey *string
}

func (p earningsCalendarParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

// EarningsCalendar retrieves the earnings announcement calendar from the /equities/earnings-calendar endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/eq-earnings-cal for API details.
func (c *Client) EarningsCalendar(ctx context.Context, req EarningsCalendarRequest) ([]EarningsCalendar, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[EarningsCalendar], error) {
		params := earningsCalendarParameters{EarningsCalendarRequest: req, PaginationKey: paginationKey}
		return getJSON[page[EarningsCalendar]](ctx, c, "/equities/earnings-calendar", params)
	})
}

// TradingBalance represents trading activity metrics for a specific investor type.
// All values are in units of 1,000 shares.
type TradingBalance struct {
	// Sales is the total sell volume.
	Sales int64
	// Purchases is the total buy volume.
	Purchases int64
	// Total is the sum of sales and purchases.
	Total int64
	// Balance is the net position (Purchases - Sales).
	Balance int64
}

func newTradingBalance(u *unmarshaler, sell, buy, total, balance nullableNumber) TradingBalance {
	return TradingBalance{
		Sales:     u.integer(sell),
		Purchases: u.integer(buy),
		Total:     u.integer(total),
		Balance:   u.integer(balance),
	}
}

// InvestorType represents weekly trading data broken down by investor category.
// It shows the buying and selling activity of different market participants.
type InvestorType struct {
	// PublishedDate is the publication date of the data.
	PublishedDate string
	// StartDate is the start of the reporting period.
	StartDate string
	// EndDate is the end of the reporting period.
	EndDate string
	// Section is the market section (e.g., "TSE1st", "TSE2nd").
	Section string
	// Proprietary is trading by securities companies for their own account.
	Proprietary TradingBalance
	// Brokerage is trading by securities companies on behalf of clients.
	Brokerage TradingBalance
	// Total is the aggregate trading across all investor types.
	Total TradingBalance
	// Individuals is trading by retail investors.
	Individuals TradingBalance
	// Foreigners is trading by foreign investors.
	Foreigners TradingBalance
	// SecuritiesCos is trading by securities companies.
	SecuritiesCos TradingBalance
	// InvestmentTrusts is trading by investment trusts.
	InvestmentTrusts TradingBalance
	// BusinessCos is trading by business corporations.
	BusinessCos TradingBalance
	// OtherCos is trading by other corporations.
	OtherCos TradingBalance
	// InsuranceCos is trading by insurance companies.
	InsuranceCos TradingBalance
	// Banks is trading by banks.
	Banks TradingBalance
	// TrustBanks is trading by trust banks.
	TrustBanks TradingBalance
	// OtherFinancialInstitutions is trading by other financial institutions.
	OtherFinancialInstitutions TradingBalance
}

func (it *InvestorType) UnmarshalJSON(b []byte) error {
	var raw struct {
		PubDate     string         `json:"PubDate"`
		StDate      string         `json:"StDate"`
		EnDate      string         `json:"EnDate"`
		Section     string         `json:"Section"`
		PropSell    nullableNumber `json:"PropSell"`
		PropBuy     nullableNumber `json:"PropBuy"`
		PropTot     nullableNumber `json:"PropTot"`
		PropBal     nullableNumber `json:"PropBal"`
		BrkSell     nullableNumber `json:"BrkSell"`
		BrkBuy      nullableNumber `json:"BrkBuy"`
		BrkTot      nullableNumber `json:"BrkTot"`
		BrkBal      nullableNumber `json:"BrkBal"`
		TotSell     nullableNumber `json:"TotSell"`
		TotBuy      nullableNumber `json:"TotBuy"`
		TotTot      nullableNumber `json:"TotTot"`
		TotBal      nullableNumber `json:"TotBal"`
		IndSell     nullableNumber `json:"IndSell"`
		IndBuy      nullableNumber `json:"IndBuy"`
		IndTot      nullableNumber `json:"IndTot"`
		IndBal      nullableNumber `json:"IndBal"`
		FrgnSell    nullableNumber `json:"FrgnSell"`
		FrgnBuy     nullableNumber `json:"FrgnBuy"`
		FrgnTot     nullableNumber `json:"FrgnTot"`
		FrgnBal     nullableNumber `json:"FrgnBal"`
		SecCoSell   nullableNumber `json:"SecCoSell"`
		SecCoBuy    nullableNumber `json:"SecCoBuy"`
		SecCoTot    nullableNumber `json:"SecCoTot"`
		SecCoBal    nullableNumber `json:"SecCoBal"`
		InvTrSell   nullableNumber `json:"InvTrSell"`
		InvTrBuy    nullableNumber `json:"InvTrBuy"`
		InvTrTot    nullableNumber `json:"InvTrTot"`
		InvTrBal    nullableNumber `json:"InvTrBal"`
		BusCoSell   nullableNumber `json:"BusCoSell"`
		BusCoBuy    nullableNumber `json:"BusCoBuy"`
		BusCoTot    nullableNumber `json:"BusCoTot"`
		BusCoBal    nullableNumber `json:"BusCoBal"`
		OthCoSell   nullableNumber `json:"OthCoSell"`
		OthCoBuy    nullableNumber `json:"OthCoBuy"`
		OthCoTot    nullableNumber `json:"OthCoTot"`
		OthCoBal    nullableNumber `json:"OthCoBal"`
		InsCoSell   nullableNumber `json:"InsCoSell"`
		InsCoBuy    nullableNumber `json:"InsCoBuy"`
		InsCoTot    nullableNumber `json:"InsCoTot"`
		InsCoBal    nullableNumber `json:"InsCoBal"`
		BankSell    nullableNumber `json:"BankSell"`
		BankBuy     nullableNumber `json:"BankBuy"`
		BankTot     nullableNumber `json:"BankTot"`
		BankBal     nullableNumber `json:"BankBal"`
		TrstBnkSell nullableNumber `json:"TrstBnkSell"`
		TrstBnkBuy  nullableNumber `json:"TrstBnkBuy"`
		TrstBnkTot  nullableNumber `json:"TrstBnkTot"`
		TrstBnkBal  nullableNumber `json:"TrstBnkBal"`
		OthFinSell  nullableNumber `json:"OthFinSell"`
		OthFinBuy   nullableNumber `json:"OthFinBuy"`
		OthFinTot   nullableNumber `json:"OthFinTot"`
		OthFinBal   nullableNumber `json:"OthFinBal"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	it.PublishedDate = raw.PubDate
	it.StartDate = raw.StDate
	it.EndDate = raw.EnDate
	it.Section = raw.Section
	u := &unmarshaler{}
	it.Proprietary = newTradingBalance(u, raw.PropSell, raw.PropBuy, raw.PropTot, raw.PropBal)
	it.Brokerage = newTradingBalance(u, raw.BrkSell, raw.BrkBuy, raw.BrkTot, raw.BrkBal)
	it.Total = newTradingBalance(u, raw.TotSell, raw.TotBuy, raw.TotTot, raw.TotBal)
	it.Individuals = newTradingBalance(u, raw.IndSell, raw.IndBuy, raw.IndTot, raw.IndBal)
	it.Foreigners = newTradingBalance(u, raw.FrgnSell, raw.FrgnBuy, raw.FrgnTot, raw.FrgnBal)
	it.SecuritiesCos = newTradingBalance(u, raw.SecCoSell, raw.SecCoBuy, raw.SecCoTot, raw.SecCoBal)
	it.InvestmentTrusts = newTradingBalance(u, raw.InvTrSell, raw.InvTrBuy, raw.InvTrTot, raw.InvTrBal)
	it.BusinessCos = newTradingBalance(u, raw.BusCoSell, raw.BusCoBuy, raw.BusCoTot, raw.BusCoBal)
	it.OtherCos = newTradingBalance(u, raw.OthCoSell, raw.OthCoBuy, raw.OthCoTot, raw.OthCoBal)
	it.InsuranceCos = newTradingBalance(u, raw.InsCoSell, raw.InsCoBuy, raw.InsCoTot, raw.InsCoBal)
	it.Banks = newTradingBalance(u, raw.BankSell, raw.BankBuy, raw.BankTot, raw.BankBal)
	it.TrustBanks = newTradingBalance(u, raw.TrstBnkSell, raw.TrstBnkBuy, raw.TrstBnkTot, raw.TrstBnkBal)
	it.OtherFinancialInstitutions = newTradingBalance(u, raw.OthFinSell, raw.OthFinBuy, raw.OthFinTot, raw.OthFinBal)
	return u.err
}

// InvestorTypeRequest specifies filter parameters for the InvestorType API.
type InvestorTypeRequest struct {
	// Section filters by market section (e.g., "TSE1st", "TSE2nd").
	Section *string
	// From specifies the start date for the query in YYYY-MM-DD format.
	From *string
	// To specifies the end date for the query in YYYY-MM-DD format.
	To *string
}

type investorTypeParameters struct {
	InvestorTypeRequest
	PaginationKey *string
}

func (p investorTypeParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Section != nil {
		v.Add("section", *p.Section)
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

// InvestorType retrieves weekly trading data by investor type from the /equities/investor-types endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/eq-investor-types for API details.
func (c *Client) InvestorType(ctx context.Context, req InvestorTypeRequest) ([]InvestorType, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (page[InvestorType], error) {
		params := investorTypeParameters{InvestorTypeRequest: req, PaginationKey: paginationKey}
		return getJSON[page[InvestorType]](ctx, c, "/equities/investor-types", params)
	})
}
