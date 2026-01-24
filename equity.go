package jquants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type IssueInformation struct {
	Date               string
	Code               string
	CompanyName        string
	CompanyNameEnglish string
	Sector17Code       int8
	Sector17Name       string
	Sector33Code       string
	Sector33Name       string
	ScaleCategory      string
	MarketCode         string
	MarketName         string
	MarginCode         *int8
	MarginName         *string
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

type IssueInformationRequest struct {
	Code *string
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

type StockPrice struct {
	Date             string
	Code             string
	Open             *json.Number
	High             *json.Number
	Low              *json.Number
	Close            *json.Number
	UpperLimit       bool
	LowerLimit       bool
	Volume           *int64
	TurnoverValue    *int64
	AdjustmentFactor json.Number
	AdjustedOpen     *json.Number
	AdjustedHigh     *json.Number
	AdjustedLow      *json.Number
	AdjustedClose    *json.Number
	AdjustedVolume   *int64
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

type StockPriceRequest struct {
	Code *string
	Date *string
	From *string
	To   *string
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

func (c *Client) StockPrice(ctx context.Context, req StockPriceRequest) ([]StockPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (stockPriceResponse, error) {
		params := stockPriceParameters{StockPriceRequest: req, PaginationKey: paginationKey}
		return c.sendStockPriceRequest(ctx, params)
	})
}

func (c *Client) StockPriceWithChannel(ctx context.Context, req StockPriceRequest, ch chan<- StockPrice) error {
	return fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, paginationKey *string) (stockPriceResponse, error) {
		params := stockPriceParameters{StockPriceRequest: req, PaginationKey: paginationKey}
		return c.sendStockPriceRequest(ctx, params)
	})
}

// Morning Session Stock Prices not implemented

type TradingBalance struct {
	Sales     int64
	Purchases int64
	Total     int64
	Balance   int64
}

type InvestorType struct {
	PublishedDate              string
	StartDate                  string
	EndDate                    string
	Section                    string
	Proprietary                TradingBalance
	Brokerage                  TradingBalance
	Total                      TradingBalance
	Individuals                TradingBalance
	Foreigners                 TradingBalance
	SecuritiesCos              TradingBalance
	InvestmentTrusts           TradingBalance
	BusinessCos                TradingBalance
	OtherCos                   TradingBalance
	InsuranceCos               TradingBalance
	Banks                      TradingBalance
	TrustBanks                 TradingBalance
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
	it.Proprietary = TradingBalance{
		Sales:     int64(raw.PropSell),
		Purchases: int64(raw.PropBuy),
		Total:     int64(raw.PropTot),
		Balance:   int64(raw.PropBal),
	}
	it.Brokerage = TradingBalance{
		Sales:     int64(raw.BrkSell),
		Purchases: int64(raw.BrkBuy),
		Total:     int64(raw.BrkTot),
		Balance:   int64(raw.BrkBal),
	}
	it.Total = TradingBalance{
		Sales:     int64(raw.TotSell),
		Purchases: int64(raw.TotBuy),
		Total:     int64(raw.TotTot),
		Balance:   int64(raw.TotBal),
	}
	it.Individuals = TradingBalance{
		Sales:     int64(raw.IndSell),
		Purchases: int64(raw.IndBuy),
		Total:     int64(raw.IndTot),
		Balance:   int64(raw.IndBal),
	}
	it.Foreigners = TradingBalance{
		Sales:     int64(raw.FrgnSell),
		Purchases: int64(raw.FrgnBuy),
		Total:     int64(raw.FrgnTot),
		Balance:   int64(raw.FrgnBal),
	}
	it.SecuritiesCos = TradingBalance{
		Sales:     int64(raw.SecCoSell),
		Purchases: int64(raw.SecCoBuy),
		Total:     int64(raw.SecCoTot),
		Balance:   int64(raw.SecCoBal),
	}
	it.InvestmentTrusts = TradingBalance{
		Sales:     int64(raw.InvTrSell),
		Purchases: int64(raw.InvTrBuy),
		Total:     int64(raw.InvTrTot),
		Balance:   int64(raw.InvTrBal),
	}
	it.BusinessCos = TradingBalance{
		Sales:     int64(raw.BusCoSell),
		Purchases: int64(raw.BusCoBuy),
		Total:     int64(raw.BusCoTot),
		Balance:   int64(raw.BusCoBal),
	}
	it.OtherCos = TradingBalance{
		Sales:     int64(raw.OthCoSell),
		Purchases: int64(raw.OthCoBuy),
		Total:     int64(raw.OthCoTot),
		Balance:   int64(raw.OthCoBal),
	}
	it.InsuranceCos = TradingBalance{
		Sales:     int64(raw.InsCoSell),
		Purchases: int64(raw.InsCoBuy),
		Total:     int64(raw.InsCoTot),
		Balance:   int64(raw.InsCoBal),
	}
	it.Banks = TradingBalance{
		Sales:     int64(raw.BankSell),
		Purchases: int64(raw.BankBuy),
		Total:     int64(raw.BankTot),
		Balance:   int64(raw.BankBal),
	}
	it.TrustBanks = TradingBalance{
		Sales:     int64(raw.TrstBnkSell),
		Purchases: int64(raw.TrstBnkBuy),
		Total:     int64(raw.TrstBnkTot),
		Balance:   int64(raw.TrstBnkBal),
	}
	it.OtherFinancialInstitutions = TradingBalance{
		Sales:     int64(raw.OthFinSell),
		Purchases: int64(raw.OthFinBuy),
		Total:     int64(raw.OthFinTot),
		Balance:   int64(raw.OthFinBal),
	}
	return nil
}

type InvestorTypeRequest struct {
	Section *string
	From    *string
	To      *string
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

// InvestorType provides trading by type of investors.
// https://jpx-jquants.com/en/spec/eq-investor-types
func (c *Client) InvestorType(ctx context.Context, req InvestorTypeRequest) ([]InvestorType, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (investorTypeResponse, error) {
		params := investorTypeParameters{InvestorTypeRequest: req, PaginationKey: paginationKey}
		return c.sendInvestorTypeRequest(ctx, params)
	})
}
