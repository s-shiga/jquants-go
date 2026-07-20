package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// FuturesPrice represents daily price data for a futures contract.
// It includes prices for the whole day, morning session, night session, and
// day session, along with volume, open interest, settlement price, and
// contract metadata.
type FuturesPrice struct {
	// Date is the trading date in YYYY-MM-DD format (JSON key "Date").
	Date string
	// Code is the futures contract code (JSON key "Code").
	Code string
	// ProductCategory is the product category code, e.g. "TOPIXF" (JSON key "ProdCat").
	ProductCategory string
	// WholeDayOpen is the opening price for the whole trading day (JSON key "O").
	WholeDayOpen *json.Number
	// WholeDayHigh is the highest price for the whole trading day (JSON key "H").
	WholeDayHigh *json.Number
	// WholeDayLow is the lowest price for the whole trading day (JSON key "L").
	WholeDayLow *json.Number
	// WholeDayClose is the closing price for the whole trading day (JSON key "C").
	WholeDayClose *json.Number
	// MorningSessionOpen is the opening price for the morning session (JSON key "MO").
	MorningSessionOpen *json.Number
	// MorningSessionHigh is the highest price for the morning session (JSON key "MH").
	MorningSessionHigh *json.Number
	// MorningSessionLow is the lowest price for the morning session (JSON key "ML").
	MorningSessionLow *json.Number
	// MorningSessionClose is the closing price for the morning session (JSON key "MC").
	MorningSessionClose *json.Number
	// NightSessionOpen is the opening price for the night session (JSON key "EO").
	NightSessionOpen *json.Number
	// NightSessionHigh is the highest price for the night session (JSON key "EH").
	NightSessionHigh *json.Number
	// NightSessionLow is the lowest price for the night session (JSON key "EL").
	NightSessionLow *json.Number
	// NightSessionClose is the closing price for the night session (JSON key "EC").
	NightSessionClose *json.Number
	// DaySessionOpen is the opening price for the day session (JSON key "AO").
	DaySessionOpen *json.Number
	// DaySessionHigh is the highest price for the day session (JSON key "AH").
	DaySessionHigh *json.Number
	// DaySessionLow is the lowest price for the day session (JSON key "AL").
	DaySessionLow *json.Number
	// DaySessionClose is the closing price for the day session (JSON key "AC").
	DaySessionClose *json.Number
	// Volume is the total trading volume in contracts. Nil when no data is available (JSON key "Vo").
	Volume *int64
	// OpenInterest is the number of outstanding contracts. Nil when no data is available (JSON key "OI").
	OpenInterest *int64
	// TurnoverValue is the total trading value in yen. Nil when no data is available (JSON key "Va").
	TurnoverValue *int64
	// ContractMonth is the contract expiration month in YYYY-MM format (JSON key "CM").
	ContractMonth string
	// VolumeOnlyAuction is the volume from auction-only trades (JSON key "VoOA").
	VolumeOnlyAuction *int64
	// EmergencyMarginTriggerDivision indicates emergency margin status
	// ("001": triggered, "002": settlement price calculation) (JSON key "EmMrgnTrgDiv").
	EmergencyMarginTriggerDivision string
	// LastTradingDay is the last trading day for this contract (JSON key "LTD").
	LastTradingDay *string
	// SpecialQuotationDay is the special quotation day (SQ day). May be blank before 2016-07-19 (JSON key "SQD").
	SpecialQuotationDay *string
	// SettlementPrice is the daily settlement price (JSON key "Settle").
	SettlementPrice *json.Number
	// CentralContractMonthFlag indicates whether this is the central contract month
	// ("1": yes, "0": no). May be blank before 2016-07-19 (JSON key "CCMFlag").
	CentralContractMonthFlag string
}

// UnmarshalJSON decodes a single futures price record, translating the
// abbreviated J-Quants JSON keys into descriptive fields and normalizing
// numeric fields that may arrive as floats, strings, or null.
func (fp *FuturesPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date                           string `json:"Date"`
		Code                           string `json:"Code"`
		ProductCategory                string `json:"ProdCat"`
		WholeDayOpen                   any    `json:"O"`
		WholeDayHigh                   any    `json:"H"`
		WholeDayLow                    any    `json:"L"`
		WholeDayClose                  any    `json:"C"`
		MorningSessionOpen             any    `json:"MO"`
		MorningSessionHigh             any    `json:"MH"`
		MorningSessionLow              any    `json:"ML"`
		MorningSessionClose            any    `json:"MC"`
		NightSessionOpen               any    `json:"EO"`
		NightSessionHigh               any    `json:"EH"`
		NightSessionLow                any    `json:"EL"`
		NightSessionClose              any    `json:"EC"`
		DaySessionOpen                 any    `json:"AO"`
		DaySessionHigh                 any    `json:"AH"`
		DaySessionLow                  any    `json:"AL"`
		DaySessionClose                any    `json:"AC"`
		Volume                         any    `json:"Vo"`
		OpenInterest                   any    `json:"OI"`
		TurnoverValue                  any    `json:"Va"`
		ContractMonth                  string `json:"CM"`
		VolumeOnlyAuction              any    `json:"VoOA"`
		EmergencyMarginTriggerDivision string `json:"EmMrgnTrgDiv"`
		LastTradingDay                 string `json:"LTD"`
		SpecialQuotationDay            string `json:"SQD"`
		SettlementPrice                any    `json:"Settle"`
		CentralContractMonthFlag       string `json:"CCMFlag"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal futures price: %w", err)
	}

	u := &unmarshaler{}

	fp.Date = raw.Date
	fp.Code = raw.Code
	fp.ProductCategory = raw.ProductCategory
	fp.WholeDayOpen = u.jsonNumber(raw.WholeDayOpen)
	fp.WholeDayHigh = u.jsonNumber(raw.WholeDayHigh)
	fp.WholeDayLow = u.jsonNumber(raw.WholeDayLow)
	fp.WholeDayClose = u.jsonNumber(raw.WholeDayClose)
	fp.MorningSessionOpen = u.jsonNumber(raw.MorningSessionOpen)
	fp.MorningSessionHigh = u.jsonNumber(raw.MorningSessionHigh)
	fp.MorningSessionLow = u.jsonNumber(raw.MorningSessionLow)
	fp.MorningSessionClose = u.jsonNumber(raw.MorningSessionClose)
	fp.NightSessionOpen = u.jsonNumber(raw.NightSessionOpen)
	fp.NightSessionHigh = u.jsonNumber(raw.NightSessionHigh)
	fp.NightSessionLow = u.jsonNumber(raw.NightSessionLow)
	fp.NightSessionClose = u.jsonNumber(raw.NightSessionClose)
	fp.DaySessionOpen = u.jsonNumber(raw.DaySessionOpen)
	fp.DaySessionHigh = u.jsonNumber(raw.DaySessionHigh)
	fp.DaySessionLow = u.jsonNumber(raw.DaySessionLow)
	fp.DaySessionClose = u.jsonNumber(raw.DaySessionClose)
	fp.Volume = u.volume(raw.Volume)
	fp.OpenInterest = u.volume(raw.OpenInterest)
	fp.TurnoverValue = u.volume(raw.TurnoverValue)
	fp.ContractMonth = raw.ContractMonth
	fp.VolumeOnlyAuction = u.volume(raw.VolumeOnlyAuction)
	fp.EmergencyMarginTriggerDivision = raw.EmergencyMarginTriggerDivision
	fp.LastTradingDay = nilIfEmpty(raw.LastTradingDay)
	fp.SpecialQuotationDay = nilIfEmpty(raw.SpecialQuotationDay)
	fp.SettlementPrice = u.jsonNumber(raw.SettlementPrice)
	fp.CentralContractMonthFlag = raw.CentralContractMonthFlag

	return u.err
}

// FuturesPriceRequest specifies filter parameters for the FuturesPrice API.
type FuturesPriceRequest struct {
	// Date is the trading date to query in YYYY-MM-DD or YYYYMMDD format. Required.
	Date string
	// Category filters by product category code, e.g. "TOPIXF". Optional.
	Category *string
	// ContractFlag filters by central contract month flag ("1" for central month only). Optional.
	ContractFlag *string
}

type futuresPriceParameters struct {
	FuturesPriceRequest
	PaginationKey *string
}

func (p futuresPriceParameters) values() (url.Values, error) {
	v := url.Values{}
	v.Add("date", p.Date)
	if p.Category != nil {
		v.Add("category", *p.Category)
	}
	if p.ContractFlag != nil {
		v.Add("contract_flag", *p.ContractFlag)
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type futuresPriceResponse struct {
	Data          []FuturesPrice `json:"data"`
	PaginationKey *string        `json:"pagination_key"`
}

func (r futuresPriceResponse) Items() []FuturesPrice { return r.Data }
func (r futuresPriceResponse) NextPageKey() *string  { return r.PaginationKey }

// FuturesPrice retrieves futures prices from the /derivatives/bars/daily/futures endpoint.
// It automatically handles pagination to fetch all matching records.
// This endpoint requires the J-Quants Premium plan.
// See https://jpx-jquants.com/en/spec/drv-bars-daily-fut for details.
func (c *Client) FuturesPrice(ctx context.Context, req FuturesPriceRequest) ([]FuturesPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (futuresPriceResponse, error) {
		params := futuresPriceParameters{FuturesPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[futuresPriceResponse](ctx, c, "/derivatives/bars/daily/futures", params)
	})
}

// FuturesPriceWithChannel retrieves futures prices and streams each record to the provided channel.
// The channel is closed when all records have been sent or an error occurs.
// On error the channel is closed and the error is returned from this method, so callers
// must check the returned error after the channel closes; ranging the channel alone will not surface it.
// This endpoint requires the J-Quants Premium plan.
// See https://jpx-jquants.com/en/spec/drv-bars-daily-fut for details.
func (c *Client) FuturesPriceWithChannel(ctx context.Context, req FuturesPriceRequest, ch chan<- FuturesPrice) error {
	return fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, paginationKey *string) (futuresPriceResponse, error) {
		params := futuresPriceParameters{FuturesPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[futuresPriceResponse](ctx, c, "/derivatives/bars/daily/futures", params)
	})
}
