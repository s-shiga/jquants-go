package jquants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"
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
			return nil, fmt.Errorf("code or date is required")
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

func (c *Client) sendStockPriceRequest(ctx context.Context, param stockPriceParameters) (*stockPriceResponse, error) {
	var r stockPriceResponse
	resp, err := c.sendRequest(ctx, "/equities/bars/daily", param)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp)
	}
	if err = decodeResponse(resp, &r); err != nil {
		return nil, fmt.Errorf("failed to decode HTTP response: %w", err)
	}
	return &r, nil
}

func (c *Client) StockPrice(ctx context.Context, req StockPriceRequest) ([]StockPrice, error) {
	var data = make([]StockPrice, 0)
	var paginationKey *string
	ctx, cancel := context.WithTimeout(ctx, c.LoopTimeout)
	defer cancel()
	for {
		params := stockPriceParameters{req, paginationKey}
		resp, err := c.sendStockPriceRequest(ctx, params)
		if err != nil {
			if errors.As(err, &InternalServerError{}) {
				slog.Warn("Retrying HTTP request", "error", err.Error())
				time.Sleep(c.RetryInterval)
				continue
			} else {
				return nil, fmt.Errorf("failed to send stock price request: %w", err)
			}
		}
		data = append(data, resp.Data...)
		paginationKey = resp.PaginationKey
		if resp.PaginationKey == nil {
			break
		}
	}
	return data, nil
}

func (c *Client) StockPriceWithChannel(ctx context.Context, req StockPriceRequest, ch chan<- StockPrice) error {
	var paginationKey *string
	ctx, cancel := context.WithTimeout(ctx, c.LoopTimeout)
	defer cancel()
	for {
		params := stockPriceParameters{StockPriceRequest: req, PaginationKey: paginationKey}
		resp, err := c.sendStockPriceRequest(ctx, params)
		if err != nil {
			if errors.As(err, &InternalServerError{}) {
				slog.Warn("Retrying HTTP request", "error", err.Error())
				time.Sleep(c.RetryInterval)
				continue
			} else {
				return fmt.Errorf("failed to send stock price request: %w", err)
			}
		}
		for _, d := range resp.Data {
			ch <- d
		}
		paginationKey = resp.PaginationKey
		if resp.PaginationKey == nil {
			break
		}
	}
	close(ch)
	return nil
}

// Morning Session Stock Prices not implemented

type TradingBalance struct {
	Sales     int64
	Purchases int64
	Total     int64
	Balance   int64
}

type StockTradingValue struct {
	PublishedDate              string
	StartDate                  string
	EndDate                    string
	Section                    string
	Proprietary                TradingBalance
	Brokerage                  TradingBalance
	Net                        TradingBalance
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

func (stv *StockTradingValue) UnmarshalJSON(b []byte) error {
	var raw struct {
		PublishedDate                       string  `json:"PublishedDate"`
		StartDate                           string  `json:"StartDate"`
		EndDate                             string  `json:"EndDate"`
		Section                             string  `json:"Section"`
		ProprietarySales                    float64 `json:"ProprietarySales"`
		ProprietaryPurchases                float64 `json:"ProprietaryPurchases"`
		ProprietaryTotal                    float64 `json:"ProprietaryTotal"`
		ProprietaryBalance                  float64 `json:"ProprietaryBalance"`
		BrokerageSales                      float64 `json:"BrokerageSales"`
		BrokeragePurchases                  float64 `json:"BrokeragePurchases"`
		BrokerageTotal                      float64 `json:"BrokerageTotal"`
		BrokerageBalance                    float64 `json:"BrokerageBalance"`
		TotalSales                          float64 `json:"TotalSales"`
		TotalPurchases                      float64 `json:"TotalPurchases"`
		TotalTotal                          float64 `json:"TotalTotal"`
		TotalBalance                        float64 `json:"TotalBalance"`
		IndividualsSales                    float64 `json:"IndividualsSales"`
		IndividualsPurchases                float64 `json:"IndividualsPurchases"`
		IndividualsTotal                    float64 `json:"IndividualsTotal"`
		IndividualsBalance                  float64 `json:"IndividualsBalance"`
		ForeignersSales                     float64 `json:"ForeignersSales"`
		ForeignersPurchases                 float64 `json:"ForeignersPurchases"`
		ForeignersTotal                     float64 `json:"ForeignersTotal"`
		ForeignersBalance                   float64 `json:"ForeignersBalance"`
		SecuritiesCosSales                  float64 `json:"SecuritiesCosSales"`
		SecuritiesCosPurchases              float64 `json:"SecuritiesCosPurchases"`
		SecuritiesCosTotal                  float64 `json:"SecuritiesCosTotal"`
		SecuritiesCosBalance                float64 `json:"SecuritiesCosBalance"`
		InvestmentTrustsSales               float64 `json:"InvestmentTrustsSales"`
		InvestmentTrustsPurchases           float64 `json:"InvestmentTrustsPurchases"`
		InvestmentTrustsTotal               float64 `json:"InvestmentTrustsTotal"`
		InvestmentTrustsBalance             float64 `json:"InvestmentTrustsBalance"`
		BusinessCosSales                    float64 `json:"BusinessCosSales"`
		BusinessCosPurchases                float64 `json:"BusinessCosPurchases"`
		BusinessCosTotal                    float64 `json:"BusinessCosTotal"`
		BusinessCosBalance                  float64 `json:"BusinessCosBalance"`
		OtherCosSales                       float64 `json:"OtherCosSales"`
		OtherCosPurchases                   float64 `json:"OtherCosPurchases"`
		OtherCosTotal                       float64 `json:"OtherCosTotal"`
		OtherCosBalance                     float64 `json:"OtherCosBalance"`
		InsuranceCosSales                   float64 `json:"InsuranceCosSales"`
		InsuranceCosPurchases               float64 `json:"InsuranceCosPurchases"`
		InsuranceCosTotal                   float64 `json:"InsuranceCosTotal"`
		InsuranceCosBalance                 float64 `json:"InsuranceCosBalance"`
		CityBKsRegionalBKsEtcSales          float64 `json:"CityBKsRegionalBKsEtcSales"`
		CityBKsRegionalBKsEtcPurchases      float64 `json:"CityBKsRegionalBKsEtcPurchases"`
		CityBKsRegionalBKsEtcTotal          float64 `json:"CityBKsRegionalBKsEtcTotal"`
		CityBKsRegionalBKsEtcBalance        float64 `json:"CityBKsRegionalBKsEtcBalance"`
		TrustBanksSales                     float64 `json:"TrustBanksSales"`
		TrustBanksPurchases                 float64 `json:"TrustBanksPurchases"`
		TrustBanksTotal                     float64 `json:"TrustBanksTotal"`
		TrustBanksBalance                   float64 `json:"TrustBanksBalance"`
		OtherFinancialInstitutionsSales     float64 `json:"OtherFinancialInstitutionsSales"`
		OtherFinancialInstitutionsPurchases float64 `json:"OtherFinancialInstitutionsPurchases"`
		OtherFinancialInstitutionsTotal     float64 `json:"OtherFinancialInstitutionsTotal"`
		OtherFinancialInstitutionsBalance   float64 `json:"OtherFinancialInstitutionsBalance"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	stv.PublishedDate = raw.PublishedDate
	stv.StartDate = raw.StartDate
	stv.EndDate = raw.EndDate
	stv.Section = raw.Section
	stv.Proprietary = TradingBalance{
		Sales:     int64(raw.ProprietarySales),
		Purchases: int64(raw.ProprietaryPurchases),
		Total:     int64(raw.ProprietaryTotal),
		Balance:   int64(raw.ProprietaryBalance),
	}
	stv.Brokerage = TradingBalance{
		Sales:     int64(raw.BrokerageSales),
		Purchases: int64(raw.BrokeragePurchases),
		Total:     int64(raw.BrokerageTotal),
		Balance:   int64(raw.BrokerageBalance),
	}
	stv.Net = TradingBalance{
		Sales:     int64(raw.TotalSales),
		Purchases: int64(raw.TotalPurchases),
		Total:     int64(raw.TotalTotal),
		Balance:   int64(raw.TotalBalance),
	}
	stv.Individuals = TradingBalance{
		Sales:     int64(raw.IndividualsSales),
		Purchases: int64(raw.IndividualsPurchases),
		Total:     int64(raw.IndividualsTotal),
		Balance:   int64(raw.IndividualsBalance),
	}
	stv.Foreigners = TradingBalance{
		Sales:     int64(raw.ForeignersSales),
		Purchases: int64(raw.ForeignersPurchases),
		Total:     int64(raw.ForeignersTotal),
		Balance:   int64(raw.ForeignersBalance),
	}
	stv.SecuritiesCos = TradingBalance{
		Sales:     int64(raw.SecuritiesCosSales),
		Purchases: int64(raw.SecuritiesCosPurchases),
		Total:     int64(raw.SecuritiesCosTotal),
		Balance:   int64(raw.SecuritiesCosBalance),
	}
	stv.InvestmentTrusts = TradingBalance{
		Sales:     int64(raw.InvestmentTrustsSales),
		Purchases: int64(raw.InvestmentTrustsPurchases),
		Total:     int64(raw.InvestmentTrustsTotal),
		Balance:   int64(raw.InvestmentTrustsBalance),
	}
	stv.BusinessCos = TradingBalance{
		Sales:     int64(raw.BusinessCosSales),
		Purchases: int64(raw.BusinessCosPurchases),
		Total:     int64(raw.BusinessCosTotal),
		Balance:   int64(raw.BusinessCosBalance),
	}
	stv.OtherCos = TradingBalance{
		Sales:     int64(raw.OtherCosSales),
		Purchases: int64(raw.OtherCosPurchases),
		Total:     int64(raw.OtherCosTotal),
		Balance:   int64(raw.OtherCosBalance),
	}
	stv.InsuranceCos = TradingBalance{
		Sales:     int64(raw.InsuranceCosSales),
		Purchases: int64(raw.InsuranceCosPurchases),
		Total:     int64(raw.InsuranceCosTotal),
		Balance:   int64(raw.InsuranceCosBalance),
	}
	stv.Banks = TradingBalance{
		Sales:     int64(raw.CityBKsRegionalBKsEtcSales),
		Purchases: int64(raw.CityBKsRegionalBKsEtcPurchases),
		Total:     int64(raw.CityBKsRegionalBKsEtcTotal),
		Balance:   int64(raw.CityBKsRegionalBKsEtcBalance),
	}
	stv.TrustBanks = TradingBalance{
		Sales:     int64(raw.TrustBanksSales),
		Purchases: int64(raw.TrustBanksPurchases),
		Total:     int64(raw.TrustBanksTotal),
		Balance:   int64(raw.TrustBanksBalance),
	}
	stv.OtherFinancialInstitutions = TradingBalance{
		Sales:     int64(raw.OtherFinancialInstitutionsSales),
		Purchases: int64(raw.OtherFinancialInstitutionsPurchases),
		Total:     int64(raw.OtherFinancialInstitutionsTotal),
		Balance:   int64(raw.OtherFinancialInstitutionsBalance),
	}
	return nil
}

type StockTradingValueRequest struct {
	Section *string
	From    *string
	To      *string
}

type stockTradingValueParameters struct {
	StockTradingValueRequest
	PaginationKey *string
}

func (p stockTradingValueParameters) values() (url.Values, error) {
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

type stockTradingValueResponse struct {
	Data          []StockTradingValue `json:"data"`
	PaginationKey *string             `json:"pagination_key"`
}

func (c *Client) sendStockTradingValueRequest(ctx context.Context, param stockTradingValueParameters) (stockTradingValueResponse, error) {
	var r stockTradingValueResponse
	resp, err := c.sendRequest(ctx, "/equities/investor-types", param)
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

// StockTradingValue provides trading by type of investors.
// https://jpx.gitbook.io/j-quants-en/api-reference/trades_spec
func (c *Client) StockTradingValue(ctx context.Context, req StockTradingValueRequest) ([]StockTradingValue, error) {
	var data = make([]StockTradingValue, 0)
	var paginationKey *string
	ctx, cancel := context.WithTimeout(ctx, c.LoopTimeout)
	defer cancel()
	for {
		params := stockTradingValueParameters{StockTradingValueRequest: req, PaginationKey: paginationKey}
		resp, err := c.sendStockTradingValueRequest(ctx, params)
		if err != nil {
			if errors.As(err, &InternalServerError{}) {
				slog.Warn("Retrying HTTP request", "error", err.Error())
				time.Sleep(c.RetryInterval)
				continue
			} else {
				return nil, fmt.Errorf("failed to send stock trading value request: %w", err)
			}
		}
		data = append(data, resp.Data...)
		paginationKey = resp.PaginationKey
		if resp.PaginationKey == nil {
			break
		}
	}
	return data, nil
}
