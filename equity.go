package jquants

import (
	"context"
	"encoding/json"
	"errors"
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
	if raw.MarginCode != nil {
		marginCode, err := strconv.ParseInt(*raw.MarginCode, 10, 8)
		if err != nil {
			return err
		}
		v := int8(marginCode)
		ii.MarginCode = &v
	}
	ii.MarginName = raw.MarginCodeName
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
}

func (p issueInformationParameters) values() (url.Values, error) {
	v := url.Values{}
	if p.Code != nil {
		v.Add("code", *p.Code)
	}
	if p.Date != nil {
		v.Add("date", *p.Date)
	}
	return v, nil
}

type issueInformationResponse struct {
	Information []IssueInformation `json:"data"`
}

// IssueInformation retrieves master data for listed securities from the /equities/master endpoint.
// It returns company information, sector classifications, and market details.
func (c *Client) IssueInformation(ctx context.Context, req IssueInformationRequest) ([]IssueInformation, error) {
	var r issueInformationResponse
	params := issueInformationParameters{req}
	resp, err := c.sendRequest(ctx, "/equities/master", params)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp)
	}
	if err = decodeResponse(resp, &r); err != nil {
		return nil, fmt.Errorf("failed to decode HTTP response: %w", err)
	}
	return r.Information, nil
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
	// AdjustedVolume is the split-adjusted trading volume.
	AdjustedVolume *int64
}

func (sp *StockPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date             string       `json:"Date"`
		Code             string       `json:"Code"`
		Open             *json.Number `json:"O"`
		High             *json.Number `json:"H"`
		Low              *json.Number `json:"L"`
		Close            *json.Number `json:"C"`
		UpperLimit       string       `json:"UL"`
		LowerLimit       string       `json:"LL"`
		Volume           *float64     `json:"Vo"`
		TurnoverValue    *float64     `json:"Va"`
		AdjustmentFactor json.Number  `json:"AdjFactor"`
		AdjustedOpen     *json.Number `json:"AdjO"`
		AdjustedHigh     *json.Number `json:"AdjH"`
		AdjustedLow      *json.Number `json:"AdjL"`
		AdjustedClose    *json.Number `json:"AdjC"`
		AdjustedVolume   *float64     `json:"AdjVo"`
	}
	var volume, turnoverValue *int64
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	upperLimit, err := unmarshalLimit(raw.UpperLimit)
	if err != nil {
		return err
	}
	lowerLimit, err := unmarshalLimit(raw.LowerLimit)
	if err != nil {
		return err
	}
	if raw.Volume != nil {
		v := int64(*raw.Volume)
		volume = &v
	}
	if raw.TurnoverValue != nil {
		v := int64(*raw.TurnoverValue)
		turnoverValue = &v
	}
	var adjustedVolume *int64
	if raw.AdjustedVolume != nil {
		v := int64(*raw.AdjustedVolume)
		adjustedVolume = &v
	}
	sp.Date = raw.Date
	sp.Code = raw.Code
	sp.Open = raw.Open
	sp.High = raw.High
	sp.Low = raw.Low
	sp.Close = raw.Close
	sp.UpperLimit = upperLimit
	sp.LowerLimit = lowerLimit
	sp.Volume = volume
	sp.TurnoverValue = turnoverValue
	sp.AdjustmentFactor = raw.AdjustmentFactor
	sp.AdjustedOpen = raw.AdjustedOpen
	sp.AdjustedHigh = raw.AdjustedHigh
	sp.AdjustedLow = raw.AdjustedLow
	sp.AdjustedClose = raw.AdjustedClose
	sp.AdjustedVolume = adjustedVolume
	return nil
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

// StockPriceRequest specifies filter parameters for the StockPrice API.
// Either Code or Date must be provided.
type StockPriceRequest struct {
	// Code filters by security code. Required if Date is not specified.
	Code *string
	// Date filters by a specific date in YYYY-MM-DD format. If specified, Code is ignored.
	Date *string
	// From specifies the start date for a date range query (used with Code).
	From *string
	// To specifies the end date for a date range query (used with Code).
	To *string
}

type stockPriceParameters struct {
	StockPriceRequest
	PaginationKey *string
}

func (p stockPriceParameters) values() (url.Values, error) {
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

type stockPriceResponse struct {
	Data          []StockPrice `json:"data"`
	PaginationKey *string      `json:"pagination_key"`
}

func (r stockPriceResponse) getData() []StockPrice   { return r.Data }
func (r stockPriceResponse) getPaginationKey() *string { return r.PaginationKey }

func (c *Client) sendStockPriceRequest(ctx context.Context, params stockPriceParameters) (stockPriceResponse, error) {
	var r stockPriceResponse
	resp, err := c.sendRequest(ctx, "/equities/bars/daily", params)
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

// StockPrice retrieves daily stock prices from the /equities/bars/daily endpoint.
// It automatically handles pagination to fetch all matching records.
func (c *Client) StockPrice(ctx context.Context, req StockPriceRequest) ([]StockPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (stockPriceResponse, error) {
		params := stockPriceParameters{StockPriceRequest: req, PaginationKey: paginationKey}
		return c.sendStockPriceRequest(ctx, params)
	})
}

// StockPriceWithChannel retrieves daily stock prices and streams each record to the provided channel.
// The channel is closed when all records have been sent or an error occurs.
func (c *Client) StockPriceWithChannel(ctx context.Context, req StockPriceRequest, ch chan<- StockPrice) error {
	return fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, paginationKey *string) (stockPriceResponse, error) {
		params := stockPriceParameters{StockPriceRequest: req, PaginationKey: paginationKey}
		return c.sendStockPriceRequest(ctx, params)
	})
}

// Morning Session Stock Prices not implemented

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

func newTradingBalance(sell, buy, total, balance float64) TradingBalance {
	return TradingBalance{
		Sales:     int64(sell),
		Purchases: int64(buy),
		Total:     int64(total),
		Balance:   int64(balance),
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
		PubDate     string  `json:"PubDate"`
		StDate      string  `json:"StDate"`
		EnDate      string  `json:"EnDate"`
		Section     string  `json:"Section"`
		PropSell    float64 `json:"PropSell"`
		PropBuy     float64 `json:"PropBuy"`
		PropTot     float64 `json:"PropTot"`
		PropBal     float64 `json:"PropBal"`
		BrkSell     float64 `json:"BrkSell"`
		BrkBuy      float64 `json:"BrkBuy"`
		BrkTot      float64 `json:"BrkTot"`
		BrkBal      float64 `json:"BrkBal"`
		TotSell     float64 `json:"TotSell"`
		TotBuy      float64 `json:"TotBuy"`
		TotTot      float64 `json:"TotTot"`
		TotBal      float64 `json:"TotBal"`
		IndSell     float64 `json:"IndSell"`
		IndBuy      float64 `json:"IndBuy"`
		IndTot      float64 `json:"IndTot"`
		IndBal      float64 `json:"IndBal"`
		FrgnSell    float64 `json:"FrgnSell"`
		FrgnBuy     float64 `json:"FrgnBuy"`
		FrgnTot     float64 `json:"FrgnTot"`
		FrgnBal     float64 `json:"FrgnBal"`
		SecCoSell   float64 `json:"SecCoSell"`
		SecCoBuy    float64 `json:"SecCoBuy"`
		SecCoTot    float64 `json:"SecCoTot"`
		SecCoBal    float64 `json:"SecCoBal"`
		InvTrSell   float64 `json:"InvTrSell"`
		InvTrBuy    float64 `json:"InvTrBuy"`
		InvTrTot    float64 `json:"InvTrTot"`
		InvTrBal    float64 `json:"InvTrBal"`
		BusCoSell   float64 `json:"BusCoSell"`
		BusCoBuy    float64 `json:"BusCoBuy"`
		BusCoTot    float64 `json:"BusCoTot"`
		BusCoBal    float64 `json:"BusCoBal"`
		OthCoSell   float64 `json:"OthCoSell"`
		OthCoBuy    float64 `json:"OthCoBuy"`
		OthCoTot    float64 `json:"OthCoTot"`
		OthCoBal    float64 `json:"OthCoBal"`
		InsCoSell   float64 `json:"InsCoSell"`
		InsCoBuy    float64 `json:"InsCoBuy"`
		InsCoTot    float64 `json:"InsCoTot"`
		InsCoBal    float64 `json:"InsCoBal"`
		BankSell    float64 `json:"BankSell"`
		BankBuy     float64 `json:"BankBuy"`
		BankTot     float64 `json:"BankTot"`
		BankBal     float64 `json:"BankBal"`
		TrstBnkSell float64 `json:"TrstBnkSell"`
		TrstBnkBuy  float64 `json:"TrstBnkBuy"`
		TrstBnkTot  float64 `json:"TrstBnkTot"`
		TrstBnkBal  float64 `json:"TrstBnkBal"`
		OthFinSell  float64 `json:"OthFinSell"`
		OthFinBuy   float64 `json:"OthFinBuy"`
		OthFinTot   float64 `json:"OthFinTot"`
		OthFinBal   float64 `json:"OthFinBal"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	it.PublishedDate = raw.PubDate
	it.StartDate = raw.StDate
	it.EndDate = raw.EnDate
	it.Section = raw.Section
	it.Proprietary = newTradingBalance(raw.PropSell, raw.PropBuy, raw.PropTot, raw.PropBal)
	it.Brokerage = newTradingBalance(raw.BrkSell, raw.BrkBuy, raw.BrkTot, raw.BrkBal)
	it.Total = newTradingBalance(raw.TotSell, raw.TotBuy, raw.TotTot, raw.TotBal)
	it.Individuals = newTradingBalance(raw.IndSell, raw.IndBuy, raw.IndTot, raw.IndBal)
	it.Foreigners = newTradingBalance(raw.FrgnSell, raw.FrgnBuy, raw.FrgnTot, raw.FrgnBal)
	it.SecuritiesCos = newTradingBalance(raw.SecCoSell, raw.SecCoBuy, raw.SecCoTot, raw.SecCoBal)
	it.InvestmentTrusts = newTradingBalance(raw.InvTrSell, raw.InvTrBuy, raw.InvTrTot, raw.InvTrBal)
	it.BusinessCos = newTradingBalance(raw.BusCoSell, raw.BusCoBuy, raw.BusCoTot, raw.BusCoBal)
	it.OtherCos = newTradingBalance(raw.OthCoSell, raw.OthCoBuy, raw.OthCoTot, raw.OthCoBal)
	it.InsuranceCos = newTradingBalance(raw.InsCoSell, raw.InsCoBuy, raw.InsCoTot, raw.InsCoBal)
	it.Banks = newTradingBalance(raw.BankSell, raw.BankBuy, raw.BankTot, raw.BankBal)
	it.TrustBanks = newTradingBalance(raw.TrstBnkSell, raw.TrstBnkBuy, raw.TrstBnkTot, raw.TrstBnkBal)
	it.OtherFinancialInstitutions = newTradingBalance(raw.OthFinSell, raw.OthFinBuy, raw.OthFinTot, raw.OthFinBal)
	return nil
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

type investorTypeResponse struct {
	Data          []InvestorType `json:"data"`
	PaginationKey *string        `json:"pagination_key"`
}

func (r investorTypeResponse) getData() []InvestorType { return r.Data }
func (r investorTypeResponse) getPaginationKey() *string { return r.PaginationKey }

func (c *Client) sendInvestorTypeRequest(ctx context.Context, params investorTypeParameters) (investorTypeResponse, error) {
	var r investorTypeResponse
	resp, err := c.sendRequest(ctx, "/equities/investor-types", params)
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

// InvestorType retrieves weekly trading data by investor type from the /equities/investor-types endpoint.
// It automatically handles pagination to fetch all matching records.
// See https://jpx-jquants.com/en/spec/eq-investor-types for API details.
func (c *Client) InvestorType(ctx context.Context, req InvestorTypeRequest) ([]InvestorType, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (investorTypeResponse, error) {
		params := investorTypeParameters{InvestorTypeRequest: req, PaginationKey: paginationKey}
		return c.sendInvestorTypeRequest(ctx, params)
	})
}
