package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// IndexOptionPrice represents daily price data for Nikkei 225 index options.
// It includes prices for whole day, night session, and day session, along with
// volume, open interest, and option Greeks.
type IndexOptionPrice struct {
	// Date is the trading date in YYYY-MM-DD format.
	Date string
	// Code is the option contract code.
	Code string
	// WholeDayOpen is the opening price for the whole trading day.
	WholeDayOpen *int32
	// WholeDayHigh is the highest price for the whole trading day.
	WholeDayHigh *int32
	// WholeDayLow is the lowest price for the whole trading day.
	WholeDayLow *int32
	// WholeDayClose is the closing price for the whole trading day.
	WholeDayClose *int32
	// NightSessionOpen is the opening price for the night session.
	NightSessionOpen *int32
	// NightSessionHigh is the highest price for the night session.
	NightSessionHigh *int32
	// NightSessionLow is the lowest price for the night session.
	NightSessionLow *int32
	// NightSessionClose is the closing price for the night session.
	NightSessionClose *int32
	// DaySessionOpen is the opening price for the day session.
	DaySessionOpen *int32
	// DaySessionHigh is the highest price for the day session.
	DaySessionHigh *int32
	// DaySessionLow is the lowest price for the day session.
	DaySessionLow *int32
	// DaySessionClose is the closing price for the day session.
	DaySessionClose *int32
	// Volume is the total trading volume in contracts. Nil when no data is available.
	Volume *int64
	// OpenInterest is the number of outstanding contracts. Nil when no data is available.
	OpenInterest *int64
	// TurnoverValue is the total trading value in yen. Nil when no data is available.
	TurnoverValue *int64
	// ContractMonth is the contract expiration month in YYYYMM format.
	ContractMonth string
	// StrikePrice is the option strike price.
	StrikePrice int32
	// VolumeOnlyAuction is the volume from auction-only trades.
	VolumeOnlyAuction *int64
	// EmergencyMarginTriggerDivision indicates emergency margin status.
	EmergencyMarginTriggerDivision string
	// PutCallDivision indicates the option type (1: Put, 2: Call).
	PutCallDivision int8
	// LastTradingDay is the last trading day for this contract.
	LastTradingDay *string
	// SpecialQuotationDay is the special quotation day (SQ day).
	SpecialQuotationDay *string
	// SettlementPrice is the daily settlement price.
	SettlementPrice *int32
	// TheoreticalPrice is the theoretical option price.
	TheoreticalPrice *json.Number
	// BaseVolatility is the base volatility used for theoretical price calculation.
	BaseVolatility *json.Number
	// UnderlyingPrice is the price of the underlying index.
	UnderlyingPrice *json.Number
	// ImpliedVolatility is the implied volatility derived from market price.
	ImpliedVolatility *json.Number
	// InterestRate is the interest rate used for pricing.
	InterestRate *json.Number
}

// unmarshaler accumulates errors during unmarshaling, allowing cleaner code flow.
type unmarshaler struct {
	err error
}

func (u *unmarshaler) price(v any) *int32 {
	if u.err != nil {
		return nil
	}
	result, err := unmarshalPrice(v)
	u.err = err
	return result
}

func (u *unmarshaler) volume(v any) *int64 {
	if u.err != nil {
		return nil
	}
	result, err := unmarshalVolume(v)
	u.err = err
	return result
}

func (u *unmarshaler) jsonNumber(v any) *json.Number {
	if u.err != nil {
		return nil
	}
	result, err := unmarshalJSONNumber(v)
	u.err = err
	return result
}

func (iop *IndexOptionPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date                           string  `json:"Date"`
		Code                           string  `json:"Code"`
		WholeDayOpen                   any     `json:"O"`
		WholeDayHigh                   any     `json:"H"`
		WholeDayLow                    any     `json:"L"`
		WholeDayClose                  any     `json:"C"`
		NightSessionOpen               any     `json:"EO"`
		NightSessionHigh               any     `json:"EH"`
		NightSessionLow                any     `json:"EL"`
		NightSessionClose              any     `json:"EC"`
		DaySessionOpen                 any     `json:"AO"`
		DaySessionHigh                 any     `json:"AH"`
		DaySessionLow                  any     `json:"AL"`
		DaySessionClose                any     `json:"AC"`
		Volume                         any     `json:"Vo"`
		OpenInterest                   any     `json:"OI"`
		TurnoverValue                  any     `json:"Va"`
		ContractMonth                  string  `json:"CM"`
		StrikePrice                    float64 `json:"Strike"`
		VolumeOnlyAuction              any     `json:"VoOA"`
		EmergencyMarginTriggerDivision string  `json:"EmMrgnTrgDiv"`
		PutCallDivision                string  `json:"PCDiv"`
		LastTradingDay                 string  `json:"LTD"`
		SpecialQuotationDay            string  `json:"SQD"`
		SettlementPrice                any     `json:"Settle"`
		TheoreticalPrice               any     `json:"Theo"`
		BaseVolatility                 any     `json:"BaseVol"`
		UnderlyingPrice                any     `json:"UnderPx"`
		ImpliedVolatility              any     `json:"IV"`
		InterestRate                   any     `json:"IR"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal index option price: %w", err)
	}
	putCallDivision, err := strconv.ParseInt(raw.PutCallDivision, 10, 8)
	if err != nil {
		return fmt.Errorf("failed to parse put/call division: %w", err)
	}

	u := &unmarshaler{}

	iop.Date = raw.Date
	iop.Code = raw.Code
	iop.WholeDayOpen = u.price(raw.WholeDayOpen)
	iop.WholeDayHigh = u.price(raw.WholeDayHigh)
	iop.WholeDayLow = u.price(raw.WholeDayLow)
	iop.WholeDayClose = u.price(raw.WholeDayClose)
	iop.NightSessionOpen = u.price(raw.NightSessionOpen)
	iop.NightSessionHigh = u.price(raw.NightSessionHigh)
	iop.NightSessionLow = u.price(raw.NightSessionLow)
	iop.NightSessionClose = u.price(raw.NightSessionClose)
	iop.DaySessionOpen = u.price(raw.DaySessionOpen)
	iop.DaySessionHigh = u.price(raw.DaySessionHigh)
	iop.DaySessionLow = u.price(raw.DaySessionLow)
	iop.DaySessionClose = u.price(raw.DaySessionClose)
	iop.Volume = u.volume(raw.Volume)
	iop.OpenInterest = u.volume(raw.OpenInterest)
	iop.TurnoverValue = u.volume(raw.TurnoverValue)
	iop.ContractMonth = raw.ContractMonth
	iop.StrikePrice = int32(raw.StrikePrice)
	iop.VolumeOnlyAuction = u.volume(raw.VolumeOnlyAuction)
	iop.EmergencyMarginTriggerDivision = raw.EmergencyMarginTriggerDivision
	iop.PutCallDivision = int8(putCallDivision)
	iop.LastTradingDay = nilIfEmpty(raw.LastTradingDay)
	iop.SpecialQuotationDay = nilIfEmpty(raw.SpecialQuotationDay)
	iop.SettlementPrice = u.price(raw.SettlementPrice)
	iop.TheoreticalPrice = u.jsonNumber(raw.TheoreticalPrice)
	iop.BaseVolatility = u.jsonNumber(raw.BaseVolatility)
	iop.UnderlyingPrice = u.jsonNumber(raw.UnderlyingPrice)
	iop.ImpliedVolatility = u.jsonNumber(raw.ImpliedVolatility)
	iop.InterestRate = u.jsonNumber(raw.InterestRate)

	return u.err
}

func unmarshalPrice(value any) (*int32, error) {
	switch v := value.(type) {
	case float64:
		i := int32(v)
		return &i, nil
	case string:
		return nil, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unmarshalPrice: unknown type %T", v)
	}
}

func unmarshalVolume(value any) (*int64, error) {
	switch v := value.(type) {
	case float64:
		i := int64(v)
		return &i, nil
	case string:
		return nil, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unmarshalVolume: unknown type %T", v)
	}
}

func unmarshalJSONNumber(value any) (*json.Number, error) {
	switch v := value.(type) {
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		n := json.Number(s)
		return &n, nil
	case string:
		return nil, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unmarshalJSONNumber: unknown type %T", v)
	}
}

func nilIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// IndexOptionPriceRequest specifies filter parameters for the IndexOptionPrice API.
type IndexOptionPriceRequest struct {
	// Date is the trading date to query in YYYY-MM-DD format. Required.
	Date string
}

type indexOptionPriceParameters struct {
	IndexOptionPriceRequest
	PaginationKey *string
}

func (p indexOptionPriceParameters) values() (url.Values, error) {
	v := url.Values{}
	v.Add("date", p.Date)
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type indexOptionPriceResponse struct {
	Data          []IndexOptionPrice `json:"data"`
	PaginationKey *string            `json:"pagination_key"`
}

func (r indexOptionPriceResponse) Items() []IndexOptionPrice { return r.Data }
func (r indexOptionPriceResponse) NextPageKey() *string      { return r.PaginationKey }

// IndexOptionPrice retrieves Nikkei 225 index option prices from the /derivatives/bars/daily/options/225 endpoint.
// It automatically handles pagination to fetch all matching records.
func (c *Client) IndexOptionPrice(ctx context.Context, req IndexOptionPriceRequest) ([]IndexOptionPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (indexOptionPriceResponse, error) {
		params := indexOptionPriceParameters{IndexOptionPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[indexOptionPriceResponse](ctx, c, "/derivatives/bars/daily/options/225", params)
	})
}

// IndexOptionPriceWithChannel retrieves Nikkei 225 index option prices and streams each record to the provided channel.
// The channel is closed when all records have been sent or an error occurs.
// On error the channel is closed and the error is returned from this method, so callers
// must check the returned error after the channel closes; ranging the channel alone will not surface it.
func (c *Client) IndexOptionPriceWithChannel(ctx context.Context, req IndexOptionPriceRequest, ch chan<- IndexOptionPrice) error {
	return fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, paginationKey *string) (indexOptionPriceResponse, error) {
		params := indexOptionPriceParameters{IndexOptionPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[indexOptionPriceResponse](ctx, c, "/derivatives/bars/daily/options/225", params)
	})
}

// OptionPrice represents daily price data for an options contract across all underlyings.
// It includes prices for the whole day, morning session, night session, and day session,
// along with volume, open interest, settlement price, option Greeks, and contract metadata.
type OptionPrice struct {
	// Date is the trading date in YYYY-MM-DD format (JSON key "Date").
	Date string
	// Code is the option contract code (JSON key "Code").
	Code string
	// ProductCategory is the product category code, e.g. "TOPIXE" or "NK225E" (JSON key "ProdCat").
	ProductCategory string
	// UnderlyingSSO is the underlying securities code for securities options, or "-" for non-securities options (JSON key "UndSSO").
	UnderlyingSSO string
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
	// StrikePrice is the option strike price (JSON key "Strike").
	StrikePrice float64
	// VolumeOnlyAuction is the volume from auction-only trades (JSON key "VoOA").
	VolumeOnlyAuction *int64
	// EmergencyMarginTriggerDivision indicates emergency margin status
	// ("001": triggered, "002": settlement price calculation) (JSON key "EmMrgnTrgDiv").
	EmergencyMarginTriggerDivision string
	// PutCallDivision indicates the option type (1: Put, 2: Call) (JSON key "PCDiv").
	PutCallDivision int8
	// LastTradingDay is the last trading day for this contract (JSON key "LTD").
	LastTradingDay *string
	// SpecialQuotationDay is the special quotation day (SQ day) (JSON key "SQD").
	SpecialQuotationDay *string
	// SettlementPrice is the daily settlement price (JSON key "Settle").
	SettlementPrice *json.Number
	// TheoreticalPrice is the theoretical option price (JSON key "Theo").
	TheoreticalPrice *json.Number
	// BaseVolatility is the base volatility used for theoretical price calculation (JSON key "BaseVol").
	BaseVolatility *json.Number
	// UnderlyingPrice is the price of the underlying asset (JSON key "UnderPx").
	UnderlyingPrice *json.Number
	// ImpliedVolatility is the implied volatility derived from market price (JSON key "IV").
	ImpliedVolatility *json.Number
	// InterestRate is the interest rate used for pricing (JSON key "IR").
	InterestRate *json.Number
	// CentralContractMonthFlag indicates whether this is the central contract month
	// ("1": yes, "0": no) (JSON key "CCMFlag").
	CentralContractMonthFlag string
}

// UnmarshalJSON decodes a single option price record, translating the
// abbreviated J-Quants JSON keys into descriptive fields and normalizing
// numeric fields that may arrive as floats, strings, or null.
func (op *OptionPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date                           string  `json:"Date"`
		Code                           string  `json:"Code"`
		ProductCategory                string  `json:"ProdCat"`
		UnderlyingSSO                  string  `json:"UndSSO"`
		WholeDayOpen                   any     `json:"O"`
		WholeDayHigh                   any     `json:"H"`
		WholeDayLow                    any     `json:"L"`
		WholeDayClose                  any     `json:"C"`
		MorningSessionOpen             any     `json:"MO"`
		MorningSessionHigh             any     `json:"MH"`
		MorningSessionLow              any     `json:"ML"`
		MorningSessionClose            any     `json:"MC"`
		NightSessionOpen               any     `json:"EO"`
		NightSessionHigh               any     `json:"EH"`
		NightSessionLow                any     `json:"EL"`
		NightSessionClose              any     `json:"EC"`
		DaySessionOpen                 any     `json:"AO"`
		DaySessionHigh                 any     `json:"AH"`
		DaySessionLow                  any     `json:"AL"`
		DaySessionClose                any     `json:"AC"`
		Volume                         any     `json:"Vo"`
		OpenInterest                   any     `json:"OI"`
		TurnoverValue                  any     `json:"Va"`
		ContractMonth                  string  `json:"CM"`
		StrikePrice                    float64 `json:"Strike"`
		VolumeOnlyAuction              any     `json:"VoOA"`
		EmergencyMarginTriggerDivision string  `json:"EmMrgnTrgDiv"`
		PutCallDivision                string  `json:"PCDiv"`
		LastTradingDay                 string  `json:"LTD"`
		SpecialQuotationDay            string  `json:"SQD"`
		SettlementPrice                any     `json:"Settle"`
		TheoreticalPrice               any     `json:"Theo"`
		BaseVolatility                 any     `json:"BaseVol"`
		UnderlyingPrice                any     `json:"UnderPx"`
		ImpliedVolatility              any     `json:"IV"`
		InterestRate                   any     `json:"IR"`
		CentralContractMonthFlag       string  `json:"CCMFlag"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal option price: %w", err)
	}
	putCallDivision, err := strconv.ParseInt(raw.PutCallDivision, 10, 8)
	if err != nil {
		return fmt.Errorf("failed to parse put/call division: %w", err)
	}

	u := &unmarshaler{}

	op.Date = raw.Date
	op.Code = raw.Code
	op.ProductCategory = raw.ProductCategory
	op.UnderlyingSSO = raw.UnderlyingSSO
	op.WholeDayOpen = u.jsonNumber(raw.WholeDayOpen)
	op.WholeDayHigh = u.jsonNumber(raw.WholeDayHigh)
	op.WholeDayLow = u.jsonNumber(raw.WholeDayLow)
	op.WholeDayClose = u.jsonNumber(raw.WholeDayClose)
	op.MorningSessionOpen = u.jsonNumber(raw.MorningSessionOpen)
	op.MorningSessionHigh = u.jsonNumber(raw.MorningSessionHigh)
	op.MorningSessionLow = u.jsonNumber(raw.MorningSessionLow)
	op.MorningSessionClose = u.jsonNumber(raw.MorningSessionClose)
	op.NightSessionOpen = u.jsonNumber(raw.NightSessionOpen)
	op.NightSessionHigh = u.jsonNumber(raw.NightSessionHigh)
	op.NightSessionLow = u.jsonNumber(raw.NightSessionLow)
	op.NightSessionClose = u.jsonNumber(raw.NightSessionClose)
	op.DaySessionOpen = u.jsonNumber(raw.DaySessionOpen)
	op.DaySessionHigh = u.jsonNumber(raw.DaySessionHigh)
	op.DaySessionLow = u.jsonNumber(raw.DaySessionLow)
	op.DaySessionClose = u.jsonNumber(raw.DaySessionClose)
	op.Volume = u.volume(raw.Volume)
	op.OpenInterest = u.volume(raw.OpenInterest)
	op.TurnoverValue = u.volume(raw.TurnoverValue)
	op.ContractMonth = raw.ContractMonth
	op.StrikePrice = raw.StrikePrice
	op.VolumeOnlyAuction = u.volume(raw.VolumeOnlyAuction)
	op.EmergencyMarginTriggerDivision = raw.EmergencyMarginTriggerDivision
	op.PutCallDivision = int8(putCallDivision)
	op.LastTradingDay = nilIfEmpty(raw.LastTradingDay)
	op.SpecialQuotationDay = nilIfEmpty(raw.SpecialQuotationDay)
	op.SettlementPrice = u.jsonNumber(raw.SettlementPrice)
	op.TheoreticalPrice = u.jsonNumber(raw.TheoreticalPrice)
	op.BaseVolatility = u.jsonNumber(raw.BaseVolatility)
	op.UnderlyingPrice = u.jsonNumber(raw.UnderlyingPrice)
	op.ImpliedVolatility = u.jsonNumber(raw.ImpliedVolatility)
	op.InterestRate = u.jsonNumber(raw.InterestRate)
	op.CentralContractMonthFlag = raw.CentralContractMonthFlag

	return u.err
}

// OptionPriceRequest specifies filter parameters for the OptionPrice API.
type OptionPriceRequest struct {
	// Date is the trading date to query in YYYY-MM-DD or YYYYMMDD format. Required.
	Date string
	// Category filters by product category code, e.g. "TOPIXE" or "NK225E". Optional.
	Category *string
	// Code filters by underlying securities code, used with the securities options category. Optional.
	Code *string
	// ContractFlag filters by central contract month flag ("1" for central month only). Optional.
	ContractFlag *string
}

type optionPriceParameters struct {
	OptionPriceRequest
	PaginationKey *string
}

func (p optionPriceParameters) values() (url.Values, error) {
	v := url.Values{}
	v.Add("date", p.Date)
	if p.Category != nil {
		v.Add("category", *p.Category)
	}
	if p.Code != nil {
		v.Add("code", *p.Code)
	}
	if p.ContractFlag != nil {
		v.Add("contract_flag", *p.ContractFlag)
	}
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type optionPriceResponse struct {
	Data          []OptionPrice `json:"data"`
	PaginationKey *string       `json:"pagination_key"`
}

func (r optionPriceResponse) Items() []OptionPrice { return r.Data }
func (r optionPriceResponse) NextPageKey() *string { return r.PaginationKey }

// OptionPrice retrieves option prices for all underlyings from the /derivatives/bars/daily/options endpoint.
// It automatically handles pagination to fetch all matching records.
// This endpoint requires the J-Quants Premium plan.
// See https://jpx-jquants.com/en/spec/drv-bars-daily-opt for details.
func (c *Client) OptionPrice(ctx context.Context, req OptionPriceRequest) ([]OptionPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (optionPriceResponse, error) {
		params := optionPriceParameters{OptionPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[optionPriceResponse](ctx, c, "/derivatives/bars/daily/options", params)
	})
}

// OptionPriceWithChannel retrieves option prices for all underlyings and streams each record to the provided channel.
// The channel is closed when all records have been sent or an error occurs.
// On error the channel is closed and the error is returned from this method, so callers
// must check the returned error after the channel closes; ranging the channel alone will not surface it.
// This endpoint requires the J-Quants Premium plan.
// See https://jpx-jquants.com/en/spec/drv-bars-daily-opt for details.
func (c *Client) OptionPriceWithChannel(ctx context.Context, req OptionPriceRequest, ch chan<- OptionPrice) error {
	return fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, paginationKey *string) (optionPriceResponse, error) {
		params := optionPriceParameters{OptionPriceRequest: req, PaginationKey: paginationKey}
		return getJSON[optionPriceResponse](ctx, c, "/derivatives/bars/daily/options", params)
	})
}
