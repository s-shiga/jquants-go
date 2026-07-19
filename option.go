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
